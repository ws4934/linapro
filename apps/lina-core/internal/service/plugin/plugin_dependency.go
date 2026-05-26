// This file connects the side-effect-free dependency resolver to plugin
// lifecycle orchestration, API projections, and upgrade validation.

package plugin

import (
	"context"
	"strings"

	"lina-core/internal/model/entity"
	"lina-core/internal/service/plugin/internal/catalog"
	plugindep "lina-core/internal/service/plugin/internal/dependency"
	"lina-core/pkg/bizerr"
)

type (
	// DependencyFrameworkCheck exposes framework compatibility for management clients.
	DependencyFrameworkCheck struct {
		// RequiredVersion is the framework semantic-version range declared by the plugin.
		RequiredVersion string
		// CurrentVersion is the current LinaPro framework version.
		CurrentVersion string
		// Status is the compatibility state returned by the resolver.
		Status string
	}

	// DependencyPluginCheck exposes one plugin dependency edge.
	DependencyPluginCheck struct {
		// OwnerID is the plugin that declares the dependency.
		OwnerID string
		// DependencyID is the depended-on plugin identifier.
		DependencyID string
		// DependencyName is the depended-on plugin display name when known.
		DependencyName string
		// RequiredVersion is the declared dependency version range.
		RequiredVersion string
		// CurrentVersion is the discovered or installed dependency version.
		CurrentVersion string
		// Installed reports whether the dependency is already installed.
		Installed bool
		// Discovered reports whether the dependency is discoverable.
		Discovered bool
		// Status is the dependency state returned by the resolver.
		Status string
		// Chain is the dependency chain leading to this edge.
		Chain []string
	}

	// DependencyBlocker exposes one hard dependency failure.
	DependencyBlocker struct {
		// Code identifies the dependency failure category.
		Code string
		// PluginID is the plugin whose lifecycle is blocked.
		PluginID string
		// DependencyID is the dependency plugin when applicable.
		DependencyID string
		// RequiredVersion is the declared version range when applicable.
		RequiredVersion string
		// CurrentVersion is the observed version when applicable.
		CurrentVersion string
		// Chain is the dependency chain associated with the blocker.
		Chain []string
		// Detail is a concise operator diagnostic.
		Detail string
	}

	// DependencyReverseDependent exposes one installed downstream hard dependency.
	DependencyReverseDependent struct {
		// PluginID is the downstream plugin ID.
		PluginID string
		// Name is the downstream plugin display name.
		Name string
		// Version is the downstream plugin version.
		Version string
		// RequiredVersion is the target version range declared by the downstream plugin.
		RequiredVersion string
	}

	// DependencyCheckResult is the management-facing dependency status snapshot.
	DependencyCheckResult struct {
		// TargetID is the checked plugin.
		TargetID string
		// Framework contains the framework compatibility result.
		Framework DependencyFrameworkCheck
		// Dependencies contains direct and transitive dependency edge checks.
		Dependencies []*DependencyPluginCheck
		// Blockers lists install-side hard failures.
		Blockers []*DependencyBlocker
		// Cycle contains the first detected dependency cycle.
		Cycle []string
		// ReverseDependents lists installed downstream hard dependencies.
		ReverseDependents []*DependencyReverseDependent
		// ReverseBlockers lists uninstall or downstream-version blockers.
		ReverseBlockers []*DependencyBlocker
	}

	// dependencyInstallContext records automatic install state for one request.
	dependencyInstallContext struct {
		// active marks target IDs already being installed in this request.
		active map[string]bool
	}

	// dependencySnapshotCache stores request-local dependency snapshots for
	// repeated read-only dependency checks during one plugin list projection.
	dependencySnapshotCache struct {
		snapshots []*plugindep.PluginSnapshot
	}
)

// dependencyInstallContextKey stores request-local dependency orchestration state.
type dependencyInstallContextKey struct{}

// dependencySnapshotCacheContextKey stores request-local dependency snapshots.
type dependencySnapshotCacheContextKey struct{}

