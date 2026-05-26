// Package sourceupgrade implements source-plugin upgrade discovery and explicit
// runtime upgrade execution for the host plugin domain.
package sourceupgrade

import (
	"context"

	"lina-core/internal/service/plugin/internal/catalog"
	"lina-core/internal/service/plugin/internal/integration"
	"lina-core/internal/service/plugin/internal/lifecycle"
	"lina-core/internal/service/plugin/internal/runtime"
)

// SourceUpgradeStatus describes one source plugin's effective version,
// discovered source version, and pending-upgrade state.
type SourceUpgradeStatus struct {
	// PluginID is the immutable plugin identifier.
	PluginID string
	// Name is the human-readable plugin display name.
	Name string
	// EffectiveVersion is the current effective version stored in sys_plugin.
	EffectiveVersion string
	// DiscoveredVersion is the version currently discovered from plugin.yaml.
	DiscoveredVersion string
	// Installed reports whether the plugin is already installed.
	Installed int
	// Enabled reports whether the plugin is currently enabled.
	Enabled int
	// NeedsUpgrade reports whether an installed plugin discovered a newer source version.
	NeedsUpgrade bool
	// DowngradeDetected reports whether the discovered source version is lower
	// than the current effective version, which is unsupported in this iteration.
	DowngradeDetected bool
}

// SourceUpgradeResult describes the outcome of one explicit source-plugin upgrade request.
type SourceUpgradeResult struct {
	// PluginID is the immutable plugin identifier.
	PluginID string
	// Name is the human-readable plugin display name.
	Name string
	// FromVersion is the effective version before the request ran.
	FromVersion string
	// ToVersion is the discovered version targeted by the request.
	ToVersion string
	// Executed reports whether upgrade work actually ran.
	Executed bool
	// Message explains the no-op or successful outcome in the effective locale.
	Message string
	// MessageKey is the runtime i18n key used to render Message.
	MessageKey string
	// MessageParams stores runtime i18n named parameters for MessageKey.
	MessageParams map[string]any
}

// Service defines the host-side source-plugin upgrade governance contract.
type Service interface {
	// ListSourceUpgradeStatuses scans source manifests and returns one
	// effective-versus-discovered upgrade-status item per source plugin.
	ListSourceUpgradeStatuses(ctx context.Context) ([]*SourceUpgradeStatus, error)
	// UpgradeSourcePlugin applies one explicit source-plugin upgrade from the
	// current effective version to the newer discovered source version.
	UpgradeSourcePlugin(ctx context.Context, pluginID string) (*SourceUpgradeResult, error)
	// ValidateSourcePluginUpgradeReadiness scans source-plugin version drift
	// without failing on pending upgrades.
	ValidateSourcePluginUpgradeReadiness(ctx context.Context) error
}

// DependencyValidator validates source-plugin upgrade candidates before the
// upgrade service runs SQL, menu sync, or release switching.
type DependencyValidator interface {
	// ValidateSourcePluginUpgradeCandidate verifies candidate dependencies and
	// reverse-dependency version safety for one source plugin upgrade.
	ValidateSourcePluginUpgradeCandidate(ctx context.Context, manifest *catalog.Manifest) error
}

// Ensure serviceImpl satisfies Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	// catalogSvc provides manifest discovery, registry, and release governance.
	catalogSvc catalog.Service
	// lifecycleSvc provides install/uninstall lifecycle orchestration.
	lifecycleSvc lifecycle.Service
	// runtimeSvc provides dynamic plugin reconciliation and route dispatch.
	runtimeSvc runtime.Service
	// integrationSvc provides host extension, menu, hook, and resource integration.
	integrationSvc integration.Service
	// i18nSvc localizes operator-facing result messages.
	i18nSvc sourceUpgradeI18nService
	// dependencyValidator checks candidate release dependency constraints before upgrade side effects.
	dependencyValidator DependencyValidator
}

// sourceUpgradeI18nService defines the narrow i18n capability needed by source upgrade.
type sourceUpgradeI18nService interface {
	// Translate returns one runtime translation key with caller-provided fallback text.
	Translate(ctx context.Context, key string, fallback string) string
}

// New creates and returns a new source-plugin upgrade governance service.
func New(
	catalogSvc catalog.Service,
	lifecycleSvc lifecycle.Service,
	runtimeSvc runtime.Service,
	integrationSvc integration.Service,
	i18nSvc sourceUpgradeI18nService,
	dependencyValidator DependencyValidator,
) Service {
	return &serviceImpl{
		catalogSvc:          catalogSvc,
		lifecycleSvc:        lifecycleSvc,
		runtimeSvc:          runtimeSvc,
		integrationSvc:      integrationSvc,
		i18nSvc:             i18nSvc,
		dependencyValidator: dependencyValidator,
	}
}
