// This file tests plugin dependency API DTO projections used by management
// endpoints and list/install responses.

package plugin

import (
	"testing"

	pluginsvc "lina-core/internal/service/plugin"
)

// TestBuildPluginDependencyCheckResultProjectsAllSections verifies the
// dependency projection preserves hard dependency checks, blockers, cycles,
// reverse dependents, and reverse blockers.
func TestBuildPluginDependencyCheckResultProjectsAllSections(t *testing.T) {
	out := buildPluginDependencyCheckResult(&pluginsvc.DependencyCheckResult{
		TargetID: "plugin-consumer",
		Framework: pluginsvc.DependencyFrameworkCheck{
			RequiredVersion: ">=0.1.0",
			CurrentVersion:  "v0.1.0",
			Status:          "satisfied",
		},
		Dependencies: []*pluginsvc.DependencyPluginCheck{
			{
				OwnerID:         "plugin-consumer",
				DependencyID:    "plugin-base",
				DependencyName:  "Plugin Base",
				RequiredVersion: ">=0.1.0",
				CurrentVersion:  "v0.1.0",
				Installed:       false,
				Discovered:      true,
				Status:          "missing",
				Chain:           []string{"plugin-consumer", "plugin-base"},
			},
		},
		Blockers: []*pluginsvc.DependencyBlocker{
			{
				Code:            "dependency_missing",
				PluginID:        "plugin-consumer",
				DependencyID:    "plugin-missing",
				RequiredVersion: ">=0.1.0",
				CurrentVersion:  "",
				Chain:           []string{"plugin-consumer", "plugin-missing"},
				Detail:          "missing",
			},
		},
		Cycle: []string{"a", "b", "a"},
		ReverseDependents: []*pluginsvc.DependencyReverseDependent{
			{
				PluginID:        "plugin-downstream",
				Name:            "Plugin Downstream",
				Version:         "v0.2.0",
				RequiredVersion: "<0.3.0",
			},
		},
		ReverseBlockers: []*pluginsvc.DependencyBlocker{
			{
				Code:         "dependency_snapshot_unknown",
				PluginID:     "plugin-unknown",
				DependencyID: "",
				Detail:       "unknown snapshot",
			},
		},
	})

	if out == nil {
		t.Fatal("expected projected dependency result")
	}
	if out.TargetId != "plugin-consumer" || out.Framework.Status != "satisfied" {
		t.Fatalf("unexpected top-level projection: %#v", out)
	}
	if len(out.Dependencies) != 1 || out.Dependencies[0].DependencyId != "plugin-base" {
		t.Fatalf("expected dependency edge projection, got %#v", out.Dependencies)
	}
	if len(out.Blockers) != 1 || out.Blockers[0].Code != "dependency_missing" {
		t.Fatalf("expected blocker projection, got %#v", out.Blockers)
	}
	if len(out.Cycle) != 3 || out.Cycle[2] != "a" {
		t.Fatalf("expected cycle projection, got %#v", out.Cycle)
	}
	if len(out.ReverseDependents) != 1 || out.ReverseDependents[0].PluginId != "plugin-downstream" {
		t.Fatalf("expected reverse dependent projection, got %#v", out.ReverseDependents)
	}
	if len(out.ReverseBlockers) != 1 || out.ReverseBlockers[0].Code != "dependency_snapshot_unknown" {
		t.Fatalf("expected reverse blocker projection, got %#v", out.ReverseBlockers)
	}
}
