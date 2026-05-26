// This file tests transport codecs shared by org and tenant capability host services.

package hostservice

import (
	"reflect"
	"testing"
)

// TestHostServiceCapabilityCodecsRoundTrip verifies primitive request and JSON
// response codecs used by organization and tenant host services.
func TestHostServiceCapabilityCodecsRoundTrip(t *testing.T) {
	userRequest := &HostServiceCapabilityUserRequest{UserID: 42}
	decodedUserRequest, err := UnmarshalHostServiceCapabilityUserRequest(
		MarshalHostServiceCapabilityUserRequest(userRequest),
	)
	if err != nil {
		t.Fatalf("decode user request failed: %v", err)
	}
	if !reflect.DeepEqual(decodedUserRequest, userRequest) {
		t.Fatalf("unexpected user request: %#v", decodedUserRequest)
	}

	usersRequest := &HostServiceCapabilityUsersRequest{UserIDs: []int{7, 8}}
	decodedUsersRequest, err := UnmarshalHostServiceCapabilityUsersRequest(
		MarshalHostServiceCapabilityUsersRequest(usersRequest),
	)
	if err != nil {
		t.Fatalf("decode users request failed: %v", err)
	}
	if !reflect.DeepEqual(decodedUsersRequest, usersRequest) {
		t.Fatalf("unexpected users request: %#v", decodedUsersRequest)
	}

	userTenantRequest := &HostServiceCapabilityUserTenantRequest{UserID: 42, TenantID: 3}
	decodedUserTenantRequest, err := UnmarshalHostServiceCapabilityUserTenantRequest(
		MarshalHostServiceCapabilityUserTenantRequest(userTenantRequest),
	)
	if err != nil {
		t.Fatalf("decode user tenant request failed: %v", err)
	}
	if !reflect.DeepEqual(decodedUserTenantRequest, userTenantRequest) {
		t.Fatalf("unexpected user tenant request: %#v", decodedUserTenantRequest)
	}

	switchRequest := &HostServiceCapabilityUserTenantSwitchRequest{UserID: 42, TargetTenantID: 3}
	decodedSwitchRequest, err := UnmarshalHostServiceCapabilityUserTenantSwitchRequest(
		MarshalHostServiceCapabilityUserTenantSwitchRequest(switchRequest),
	)
	if err != nil {
		t.Fatalf("decode tenant switch request failed: %v", err)
	}
	if !reflect.DeepEqual(decodedSwitchRequest, switchRequest) {
		t.Fatalf("unexpected tenant switch request: %#v", decodedSwitchRequest)
	}

	response := &HostServiceCapabilityJSONResponse{Value: []byte(`{"ok":true}`)}
	decodedResponse, err := UnmarshalHostServiceCapabilityJSONResponse(
		MarshalHostServiceCapabilityJSONResponse(response),
	)
	if err != nil {
		t.Fatalf("decode JSON response failed: %v", err)
	}
	if !reflect.DeepEqual(decodedResponse, response) {
		t.Fatalf("unexpected JSON response: %#v", decodedResponse)
	}
}
