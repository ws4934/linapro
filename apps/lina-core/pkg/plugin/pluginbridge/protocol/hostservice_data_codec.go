// hostservice_data_codec.go exposes data host service payload codecs.
// The declarations are direct aliases and must not add data-access behavior or authorization decisions.

package protocol

import "lina-core/pkg/plugin/pluginbridge/internal/hostservice"

var (
	MarshalHostServiceDataListRequest           = hostservice.MarshalHostServiceDataListRequest
	UnmarshalHostServiceDataListRequest         = hostservice.UnmarshalHostServiceDataListRequest
	MarshalHostServiceDataListResponse          = hostservice.MarshalHostServiceDataListResponse
	UnmarshalHostServiceDataListResponse        = hostservice.UnmarshalHostServiceDataListResponse
	MarshalHostServiceDataGetRequest            = hostservice.MarshalHostServiceDataGetRequest
	UnmarshalHostServiceDataGetRequest          = hostservice.UnmarshalHostServiceDataGetRequest
	MarshalHostServiceDataGetResponse           = hostservice.MarshalHostServiceDataGetResponse
	UnmarshalHostServiceDataGetResponse         = hostservice.UnmarshalHostServiceDataGetResponse
	MarshalHostServiceDataMutationRequest       = hostservice.MarshalHostServiceDataMutationRequest
	UnmarshalHostServiceDataMutationRequest     = hostservice.UnmarshalHostServiceDataMutationRequest
	MarshalHostServiceDataMutationResponse      = hostservice.MarshalHostServiceDataMutationResponse
	UnmarshalHostServiceDataMutationResponse    = hostservice.UnmarshalHostServiceDataMutationResponse
	MarshalHostServiceDataTransactionRequest    = hostservice.MarshalHostServiceDataTransactionRequest
	UnmarshalHostServiceDataTransactionRequest  = hostservice.UnmarshalHostServiceDataTransactionRequest
	MarshalHostServiceDataTransactionResponse   = hostservice.MarshalHostServiceDataTransactionResponse
	UnmarshalHostServiceDataTransactionResponse = hostservice.UnmarshalHostServiceDataTransactionResponse
)
