package main

import (
	"context"
	"maps"
	"reflect"
	"testing"

	"github.com/nais/dependencytrack/pkg/dependencytrack"
	"github.com/nais/dependencytrack/pkg/dependencytracktest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func findUpdate(t *testing.T, updates []extensionUpdate, extension string) extensionUpdate {
	t.Helper()
	for _, u := range updates {
		if u.extension == extension {
			return u
		}
	}
	t.Fatalf("no update for extension %q", extension)
	return extensionUpdate{}
}

func hasUpdate(updates []extensionUpdate, extension string) bool {
	for _, u := range updates {
		if u.extension == extension {
			return true
		}
	}
	return false
}

func applied(u extensionUpdate, current dependencytrack.ExtensionConfig) dependencytrack.ExtensionConfig {
	cfg := dependencytrack.ExtensionConfig{}
	maps.Copy(cfg, current)
	u.apply(cfg)
	return cfg
}

func TestExtensionUpdatesTrivy(t *testing.T) {
	c := &Config{TrivyApiToken: "token", TrivyBaseURL: "http://trivy:4954", TrivyIgnoreUnfixed: true}

	u := findUpdate(t, extensionUpdates(c), "trivy")
	if u.extensionPoint != dependencytrack.ExtensionPointVulnAnalyzer {
		t.Fatalf("extension point = %q, want %q", u.extensionPoint, dependencytrack.ExtensionPointVulnAnalyzer)
	}
	if !reflect.DeepEqual(u.secrets, map[string]string{secretTrivyApiToken: "token"}) {
		t.Fatalf("secrets = %v", u.secrets)
	}

	got := applied(u, dependencytrack.ExtensionConfig{"enabled": false, "scanOs": false})
	want := dependencytrack.ExtensionConfig{
		"enabled":       true,
		"apiToken":      secretTrivyApiToken,
		"apiUrl":        "http://trivy:4954",
		"ignoreUnfixed": true,
		"scanOs":        true,
		"scanLibrary":   true,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("config = %v, want %v", got, want)
	}
}

func TestExtensionUpdatesOsvEcosystems(t *testing.T) {
	stored := dependencytrack.ExtensionConfig{"enabled": false, "ecosystems": []any{"Maven", "npm"}}

	t.Run("keeps stored ecosystems when none are configured", func(t *testing.T) {
		got := applied(findUpdate(t, extensionUpdates(&Config{GoogleOSVEnabled: true}), "osv"), stored)
		if !reflect.DeepEqual(got["ecosystems"], []any{"Maven", "npm"}) {
			t.Fatalf("ecosystems = %v, want stored list", got["ecosystems"])
		}
		if got["enabled"] != true {
			t.Fatalf("enabled = %v, want true", got["enabled"])
		}
	})

	t.Run("replaces ecosystems when configured", func(t *testing.T) {
		got := applied(findUpdate(t, extensionUpdates(&Config{GoogleOSVEnabled: true, OsvEcosystems: "npm;Go;npm"}), "osv"), stored)
		if !reflect.DeepEqual(got["ecosystems"], []string{"Go", "npm"}) {
			t.Fatalf("ecosystems = %v, want [Go npm]", got["ecosystems"])
		}
	})
}

func TestExtensionUpdatesOptionalSources(t *testing.T) {
	updates := extensionUpdates(&Config{})

	for _, extension := range []string{"github", "osv", "trivy"} {
		if hasUpdate(updates, extension) {
			t.Errorf("unexpected update for %s without its token or flag", extension)
		}
	}
	if got := applied(findUpdate(t, updates, "nvd"), nil); got["enabled"] != true {
		t.Errorf("nvd enabled = %v, want true", got["enabled"])
	}
	if got := applied(findUpdate(t, updates, "oss-index"), nil); got["enabled"] != false {
		t.Errorf("oss-index enabled = %v, want false without credentials", got["enabled"])
	}
}

func TestConfigureExtensions(t *testing.T) {
	ctx := context.Background()
	c := dependencytracktest.NewMockManagementClient(t)
	updates := extensionUpdates(&Config{TrivyApiToken: "token"})

	c.EXPECT().EnsureSecret(ctx, secretTrivyApiToken, "token").Return(nil).Once()
	c.EXPECT().GetExtensionConfig(ctx, mock.Anything, mock.Anything).Return(dependencytrack.ExtensionConfig{"enabled": false}, nil)
	c.EXPECT().UpdateExtensionConfig(ctx, dependencytrack.ExtensionPointVulnAnalyzer, "trivy",
		mock.MatchedBy(func(cfg dependencytrack.ExtensionConfig) bool {
			return cfg["enabled"] == true && cfg["scanOs"] == true && cfg["apiToken"] == secretTrivyApiToken
		})).Return(true, nil).Once()
	c.EXPECT().UpdateExtensionConfig(ctx, mock.Anything, mock.Anything, mock.Anything).Return(false, nil)

	if err := configureExtensions(ctx, c, updates, logrus.New()); err != nil {
		t.Fatalf("configureExtensions: %v", err)
	}
}
