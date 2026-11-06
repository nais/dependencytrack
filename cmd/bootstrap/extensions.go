package main

import (
	"context"
	"fmt"

	"github.com/nais/dependencytrack/pkg/dependencytrack"
	"github.com/sirupsen/logrus"
)

// Names of the managed secrets that v5 extension configs reference. They are
// identifiers, not credentials.
const (
	secretGithubAdvisoryToken = "GITHUB_ADVISORY_TOKEN" // #nosec G101 -- secret name
	secretTrivyApiToken       = "TRIVY_API_TOKEN"       // #nosec G101 -- secret name
	secretOssIndexApiToken    = "OSS_INDEX_API_TOKEN"   // #nosec G101 -- secret name
)

// extensionUpdate describes the desired runtime config of one v5 extension.
// Dependency-Track v5 replaced the v4 config properties for vulnerability data
// sources and analyzers with per-extension configs, stored in the database and
// only writable through the v2 API.
type extensionUpdate struct {
	extensionPoint string
	extension      string
	// secrets are created or updated before the config is written, keyed by
	// secret name. Config fields marked x-secret-ref hold the secret name.
	secrets map[string]string
	apply   func(dependencytrack.ExtensionConfig)
}

func extensionUpdates(c *Config) []extensionUpdate {
	updates := []extensionUpdate{
		{
			// The v4 bootstrap mirrored NVD through its feeds; v5 mirrors feeds only.
			extensionPoint: dependencytrack.ExtensionPointVulnDataSource,
			extension:      "nvd",
			apply: func(cfg dependencytrack.ExtensionConfig) {
				cfg["enabled"] = true
			},
		},
	}

	if c.GithubAdvisoryToken != "" {
		updates = append(updates, extensionUpdate{
			extensionPoint: dependencytrack.ExtensionPointVulnDataSource,
			extension:      "github",
			secrets:        map[string]string{secretGithubAdvisoryToken: c.GithubAdvisoryToken},
			apply: func(cfg dependencytrack.ExtensionConfig) {
				cfg["enabled"] = true
				cfg["apiToken"] = secretGithubAdvisoryToken
			},
		})
	}

	if c.GoogleOSVEnabled {
		ecosystems := parseEcosystemList(c.OsvEcosystems)
		updates = append(updates, extensionUpdate{
			extensionPoint: dependencytrack.ExtensionPointVulnDataSource,
			extension:      "osv",
			apply: func(cfg dependencytrack.ExtensionConfig) {
				cfg["enabled"] = true
				// v5 rejects an empty list, so keep the stored ecosystems
				// unless an explicit list is configured.
				if len(ecosystems) > 0 {
					cfg["ecosystems"] = ecosystems
				}
			},
		})
	}

	if c.TrivyApiToken != "" {
		updates = append(updates, extensionUpdate{
			extensionPoint: dependencytrack.ExtensionPointVulnAnalyzer,
			extension:      "trivy",
			secrets:        map[string]string{secretTrivyApiToken: c.TrivyApiToken},
			apply: func(cfg dependencytrack.ExtensionConfig) {
				cfg["enabled"] = true
				cfg["apiToken"] = secretTrivyApiToken
				if c.TrivyBaseURL != "" {
					cfg["apiUrl"] = c.TrivyBaseURL
				}
				cfg["ignoreUnfixed"] = c.TrivyIgnoreUnfixed
				// v4 scanned OS packages by default; v5 does not. Without this,
				// apk and deb findings disappear after the migration.
				cfg["scanOs"] = true
				cfg["scanLibrary"] = true
			},
		})
	}

	ossIndex := extensionUpdate{
		extensionPoint: dependencytrack.ExtensionPointVulnAnalyzer,
		extension:      "oss-index",
		apply: func(cfg dependencytrack.ExtensionConfig) {
			cfg["enabled"] = false
		},
	}
	if c.OssIndexApiUsername != "" && c.OssIndexApiToken != "" {
		ossIndex.secrets = map[string]string{secretOssIndexApiToken: c.OssIndexApiToken}
		ossIndex.apply = func(cfg dependencytrack.ExtensionConfig) {
			cfg["enabled"] = true
			cfg["username"] = c.OssIndexApiUsername
			cfg["apiToken"] = secretOssIndexApiToken
		}
	}
	updates = append(updates, ossIndex)

	return updates
}

func configureExtensions(ctx context.Context, c dependencytrack.ManagementClient, updates []extensionUpdate, log logrus.FieldLogger) error {
	for _, u := range updates {
		name := u.extensionPoint + "/" + u.extension

		for secret, value := range u.secrets {
			if err := c.EnsureSecret(ctx, secret, value); err != nil {
				return fmt.Errorf("ensure secret %s for %s: %w", secret, name, err)
			}
		}

		cfg, err := c.GetExtensionConfig(ctx, u.extensionPoint, u.extension)
		if err != nil {
			return fmt.Errorf("get config for %s: %w", name, err)
		}
		if cfg == nil {
			cfg = dependencytrack.ExtensionConfig{}
		}
		u.apply(cfg)

		changed, err := c.UpdateExtensionConfig(ctx, u.extensionPoint, u.extension, cfg)
		if err != nil {
			return fmt.Errorf("update config for %s: %w", name, err)
		}
		if changed {
			log.Infof("updated: %s", name)
		} else {
			log.Infof("%s already up to date", name)
		}
	}
	return nil
}
