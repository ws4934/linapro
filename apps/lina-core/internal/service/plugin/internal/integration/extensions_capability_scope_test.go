// This file verifies source-plugin callback flows receive plugin-scoped
// capability services through the public scoping mechanism.

package integration_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"lina-core/internal/service/plugin/internal/testutil"
	"lina-core/pkg/plugin/capability"
	"lina-core/pkg/plugin/capability/contract"
	capabilityorgcap "lina-core/pkg/plugin/capability/orgcap"
	tenantcapsvc "lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/pluginhost"
)

// capabilityScopeRecorder records plugin-scope binding performed by
// capability.ServicesForPlugin.
type capabilityScopeRecorder struct {
	emptySourceServicesDirectory

	mu     sync.Mutex
	scopes []string
}

var _ capability.Services = (*capabilityScopeRecorder)(nil)
var _ capability.ScopedServicesFactory = (*capabilityScopeRecorder)(nil)

// ForPlugin returns a source-plugin capability view bound to pluginID.
func (r *capabilityScopeRecorder) ForPlugin(pluginID string) capability.Services {
	normalizedPluginID := strings.TrimSpace(pluginID)
	r.mu.Lock()
	r.scopes = append(r.scopes, normalizedPluginID)
	r.mu.Unlock()
	return &scopedSourceServicesDirectory{pluginID: normalizedPluginID}
}

// seenScope reports whether ServicesForPlugin bound pluginID at least once.
func (r *capabilityScopeRecorder) seenScope(pluginID string) bool {
	normalizedPluginID := strings.TrimSpace(pluginID)
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, scope := range r.scopes {
		if scope == normalizedPluginID {
			return true
		}
	}
	return false
}

// scopedSourceServicesDirectory is the plugin-bound Services view
// returned to source-plugin callback code.
type scopedSourceServicesDirectory struct {
	emptySourceServicesDirectory

	pluginID string
}

var _ pluginhost.Services = (*scopedSourceServicesDirectory)(nil)

// scopedPluginID exposes the test-only scope marker for assertions.
func (d *scopedSourceServicesDirectory) scopedPluginID() string {
	if d == nil {
		return ""
	}
	return d.pluginID
}

// emptySourceServicesDirectory is a minimal Services test double.
type emptySourceServicesDirectory struct{}

var _ pluginhost.Services = (*emptySourceServicesDirectory)(nil)

// APIDoc returns no API-doc service for this capability-scope test.
func (emptySourceServicesDirectory) APIDoc() contract.APIDocService { return nil }

// Auth returns no auth service for this capability-scope test.
func (emptySourceServicesDirectory) Auth() contract.AuthService { return nil }

// BizCtx returns no business-context service for this capability-scope test.
func (emptySourceServicesDirectory) BizCtx() contract.BizCtxService { return nil }

// Cache returns no cache service for this capability-scope test.
func (emptySourceServicesDirectory) Cache() contract.CacheService { return nil }

// Config returns no config service for this capability-scope test.
func (emptySourceServicesDirectory) Config() contract.ConfigService { return nil }

// HostConfig returns no host-config service for this capability-scope test.
func (emptySourceServicesDirectory) HostConfig() contract.HostConfigService { return nil }

// I18n returns no i18n service for this capability-scope test.
func (emptySourceServicesDirectory) I18n() contract.I18nService { return nil }

// Manifest returns no manifest service for this capability-scope test.
func (emptySourceServicesDirectory) Manifest() contract.ManifestService { return nil }

// Notify returns no notification service for this capability-scope test.
func (emptySourceServicesDirectory) Notify() contract.NotifyService { return nil }

// Org returns no organization capability for this capability-scope test.
func (emptySourceServicesDirectory) Org() capabilityorgcap.Service { return nil }

// PluginLifecycle returns no plugin lifecycle service for this capability-scope test.
func (emptySourceServicesDirectory) PluginLifecycle() contract.PluginLifecycleService { return nil }

// PluginState returns no plugin-state service for this capability-scope test.
func (emptySourceServicesDirectory) PluginState() contract.PluginStateService { return nil }

// Route returns no route service for this capability-scope test.
func (emptySourceServicesDirectory) Route() contract.RouteService { return nil }

// Session returns no session service for this capability-scope test.
func (emptySourceServicesDirectory) Session() contract.SessionService { return nil }

// Tenant returns no tenant capability for this capability-scope test.
func (emptySourceServicesDirectory) Tenant() tenantcapsvc.Service { return nil }

// TenantFilter returns no tenant-filter service for this capability-scope test.
func (emptySourceServicesDirectory) TenantFilter() contract.TenantFilterService { return nil }

// scopedCapabilityView is implemented by test doubles returned from ForPlugin.
type scopedCapabilityView interface {
	scopedPluginID() string
}

