// This file verifies plugin-scoped manifest resource reads.

package manifest

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testMetadata captures manifest metadata fixture values.
type testMetadata struct {
	// Name is the fixture metadata name.
	Name string `yaml:"name"`
	// Enabled is the fixture metadata switch.
	Enabled bool `yaml:"enabled"`
}

// TestManifestReadsDevelopmentMetadata verifies a plugin can read its own
// manifest/metadata.yaml from the development source tree.
func TestManifestReadsDevelopmentMetadata(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifestFile(t, repoRoot, "plugin-a", "metadata.yaml", "name: alpha\nenabled: true\n")

	content, err := NewFactory(repoRoot).ForPlugin("plugin-a").Get(context.Background(), "metadata.yaml")
	if err != nil {
		t.Fatalf("read metadata: %v", err)
	}
	if !strings.Contains(string(content), "alpha") {
		t.Fatalf("expected metadata content, got %s", content)
	}
}

// TestManifestScansYAML verifies YAML documents scan into caller-owned structs.
func TestManifestScansYAML(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifestFile(t, repoRoot, "plugin-a", "metadata.yaml", "name: alpha\nenabled: true\n")

	target := &testMetadata{}
	err := NewFactory(repoRoot).ForPlugin("plugin-a").Scan(context.Background(), "metadata.yaml", "", target)
	if err != nil {
		t.Fatalf("scan metadata: %v", err)
	}
	if target.Name != "alpha" || !target.Enabled {
		t.Fatalf("unexpected target: %#v", target)
	}
}

// TestManifestRejectsUnsafePaths verifies path governance rejects unsafe inputs.
func TestManifestRejectsUnsafePaths(t *testing.T) {
	svc := NewFactory(t.TempDir()).ForPlugin("plugin-a")
	for _, path := range []string{
		"",
		".",
		"../plugin-b/manifest/metadata.yaml",
		"/etc/passwd",
		"C:\\secret.yaml",
		"http://example.com/config.yaml",
		"manifest/metadata.yaml",
		"config/config.yaml",
		"sql/001-schema.sql",
		"i18n/zh-CN/plugin.json",
	} {
		if _, err := svc.Get(context.Background(), path); err == nil {
			t.Fatalf("expected path %q to be rejected", path)
		}
	}
}

// TestManifestDoesNotCrossPluginScope verifies another plugin's manifest file
// cannot be reached through relative path traversal.
func TestManifestDoesNotCrossPluginScope(t *testing.T) {
	repoRoot := t.TempDir()
	writeManifestFile(t, repoRoot, "plugin-b", "metadata.yaml", "name: beta\n")

	svc := NewFactory(repoRoot).ForPlugin("plugin-a")
	if _, err := svc.Get(context.Background(), "../plugin-b/manifest/metadata.yaml"); err == nil {
		t.Fatal("expected cross-plugin traversal to fail")
	}
	exists, err := svc.Exists(context.Background(), "metadata.yaml")
	if err != nil {
		t.Fatalf("check missing own metadata: %v", err)
	}
	if exists {
		t.Fatal("expected plugin-a metadata to be absent")
	}
}

// writeManifestFile writes one plugin manifest fixture file.
func writeManifestFile(t *testing.T, repoRoot string, pluginID string, resourcePath string, content string) {
	t.Helper()
	filePath := filepath.Join(repoRoot, "apps", "lina-plugins", pluginID, "manifest", filepath.FromSlash(resourcePath))
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatalf("create manifest fixture dir: %v", err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest fixture: %v", err)
	}
}
