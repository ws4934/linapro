// This file owns the process-local plugin management list read-model cache.
// It keeps the cache scoped to the root plugin facade so lifecycle mutations can
// invalidate the complete list without changing API response contracts.

package plugin

import (
	"context"
	"strconv"
	"sync"

	i18nsvc "lina-core/internal/service/i18n"
	"lina-core/internal/service/plugin/internal/catalog"
	"lina-core/pkg/plugin/pluginbridge/protocol"
)

// pluginManagementListCache stores one complete unfiltered plugin management
// list read model. Filtered API requests derive their page data from this
// complete projection so modal-dependent fields remain available.
type pluginManagementListCache struct {
	mu     sync.RWMutex
	values map[string]*ListOutput
}

// get returns a defensive copy of the cached list, if available.
func (c *pluginManagementListCache) get(key pluginManagementListCacheKey) (*ListOutput, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.values == nil {
		return nil, false
	}
	value := c.values[key.String()]
	if value == nil {
		return nil, false
	}
	return cloneListOutput(value), true
}

// store replaces the cached list with a defensive copy.
func (c *pluginManagementListCache) store(key pluginManagementListCacheKey, value *ListOutput) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if value == nil {
		return
	}
	if c.values == nil {
		c.values = make(map[string]*ListOutput)
	}
	for existingKey := range c.values {
		if pluginManagementListCacheKeyLocale(existingKey) == key.Locale && existingKey != key.String() {
			delete(c.values, existingKey)
		}
	}
	c.values[key.String()] = cloneListOutput(value)
}

// invalidate clears the cached list after plugin governance changes.
func (c *pluginManagementListCache) invalidate() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values = nil
}

// PrewarmManagementList builds the complete plugin management list read model
// so the first administrator request can reuse hot discovery and dependency
// projections. Failures are returned to foreground callers and logged by
// asynchronous startup callers.
func (s *serviceImpl) PrewarmManagementList(ctx context.Context) error {
	if _, err := s.managementList(ctx); err != nil {
		return err
	}
	return nil
}

// managementList returns the complete unfiltered plugin management read model.
func (s *serviceImpl) managementList(ctx context.Context) (*ListOutput, error) {
	if err := s.ensureRuntimeCacheFresh(ctx); err != nil {
		return nil, err
	}
	cacheKey := s.managementListCacheKey(ctx)
	if cached, ok := s.managementListCache.get(cacheKey); ok {
		return cached, nil
	}
	out, err := s.buildManagementList(ctx)
	if err != nil {
		return nil, err
	}
	s.managementListCache.store(s.managementListCacheKey(ctx), out)
	return cloneListOutput(out), nil
}

// InvalidateManagementListCache clears this process-local read model. Cluster
// peers observe the same plugin-runtime revision and invalidate through the
// root runtime-cache refresh callback.
func (s *serviceImpl) InvalidateManagementListCache(_ context.Context, _ string) {
	if s == nil {
		return
	}
	s.managementListCache.invalidate()
}

// managementListCacheKey returns the current cache partition because plugin
// display metadata is localized during projection and can change when the
// runtime translation bundle version changes.
func (s *serviceImpl) managementListCacheKey(ctx context.Context) pluginManagementListCacheKey {
	if s == nil || s.i18nSvc == nil {
		return pluginManagementListCacheKey{Locale: i18nsvc.DefaultLocale}
	}
	locale := normalizeManagementListCacheLocale(s.i18nSvc.GetLocale(ctx))
	return pluginManagementListCacheKey{
		Locale:               locale,
		RuntimeBundleVersion: s.i18nSvc.BundleVersion(locale),
	}
}

// normalizeManagementListCacheLocale keeps cache keys stable for detached
// startup contexts and tests that do not carry business locale metadata.
func normalizeManagementListCacheLocale(locale string) string {
	if locale == "" {
		return i18nsvc.DefaultLocale
	}
	return locale
}

// pluginManagementListCacheKey identifies one localized management list read model.
type pluginManagementListCacheKey struct {
	Locale               string
	RuntimeBundleVersion uint64
}

// String returns a stable map key for the localized read-model cache.
func (k pluginManagementListCacheKey) String() string {
	return normalizeManagementListCacheLocale(k.Locale) + "@" + strconv.FormatUint(k.RuntimeBundleVersion, 10)
}

// pluginManagementListCacheKeyLocale extracts the locale prefix from one cache key.
func pluginManagementListCacheKeyLocale(key string) string {
	for index, char := range key {
		if char == '@' {
			return normalizeManagementListCacheLocale(key[:index])
		}
	}
	return normalizeManagementListCacheLocale(key)
}

// withManifestSnapshot stores one already-scanned manifest list in context so
// dependency checks inside the same list build do not rescan source plugins and
// dynamic artifacts for every plugin row.
func withManifestSnapshot(ctx context.Context, manifests []*catalog.Manifest) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if manifestSnapshotFromContext(ctx) != nil {
		return ctx
	}
	return context.WithValue(ctx, manifestSnapshotContextKey{}, cloneManifestSlice(manifests))
}

// manifestSnapshotFromContext returns the request-local manifest list, if set.
func manifestSnapshotFromContext(ctx context.Context) []*catalog.Manifest {
	if ctx == nil {
		return nil
	}
	manifests, ok := ctx.Value(manifestSnapshotContextKey{}).([]*catalog.Manifest)
	if !ok || manifests == nil {
		return nil
	}
	return cloneManifestSlice(manifests)
}

