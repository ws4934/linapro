// This file implements transport-only codecs for organization and tenant
// capability host-service calls. Capability DTO ownership stays in
// capability/orgcap, capability/tenantcap, and capability/contract; the bridge
// only carries primitive request fields and compact JSON projections.

package hostservice

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/protobuf/encoding/protowire"
)

// HostServiceCapabilityUserRequest carries one user identifier.
type HostServiceCapabilityUserRequest struct {
	// UserID is the target host user identifier.
	UserID int `json:"userId"`
}

// HostServiceCapabilityUsersRequest carries multiple user identifiers.
type HostServiceCapabilityUsersRequest struct {
	// UserIDs are the target host user identifiers.
	UserIDs []int `json:"userIds"`
}

// HostServiceCapabilityTenantRequest carries one tenant identifier.
type HostServiceCapabilityTenantRequest struct {
	// TenantID is the target tenant identifier.
	TenantID int `json:"tenantId"`
}

// HostServiceCapabilityUserTenantRequest carries one user and tenant pair.
type HostServiceCapabilityUserTenantRequest struct {
	// UserID is the target host user identifier.
	UserID int `json:"userId"`
	// TenantID is the target tenant identifier.
	TenantID int `json:"tenantId"`
}

// HostServiceCapabilityUserTenantSwitchRequest carries one tenant switch check.
type HostServiceCapabilityUserTenantSwitchRequest struct {
	// UserID is the target host user identifier.
	UserID int `json:"userId"`
	// TargetTenantID is the tenant switch target.
	TargetTenantID int `json:"targetTenantId"`
}

// HostServiceCapabilityJSONResponse carries one JSON-encoded capability projection.
type HostServiceCapabilityJSONResponse struct {
	// Value is the compact JSON projection returned by the capability service.
	Value []byte `json:"value"`
}

// MarshalHostServiceCapabilityUserRequest encodes one user request.
func MarshalHostServiceCapabilityUserRequest(req *HostServiceCapabilityUserRequest) []byte {
	var content []byte
	if req == nil {
		return content
	}
	if req.UserID != 0 {
		content = appendVarintField(content, 1, uint64(req.UserID))
	}
	return content
}

// UnmarshalHostServiceCapabilityUserRequest decodes one user request.
func UnmarshalHostServiceCapabilityUserRequest(data []byte) (*HostServiceCapabilityUserRequest, error) {
	out := &HostServiceCapabilityUserRequest{}
	content := data
	for len(content) > 0 {
		fieldNumber, wireType, length := protowire.ConsumeTag(content)
		if length < 0 {
			return nil, gerror.New("failed to decode capability user request tag")
		}
		content = content[length:]
		switch fieldNumber {
		case 1:
			value, size := protowire.ConsumeVarint(content)
			if size < 0 {
				return nil, gerror.New("failed to decode capability user request userId")
			}
			out.UserID = int(value)
			content = content[size:]
		default:
			size := protowire.ConsumeFieldValue(fieldNumber, wireType, content)
			if size < 0 {
				return nil, gerror.New("failed to skip unknown capability user request field")
			}
			content = content[size:]
		}
	}
	return out, nil
}

// MarshalHostServiceCapabilityUsersRequest encodes one users request.
func MarshalHostServiceCapabilityUsersRequest(req *HostServiceCapabilityUsersRequest) []byte {
	var content []byte
	if req == nil {
		return content
	}
	for _, userID := range req.UserIDs {
		content = appendVarintField(content, 1, uint64(userID))
	}
	return content
}

// UnmarshalHostServiceCapabilityUsersRequest decodes one users request.
func UnmarshalHostServiceCapabilityUsersRequest(data []byte) (*HostServiceCapabilityUsersRequest, error) {
	out := &HostServiceCapabilityUsersRequest{}
	content := data
	for len(content) > 0 {
		fieldNumber, wireType, length := protowire.ConsumeTag(content)
		if length < 0 {
			return nil, gerror.New("failed to decode capability users request tag")
		}
		content = content[length:]
		switch fieldNumber {
		case 1:
			value, size := protowire.ConsumeVarint(content)
			if size < 0 {
				return nil, gerror.New("failed to decode capability users request userId")
			}
			out.UserIDs = append(out.UserIDs, int(value))
			content = content[size:]
		default:
			size := protowire.ConsumeFieldValue(fieldNumber, wireType, content)
			if size < 0 {
				return nil, gerror.New("failed to skip unknown capability users request field")
			}
			content = content[size:]
		}
	}
	return out, nil
}

// MarshalHostServiceCapabilityTenantRequest encodes one tenant request.
func MarshalHostServiceCapabilityTenantRequest(req *HostServiceCapabilityTenantRequest) []byte {
	var content []byte
	if req == nil {
		return content
	}
	if req.TenantID != 0 {
		content = appendVarintField(content, 1, uint64(req.TenantID))
	}
	return content
}

// UnmarshalHostServiceCapabilityTenantRequest decodes one tenant request.
func UnmarshalHostServiceCapabilityTenantRequest(data []byte) (*HostServiceCapabilityTenantRequest, error) {
	out := &HostServiceCapabilityTenantRequest{}
	content := data
	for len(content) > 0 {
		fieldNumber, wireType, length := protowire.ConsumeTag(content)
		if length < 0 {
			return nil, gerror.New("failed to decode capability tenant request tag")
		}
		content = content[length:]
		switch fieldNumber {
		case 1:
			value, size := protowire.ConsumeVarint(content)
			if size < 0 {
				return nil, gerror.New("failed to decode capability tenant request tenantId")
			}
			out.TenantID = int(value)
			content = content[size:]
		default:
			size := protowire.ConsumeFieldValue(fieldNumber, wireType, content)
			if size < 0 {
				return nil, gerror.New("failed to skip unknown capability tenant request field")
			}
			content = content[size:]
		}
	}
	return out, nil
}

