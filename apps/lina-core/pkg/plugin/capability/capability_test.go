// This file verifies the public capability services stay limited to
// ordinary plugin consumption and does not publish host-internal seams.

package capability

import (
	"reflect"
	"testing"

	"lina-core/pkg/plugin/capability/orgcap"
	"lina-core/pkg/plugin/capability/tenantcap"
)

// TestServicesDoesNotExposeTenantFilter verifies database-query-builder
// tenant filtering stays out of the ordinary capability services.
func TestServicesDoesNotExposeTenantFilter(t *testing.T) {
	servicesType := reflect.TypeOf((*Services)(nil)).Elem()
	if _, ok := servicesType.MethodByName("TenantFilter"); ok {
		t.Fatal("capability.Services must not expose TenantFilter")
	}
}

// TestOrgServiceDoesNotExposeWorkspaceProjections verifies ordinary
// organization capability consumption stays independent from user-management
// workspace projections.
func TestOrgServiceDoesNotExposeWorkspaceProjections(t *testing.T) {
	serviceType := reflect.TypeOf((*orgcap.Service)(nil)).Elem()
	for _, method := range []string{"UserDeptTree", "ListPostOptions"} {
		if _, ok := serviceType.MethodByName(method); ok {
			t.Fatalf("orgcap.Service must not expose host workspace projection method %s", method)
		}
	}
}

// TestTenantRuntimeServiceDoesNotExposeFallback verifies unused platform
// fallback helpers stay out of the broad runtime combination interface.
func TestTenantRuntimeServiceDoesNotExposeFallback(t *testing.T) {
	serviceType := reflect.TypeOf((*tenantcap.RuntimeService)(nil)).Elem()
	if _, ok := serviceType.MethodByName("ReadWithPlatformFallback"); ok {
		t.Fatal("tenantcap.RuntimeService must not expose ReadWithPlatformFallback")
	}
}