// WithDependencySnapshotCache returns a child context that can reuse dependency
// snapshots across repeated read-only dependency checks in one request.
func (s *serviceImpl) WithDependencySnapshotCache(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if dependencySnapshotCacheFromContext(ctx) != nil {
		return ctx
	}
	return context.WithValue(ctx, dependencySnapshotCacheContextKey{}, &dependencySnapshotCache{})
}

// CheckPluginDependencies evaluates install and uninstall dependency status for one plugin.
func (s *serviceImpl) CheckPluginDependencies(ctx context.Context, pluginID string) (*DependencyCheckResult, error) {
	installResult, err := s.resolveInstallDependencies(ctx, pluginID)
	if err != nil {
		return nil, err
	}
	reverseResult, err := s.resolveReverseDependencies(ctx, pluginID, "")
	if err != nil {
		return nil, err
	}
	result := toDependencyCheckResult(installResult)
	result.ReverseDependents = toDependencyReverseDependents(reverseResult.Dependents)
	result.ReverseBlockers = toDependencyBlockers(reverseResult.Blockers)
	return result, nil
}

// prepareInstallDependencies verifies a target before lifecycle side effects.
func (s *serviceImpl) prepareInstallDependencies(
	ctx context.Context,
	pluginID string,
	options InstallOptions,
) (*DependencyCheckResult, context.Context, error) {
	depCtx := dependencyContextFrom(ctx)
	normalizedID := strings.TrimSpace(pluginID)
	if normalizedID == "" {
		return nil, ctx, nil
	}
	if depCtx.active[normalizedID] {
		return nil, ctx, nil
	}

	depCtx.active[normalizedID] = true
	defer delete(depCtx.active, normalizedID)

	nextCtx := context.WithValue(ctx, dependencyInstallContextKey{}, depCtx)
	check, err := s.resolveInstallDependencies(nextCtx, normalizedID)
	if err != nil {
		return nil, nextCtx, err
	}
	result := toDependencyCheckResult(check)
	if hasDependencyBlockers(check.Blockers) {
		return result, nextCtx, s.buildDependencyBlockedError(normalizedID, check.Blockers)
	}
	return result, nextCtx, nil
}

// ensureNoReverseDependencies blocks uninstall when installed downstream plugins depend on target.
func (s *serviceImpl) ensureNoReverseDependencies(ctx context.Context, pluginID string) error {
	result, err := s.resolveReverseDependencies(ctx, pluginID, "")
	if err != nil {
		return err
	}
	if !hasDependencyBlockers(result.Blockers) {
		return nil
	}
	return s.buildReverseDependencyBlockedError(pluginID, result)
}

// ValidateSourcePluginUpgradeCandidate validates a source upgrade target before side effects.
func (s *serviceImpl) ValidateSourcePluginUpgradeCandidate(ctx context.Context, manifest *catalog.Manifest) error {
	return s.validateUpgradeCandidateDependencies(ctx, manifest)
}

// ValidateDynamicPluginCandidate validates a dynamic release candidate before side effects.
func (s *serviceImpl) ValidateDynamicPluginCandidate(ctx context.Context, manifest *catalog.Manifest) error {
	return s.validateUpgradeCandidateDependencies(ctx, manifest)
}

// validateUpgradeCandidateDependencies checks candidate dependencies and downstream version safety.
func (s *serviceImpl) validateUpgradeCandidateDependencies(ctx context.Context, manifest *catalog.Manifest) error {
	if manifest == nil {
		return nil
	}
	installResult, err := s.resolveInstallDependenciesForManifest(ctx, manifest)
	if err != nil {
		return err
	}
	if hasDependencyBlockers(installResult.Blockers) {
		return s.buildDependencyBlockedError(manifest.ID, installResult.Blockers)
	}

	if !s.dependencyTargetAlreadyInstalled(ctx, manifest.ID) {
		return nil
	}
	reverseResult, err := s.resolveReverseDependencies(ctx, manifest.ID, manifest.Version)
	if err != nil {
		return err
	}
	if hasDependencyBlockers(reverseResult.Blockers) {
		return s.buildReverseDependencyBlockedError(manifest.ID, reverseResult)
	}
	return nil
}

