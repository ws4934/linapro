// Package manifestview implements immutable source-plugin manifest snapshots
// returned through the public pluginhost upgrade callback contract.
package manifestview

import bridgecontract "lina-core/pkg/plugin/pluginbridge/contract"

// Snapshot is the host-owned immutable manifest snapshot view.
type Snapshot struct {
	value bridgecontract.ManifestSnapshotV1
}

// NewSnapshot creates one immutable manifest snapshot view.
func NewSnapshot(value *bridgecontract.ManifestSnapshotV1) *Snapshot {
	if value == nil {
		return nil
	}
	return &Snapshot{
		value: *value,
	}
}

// ID returns the plugin identifier recorded in the manifest snapshot.
func (s *Snapshot) ID() string {
	if s == nil {
		return ""
	}
	return s.value.ID
}

// Name returns the plugin display name recorded in the manifest snapshot.
func (s *Snapshot) Name() string {
	if s == nil {
		return ""
	}
	return s.value.Name
}

// Version returns the plugin version recorded in the manifest snapshot.
func (s *Snapshot) Version() string {
	if s == nil {
		return ""
	}
	return s.value.Version
}

// Type returns the plugin type recorded in the manifest snapshot.
func (s *Snapshot) Type() string {
	if s == nil {
		return ""
	}
	return s.value.Type
}

// Values returns a copy of the shared typed manifest snapshot contract.
func (s *Snapshot) Values() *bridgecontract.ManifestSnapshotV1 {
	if s == nil {
		return nil
	}
	value := s.value
	return &value
}
