// This file applies typed data capability query plans to governed data service
// requests.

package datahost

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/internal/service/plugin/internal/catalog"
	plugindata "lina-core/pkg/plugin/capability/data"
	"lina-core/pkg/plugin/pluginbridge/protocol"
)

// decodeDataListPlan restores a typed data capability list plan.
func decodeDataListPlan(table string, request *protocol.HostServiceDataListRequest) (*plugindata.DataQueryPlan, error) {
	if request == nil || len(request.PlanJSON) == 0 {
		return nil, gerror.New("data list request requires planJson")
	}
	requestPlan, err := plugindata.UnmarshalQueryPlanJSON(request.PlanJSON)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(requestPlan.Table) == "" {
		requestPlan.Table = strings.TrimSpace(table)
	}
	if strings.TrimSpace(requestPlan.Table) != strings.TrimSpace(table) {
		return nil, gerror.Newf("data capability query plan table mismatch: %s != %s", requestPlan.Table, table)
	}
	if requestPlan.Action == "" {
		requestPlan.Action = plugindata.DataPlanActionList
	}
	if requestPlan.Action != plugindata.DataPlanActionList && requestPlan.Action != plugindata.DataPlanActionCount {
		return nil, gerror.Newf("data capability list request action is invalid: %s", requestPlan.Action)
	}
	if requestPlan.Action == plugindata.DataPlanActionList {
		if requestPlan.Page == nil {
			requestPlan.Page = &plugindata.DataPagination{PageNum: defaultDataListPageNum, PageSize: defaultDataListPageSize}
		}
		if requestPlan.Page.PageNum <= 0 {
			requestPlan.Page.PageNum = defaultDataListPageNum
		}
		if requestPlan.Page.PageSize <= 0 {
			requestPlan.Page.PageSize = defaultDataListPageSize
		}
		if requestPlan.Page.PageSize > maxDataListPageSize {
			requestPlan.Page.PageSize = maxDataListPageSize
		}
	}
	return requestPlan, plugindata.ValidateDataQueryPlan(requestPlan)
}

// decodeDataGetPlan restores a typed data capability get plan.
func decodeDataGetPlan(table string, request *protocol.HostServiceDataGetRequest) (*plugindata.DataQueryPlan, error) {
	if request == nil || len(request.PlanJSON) == 0 {
		return nil, gerror.New("data get request requires planJson")
	}
	requestPlan, err := plugindata.UnmarshalQueryPlanJSON(request.PlanJSON)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(requestPlan.Table) == "" {
		requestPlan.Table = strings.TrimSpace(table)
	}
	if strings.TrimSpace(requestPlan.Table) != strings.TrimSpace(table) {
		return nil, gerror.Newf("data capability get request table mismatch: %s != %s", requestPlan.Table, table)
	}
	if requestPlan.Action == "" {
		requestPlan.Action = plugindata.DataPlanActionGet
	}
	if requestPlan.Action != plugindata.DataPlanActionGet {
		return nil, gerror.Newf("data capability get request action is invalid: %s", requestPlan.Action)
	}
	if len(requestPlan.KeyJSON) == 0 {
		return nil, gerror.New("data key cannot be empty")
	}
	return requestPlan, plugindata.ValidateDataQueryPlan(requestPlan)
}

