// This file keeps platform-context checks for global menu-governance writes in
// one place. sys_menu remains the host-wide permission topology, so tenant and
// impersonation contexts must fail before any mutation or topology revision.

package menu

import (
	"context"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
)

// platformMenuTenantService is the tenant-capability slice required by
// global menu governance guards.
type platformMenuTenantService interface {
	// Available reports whether multi-tenancy governance is active.
	Available(ctx context.Context) bool
	// PlatformBypass reports whether the request is a platform all-data context.
	PlatformBypass(ctx context.Context) bool
}

// ensurePlatformMenuGovernance verifies the current request can mutate the
// global menu topology.
func (s *serviceImpl) ensurePlatformMenuGovernance(ctx context.Context) error {
	return ensurePlatformMenuGovernanceContext(ctx, s)
}

// ensurePlatformMenuGovernanceContext applies platform-menu checks without
// coupling tests to the full tenantcap service interface.
func ensurePlatformMenuGovernanceContext(ctx context.Context, holder interface {
	platformMenuTenantService() platformMenuTenantService
}) error {
	if holder == nil {
		return nil
	}
	tenantSvc := holder.platformMenuTenantService()
	if tenantSvc == nil || !tenantSvc.Available(ctx) || tenantSvc.PlatformBypass(ctx) {
		return nil
	}
	return bizerr.NewCode(tenantcap.CodePlatformPermissionRequired)
}

// platformMenuTenantService returns the tenant capability used by the menu
// governance guard.
func (s *serviceImpl) platformMenuTenantService() platformMenuTenantService {
	if s == nil {
		return nil
	}
	return s.tenantSvc
}
