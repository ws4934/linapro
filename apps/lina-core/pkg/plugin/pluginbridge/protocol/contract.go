// contract.go exposes bridge contract aliases, constants, and validators through the public protocol facade.
// Keep these declarations as direct one-to-one aliases so protocol callers do not depend on internal bridge subpackages.

package protocol

import "lina-core/pkg/plugin/pluginbridge/contract"

type (
	BridgeFailureV1          = contract.BridgeFailureV1
	BridgeRequestEnvelopeV1  = contract.BridgeRequestEnvelopeV1
	BridgeResponseEnvelopeV1 = contract.BridgeResponseEnvelopeV1
	BridgeSpec               = contract.BridgeSpec
	CronConcurrency          = contract.CronConcurrency
	CronContract             = contract.CronContract
	CronScope                = contract.CronScope
	ExecutionSource          = contract.ExecutionSource
	HTTPRequestSnapshotV1    = contract.HTTPRequestSnapshotV1
	IdentitySnapshotV1       = contract.IdentitySnapshotV1
	LifecycleContract        = contract.LifecycleContract
	LifecycleDecision        = contract.LifecycleDecision
	LifecycleOperation       = contract.LifecycleOperation
	LifecycleRequest         = contract.LifecycleRequest
	ManifestSnapshotV1       = contract.ManifestSnapshotV1
	RouteContract            = contract.RouteContract
	RouteMatchSnapshotV1     = contract.RouteMatchSnapshotV1
)

const (
	CodecProtobuf                 = contract.CodecProtobuf
	AccessPublic                  = contract.AccessPublic
	AccessLogin                   = contract.AccessLogin
	RuntimeKindWasm               = contract.RuntimeKindWasm
	ABIVersionV1                  = contract.ABIVersionV1
	SupportedABIVersion           = contract.SupportedABIVersion
	DefaultGuestAllocExport       = contract.DefaultGuestAllocExport
	DefaultGuestExecuteExport     = contract.DefaultGuestExecuteExport
	BridgeFailureCodeUnauthorized = contract.BridgeFailureCodeUnauthorized
	BridgeFailureCodeForbidden    = contract.BridgeFailureCodeForbidden
	BridgeFailureCodeBadRequest   = contract.BridgeFailureCodeBadRequest
	BridgeFailureCodeNotFound     = contract.BridgeFailureCodeNotFound
	BridgeFailureCodeInternal     = contract.BridgeFailureCodeInternal

	DefaultCronContractTimezone          = contract.DefaultCronContractTimezone
	DefaultCronContractTimeoutSeconds    = contract.DefaultCronContractTimeoutSeconds
	DeclaredCronRouteBasePath            = contract.DeclaredCronRouteBasePath
	DeclaredCronRegistrationInternalPath = contract.DeclaredCronRegistrationInternalPath
	DeclaredCronRegistrationRequestType  = contract.DeclaredCronRegistrationRequestType
	CronScopeMasterOnly                  = contract.CronScopeMasterOnly
	CronScopeAllNode                     = contract.CronScopeAllNode
	CronConcurrencySingleton             = contract.CronConcurrencySingleton
	CronConcurrencyParallel              = contract.CronConcurrencyParallel

	ExecutionSourceRoute         = contract.ExecutionSourceRoute
	ExecutionSourceHook          = contract.ExecutionSourceHook
	ExecutionSourceCron          = contract.ExecutionSourceCron
	ExecutionSourceCronDiscovery = contract.ExecutionSourceCronDiscovery
	ExecutionSourceLifecycle     = contract.ExecutionSourceLifecycle

	LifecycleOperationBeforeInstall           = contract.LifecycleOperationBeforeInstall
	LifecycleOperationAfterInstall            = contract.LifecycleOperationAfterInstall
	LifecycleOperationBeforeUpgrade           = contract.LifecycleOperationBeforeUpgrade
	LifecycleOperationUpgrade                 = contract.LifecycleOperationUpgrade
	LifecycleOperationAfterUpgrade            = contract.LifecycleOperationAfterUpgrade
	LifecycleOperationBeforeDisable           = contract.LifecycleOperationBeforeDisable
	LifecycleOperationAfterDisable            = contract.LifecycleOperationAfterDisable
	LifecycleOperationBeforeUninstall         = contract.LifecycleOperationBeforeUninstall
	LifecycleOperationUninstall               = contract.LifecycleOperationUninstall
	LifecycleOperationAfterUninstall          = contract.LifecycleOperationAfterUninstall
	LifecycleOperationBeforeTenantDisable     = contract.LifecycleOperationBeforeTenantDisable
	LifecycleOperationAfterTenantDisable      = contract.LifecycleOperationAfterTenantDisable
	LifecycleOperationBeforeTenantDelete      = contract.LifecycleOperationBeforeTenantDelete
	LifecycleOperationAfterTenantDelete       = contract.LifecycleOperationAfterTenantDelete
	LifecycleOperationBeforeInstallModeChange = contract.LifecycleOperationBeforeInstallModeChange
	LifecycleOperationAfterInstallModeChange  = contract.LifecycleOperationAfterInstallModeChange
)

var (
	ValidateRouteContracts        = contract.ValidateRouteContracts
	NormalizeBridgeSpec           = contract.NormalizeBridgeSpec
	ValidateBridgeSpec            = contract.ValidateBridgeSpec
	NormalizeLifecycleContract    = contract.NormalizeLifecycleContract
	ValidateLifecycleContracts    = contract.ValidateLifecycleContracts
	IsSupportedLifecycleOperation = contract.IsSupportedLifecycleOperation
	NormalizeCronScope            = contract.NormalizeCronScope
	NormalizeCronConcurrency      = contract.NormalizeCronConcurrency
	NormalizeCronContract         = contract.NormalizeCronContract
	BuildPluginCronHandlerRef     = contract.BuildPluginCronHandlerRef
	BuildDeclaredCronRoutePath    = contract.BuildDeclaredCronRoutePath
	ValidateCronContracts         = contract.ValidateCronContracts
	NormalizeExecutionSource      = contract.NormalizeExecutionSource
)
