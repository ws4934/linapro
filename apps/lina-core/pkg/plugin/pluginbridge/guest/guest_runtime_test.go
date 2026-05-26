// This file tests guest runtime request execution through encoded bridge
// envelopes.

package guest

import (
	"testing"

	"lina-core/pkg/plugin/pluginbridge/protocol"
)

// TestGuestRuntimeRoundTrip verifies the guest runtime allocator and execute
// path expose one decodable bridge response.
func TestGuestRuntimeRoundTrip(t *testing.T) {
	runtime := NewGuestRuntime(func(request *protocol.BridgeRequestEnvelopeV1) (*protocol.BridgeResponseEnvelopeV1, error) {
		return protocol.NewJSONResponse(200, []byte(`{"ok":true}`)), nil
	})

	requestContent, err := protocol.EncodeRequestEnvelope(&protocol.BridgeRequestEnvelopeV1{
		PluginID: "linapro-demo-dynamic",
	})
	if err != nil {
		t.Fatalf("expected request encode to succeed, got error: %v", err)
	}

	pointer := runtime.Alloc(uint32(len(requestContent)))
	if pointer == 0 {
		t.Fatal("expected guest alloc to return non-zero pointer")
	}
	copy(runtime.RequestBuffer(), requestContent)

	responsePointer, responseLength, err := runtime.Execute(uint32(len(requestContent)))
	if err != nil {
		t.Fatalf("expected guest execute to succeed, got error: %v", err)
	}
	if responsePointer == 0 || responseLength == 0 {
		t.Fatal("expected guest execute to expose one encoded response")
	}

	response, err := protocol.DecodeResponseEnvelope(runtime.ResponseBuffer())
	if err != nil {
		t.Fatalf("expected response decode to succeed, got error: %v", err)
	}
	if response.StatusCode != 200 || string(response.Body) != `{"ok":true}` {
		t.Fatalf("unexpected guest response: %#v", response)
	}
}