// cloneListOutput copies one list output and the plugin rows it owns.
func cloneListOutput(in *ListOutput) *ListOutput {
	if in == nil {
		return nil
	}
	out := &ListOutput{
		List:  make([]*PluginItem, 0, len(in.List)),
		Total: in.Total,
	}
	for _, item := range in.List {
		out.List = append(out.List, clonePluginItem(item))
	}
	return out
}

// clonePluginItem copies one plugin item while preserving immutable nested
// projections by value where callers may otherwise mutate list rows.
func clonePluginItem(in *PluginItem) *PluginItem {
	if in == nil {
		return nil
	}
	out := *in
	if in.LastUpgradeFailure != nil {
		lastUpgradeFailure := *in.LastUpgradeFailure
		out.LastUpgradeFailure = &lastUpgradeFailure
	}
	out.RequestedHostServices = cloneHostServiceSpecs(in.RequestedHostServices)
	out.AuthorizedHostServices = cloneHostServiceSpecs(in.AuthorizedHostServices)
	out.DeclaredRoutes = cloneRouteContractsForList(in.DeclaredRoutes)
	out.DependencyCheck = cloneDependencyCheckResult(in.DependencyCheck)
	return &out
}

// cloneHostServiceSpecs deep-copies host-service declarations because list
// consumers may reuse rows while action modals are open.
func cloneHostServiceSpecs(in []*protocol.HostServiceSpec) []*protocol.HostServiceSpec {
	if len(in) == 0 {
		return nil
	}
	out := make([]*protocol.HostServiceSpec, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		out = append(out, &protocol.HostServiceSpec{
			Service:   item.Service,
			Methods:   append([]string(nil), item.Methods...),
			Paths:     append([]string(nil), item.Paths...),
			Tables:    append([]string(nil), item.Tables...),
			Keys:      append([]string(nil), item.Keys...),
			Resources: cloneHostServiceResources(item.Resources),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// cloneHostServiceResources deep-copies governed host-service resource specs.
func cloneHostServiceResources(in []*protocol.HostServiceResourceSpec) []*protocol.HostServiceResourceSpec {
	if len(in) == 0 {
		return nil
	}
	out := make([]*protocol.HostServiceResourceSpec, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		out = append(out, &protocol.HostServiceResourceSpec{
			Ref:             item.Ref,
			AllowMethods:    append([]string(nil), item.AllowMethods...),
			HeaderAllowList: append([]string(nil), item.HeaderAllowList...),
			TimeoutMs:       item.TimeoutMs,
			MaxBodyBytes:    item.MaxBodyBytes,
			Attributes:      cloneStringMapForList(item.Attributes),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// cloneRouteContractsForList deep-copies route review declarations.
func cloneRouteContractsForList(in []*protocol.RouteContract) []*protocol.RouteContract {
	if len(in) == 0 {
		return nil
	}
	out := make([]*protocol.RouteContract, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		out = append(out, &protocol.RouteContract{
			Path:        item.Path,
			Method:      item.Method,
			Tags:        append([]string(nil), item.Tags...),
			Summary:     item.Summary,
			Description: item.Description,
			Access:      item.Access,
			Permission:  item.Permission,
			Meta:        cloneStringMapForList(item.Meta),
			RequestType: item.RequestType,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// cloneDependencyCheckResult deep-copies the dependency result attached to list rows.
func cloneDependencyCheckResult(in *DependencyCheckResult) *DependencyCheckResult {
	if in == nil {
		return nil
	}
	out := *in
	out.Dependencies = cloneDependencyPluginChecks(in.Dependencies)
	out.Blockers = cloneDependencyBlockers(in.Blockers)
	out.Cycle = append([]string(nil), in.Cycle...)
	out.ReverseDependents = cloneDependencyReverseDependents(in.ReverseDependents)
	out.ReverseBlockers = cloneDependencyBlockers(in.ReverseBlockers)
	return &out
}

// cloneDependencyPluginChecks deep-copies dependency edge checks.
func cloneDependencyPluginChecks(in []*DependencyPluginCheck) []*DependencyPluginCheck {
	if len(in) == 0 {
		return nil
	}
	out := make([]*DependencyPluginCheck, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		clone := *item
		clone.Chain = append([]string(nil), item.Chain...)
		out = append(out, &clone)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// cloneDependencyBlockers deep-copies dependency blocker diagnostics.
func cloneDependencyBlockers(in []*DependencyBlocker) []*DependencyBlocker {
	if len(in) == 0 {
		return nil
	}
	out := make([]*DependencyBlocker, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		clone := *item
		clone.Chain = append([]string(nil), item.Chain...)
		out = append(out, &clone)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// cloneDependencyReverseDependents deep-copies reverse dependency summaries.
func cloneDependencyReverseDependents(in []*DependencyReverseDependent) []*DependencyReverseDependent {
	if len(in) == 0 {
		return nil
	}
	out := make([]*DependencyReverseDependent, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		clone := *item
		out = append(out, &clone)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// cloneStringMapForList copies string maps used by cached list projections.
func cloneStringMapForList(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

// cloneManifestSlice copies the manifest slice header so callers cannot mutate
// the request-local list ordering.
func cloneManifestSlice(in []*catalog.Manifest) []*catalog.Manifest {
	if in == nil {
		return nil
	}
	out := make([]*catalog.Manifest, len(in))
	copy(out, in)
	return out
}

// manifestSnapshotContextKey stores one request-local manifest discovery result.
type manifestSnapshotContextKey struct{}
