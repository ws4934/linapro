// This file adapts dynamic-plugin organization host-service calls to the
// ordinary orgcap.Service consumer contract. The dispatcher intentionally keeps
// host-internal scope, assignment, workspace projection, and database query
// builder seams out of the dynamic-plugin protocol.

package wasm

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/plugin/capability"
	"lina-core/pkg/plugin/capability/orgcap"
	bridgehostcall "lina-core/pkg/plugin/pluginbridge/protocol"
	bridgehostservice "lina-core/pkg/plugin/pluginbridge/protocol"
)

// orgHostServices stores the runtime-owned capability services used by
// organization host-service dispatch.
var orgHostServices capability.Services

// ConfigureOrgHostService replaces the organization capability service
// directory used by dynamic-plugin host calls.
func ConfigureOrgHostService(services capability.Services) error {
	if services == nil {
		return gerror.New("org host services directory is nil")
	}
	orgHostServices = services
	return nil
}

// dispatchOrgHostService routes one organization host-service method to the
// same ordinary orgcap.Service surface exposed to source plugins.
func dispatchOrgHostService(
	ctx context.Context,
	hcc *hostCallContext,
	method string,
	payload []byte,
) *bridgehostcall.HostCallResponseEnvelope {
	service := orgServiceForHostCall(hcc)
	if service == nil {
		return bridgehostcall.NewHostCallErrorResponse(
			bridgehostcall.HostCallStatusInternalError,
			"org host service is not scoped",
		)
	}

	switch method {
	case bridgehostservice.HostServiceMethodOrgAvailable:
		return capabilityJSONResponse(service.Available(ctx))
	case bridgehostservice.HostServiceMethodOrgStatus:
		return capabilityJSONResponse(service.Status(ctx))
	case bridgehostservice.HostServiceMethodOrgListUserDeptAssignments:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityUsersRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		assignments, err := service.ListUserDeptAssignments(ctx, request.UserIDs)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(assignments)
	case bridgehostservice.HostServiceMethodOrgGetUserDeptInfo:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityUserRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		deptID, deptName, err := service.GetUserDeptInfo(ctx, request.UserID)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(orgUserDeptInfoResponse{DeptID: deptID, DeptName: deptName})
	case bridgehostservice.HostServiceMethodOrgGetUserDeptName:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityUserRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		name, err := service.GetUserDeptName(ctx, request.UserID)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(name)
	case bridgehostservice.HostServiceMethodOrgGetUserDeptIDs:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityUserRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		deptIDs, err := service.GetUserDeptIDs(ctx, request.UserID)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(deptIDs)
	case bridgehostservice.HostServiceMethodOrgGetUserPostIDs:
		request, err := bridgehostservice.UnmarshalHostServiceCapabilityUserRequest(payload)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		postIDs, err := service.GetUserPostIDs(ctx, request.UserID)
		if err != nil {
			return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInvalidRequest, err.Error())
		}
		return capabilityJSONResponse(postIDs)
	default:
		return bridgehostcall.NewHostCallErrorResponse(
			bridgehostcall.HostCallStatusNotFound,
			"org host service method not implemented: "+method,
		)
	}
}

// orgServiceForPlugin returns the runtime-owned organization capability bound
// to one dynamic plugin. Data host service uses this for department data-scope
// filtering without depending on an org host-service authorization declaration.
func orgServiceForPlugin(pluginID string) orgcap.Service {
	if orgHostServices == nil {
		return nil
	}
	services := capability.ServicesForPlugin(orgHostServices, pluginID)
	if services == nil {
		return nil
	}
	return services.Org()
}

// orgServiceForHostCall resolves the organization service for one host call.
func orgServiceForHostCall(hcc *hostCallContext) orgcap.Service {
	if hcc == nil {
		return nil
	}
	return orgServiceForPlugin(hcc.pluginID)
}

// orgUserDeptInfoResponse carries the tuple returned by orgcap.Service.GetUserDeptInfo.
type orgUserDeptInfoResponse struct {
	// DeptID is the department identifier.
	DeptID int `json:"deptId"`
	// DeptName is the department display name.
	DeptName string `json:"deptName"`
}

// capabilityJSONResponse encodes one capability result as a transport-owned
// JSON response without making pluginbridge own capability DTO definitions.
func capabilityJSONResponse(value any) *bridgehostcall.HostCallResponseEnvelope {
	content, err := json.Marshal(value)
	if err != nil {
		return bridgehostcall.NewHostCallErrorResponse(bridgehostcall.HostCallStatusInternalError, err.Error())
	}
	payload := bridgehostservice.MarshalHostServiceCapabilityJSONResponse(
		&bridgehostservice.HostServiceCapabilityJSONResponse{Value: content},
	)
	return bridgehostcall.NewHostCallSuccessResponse(payload)
}
