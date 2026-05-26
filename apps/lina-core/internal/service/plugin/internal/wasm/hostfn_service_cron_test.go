// This file tests cron host service discovery registration flows.

package wasm

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/plugin/pluginbridge/protocol"
)

// testCronRegistrationCollector stores discovered contracts for assertions.
type testCronRegistrationCollector struct {
	items []*protocol.CronContract
	err   error
}

// Register stores one received contract or returns the configured error.
func (c *testCronRegistrationCollector) Register(contract *protocol.CronContract) error {
	if c.err != nil {
		return c.err
	}
	if contract == nil {
		return gerror.New("nil cron contract")
	}
	contractSnapshot := *contract
	c.items = append(c.items, &contractSnapshot)
	return nil
}

// TestHandleHostServiceInvokeCronRegister verifies cron discovery executions
// can register one normalized cron contract through the host service.
func TestHandleHostServiceInvokeCronRegister(t *testing.T) {
	collector := &testCronRegistrationCollector{}
	hcc := &hostCallContext{
		pluginID: "linapro-demo-dynamic",
		capabilities: map[string]struct{}{
			protocol.CapabilityCron: {},
		},
		hostServices: []*protocol.HostServiceSpec{{
			Service: protocol.HostServiceCron,
			Methods: []string{protocol.HostServiceMethodCronRegister},
		}},
		executionSource: protocol.ExecutionSourceCronDiscovery,
		cronCollector:   collector,
	}

	response := invokeCronHostService(
		t,
		hcc,
		protocol.MarshalHostServiceCronRegisterRequest(&protocol.HostServiceCronRegisterRequest{
			Contract: &protocol.CronContract{
				Name:         "heartbeat",
				Pattern:      "# */10 * * * *",
				RequestType:  "CronHeartbeatReq",
				InternalPath: "cron-heartbeat",
			},
		}),
	)
	if response.Status != protocol.HostCallStatusSuccess {
		t.Fatalf("expected cron.register success, got status=%d payload=%s", response.Status, string(response.Payload))
	}
	if len(collector.items) != 1 {
		t.Fatalf("expected one registered cron contract, got %#v", collector.items)
	}
	if collector.items[0].InternalPath != "/cron-heartbeat" || collector.items[0].Scope != protocol.CronScopeAllNode {
		t.Fatalf("expected registered cron contract to be normalized, got %#v", collector.items[0])
	}
}

// TestHandleHostServiceInvokeCronRegisterRejectsWrongExecutionSource verifies
// cron registration is only available during discovery executions.
func TestHandleHostServiceInvokeCronRegisterRejectsWrongExecutionSource(t *testing.T) {
	hcc := &hostCallContext{
		pluginID: "linapro-demo-dynamic",
		capabilities: map[string]struct{}{
			protocol.CapabilityCron: {},
		},
		hostServices: []*protocol.HostServiceSpec{{
			Service: protocol.HostServiceCron,
			Methods: []string{protocol.HostServiceMethodCronRegister},
		}},
		executionSource: protocol.ExecutionSourceCron,
		cronCollector:   &testCronRegistrationCollector{},
	}

	response := invokeCronHostService(
		t,
		hcc,
		protocol.MarshalHostServiceCronRegisterRequest(&protocol.HostServiceCronRegisterRequest{
			Contract: &protocol.CronContract{
				Name:         "heartbeat",
				Pattern:      "# */10 * * * *",
				RequestType:  "CronHeartbeatReq",
				InternalPath: "/cron-heartbeat",
			},
		}),
	)
	if response.Status != protocol.HostCallStatusCapabilityDenied {
		t.Fatalf("expected cron.register to reject non-discovery execution, got status=%d", response.Status)
	}
}

// TestHandleHostServiceInvokeCronRegisterRejectsCollectorErrors verifies
// collector validation failures are surfaced as invalid requests.
func TestHandleHostServiceInvokeCronRegisterRejectsCollectorErrors(t *testing.T) {
	hcc := &hostCallContext{
		pluginID: "linapro-demo-dynamic",
		capabilities: map[string]struct{}{
			protocol.CapabilityCron: {},
		},
		hostServices: []*protocol.HostServiceSpec{{
			Service: protocol.HostServiceCron,
			Methods: []string{protocol.HostServiceMethodCronRegister},
		}},
		executionSource: protocol.ExecutionSourceCronDiscovery,
		cronCollector: &testCronRegistrationCollector{
			err: gerror.New("duplicate cron name"),
		},
	}

	response := invokeCronHostService(
		t,
		hcc,
		protocol.MarshalHostServiceCronRegisterRequest(&protocol.HostServiceCronRegisterRequest{
			Contract: &protocol.CronContract{
				Name:         "heartbeat",
				Pattern:      "# */10 * * * *",
				RequestType:  "CronHeartbeatReq",
				InternalPath: "/cron-heartbeat",
			},
		}),
	)
	if response.Status != protocol.HostCallStatusInvalidRequest {
		t.Fatalf("expected cron.register collector error to map to invalid request, got status=%d", response.Status)
	}
}

// TestHandleHostServiceInvokeCronDiscoveryBlocksNonCronServices verifies the
// reserved discovery execution cannot call other host-service families.
func TestHandleHostServiceInvokeCronDiscoveryBlocksNonCronServices(t *testing.T) {
	hcc := &hostCallContext{
		pluginID: "linapro-demo-dynamic",
		capabilities: map[string]struct{}{
			protocol.CapabilityRuntime: {},
		},
		hostServices: []*protocol.HostServiceSpec{{
			Service: protocol.HostServiceRuntime,
			Methods: []string{protocol.HostServiceMethodRuntimeInfoNow},
		}},
		executionSource: protocol.ExecutionSourceCronDiscovery,
	}

	request := &protocol.HostServiceRequestEnvelope{
		Service: protocol.HostServiceRuntime,
		Method:  protocol.HostServiceMethodRuntimeInfoNow,
	}
	response := handleHostServiceInvoke(context.Background(), hcc, protocol.MarshalHostServiceRequestEnvelope(request))
	if response.Status != protocol.HostCallStatusCapabilityDenied {
		t.Fatalf("expected non-cron host service to be blocked during cron discovery, got status=%d", response.Status)
	}
}

// invokeCronHostService dispatches one cron host-service request and returns
// the raw response envelope for assertions.
func invokeCronHostService(
	t *testing.T,
	hcc *hostCallContext,
	payload []byte,
) *protocol.HostCallResponseEnvelope {
	t.Helper()

	request := &protocol.HostServiceRequestEnvelope{
		Service: protocol.HostServiceCron,
		Method:  protocol.HostServiceMethodCronRegister,
		Payload: payload,
	}
	return handleHostServiceInvoke(context.Background(), hcc, protocol.MarshalHostServiceRequestEnvelope(request))
}