// resolveInstallDependencies evaluates dependency status for one discovered target.
func (s *serviceImpl) resolveInstallDependencies(ctx context.Context, pluginID string) (*plugindep.InstallCheckResult, error) {
	normalizedPluginID := strings.TrimSpace(pluginID)
	if manifest := manifestByIDFromContext(ctx, normalizedPluginID); manifest != nil {
		return s.resolveInstallDependenciesForManifest(ctx, manifest)
	}
	manifest, err := s.catalogSvc.GetDesiredManifest(normalizedPluginID)
	if err != nil {
		return nil, err
	}
	return s.resolveInstallDependenciesForManifest(ctx, manifest)
}

// resolveInstallDependenciesForManifest evaluates dependency status using a candidate manifest override.
func (s *serviceImpl) resolveInstallDependenciesForManifest(
	ctx context.Context,
	manifest *catalog.Manifest,
) (*plugindep.InstallCheckResult, error) {
	snapshots, err := s.buildDependencySnapshots(ctx, manifest)
	if err != nil {
		return nil, err
	}
	resolver := plugindep.New()
	return resolver.CheckInstall(plugindep.InstallCheckInput{
		TargetID:         strings.TrimSpace(manifest.ID),
		FrameworkVersion: s.frameworkVersion(ctx),
		Plugins:          snapshots,
	}), nil
}

// resolveReverseDependencies evaluates installed downstream dependencies for one target.
func (s *serviceImpl) resolveReverseDependencies(
	ctx context.Context,
	pluginID string,
	candidateVersion string,
) (*plugindep.ReverseCheckResult, error) {
	snapshots, err := s.buildDependencySnapshots(ctx, nil)
	if err != nil {
		return nil, err
	}
	resolver := plugindep.New()
	return resolver.CheckReverse(plugindep.ReverseCheckInput{
		TargetID:         strings.TrimSpace(pluginID),
		CandidateVersion: strings.TrimSpace(candidateVersion),
		Plugins:          snapshots,
	}), nil
}

// buildDependencySnapshots combines discovered manifests with installed registry release snapshots.
func (s *serviceImpl) buildDependencySnapshots(
	ctx context.Context,
	candidate *catalog.Manifest,
) ([]*plugindep.PluginSnapshot, error) {
	if candidate == nil {
		if cache := dependencySnapshotCacheFromContext(ctx); cache != nil && cache.snapshots != nil {
			return cloneDependencySnapshots(cache.snapshots), nil
		}
	}
	manifests := manifestSnapshotFromContext(ctx)
	if manifests == nil {
		var err error
		manifests, err = s.catalogSvc.ScanManifests()
		if err != nil {
			return nil, err
		}
	}
	snapshotByID := make(map[string]*plugindep.PluginSnapshot, len(manifests)+1)
	for _, manifest := range manifests {
		if manifest == nil || strings.TrimSpace(manifest.ID) == "" {
			continue
		}
		snapshotByID[manifest.ID] = &plugindep.PluginSnapshot{
			ID:           strings.TrimSpace(manifest.ID),
			Name:         strings.TrimSpace(manifest.Name),
			Version:      strings.TrimSpace(manifest.Version),
			Manifest:     manifest,
			Dependencies: catalog.CloneDependencySpec(manifest.Dependencies),
		}
	}
	if candidate != nil && strings.TrimSpace(candidate.ID) != "" {
		snapshotByID[candidate.ID] = &plugindep.PluginSnapshot{
			ID:           strings.TrimSpace(candidate.ID),
			Name:         strings.TrimSpace(candidate.Name),
			Version:      strings.TrimSpace(candidate.Version),
			Manifest:     candidate,
			Dependencies: catalog.CloneDependencySpec(candidate.Dependencies),
		}
	}

	readCtx, err := s.catalogSvc.WithStartupDataSnapshot(ctx)
	if err != nil {
		return nil, err
	}
	registries, err := s.catalogSvc.ListAllRegistries(readCtx)
	if err != nil {
		return nil, err
	}
	candidateID := ""
	if candidate != nil {
		candidateID = strings.TrimSpace(candidate.ID)
	}
	for _, registry := range registries {
		if registry == nil {
			continue
		}
		registryPluginID := strings.TrimSpace(registry.PluginId)
		if registryPluginID == "" {
			continue
		}
		snapshot := snapshotByID[registryPluginID]
		if snapshot == nil {
			if registry.ReleaseId <= 0 {
				continue
			}
			snapshot = &plugindep.PluginSnapshot{ID: registryPluginID}
			snapshotByID[registryPluginID] = snapshot
		}
		if registryPluginID == candidateID {
			snapshot.Installed = registry.Installed == catalog.InstalledYes
			continue
		}
		applyRegistryDependencySnapshot(readCtx, s.catalogSvc, snapshot, registry)
	}

	out := make([]*plugindep.PluginSnapshot, 0, len(snapshotByID))
	for _, snapshot := range snapshotByID {
		out = append(out, snapshot)
	}
	if candidate == nil {
		if cache := dependencySnapshotCacheFromContext(ctx); cache != nil {
			cache.snapshots = cloneDependencySnapshots(out)
		}
	}
	return out, nil
}

