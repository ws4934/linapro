// hostservice_storage_codec.go exposes storage host service payload codecs.
// The declarations are direct aliases and must stay behavior-free so storage protocol changes remain owned by hostservice.

package protocol

import "lina-core/pkg/plugin/pluginbridge/internal/hostservice"

var (
	MarshalHostServiceStoragePutRequest      = hostservice.MarshalHostServiceStoragePutRequest
	UnmarshalHostServiceStoragePutRequest    = hostservice.UnmarshalHostServiceStoragePutRequest
	MarshalHostServiceStoragePutResponse     = hostservice.MarshalHostServiceStoragePutResponse
	UnmarshalHostServiceStoragePutResponse   = hostservice.UnmarshalHostServiceStoragePutResponse
	MarshalHostServiceStorageGetRequest      = hostservice.MarshalHostServiceStorageGetRequest
	UnmarshalHostServiceStorageGetRequest    = hostservice.UnmarshalHostServiceStorageGetRequest
	MarshalHostServiceStorageGetResponse     = hostservice.MarshalHostServiceStorageGetResponse
	UnmarshalHostServiceStorageGetResponse   = hostservice.UnmarshalHostServiceStorageGetResponse
	MarshalHostServiceStorageDeleteRequest   = hostservice.MarshalHostServiceStorageDeleteRequest
	UnmarshalHostServiceStorageDeleteRequest = hostservice.UnmarshalHostServiceStorageDeleteRequest
	MarshalHostServiceStorageListRequest     = hostservice.MarshalHostServiceStorageListRequest
	UnmarshalHostServiceStorageListRequest   = hostservice.UnmarshalHostServiceStorageListRequest
	MarshalHostServiceStorageListResponse    = hostservice.MarshalHostServiceStorageListResponse
	UnmarshalHostServiceStorageListResponse  = hostservice.UnmarshalHostServiceStorageListResponse
	MarshalHostServiceStorageStatRequest     = hostservice.MarshalHostServiceStorageStatRequest
	UnmarshalHostServiceStorageStatRequest   = hostservice.UnmarshalHostServiceStorageStatRequest
	MarshalHostServiceStorageStatResponse    = hostservice.MarshalHostServiceStorageStatResponse
	UnmarshalHostServiceStorageStatResponse  = hostservice.UnmarshalHostServiceStorageStatResponse
)