// applyPlanFilters applies typed data capability filters against authorized resource fields.
func applyPlanFilters(model *gdb.Model, resource *catalog.ResourceSpec, filters []*plugindata.DataFilter) (*gdb.Model, error) {
	if model == nil || resource == nil || len(filters) == 0 {
		return model, nil
	}
	for _, filter := range filters {
		if err := plugindata.ValidateDataFilter(filter); err != nil {
			return nil, err
		}
		column := resolveResourceFieldColumn(resource, filter.Field)
		if column == "" {
			return nil, gerror.Newf("data capability filter field is not authorized: %s", filter.Field)
		}
		switch filter.Operator {
		case plugindata.DataFilterOperatorEQ:
			value, err := plugindata.UnmarshalValueJSON(filter.ValueJSON)
			if err != nil {
				return nil, err
			}
			model = model.Where(column, value)
		case plugindata.DataFilterOperatorIN:
			values, err := plugindata.UnmarshalValuesJSON(filter.ValuesJSON)
			if err != nil {
				return nil, err
			}
			if len(values) == 0 {
				return nil, gerror.Newf("data capability filter %s requires at least one value", filter.Operator)
			}
			model = model.WhereIn(column, values)
		case plugindata.DataFilterOperatorLike:
			value, err := plugindata.UnmarshalValueJSON(filter.ValueJSON)
			if err != nil {
				return nil, err
			}
			model = model.WhereLike(column, "%"+fmt.Sprint(value)+"%")
		default:
			return nil, gerror.Newf("data capability filter operator is not supported: %s", filter.Operator)
		}
	}
	return model, nil
}

// buildPlanFieldArgs builds select expressions for the requested field subset.
func buildPlanFieldArgs(resource *catalog.ResourceSpec, selected []string) ([]any, error) {
	if len(selected) == 0 {
		return buildResourceFieldArgs(resource), nil
	}
	fields := make([]any, 0, len(selected))
	seen := make(map[string]struct{}, len(selected))
	for _, fieldName := range selected {
		normalizedField := strings.TrimSpace(fieldName)
		if normalizedField == "" {
			return nil, gerror.New("data capability selected field cannot be empty")
		}
		if _, ok := seen[normalizedField]; ok {
			continue
		}
		seen[normalizedField] = struct{}{}
		column := resolveResourceFieldColumn(resource, normalizedField)
		if column == "" {
			return nil, gerror.Newf("data capability selected field is not authorized: %s", normalizedField)
		}
		fields = append(fields, fmt.Sprintf("%s AS %s", column, quoteResourceAlias(normalizedField)))
	}
	return fields, nil
}

// buildPlanOrderBy builds the ORDER BY clause for the typed query plan.
func buildPlanOrderBy(resource *catalog.ResourceSpec, orders []*plugindata.DataOrder) (string, error) {
	if len(orders) == 0 {
		return buildResourceOrderBy(resource), nil
	}
	parts := make([]string, 0, len(orders))
	for _, order := range orders {
		if err := plugindata.ValidateDataOrder(order); err != nil {
			return "", err
		}
		column := resolveResourceFieldColumn(resource, order.Field)
		if column == "" {
			return "", gerror.Newf("data capability order field is not authorized: %s", order.Field)
		}
		direction := "ASC"
		if order.Direction == plugindata.DataOrderDirectionDESC {
			direction = "DESC"
		}
		parts = append(parts, column+" "+direction)
	}
	return strings.Join(parts, ", "), nil
}

// buildResourceRecordWithSelection projects only the selected logical fields from one row.
func buildResourceRecordWithSelection(recordMap map[string]interface{}, resource *catalog.ResourceSpec, selected []string) map[string]interface{} {
	if len(selected) == 0 {
		return buildResourceRecord(recordMap, resource)
	}
	row := make(map[string]interface{}, len(selected))
	seen := make(map[string]struct{}, len(selected))
	for _, fieldName := range selected {
		normalizedField := strings.TrimSpace(fieldName)
		if normalizedField == "" {
			continue
		}
		if _, ok := seen[normalizedField]; ok {
			continue
		}
		seen[normalizedField] = struct{}{}
		field := findResourceField(resource, normalizedField)
		if field == nil {
			continue
		}
		row[normalizedField] = normalizeResourceValue(resolveResourceRecordValue(recordMap, field))
	}
	return row
}

// findResourceField returns the declared resource field for one logical name.
func findResourceField(resource *catalog.ResourceSpec, fieldName string) *catalog.ResourceField {
	if resource == nil {
		return nil
	}
	targetFieldName := strings.TrimSpace(fieldName)
	for _, field := range resource.Fields {
		if field != nil && field.Name == targetFieldName {
			return field
		}
	}
	return nil
}
