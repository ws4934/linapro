// This file projects validated dynamic artifact manifest/config resources into
// the per-execution views consumed by plugin config and manifest host services.

package runtime

import (
	"strings"

	"lina-core/internal/service/plugin/internal/catalog"
	capabilityconfig "lina-core/pkg/plugin/capability/config"
)

// buildArtifactDefaultConfig returns the active-release default config content
// from manifest/config/config.yaml. The template config.example.yaml is never
// exposed as runtime defaults.
func buildArtifactDefaultConfig(manifest *catalog.Manifest) []byte {
	if manifest == nil || manifest.RuntimeArtifact == nil {
		return nil
	}
	for _, resource := range manifest.RuntimeArtifact.ManifestResources {
		if resource == nil {
			continue
		}
		if strings.TrimSpace(resource.Path) == "manifest/config/"+capabilityconfig.RuntimeConfigFileName {
			return append([]byte(nil), resource.Content...)
		}
	}
	return nil
}

// buildArtifactManifestResources returns declaration resources keyed relative
// to manifest/. Config, SQL, and i18n resources stay on their dedicated
// pipelines and are intentionally omitted from HostServices.Manifest().
func buildArtifactManifestResources(manifest *catalog.Manifest) map[string][]byte {
	if manifest == nil || manifest.RuntimeArtifact == nil {
		return nil
	}
	resources := make(map[string][]byte)
	for _, resource := range manifest.RuntimeArtifact.ManifestResources {
		if resource == nil {
			continue
		}
		relativePath := strings.TrimPrefix(strings.TrimSpace(resource.Path), "manifest/")
		if relativePath == "" || relativePath == resource.Path {
			continue
		}
		if isDedicatedManifestPipelinePath(relativePath) {
			continue
		}
		resources[relativePath] = append([]byte(nil), resource.Content...)
	}
	if len(resources) == 0 {
		return nil
	}
	return resources
}

// isDedicatedManifestPipelinePath reports whether a manifest-relative path
// belongs to config, SQL, or i18n resource governance instead of Manifest().
func isDedicatedManifestPipelinePath(relativePath string) bool {
	for _, reserved := range []string{"config", "sql", "i18n"} {
		if relativePath == reserved || strings.HasPrefix(relativePath, reserved+"/") {
			return true
		}
	}
	return false
}
