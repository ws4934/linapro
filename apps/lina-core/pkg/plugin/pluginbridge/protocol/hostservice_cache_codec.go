// hostservice_cache_codec.go exposes cache host service payload codecs.
// Cache consistency and execution semantics remain outside this facade; this file only preserves protocol aliases.

package protocol

import "lina-core/pkg/plugin/pluginbridge/internal/hostservice"

var (
	MarshalHostServiceCacheGetRequest       = hostservice.MarshalHostServiceCacheGetRequest
	UnmarshalHostServiceCacheGetRequest     = hostservice.UnmarshalHostServiceCacheGetRequest
	MarshalHostServiceCacheGetResponse      = hostservice.MarshalHostServiceCacheGetResponse
	UnmarshalHostServiceCacheGetResponse    = hostservice.UnmarshalHostServiceCacheGetResponse
	MarshalHostServiceCacheSetRequest       = hostservice.MarshalHostServiceCacheSetRequest
	UnmarshalHostServiceCacheSetRequest     = hostservice.UnmarshalHostServiceCacheSetRequest
	MarshalHostServiceCacheSetResponse      = hostservice.MarshalHostServiceCacheSetResponse
	UnmarshalHostServiceCacheSetResponse    = hostservice.UnmarshalHostServiceCacheSetResponse
	MarshalHostServiceCacheDeleteRequest    = hostservice.MarshalHostServiceCacheDeleteRequest
	UnmarshalHostServiceCacheDeleteRequest  = hostservice.UnmarshalHostServiceCacheDeleteRequest
	MarshalHostServiceCacheIncrRequest      = hostservice.MarshalHostServiceCacheIncrRequest
	UnmarshalHostServiceCacheIncrRequest    = hostservice.UnmarshalHostServiceCacheIncrRequest
	MarshalHostServiceCacheIncrResponse     = hostservice.MarshalHostServiceCacheIncrResponse
	UnmarshalHostServiceCacheIncrResponse   = hostservice.UnmarshalHostServiceCacheIncrResponse
	MarshalHostServiceCacheExpireRequest    = hostservice.MarshalHostServiceCacheExpireRequest
	UnmarshalHostServiceCacheExpireRequest  = hostservice.UnmarshalHostServiceCacheExpireRequest
	MarshalHostServiceCacheExpireResponse   = hostservice.MarshalHostServiceCacheExpireResponse
	UnmarshalHostServiceCacheExpireResponse = hostservice.UnmarshalHostServiceCacheExpireResponse
)
