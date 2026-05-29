// This file exposes raw host configuration reads for trusted internal
// adapters that need business-neutral access without expanding Service.

package config

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// GetRaw returns the raw GoFrame host configuration value for key. Empty key
// and "." follow GoFrame semantics and return the full configuration snapshot.
func (s *serviceImpl) GetRaw(ctx context.Context, key string) (*gvar.Var, error) {
	value, err := g.Cfg().Get(ctx, key)
	if err != nil {
		return nil, gerror.Wrapf(err, "read host config key failed key=%s", key)
	}
	return value, nil
}
