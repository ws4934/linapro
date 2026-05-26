// This file tests storage host service authorization, path isolation, and
// logical path prefix matching.

package wasm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"lina-core/internal/service/config"
	"lina-core/pkg/plugin/pluginbridge/protocol"
)

// trackingStorageConfig records dynamic storage path reads for shared-instance
// wiring tests.
type trackingStorageConfig struct {
	rootPath string
	calls    int
}

// GetPluginDynamicStoragePath records and returns the configured root path.
func (s *trackingStorageConfig) GetPluginDynamicStoragePath(context.Context) string {
	s.calls++
	return s.rootPath
}

// TestHandleHostServiceInvokeStorageLifecycle verifies storage put/get/list/
// delete/stat behavior against the plugin-scoped storage root.
func TestHandleHostServiceInvokeStorageLifecycle(t *testing.T) {
	storageRoot := t.TempDir()
	config.SetPluginDynamicStoragePathOverride(storageRoot)
	t.Cleanup(func() {
		config.SetPluginDynamicStoragePathOverride("")
	})

	authorizedPath := "reports/"
	hcc := newStorageHostCallContext([]string{authorizedPath})

	putResponse := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStoragePut,
		"reports/demo.json",
		protocol.MarshalHostServiceStoragePutRequest(&protocol.HostServiceStoragePutRequest{
			Path:        "reports/demo.json",
			Body:        []byte(`{"ok":true}`),
			ContentType: "application/json",
			Overwrite:   false,
		}),
	)
	if putResponse.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("put: expected success, got status=%d payload=%s", putResponse.Status, string(putResponse.Payload))
	}
	putPayload, err := protocol.UnmarshalHostServiceStoragePutResponse(putResponse.Payload)
	if err != nil {
		t.Fatalf("put payload decode failed: %v", err)
	}
	if putPayload.Object == nil || putPayload.Object.Path != "reports/demo.json" {
		t.Fatalf("put object: got %#v", putPayload.Object)
	}

	absolutePath := filepath.Join(
		storageRoot,
		storageHostServiceRootDirName,
		storageHostServiceDirName,
		hcc.pluginID,
		"reports",
		"demo.json",
	)
	content, err := os.ReadFile(absolutePath)
	if err != nil {
		t.Fatalf("expected written file to exist, got error: %v", err)
	}
	if string(content) != `{"ok":true}` {
		t.Fatalf("written content: got %q", content)
	}

	getResponse := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStorageGet,
		"reports/demo.json",
		protocol.MarshalHostServiceStorageGetRequest(&protocol.HostServiceStorageGetRequest{Path: "reports/demo.json"}),
	)
	if getResponse.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("get: expected success, got status=%d payload=%s", getResponse.Status, string(getResponse.Payload))
	}
	getPayload, err := protocol.UnmarshalHostServiceStorageGetResponse(getResponse.Payload)
	if err != nil {
		t.Fatalf("get payload decode failed: %v", err)
	}
	if !getPayload.Found || string(getPayload.Body) != `{"ok":true}` {
		t.Fatalf("get payload: got %#v", getPayload)
	}

	listResponse := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStorageList,
		"reports",
		protocol.MarshalHostServiceStorageListRequest(&protocol.HostServiceStorageListRequest{
			Prefix: "reports",
			Limit:  10,
		}),
	)
	if listResponse.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("list: expected success, got status=%d payload=%s", listResponse.Status, string(listResponse.Payload))
	}
	listPayload, err := protocol.UnmarshalHostServiceStorageListResponse(listResponse.Payload)
	if err != nil {
		t.Fatalf("list payload decode failed: %v", err)
	}
	if len(listPayload.Objects) != 1 || listPayload.Objects[0].Path != "reports/demo.json" {
		t.Fatalf("list payload: got %#v", listPayload.Objects)
	}

	deleteResponse := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStorageDelete,
		"reports/demo.json",
		protocol.MarshalHostServiceStorageDeleteRequest(&protocol.HostServiceStorageDeleteRequest{Path: "reports/demo.json"}),
	)
	if deleteResponse.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("delete: expected success, got status=%d payload=%s", deleteResponse.Status, string(deleteResponse.Payload))
	}

	statResponse := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStorageStat,
		"reports/demo.json",
		protocol.MarshalHostServiceStorageStatRequest(&protocol.HostServiceStorageStatRequest{Path: "reports/demo.json"}),
	)
	if statResponse.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("stat: expected success, got status=%d payload=%s", statResponse.Status, string(statResponse.Payload))
	}
	statPayload, err := protocol.UnmarshalHostServiceStorageStatResponse(statResponse.Payload)
	if err != nil {
		t.Fatalf("stat payload decode failed: %v", err)
	}
	if statPayload.Found {
		t.Fatalf("stat: expected object to be deleted, got %#v", statPayload.Object)
	}
}

