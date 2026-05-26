// This file exposes the plugin-service facade for source-plugin runtime upgrade
// governance by delegating to the dedicated internal sourceupgrade component.

package plugin

import (
	"context"

	"lina-core/internal/service/plugin/internal/catalog"
	sourceupgradeinternal "lina-core/internal/service/plugin/internal/sourceupgrade"
	"lina-core/pkg/logger"
)

type (
	// SourceUpgradeStatus aliases the internal source-plugin upgrade status contract.
	SourceUpgradeStatus = sourceupgradeinternal.SourceUpgradeStatus

	// SourceUpgradeResult aliases the internal explicit source-plugin upgrade result contract.
	SourceUpgradeResult = sourceupgradeinternal.SourceUpgradeResult
)

// ListSourceUpgradeStatuses scans source manifests and returns one
// effective-versus-discovered upgrade-status item per source plugin.
func (s *serviceImpl) ListSourceUpgradeStatuses(ctx context.Context) ([]*SourceUpgradeStatus, error) {
	return s.sourceUpgradeSvc.ListSourceUpgradeStatuses(ctx)
}

// UpgradeSourcePlugin applies one explicit source-plugin runtime upgrade from
// the current effective version to the newer discovered source version.
func (s *serviceImpl) UpgradeSourcePlugin(ctx context.Context, pluginID string) (*SourceUpgradeResult, error) {
	if err := s.ensurePlatformGovernance(ctx); err != nil {
		return nil, err
	}
	result, err := s.sourceUpgradeSvc.UpgradeSourcePlugin(ctx, pluginID)
	if err != nil {
		s.invalidateRuntimeUpgradeCaches(ctx, pluginID, catalog.TypeSource.String(), "source_plugin_upgrade_failed")
		if _, markErr := s.markRuntimeCacheChanged(ctx, "source_plugin_upgrade_failed"); markErr != nil {
			logger.Warningf(
				ctx,
				"mark runtime cache changed after source upgrade failure failed plugin=%s err=%v",
				pluginID,
				markErr,
			)
		}
		return nil, err
	}
	if result != nil && result.Executed {
		s.invalidateRuntimeUpgradeCaches(ctx, pluginID, catalog.TypeSource.String(), "source_plugin_upgraded")
		if err = s.syncEnabledSnapshotAndPublishRuntimeChange(
			ctx,
			pluginID,
			"source_plugin_upgraded",
		); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ValidateSourcePluginUpgradeReadiness scans source-plugin version drift
// without failing on pending upgrades.
func (s *serviceImpl) ValidateSourcePluginUpgradeReadiness(ctx context.Context) error {
	return s.sourceUpgradeSvc.ValidateSourcePluginUpgradeReadiness(ctx)
}

// Ensure the plugin facade keeps delegating to the dedicated source-upgrade component.
var _ sourceupgradeinternal.Service = (*serviceImpl)(nil)
