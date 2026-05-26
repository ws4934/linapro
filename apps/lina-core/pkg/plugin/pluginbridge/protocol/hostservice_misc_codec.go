// hostservice_misc_codec.go exposes network, manifest, notify, and cron host service payload codecs.
// These smaller protocol areas share a file to keep the facade concise without mixing in implementation logic.

package protocol

import "lina-core/pkg/plugin/pluginbridge/internal/hostservice"

var (
	MarshalHostServiceNetworkRequest        = hostservice.MarshalHostServiceNetworkRequest
	UnmarshalHostServiceNetworkRequest      = hostservice.UnmarshalHostServiceNetworkRequest
	MarshalHostServiceNetworkResponse       = hostservice.MarshalHostServiceNetworkResponse
	UnmarshalHostServiceNetworkResponse     = hostservice.UnmarshalHostServiceNetworkResponse
	MarshalHostServiceManifestGetRequest    = hostservice.MarshalHostServiceManifestGetRequest
	UnmarshalHostServiceManifestGetRequest  = hostservice.UnmarshalHostServiceManifestGetRequest
	MarshalHostServiceManifestGetResponse   = hostservice.MarshalHostServiceManifestGetResponse
	UnmarshalHostServiceManifestGetResponse = hostservice.UnmarshalHostServiceManifestGetResponse
	MarshalHostServiceNotifySendRequest     = hostservice.MarshalHostServiceNotifySendRequest
	UnmarshalHostServiceNotifySendRequest   = hostservice.UnmarshalHostServiceNotifySendRequest
	MarshalHostServiceNotifySendResponse    = hostservice.MarshalHostServiceNotifySendResponse
	UnmarshalHostServiceNotifySendResponse  = hostservice.UnmarshalHostServiceNotifySendResponse
	MarshalHostServiceCronRegisterRequest   = hostservice.MarshalHostServiceCronRegisterRequest
	UnmarshalHostServiceCronRegisterRequest = hostservice.UnmarshalHostServiceCronRegisterRequest
)
