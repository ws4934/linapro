// This file verifies plugin lifecycle transitions expose provider availability
// through the platform plugin enabled snapshot used by capability capabilities.

package plugin

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/gogf/gf/v2/net/ghttp"

	"lina-core/pkg/plugin/capability/contract"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/pluginhost"
)

// TestSourceProviderAvailabilityFollowsEnabledSnapshot verifies provider
// declarations remain inert until their owning source plugin is platform-enabled.
func TestSourceProviderAvailabilityFollowsEnabledSnapshot(t *testing.T) {
	var (
		ctx      = contract.WithCurrentContext(context.Background(), contract.CurrentContext{TenantID: 0, PlatformBypass: true})
		pluginID = "plugin-dev-source-capability-revision"
		service  = newTestServiceWithTopology(&testTopology{
			enabled: true,
			primary: true,
			nodeID:  "capability-revision-node",
		})
	)
	cleanupTestPluginIDs(t, ctx, pluginID)

	plugin := pluginhost.NewSourcePlugin(pluginID)
	plugin.Assets().UseEmbeddedFiles(fstest.MapFS{
		"plugin.yaml": &fstest.MapFile{Data: []byte(
			"id: " + pluginID + "\n" +
				"name: Runtime Revision Provider\n" +
				"version: v0.1.0\n" +
				"type: source\n" +
				"scope_nature: tenant_aware\n" +
				"supports_multi_tenant: false\n" +
				"default_install_mode: global\n",
		)},
	})
	cleanup, err := pluginhost.RegisterSourcePluginForTest(plugin)
	if err != nil {
		t.Fatalf("register source plugin fixture failed: %v", err)
	}
	t.Cleanup(cleanup)

	if err = tenantcap.Provide(pluginID, func(
		context.Context,
		tenantcap.ProviderEnv,
	) (tenantcap.Provider, error) {
		return capabilityRevisionProvider{}, nil
	}); err != nil {
		t.Fatalf("register tenant provider factory failed: %v", err)
	}

	if _, err = service.Install(ctx, pluginID, InstallOptions{}); err != nil {
		t.Fatalf("install source provider plugin failed: %v", err)
	}
	tenantSvc := tenantcap.New(service, nil)
	status := tenantSvc.Status(ctx)
	if status.Available || status.ActiveProvider == pluginID {
		t.Fatalf("expected installed-but-disabled provider unavailable, got %#v", status)
	}

	if err = service.Enable(ctx, pluginID); err != nil {
		t.Fatalf("enable source provider plugin failed: %v", err)
	}
	status = tenantSvc.Status(ctx)
	if !status.Available || status.ActiveProvider != pluginID {
		t.Fatalf("expected tenant provider active for %s, got %#v", pluginID, status)
	}

	if err = service.Disable(ctx, pluginID); err != nil {
		t.Fatalf("disable source provider plugin failed: %v", err)
	}
	status = tenantSvc.Status(ctx)
	if status.Available || status.ActiveProvider == pluginID {
		t.Fatalf("expected disabled provider unavailable, got %#v", status)
	}
}

// capabilityRevisionProvider is a no-op tenant provider used by the
// lifecycle/runtime-revision integration test.
type capabilityRevisionProvider struct{}

// ResolveTenant returns the platform tenant.
func (capabilityRevisionProvider) ResolveTenant(
	context.Context,
	*ghttp.Request,
) (*tenantcap.ResolverResult, error) {
	return &tenantcap.ResolverResult{
		TenantID: tenantcap.PLATFORM,
		Matched:  true,
	}, nil
}

// ValidateUserInTenant accepts every tenant validation request.
func (capabilityRevisionProvider) ValidateUserInTenant(
	context.Context,
	int,
	tenantcap.TenantID,
) error {
	return nil
}

// ListUserTenants returns no tenant memberships.
func (capabilityRevisionProvider) ListUserTenants(
	context.Context,
	int,
) ([]tenantcap.TenantInfo, error) {
	return []tenantcap.TenantInfo{}, nil
}

// SwitchTenant accepts every tenant switch request.
func (capabilityRevisionProvider) SwitchTenant(
	context.Context,
	int,
	tenantcap.TenantID,
) error {
	return nil
}
