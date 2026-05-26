// This file adapts dynamic-plugin tenant host-service calls to the ordinary
// tenantcap.Service consumer contract. The dispatcher intentionally excludes
// host-internal HTTP request resolution, database query builders, membership
// write seams, and lifecycle governance services.

package wasm

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/plugin/capability"
	"lina-core/pkg/plugin/capability/tenantcap"
	bridgehostcall "lina-core/pkg/plugin/pluginbridge/protocol"
	bridgehostservice "lina-core/pkg/plugin/pluginbridge/protocol"
)

// tenantHostServices stores the runtime-owned capability services used by
// tenant host-service dispatch.
var tenantHostServices capability.Services

// ConfigureTenantHostService replaces the tenant capability service directory
// used by dynamic-plugin host calls.
func ConfigureTenantHostService(services capability.Services) error {
	if services == nil {
		return gerror.New("tenant host services directory is nil")
	}
	tenantHostServices = services
	return nil
}

// dispatchTenantHostService routes one tenant host-service method to the same
// ordinary tenantcap.Service surface exposed to source plugins.
func dispatchTenantHostService(
	ctx context.Context,
	hcc *hostCallContext,
	method string,
	payload []byte,
) *bridgehostcall.HostCallResponseEnvelope {
	service := tenantServiceForHostCall(hcc)
	if service == nil {
		return bridgehostcall.NewHostCallErrorResponse(
			bridgehostcall.HostCallStatusInternalError,
			"tenant host service is not scoped",
		)
	}

	switch method {
	case bridgehostservice.HostServiceMethodTenantAvailable:
		return capabilityJSONResponse(service.Available(ctx))
	case bridgehostservice.HostServiceMethodTenantStatus:
		return capabilityJSONResponse(service.Status(ctx))
	case bridgehostservice.HostServiceMethodTenantCurrent:
		return capabilityJSONResponse(service.Current(ctx))
	case bridgehostservice.HostServiceMethodTenantPlatformBypass:
		return capabilityJSONResponse(service.PlatformBypass(ctx))
	case bridgehostservice.HostServiceMethodTenantEnsureVisible:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityTenantRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		if err = service.EnsureTenantVisible(ctx, tenantcap.TenantID(request.TenantID)); err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(true)
	case bridgehostservice.HostServiceMethodTenantValidateUserInTenant:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityUserTenantRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		if err = service.ValidateUserInTenant(ctx, request.UserID, tenantcap.TenantID(request.TenantID)); err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(true)
	case bridgehostservice.HostServiceMethodTenantListUserTenants:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityUserRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		tenants, err := service.ListUserTenants(ctx, request.UserID)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(tenants)
	case bridgehostservice.HostServiceMethodTenantValidateSwitch:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityUserTenantSwitchRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		if err = service.SwitchTenant(ctx, request.UserID, tenantcap.TenantID(request.TargetTenantID)); err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(true)
	default:
		return bridgehostcall.NewHostCallErrorResponse(
			bridgehostcall.HostCallStatusNotFound,
			"tenant host service method not implemented: "+method,
		)
	}
}

// tenantServiceForHostCall resolves the tenant service for one host call.
func tenantServiceForHostCall(hcc *hostCallContext) tenantcap.Service {
	if hcc == nil || tenantHostServices == nil {
		return nil
	}
	services := capability.ServicesForPlugin(tenantHostServices, hcc.pluginID)
	if services == nil {
		return nil
	}
	return services.Tenant()
}
