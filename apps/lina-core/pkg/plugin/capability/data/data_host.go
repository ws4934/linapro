// This file exposes structured data query-plan contracts shared by guest
// builders and host-side data service execution.

package data

import dataplan "lina-core/pkg/plugin/capability/data/internal/plan"

// DataQueryPlan represents one governed typed data query plan.
type DataQueryPlan = dataplan.DataQueryPlan

// DataFilter represents one governed typed filter clause.
type DataFilter = dataplan.DataFilter

// DataOrder represents one governed typed order clause.
type DataOrder = dataplan.DataOrder

// DataPagination represents one governed typed page window.
type DataPagination = dataplan.DataPagination

// DataPlanAction represents one governed data plan action.
type DataPlanAction = dataplan.DataPlanAction

// DataFilterOperator represents one governed filter operator.
type DataFilterOperator = dataplan.DataFilterOperator

// DataOrderDirection represents one governed order direction.
type DataOrderDirection = dataplan.DataOrderDirection

const (
	// DataPlanActionList lists records from one authorized table.
	DataPlanActionList = dataplan.DataPlanActionList
	// DataPlanActionGet reads one record by key from one authorized table.
	DataPlanActionGet = dataplan.DataPlanActionGet
	// DataPlanActionCount counts records from one authorized table.
	DataPlanActionCount = dataplan.DataPlanActionCount
	// DataOrderDirectionDESC orders records in descending order.
	DataOrderDirectionDESC = dataplan.DataOrderDirectionDESC
	// DataFilterOperatorEQ compares one field by equality.
	DataFilterOperatorEQ = dataplan.DataFilterOperatorEQ
	// DataFilterOperatorIN compares one field against a value list.
	DataFilterOperatorIN = dataplan.DataFilterOperatorIN
	// DataFilterOperatorLike compares one field by wildcard matching.
	DataFilterOperatorLike = dataplan.DataFilterOperatorLike
)

// UnmarshalQueryPlanJSON decodes one governed typed query plan.
func UnmarshalQueryPlanJSON(data []byte) (*DataQueryPlan, error) {
	return dataplan.UnmarshalQueryPlanJSON(data)
}

// ValidateDataQueryPlan validates one governed typed query plan.
func ValidateDataQueryPlan(plan *DataQueryPlan) error {
	return dataplan.ValidateDataQueryPlan(plan)
}

// ValidateDataFilter validates one governed typed filter clause.
func ValidateDataFilter(filter *DataFilter) error {
	return dataplan.ValidateDataFilter(filter)
}

// ValidateDataOrder validates one governed typed order clause.
func ValidateDataOrder(order *DataOrder) error {
	return dataplan.ValidateDataOrder(order)
}

// UnmarshalValueJSON decodes one JSON-encoded scalar or object value.
func UnmarshalValueJSON(data []byte) (any, error) {
	return dataplan.UnmarshalValueJSON(data)
}

// UnmarshalValuesJSON decodes one list of JSON-encoded values.
func UnmarshalValuesJSON(items [][]byte) ([]any, error) {
	return dataplan.UnmarshalValuesJSON(items)
}