// manifestByIDFromContext returns a manifest from the request-local discovery
// snapshot without triggering another scan.
func manifestByIDFromContext(ctx context.Context, pluginID string) *catalog.Manifest {
	normalizedPluginID := strings.TrimSpace(pluginID)
	if normalizedPluginID == "" {
		return nil
	}
	for _, manifest := range manifestSnapshotFromContext(ctx) {
		if manifest != nil && strings.TrimSpace(manifest.ID) == normalizedPluginID {
			return manifest
		}
	}
	return nil
}

// applyRegistryDependencySnapshot prefers installed release snapshots for effective dependency metadata.
func applyRegistryDependencySnapshot(
	ctx context.Context,
	catalogSvc catalog.Service,
	snapshot *plugindep.PluginSnapshot,
	registry *entity.SysPlugin,
) {
	if snapshot == nil || registry == nil {
		return
	}
	if strings.TrimSpace(registry.Name) != "" {
		snapshot.Name = strings.TrimSpace(registry.Name)
	}
	if strings.TrimSpace(registry.Version) != "" {
		snapshot.Version = strings.TrimSpace(registry.Version)
	}
	snapshot.Installed = registry.Installed == catalog.InstalledYes
	if !snapshot.Installed {
		return
	}
	release, err := catalogSvc.GetRegistryRelease(ctx, registry)
	if err != nil || release == nil {
		snapshot.DependencySnapshotUnknown = true
		return
	}
	releaseSnapshot, err := catalogSvc.ParseManifestSnapshot(release.ManifestSnapshot)
	if err != nil || releaseSnapshot == nil {
		snapshot.DependencySnapshotUnknown = true
		return
	}
	if strings.TrimSpace(releaseSnapshot.Name) != "" {
		snapshot.Name = strings.TrimSpace(releaseSnapshot.Name)
	}
	if strings.TrimSpace(releaseSnapshot.Version) != "" {
		snapshot.Version = strings.TrimSpace(releaseSnapshot.Version)
	}
	snapshot.Dependencies = catalog.CloneDependencySpec(releaseSnapshot.Dependencies)
}

// frameworkVersion returns the current LinaPro framework version authority.
func (s *serviceImpl) frameworkVersion(ctx context.Context) string {
	if s == nil || s.configSvc == nil {
		return ""
	}
	metadata := s.configSvc.GetMetadata(ctx)
	if metadata == nil {
		return ""
	}
	return strings.TrimSpace(metadata.Framework.Version)
}