// TestSourcePluginCallbacksUsePluginScopedServices verifies route, cron, hook,
// and managed-cron integration flows all bind runtime services through
// capability.ServicesForPlugin before exposing them to a source plugin.
func TestSourcePluginCallbacksUsePluginScopedServices(t *testing.T) {
	const pluginID = "plugin-dev-source-capability-scope"

	services := testutil.NewServices()
	recorder := &capabilityScopeRecorder{}
	services.Integration.SetCapabilities(recorder)

	observed := make(map[string]string)
	currentPhase := ""
	recordServices := func(label string, services pluginhost.Services) error {
		if services == nil {
			return fmt.Errorf("%s services are nil", label)
		}
		scoped, ok := services.(scopedCapabilityView)
		if !ok {
			return fmt.Errorf("%s services were not plugin-scoped: %T", label, services)
		}
		if got := scoped.scopedPluginID(); got != pluginID {
			return fmt.Errorf("%s services scope = %q, want %q", label, got, pluginID)
		}
		observed[label] = scoped.scopedPluginID()
		return nil
	}

	sourcePlugin := pluginhost.NewSourcePlugin(pluginID)
	sourcePlugin.Assets().UseEmbeddedFiles(fstest.MapFS{
		"plugin.yaml": &fstest.MapFile{Data: []byte(
			"id: " + pluginID + "\n" +
				"name: Source Capability Scope Plugin\n" +
				"version: 0.1.0\n" +
				"type: source\n" +
				"scope_nature: tenant_aware\n" +
				"supports_multi_tenant: true\n" +
				"default_install_mode: tenant_scoped\n",
		)},
		"backend/plugin.go":             &fstest.MapFile{Data: []byte("package backend\n")},
		"frontend/pages/main-entry.vue": &fstest.MapFile{Data: []byte("<template><div /></template>\n")},
	})

	if err := sourcePlugin.HTTP().RegisterRoutes(
		pluginhost.ExtensionPointHTTPRouteRegister,
		pluginhost.CallbackExecutionModeBlocking,
		func(_ context.Context, registrar pluginhost.HTTPRegistrar) error {
			return recordServices("route", registrar.Services())
		},
	); err != nil {
		t.Fatalf("failed to register source route handler: %v", err)
	}
	if err := sourcePlugin.Cron().RegisterCron(
		pluginhost.ExtensionPointCronRegister,
		pluginhost.CallbackExecutionModeBlocking,
		func(_ context.Context, registrar pluginhost.CronRegistrar) error {
			if currentPhase == "" {
				return fmt.Errorf("cron registration phase was not set")
			}
			return recordServices(currentPhase, registrar.Services())
		},
	); err != nil {
		t.Fatalf("failed to register source cron handler: %v", err)
	}
	if err := sourcePlugin.Hooks().RegisterHook(
		pluginhost.ExtensionPointPluginInstalled,
		pluginhost.CallbackExecutionModeBlocking,
		func(_ context.Context, payload pluginhost.HookPayload) error {
			return recordServices("hook", payload.Services())
		},
	); err != nil {
		t.Fatalf("failed to register source hook handler: %v", err)
	}

	cleanup, err := pluginhost.RegisterSourcePluginForTest(sourcePlugin)
	if err != nil {
		t.Fatalf("failed to register source plugin fixture: %v", err)
	}
	t.Cleanup(cleanup)

	ctx := context.Background()
	server := g.Server("integration-source-capability-scope")
	server.SetDumpRouterMap(false)
	var rootGroup *ghttp.RouterGroup
	server.Group("/", func(group *ghttp.RouterGroup) {
		rootGroup = group
	})

	if err = services.Integration.RegisterHTTPRoutes(ctx, server, rootGroup, nil); err != nil {
		t.Fatalf("expected route registration to receive scoped services, got error: %v", err)
	}

	currentPhase = "cron"
	if err = services.Integration.RegisterCrons(ctx); err != nil {
		t.Fatalf("expected cron registration to receive scoped services, got error: %v", err)
	}

	if err = services.Integration.DispatchPluginHookEvent(
		ctx,
		pluginhost.ExtensionPointPluginInstalled,
		pluginhost.BuildPluginLifecycleHookPayloadValues(pluginhost.PluginLifecycleHookPayloadInput{
			PluginID: pluginID,
			Name:     "Source Capability Scope Plugin",
			Version:  "0.1.0",
		}),
	); err != nil {
		t.Fatalf("expected hook dispatch to receive scoped services, got error: %v", err)
	}

	currentPhase = "managed-cron"
	if _, err = services.Integration.ListCronDeclarationsByPlugin(ctx, pluginID); err != nil {
		t.Fatalf("expected managed cron collection to receive scoped services, got error: %v", err)
	}

	for _, label := range []string{"route", "cron", "hook", "managed-cron"} {
		if got := observed[label]; got != pluginID {
			t.Fatalf("expected %s callback to receive plugin-scoped services for %q, got %q", label, pluginID, got)
		}
	}
	if !recorder.seenScope(pluginID) {
		t.Fatalf("expected capability.ServicesForPlugin to bind %q", pluginID)
	}
}