// TestHandleHostServiceInvokeStorageLifecycleWithRelativeStorageRoot verifies
// relative storage root overrides resolve correctly before file operations.
func TestHandleHostServiceInvokeStorageLifecycleWithRelativeStorageRoot(t *testing.T) {
	storageRoot := t.TempDir()
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("resolve working directory failed: %v", err)
	}
	relativeStorageRoot, err := filepath.Rel(workingDir, storageRoot)
	if err != nil {
		t.Fatalf("build relative storage root failed: %v", err)
	}
	config.SetPluginDynamicStoragePathOverride(relativeStorageRoot)
	t.Cleanup(func() {
		config.SetPluginDynamicStoragePathOverride("")
	})

	hcc := newStorageHostCallContext([]string{"reports/"})
	putResponse := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStoragePut,
		"reports/demo.json",
		protocol.MarshalHostServiceStoragePutRequest(&protocol.HostServiceStoragePutRequest{
			Path:        "reports/demo.json",
			Body:        []byte(`{"ok":true}`),
			ContentType: "application/json",
		}),
	)
	if putResponse.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("put with relative storage root: expected success, got status=%d payload=%s", putResponse.Status, string(putResponse.Payload))
	}

	listResponse := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStorageList,
		"reports",
		protocol.MarshalHostServiceStorageListRequest(&protocol.HostServiceStorageListRequest{Prefix: "reports", Limit: 10}),
	)
	if listResponse.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("list with relative storage root: expected success, got status=%d payload=%s", listResponse.Status, string(listResponse.Payload))
	}
	listPayload, err := protocol.UnmarshalHostServiceStorageListResponse(listResponse.Payload)
	if err != nil {
		t.Fatalf("list payload decode failed: %v", err)
	}
	if len(listPayload.Objects) != 1 || listPayload.Objects[0].Path != "reports/demo.json" {
		t.Fatalf("list payload: got %#v", listPayload.Objects)
	}
}

// TestHandleHostServiceInvokeStorageRejectsUnauthorizedPath verifies requests
// outside the authorized logical path set are denied.
func TestHandleHostServiceInvokeStorageRejectsUnauthorizedPath(t *testing.T) {
	config.SetPluginDynamicStoragePathOverride(t.TempDir())
	t.Cleanup(func() {
		config.SetPluginDynamicStoragePathOverride("")
	})

	hcc := newStorageHostCallContext([]string{"reports/"})
	response := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStoragePut,
		"private/escape.txt",
		protocol.MarshalHostServiceStoragePutRequest(&protocol.HostServiceStoragePutRequest{
			Path: "private/escape.txt",
			Body: []byte("blocked"),
		}),
	)
	if response.Status != protocol.HostCallStatusCapabilityDenied {
		t.Fatalf("expected capability denied for unauthorized path, got status=%d payload=%s", response.Status, string(response.Payload))
	}
}

// TestHandleHostServiceInvokeStorageRejectsTargetMismatch verifies the request
// payload path must match the declared target resource reference.
func TestHandleHostServiceInvokeStorageRejectsTargetMismatch(t *testing.T) {
	config.SetPluginDynamicStoragePathOverride(t.TempDir())
	t.Cleanup(func() {
		config.SetPluginDynamicStoragePathOverride("")
	})

	hcc := newStorageHostCallContext([]string{"reports/"})
	response := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStoragePut,
		"reports/demo.json",
		protocol.MarshalHostServiceStoragePutRequest(&protocol.HostServiceStoragePutRequest{
			Path: "reports/other.json",
			Body: []byte("blocked"),
		}),
	)
	if response.Status != protocol.HostCallStatusInvalidRequest {
		t.Fatalf("expected invalid request for target mismatch, got status=%d payload=%s", response.Status, string(response.Payload))
	}
}

