//go:build !wasip1

// This file provides non-wasm stubs for the data capability guest facade.

package data

import (
	"github.com/gogf/gf/v2/errors/gerror"

	dataplan "lina-core/pkg/plugin/capability/data/internal/plan"
)

// One is unavailable outside wasip1 builds.
func (q *Query) One() (map[string]any, bool, error) {
	return nil, false, gerror.New("data capability guest execution is only available for wasip1 builds")
}

// All is unavailable outside wasip1 builds.
func (q *Query) All() ([]map[string]any, int32, error) {
	return nil, 0, gerror.New("data capability guest execution is only available for wasip1 builds")
}

// Count is unavailable outside wasip1 builds.
func (q *Query) Count() (int32, error) {
	return 0, gerror.New("data capability guest execution is only available for wasip1 builds")
}

// Insert is unavailable outside wasip1 builds.
func (q *Query) Insert(record map[string]any) (*MutationResult, error) {
	return nil, gerror.New("data capability guest execution is only available for wasip1 builds")
}

// Update is unavailable outside wasip1 builds.
func (q *Query) Update(record map[string]any) (*MutationResult, error) {
	return nil, gerror.New("data capability guest execution is only available for wasip1 builds")
}

// Delete is unavailable outside wasip1 builds.
func (q *Query) Delete() (*MutationResult, error) {
	return nil, gerror.New("data capability guest execution is only available for wasip1 builds")
}

// Transaction is unavailable outside wasip1 builds.
func (db *DB) Transaction(_ func(tx *Tx) error) error {
	return gerror.New("data capability guest execution is only available for wasip1 builds")
}

// Compile-time anchor to keep the internal plan package referenced in stub builds.
var _ dataplan.DataPlanAction
