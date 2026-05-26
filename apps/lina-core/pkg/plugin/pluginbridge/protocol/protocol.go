// Package protocol exposes the low-level dynamic plugin bridge wire protocol.
// It is the single public owner for bridge envelopes, ABI constants, host-call
// payloads, host-service payloads, and protocol codecs. Higher-level packages
// may use these entries in method signatures and implementations, but they
// must not re-export them as a second public protocol surface.
package protocol