// TestHandleHostServiceInvokeStorageUsesConfiguredSharedConfig verifies
// storage host service dispatch reuses the explicitly configured config reader.
func TestHandleHostServiceInvokeStorageUsesConfiguredSharedConfig(t *testing.T) {
	configSvc := &trackingStorageConfig{rootPath: t.TempDir()}
	previousConfigSvc := storageConfigSvc
	if err := ConfigureStorageHostService(configSvc); err != nil {
		t.Fatalf("configure storage host service failed: %v", err)
	}
	t.Cleanup(func() {
		storageConfigSvc = previousConfigSvc
	})

	hcc := newStorageHostCallContext([]string{"reports/"})
	response := invokeStorageHostService(
		t,
		hcc,
		protocol.HostServiceMethodStoragePut,
		"reports/demo.json",
		protocol.MarshalHostServiceStoragePutRequest(&protocol.HostServiceStoragePutRequest{
			Path: "reports/demo.json",
			Body: []byte("shared"),
		}),
	)
	if response.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("put through shared storage config: expected success, got status=%d payload=%s", response.Status, string(response.Payload))
	}
	if configSvc.calls != 1 {
		t.Fatalf("expected shared storage config to be read once, got %d", configSvc.calls)
	}
	absolutePath := filepath.Join(
		configSvc.rootPath,
		storageHostServiceRootDirName,
		storageHostServiceDirName,
		hcc.pluginID,
		"reports",
		"demo.json",
	)
	content, err := os.ReadFile(absolutePath)
	if err != nil {
		t.Fatalf("expected object under shared storage root: %v", err)
	}
	if string(content) != "shared" {
		t.Fatalf("expected shared storage content, got %q", content)
	}
}

// TestConfigureStorageHostServiceRejectsNil verifies missing runtime config
// reader injection returns an error instead of silently constructing an isolated config service.
func TestConfigureStorageHostServiceRejectsNil(t *testing.T) {
	if err := ConfigureStorageHostService(nil); err == nil {
		t.Fatal("expected nil storage host service to return an error")
	}
}

// TestMatchAuthorizedStoragePath verifies logical prefix and exact-file path
// matching for authorized storage resources.
func TestMatchAuthorizedStoragePath(t *testing.T) {
	specs := []*protocol.HostServiceSpec{{
		Service: protocol.HostServiceStorage,
		Methods: []string{protocol.HostServiceMethodStorageGet},
		Paths:   []string{"reports/", "exports/daily.json"},
	}}

	if matched := matchAuthorizedStoragePath(specs, "reports/2026/summary.json"); matched != "reports/" {
		t.Fatalf("expected reports/ prefix to match, got %q", matched)
	}
	if matched := matchAuthorizedStoragePath(specs, "exports/daily.json"); matched != "exports/daily.json" {
		t.Fatalf("expected exact file path to match, got %q", matched)
	}
	if matched := matchAuthorizedStoragePath(specs, "reports-v2/demo.json"); matched != "" {
		t.Fatalf("expected sibling prefix to be rejected, got %q", matched)
	}
}

// newStorageHostCallContext constructs a storage-capable host call context for
// the provided authorized logical paths.
func newStorageHostCallContext(paths []string) *hostCallContext {
	return &hostCallContext{
		pluginID: "test-plugin-storage",
		capabilities: map[string]struct{}{
			protocol.CapabilityStorage: {},
		},
		hostServices: []*protocol.HostServiceSpec{{
			Service: protocol.HostServiceStorage,
			Methods: []string{
				protocol.HostServiceMethodStorageDelete,
				protocol.HostServiceMethodStorageGet,
				protocol.HostServiceMethodStorageList,
				protocol.HostServiceMethodStoragePut,
				protocol.HostServiceMethodStorageStat,
			},
			Paths: paths,
		}},
	}
}

// invokeStorageHostService dispatches a storage host-service request through
// the shared handler and returns the raw response envelope.
func invokeStorageHostService(
	t *testing.T,
	hcc *hostCallContext,
	method string,
	targetPath string,
	payload []byte,
) *protocol.HostCallResponseEnvelope {
	t.Helper()

	request := &protocol.HostServiceRequestEnvelope{
		Service:     protocol.HostServiceStorage,
		Method:      method,
		ResourceRef: targetPath,
		Payload:     payload,
	}
	return handleHostServiceInvoke(
		context.Background(),
		hcc,
		protocol.MarshalHostServiceRequestEnvelope(request),
	)
}
