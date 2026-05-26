// This file implements whitelisted public host config key access.

package hostconfig

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Get returns the raw public host config value for a whitelisted key.
func (s *serviceAdapter) Get(ctx context.Context, key string) (*gvar.Var, error) {
	normalizedKey, err := normalizeHostConfigKey(key)
	if err != nil {
		return nil, err
	}
	value, err := s.valueForKey(ctx, normalizedKey)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return gvar.New(value), nil
}

// Exists reports whether a whitelisted public host config key is available.
func (s *serviceAdapter) Exists(ctx context.Context, key string) (bool, error) {
	value, err := s.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return value != nil && !value.IsNil(), nil
}

// String reads a public host config string value or returns defaultValue when the key is absent or blank.
func (s *serviceAdapter) String(ctx context.Context, key string, defaultValue string) (string, error) {
	value, err := s.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if value == nil || value.IsNil() || strings.TrimSpace(value.String()) == "" {
		return defaultValue, nil
	}
	return value.String(), nil
}

// Bool reads a public host config bool value or returns defaultValue when the key is absent.
func (s *serviceAdapter) Bool(ctx context.Context, key string, defaultValue bool) (bool, error) {
	value, err := s.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if value == nil || value.IsNil() {
		return defaultValue, nil
	}
	return value.Bool(), nil
}

// Int reads a public host config int value or returns defaultValue when the key is absent.
func (s *serviceAdapter) Int(ctx context.Context, key string, defaultValue int) (int, error) {
	value, err := s.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if value == nil || value.IsNil() {
		return defaultValue, nil
	}
	return value.Int(), nil
}

// Duration reads a public host config duration value or returns defaultValue when the key is absent or blank.
func (s *serviceAdapter) Duration(ctx context.Context, key string, defaultValue time.Duration) (time.Duration, error) {
	value, err := s.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if value == nil || value.IsNil() {
		return defaultValue, nil
	}
	raw := strings.TrimSpace(value.String())
	if raw == "" {
		return defaultValue, nil
	}
	duration, err := time.ParseDuration(raw)
	if err != nil {
		return 0, gerror.Wrapf(err, "parse host config %s duration failed", key)
	}
	return duration, nil
}

// normalizeHostConfigKey rejects root and blank lookups.
func normalizeHostConfigKey(key string) (string, error) {
	normalized := strings.TrimSpace(key)
	if normalized == "" || normalized == "." {
		return "", gerror.New("host config key cannot be empty or root")
	}
	return normalized, nil
}

// valueForKey returns a whitelisted public host config value.
func (s *serviceAdapter) valueForKey(ctx context.Context, key string) (any, error) {
	if s == nil || s.configSvc == nil {
		return nil, gerror.New("host config service is not configured")
	}
	switch key {
	case "workspace.basePath":
		return s.configSvc.GetWorkspaceBasePath(ctx), nil
	case "i18n.default":
		cfg := s.configSvc.GetI18n(ctx)
		if cfg == nil {
			return nil, nil
		}
		return cfg.Default, nil
	case "i18n.enabled":
		cfg := s.configSvc.GetI18n(ctx)
		if cfg == nil {
			return nil, nil
		}
		return cfg.Enabled, nil
	default:
		return nil, gerror.Newf("host config key is not public: %s", key)
	}
}