// MarshalHostServiceCapabilityUserTenantRequest encodes one user-tenant request.
func MarshalHostServiceCapabilityUserTenantRequest(req *HostServiceCapabilityUserTenantRequest) []byte {
	var content []byte
	if req == nil {
		return content
	}
	if req.UserID != 0 {
		content = appendVarintField(content, 1, uint64(req.UserID))
	}
	if req.TenantID != 0 {
		content = appendVarintField(content, 2, uint64(req.TenantID))
	}
	return content
}

// UnmarshalHostServiceCapabilityUserTenantRequest decodes one user-tenant request.
func UnmarshalHostServiceCapabilityUserTenantRequest(data []byte) (*HostServiceCapabilityUserTenantRequest, error) {
	out := &HostServiceCapabilityUserTenantRequest{}
	content := data
	for len(content) > 0 {
		fieldNumber, wireType, length := protowire.ConsumeTag(content)
		if length < 0 {
			return nil, gerror.New("failed to decode capability user tenant request tag")
		}
		content = content[length:]
		switch fieldNumber {
		case 1:
			value, size := protowire.ConsumeVarint(content)
			if size < 0 {
				return nil, gerror.New("failed to decode capability user tenant request userId")
			}
			out.UserID = int(value)
			content = content[size:]
		case 2:
			value, size := protowire.ConsumeVarint(content)
			if size < 0 {
				return nil, gerror.New("failed to decode capability user tenant request tenantId")
			}
			out.TenantID = int(value)
			content = content[size:]
		default:
			size := protowire.ConsumeFieldValue(fieldNumber, wireType, content)
			if size < 0 {
				return nil, gerror.New("failed to skip unknown capability user tenant request field")
			}
			content = content[size:]
		}
	}
	return out, nil
}

// MarshalHostServiceCapabilityUserTenantSwitchRequest encodes one tenant switch request.
func MarshalHostServiceCapabilityUserTenantSwitchRequest(req *HostServiceCapabilityUserTenantSwitchRequest) []byte {
	var content []byte
	if req == nil {
		return content
	}
	if req.UserID != 0 {
		content = appendVarintField(content, 1, uint64(req.UserID))
	}
	if req.TargetTenantID != 0 {
		content = appendVarintField(content, 2, uint64(req.TargetTenantID))
	}
	return content
}

// UnmarshalHostServiceCapabilityUserTenantSwitchRequest decodes one tenant switch request.
func UnmarshalHostServiceCapabilityUserTenantSwitchRequest(data []byte) (*HostServiceCapabilityUserTenantSwitchRequest, error) {
	out := &HostServiceCapabilityUserTenantSwitchRequest{}
	content := data
	for len(content) > 0 {
		fieldNumber, wireType, length := protowire.ConsumeTag(content)
		if length < 0 {
			return nil, gerror.New("failed to decode capability tenant switch request tag")
		}
		content = content[length:]
		switch fieldNumber {
		case 1:
			value, size := protowire.ConsumeVarint(content)
			if size < 0 {
				return nil, gerror.New("failed to decode capability tenant switch request userId")
			}
			out.UserID = int(value)
			content = content[size:]
		case 2:
			value, size := protowire.ConsumeVarint(content)
			if size < 0 {
				return nil, gerror.New("failed to decode capability tenant switch request targetTenantId")
			}
			out.TargetTenantID = int(value)
			content = content[size:]
		default:
			size := protowire.ConsumeFieldValue(fieldNumber, wireType, content)
			if size < 0 {
				return nil, gerror.New("failed to skip unknown capability tenant switch request field")
			}
			content = content[size:]
		}
	}
	return out, nil
}

// MarshalHostServiceCapabilityJSONResponse encodes one JSON value response.
func MarshalHostServiceCapabilityJSONResponse(resp *HostServiceCapabilityJSONResponse) []byte {
	var content []byte
	if resp == nil {
		return content
	}
	if len(resp.Value) > 0 {
		content = appendBytesField(content, 1, resp.Value)
	}
	return content
}

// UnmarshalHostServiceCapabilityJSONResponse decodes one JSON value response.
func UnmarshalHostServiceCapabilityJSONResponse(data []byte) (*HostServiceCapabilityJSONResponse, error) {
	out := &HostServiceCapabilityJSONResponse{}
	content := data
	for len(content) > 0 {
		fieldNumber, wireType, length := protowire.ConsumeTag(content)
		if length < 0 {
			return nil, gerror.New("failed to decode capability JSON response tag")
		}
		content = content[length:]
		switch fieldNumber {
		case 1:
			value, size := protowire.ConsumeBytes(content)
			if size < 0 {
				return nil, gerror.New("failed to decode capability JSON response value")
			}
			out.Value = append([]byte(nil), value...)
			content = content[size:]
		default:
			size := protowire.ConsumeFieldValue(fieldNumber, wireType, content)
			if size < 0 {
				return nil, gerror.New("failed to skip unknown capability JSON response field")
			}
			content = content[size:]
		}
	}
	return out, nil
}
