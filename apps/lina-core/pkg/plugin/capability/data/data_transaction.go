// This file implements guest-side governed data transaction builder methods.

package data

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	dataplan "lina-core/pkg/plugin/capability/data/internal/plan"
)

// Table selects the single transaction table and returns one mutation builder.
func (tx *Tx) Table(table string) *TxQuery {
	normalizedTable := strings.TrimSpace(table)
	if tx.err == nil {
		switch {
		case normalizedTable == "":
			tx.err = gerror.New("data capability transaction table cannot be empty")
		case tx.table == "":
			tx.table = normalizedTable
		case tx.table != normalizedTable:
			tx.err = gerror.Newf("data capability transaction only supports one table per transaction: %s != %s", tx.table, normalizedTable)
		}
	}
	return &TxQuery{tx: tx, table: normalizedTable, err: tx.err}
}

// WhereKey sets the key used by update/delete in a transaction.
func (q *TxQuery) WhereKey(key any) *TxQuery {
	if q.err != nil {
		return q
	}
	keyJSON, err := dataplan.MarshalValueJSON(key)
	if err != nil {
		q.err = err
		q.tx.err = err
		return q
	}
	q.keyJSON = keyJSON
	return q
}

// Insert appends one insert mutation to the transaction.
func (q *TxQuery) Insert(record map[string]any) (*MutationResult, error) {
	return q.enqueueMutation(dataplan.DataMutationActionCreate, nil, record)
}

// Update appends one update mutation to the transaction.
func (q *TxQuery) Update(record map[string]any) (*MutationResult, error) {
	return q.enqueueMutation(dataplan.DataMutationActionUpdate, q.keyJSON, record)
}

// Delete appends one delete mutation to the transaction.
func (q *TxQuery) Delete() (*MutationResult, error) {
	return q.enqueueMutation(dataplan.DataMutationActionDelete, q.keyJSON, nil)
}

// enqueueMutation validates and appends one structured mutation plan to the
// surrounding single-table transaction builder.
func (q *TxQuery) enqueueMutation(action dataplan.DataMutationAction, keyJSON []byte, record map[string]any) (*MutationResult, error) {
	if q.err != nil {
		return nil, q.err
	}
	if q.tx == nil {
		return nil, gerror.New("data capability transaction query is not initialized")
	}
	if !action.IsValid() {
		return nil, gerror.Newf("data capability mutation action is invalid: %s", action)
	}
	if (action == dataplan.DataMutationActionUpdate || action == dataplan.DataMutationActionDelete) && len(keyJSON) == 0 {
		return nil, gerror.New("data capability update/delete in transaction requires WhereKey")
	}
	recordJSON, err := dataplan.MarshalValueJSON(record)
	if err != nil {
		q.err = err
		q.tx.err = err
		return nil, err
	}
	q.tx.operations = append(q.tx.operations, &dataplan.DataMutationPlan{
		Action:     action,
		KeyJSON:    append([]byte(nil), keyJSON...),
		RecordJSON: recordJSON,
	})
	return &MutationResult{}, nil
}
