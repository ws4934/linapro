// This file implements guest-side governed data query builder methods.

package data

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	dataplan "lina-core/pkg/plugin/capability/data/internal/plan"
)

// Table starts one single-table governed query builder.
func (db *DB) Table(table string) *Query {
	return &Query{
		table: strings.TrimSpace(table),
		plan:  &dataplan.DataQueryPlan{Table: strings.TrimSpace(table)},
	}
}

// Fields requests one field projection.
func (q *Query) Fields(fields ...string) *Query {
	if q.err != nil {
		return q
	}
	for _, field := range fields {
		normalized := strings.TrimSpace(field)
		if normalized == "" {
			q.err = gerror.New("data capability fields contains an empty field name")
			return q
		}
		q.plan.Fields = append(q.plan.Fields, normalized)
	}
	return q
}

// where appends one typed filter clause.
func (q *Query) where(field string, operator dataplan.DataFilterOperator, value any) *Query {
	if q.err != nil {
		return q
	}
	normalizedField := strings.TrimSpace(field)
	if normalizedField == "" {
		q.err = gerror.New("data capability where field cannot be empty")
		return q
	}
	if !operator.IsValid() {
		q.err = gerror.Newf("data capability where operator is invalid: %s", operator)
		return q
	}
	var (
		filter *dataplan.DataFilter
		err    error
	)
	switch operator {
	case dataplan.DataFilterOperatorEQ:
		filter, err = dataplan.NewEQFilter(normalizedField, value)
	case dataplan.DataFilterOperatorIN:
		filter, err = dataplan.NewINFilter(normalizedField, value)
	case dataplan.DataFilterOperatorLike:
		filter, err = dataplan.NewLikeFilter(normalizedField, value)
	default:
		err = gerror.Newf("data capability where operator is unsupported: %s", operator)
	}
	if err != nil {
		q.err = err
		return q
	}
	q.plan.Filters = append(q.plan.Filters, filter)
	return q
}

// WhereEq appends one equality filter.
func (q *Query) WhereEq(field string, value any) *Query {
	return q.where(field, dataplan.DataFilterOperatorEQ, value)
}

// WhereIn appends one list-membership filter.
func (q *Query) WhereIn(field string, values any) *Query {
	return q.where(field, dataplan.DataFilterOperatorIN, values)
}

// WhereLike appends one wildcard filter.
func (q *Query) WhereLike(field string, value any) *Query {
	return q.where(field, dataplan.DataFilterOperatorLike, value)
}

// WhereKey sets the key used by get/update/delete operations.
func (q *Query) WhereKey(key any) *Query {
	if q.err != nil {
		return q
	}
	keyJSON, err := dataplan.MarshalValueJSON(key)
	if err != nil {
		q.err = err
		return q
	}
	q.plan.KeyJSON = keyJSON
	return q
}

// order appends one typed order clause.
func (q *Query) order(field string, direction dataplan.DataOrderDirection) *Query {
	if q.err != nil {
		return q
	}
	normalizedField := strings.TrimSpace(field)
	if normalizedField == "" {
		q.err = gerror.New("data capability order field cannot be empty")
		return q
	}
	if !direction.IsValid() {
		q.err = gerror.Newf("data capability order direction is invalid: %s", direction)
		return q
	}
	q.plan.Orders = append(q.plan.Orders, &dataplan.DataOrder{Field: normalizedField, Direction: direction})
	return q
}

// OrderAsc appends one ascending order clause.
func (q *Query) OrderAsc(field string) *Query {
	return q.order(field, dataplan.DataOrderDirectionASC)
}

// OrderDesc appends one descending order clause.
func (q *Query) OrderDesc(field string) *Query {
	return q.order(field, dataplan.DataOrderDirectionDESC)
}

// Page applies one paging window.
func (q *Query) Page(pageNum int32, pageSize int32) *Query {
	if q.err != nil {
		return q
	}
	q.plan.Page = &dataplan.DataPagination{PageNum: pageNum, PageSize: pageSize}
	return q
}