// dependencyContextFrom returns one request-local dependency orchestration context.
func dependencyContextFrom(ctx context.Context) *dependencyInstallContext {
	if ctx != nil {
		if value, ok := ctx.Value(dependencyInstallContextKey{}).(*dependencyInstallContext); ok && value != nil {
			if value.active == nil {
				value.active = make(map[string]bool)
			}
			return value
		}
	}
	return &dependencyInstallContext{active: make(map[string]bool)}
}

// dependencySnapshotCacheFromContext returns the request-local dependency
// snapshot cache, if the current read path enabled one.
func dependencySnapshotCacheFromContext(ctx context.Context) *dependencySnapshotCache {
	if ctx == nil {
		return nil
	}
	value, ok := ctx.Value(dependencySnapshotCacheContextKey{}).(*dependencySnapshotCache)
	if !ok {
		return nil
	}
	return value
}

// cloneDependencySnapshots returns a detached copy so callers cannot mutate the
// cached dependency snapshot slice for later checks in the same request.
func cloneDependencySnapshots(items []*plugindep.PluginSnapshot) []*plugindep.PluginSnapshot {
	out := make([]*plugindep.PluginSnapshot, 0, len(items))
	for _, item := range items {
		if item == nil {
			out = append(out, nil)
			continue
		}
		cloned := *item
		cloned.Dependencies = catalog.CloneDependencySpec(item.Dependencies)
		out = append(out, &cloned)
	}
	return out
}

// dependencyTargetAlreadyInstalled reports whether the target is already installed.
func (s *serviceImpl) dependencyTargetAlreadyInstalled(ctx context.Context, pluginID string) bool {
	registry, err := s.catalogSvc.GetRegistry(ctx, pluginID)
	if err != nil || registry == nil {
		return false
	}
	return registry.Installed == catalog.InstalledYes
}

// hasDependencyBlockers reports whether resolver blockers contain any hard failure.
func hasDependencyBlockers(blockers []*plugindep.Blocker) bool {
	return len(blockers) > 0
}

// buildDependencyBlockedError converts resolver blockers into one structured business error.
func (s *serviceImpl) buildDependencyBlockedError(pluginID string, blockers []*plugindep.Blocker) error {
	dependencyID, requiredVersion, currentVersion := firstDependencyBlockerFields(blockers)
	return bizerr.NewCode(
		CodePluginDependencyBlocked,
		bizerr.P("pluginId", strings.TrimSpace(pluginID)),
		bizerr.P("dependencyId", dependencyID),
		bizerr.P("requiredVersion", requiredVersion),
		bizerr.P("currentVersion", currentVersion),
		bizerr.P("chain", firstDependencyBlockerChain(blockers)),
		bizerr.P("blockers", formatDependencyBlockers(blockers)),
	)
}

// buildReverseDependencyBlockedError converts reverse dependency blockers into one structured error.
func (s *serviceImpl) buildReverseDependencyBlockedError(
	pluginID string,
	result *plugindep.ReverseCheckResult,
) error {
	dependents := toDependencyReverseDependents(result.Dependents)
	dependencyID, requiredVersion, currentVersion := firstDependencyBlockerFields(result.Blockers)
	return bizerr.NewCode(
		CodePluginReverseDependencyBlocked,
		bizerr.P("pluginId", strings.TrimSpace(pluginID)),
		bizerr.P("dependencyId", dependencyID),
		bizerr.P("requiredVersion", requiredVersion),
		bizerr.P("currentVersion", currentVersion),
		bizerr.P("dependents", strings.Join(reverseDependentIDs(dependents), ",")),
		bizerr.P("blockers", formatDependencyBlockers(result.Blockers)),
	)
}

// toDependencyCheckResult converts resolver install output into a service DTO.
func toDependencyCheckResult(result *plugindep.InstallCheckResult) *DependencyCheckResult {
	if result == nil {
		return &DependencyCheckResult{}
	}
	return &DependencyCheckResult{
		TargetID: strings.TrimSpace(result.TargetID),
		Framework: DependencyFrameworkCheck{
			RequiredVersion: result.Framework.RequiredVersion,
			CurrentVersion:  result.Framework.CurrentVersion,
			Status:          string(result.Framework.Status),
		},
		Dependencies: toDependencyPluginChecks(result.Dependencies),
		Blockers:     toDependencyBlockers(result.Blockers),
		Cycle:        cloneStringSlice(result.Cycle),
	}
}

// toDependencyPluginChecks converts resolver dependency edges.
func toDependencyPluginChecks(items []*plugindep.PluginDependencyCheck) []*DependencyPluginCheck {
	out := make([]*DependencyPluginCheck, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		out = append(out, &DependencyPluginCheck{
			OwnerID:         item.OwnerID,
			DependencyID:    item.DependencyID,
			DependencyName:  item.DependencyName,
			RequiredVersion: item.RequiredVersion,
			CurrentVersion:  item.CurrentVersion,
			Installed:       item.Installed,
			Discovered:      item.Discovered,
			Status:          string(item.Status),
			Chain:           cloneStringSlice(item.Chain),
		})
	}
	return out
}

// toDependencyBlockers converts resolver blockers.
func toDependencyBlockers(items []*plugindep.Blocker) []*DependencyBlocker {
	out := make([]*DependencyBlocker, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		out = append(out, &DependencyBlocker{
			Code:            string(item.Code),
			PluginID:        item.PluginID,
			DependencyID:    item.DependencyID,
			RequiredVersion: item.RequiredVersion,
			CurrentVersion:  item.CurrentVersion,
			Chain:           cloneStringSlice(item.Chain),
			Detail:          item.Detail,
		})
	}
	return out
}

// toDependencyReverseDependents converts resolver reverse-dependency results.
func toDependencyReverseDependents(items []*plugindep.ReverseDependent) []*DependencyReverseDependent {
	out := make([]*DependencyReverseDependent, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		out = append(out, &DependencyReverseDependent{
			PluginID:        item.PluginID,
			Name:            item.Name,
			Version:         item.Version,
			RequiredVersion: item.RequiredVersion,
		})
	}
	return out
}

// formatDependencyBlockers renders a compact deterministic blocker summary for fallback messages.
func formatDependencyBlockers(blockers []*plugindep.Blocker) string {
	if len(blockers) == 0 {
		return ""
	}
	parts := make([]string, 0, len(blockers))
	for _, blocker := range blockers {
		if blocker == nil {
			continue
		}
		parts = append(parts, strings.Join([]string{
			string(blocker.Code),
			strings.TrimSpace(blocker.PluginID),
			strings.TrimSpace(blocker.DependencyID),
			strings.TrimSpace(blocker.RequiredVersion),
			strings.TrimSpace(blocker.CurrentVersion),
			strings.Join(blocker.Chain, ">"),
		}, "|"))
	}
	return strings.Join(parts, ";")
}

// firstDependencyBlockerFields returns the first dependency/version tuple for error params.
func firstDependencyBlockerFields(blockers []*plugindep.Blocker) (string, string, string) {
	for _, blocker := range blockers {
		if blocker == nil {
			continue
		}
		return blocker.DependencyID, blocker.RequiredVersion, blocker.CurrentVersion
	}
	return "", "", ""
}

// firstDependencyBlockerChain returns the first blocker chain for structured errors.
func firstDependencyBlockerChain(blockers []*plugindep.Blocker) string {
	for _, blocker := range blockers {
		if blocker == nil {
			continue
		}
		return strings.Join(blocker.Chain, ">")
	}
	return ""
}

// reverseDependentIDs extracts downstream plugin IDs.
func reverseDependentIDs(items []*DependencyReverseDependent) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item == nil || strings.TrimSpace(item.PluginID) == "" {
			continue
		}
		out = append(out, strings.TrimSpace(item.PluginID))
	}
	return out
}

// cloneStringSlice returns a copy of string values for DTO exposure.
func cloneStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, len(values))
	copy(out, values)
	return out
}
