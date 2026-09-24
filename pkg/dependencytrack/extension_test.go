package dependencytrack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
	"github.com/sirupsen/logrus"
)

type staticTokenAuth struct{}

func (staticTokenAuth) AuthContext(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, client.ContextAccessToken, "test-token"), nil
}

func (staticTokenAuth) Login(context.Context, string, string) (string, error) {
	return "test-token", nil
}

type recordedRequest struct {
	method string
	path   string
	body   map[string]any
}

func newTestManagementClient(t *testing.T, handler func(w http.ResponseWriter, r recordedRequest)) (*managementClient, *[]recordedRequest) {
	t.Helper()
	var requests []recordedRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer test-token")
		}
		req := recordedRequest{method: r.Method, path: r.URL.EscapedPath()}
		if b, _ := io.ReadAll(r.Body); len(b) > 0 {
			if err := json.Unmarshal(b, &req.body); err != nil {
				t.Errorf("decode request body: %v", err)
			}
		}
		requests = append(requests, req)
		handler(w, req)
	}))
	t.Cleanup(server.Close)

	apiClient := client.NewAPIClient(setupConfig(server.URL+"/api", &Options{}))
	return &managementClient{
		client: apiClient,
		auth:   staticTokenAuth{},
		log:    logrus.New(),
	}, &requests
}

func TestGetExtensionConfig(t *testing.T) {
	c, requests := newTestManagementClient(t, func(w http.ResponseWriter, _ recordedRequest) {
		_, _ = w.Write([]byte(`{"config":{"enabled":false,"scanOs":false}}`))
	})

	cfg, err := c.GetExtensionConfig(context.Background(), ExtensionPointVulnAnalyzer, "trivy")
	if err != nil {
		t.Fatalf("GetExtensionConfig: %v", err)
	}

	if want := "/api/v2/extension-points/vuln-analyzer/extensions/trivy/config"; (*requests)[0].path != want {
		t.Fatalf("path = %q, want %q", (*requests)[0].path, want)
	}
	if cfg["enabled"] != false || cfg["scanOs"] != false {
		t.Fatalf("config = %v, want enabled=false scanOs=false", cfg)
	}
}

func TestUpdateExtensionConfig(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		wantChanged bool
		wantErr     bool
	}{
		{name: "updated", status: http.StatusNoContent, wantChanged: true},
		{name: "unchanged", status: http.StatusNotModified, wantChanged: false},
		{name: "invalid config", status: http.StatusBadRequest, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, requests := newTestManagementClient(t, func(w http.ResponseWriter, _ recordedRequest) {
				w.WriteHeader(tt.status)
			})

			changed, err := c.UpdateExtensionConfig(context.Background(), ExtensionPointVulnDataSource, "osv", ExtensionConfig{"enabled": true})
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateExtensionConfig error = %v, wantErr %v", err, tt.wantErr)
			}
			if changed != tt.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, tt.wantChanged)
			}

			req := (*requests)[0]
			if req.method != http.MethodPut {
				t.Fatalf("method = %s, want PUT", req.method)
			}
			config, ok := req.body["config"].(map[string]any)
			if !ok || config["enabled"] != true {
				t.Fatalf("body = %v, want config.enabled=true", req.body)
			}
		})
	}
}

func TestEnsureSecret(t *testing.T) {
	t.Run("creates missing secret", func(t *testing.T) {
		c, requests := newTestManagementClient(t, func(w http.ResponseWriter, _ recordedRequest) {
			w.WriteHeader(http.StatusCreated)
		})

		if err := c.EnsureSecret(context.Background(), "TRIVY_API_TOKEN", "s3cret"); err != nil {
			t.Fatalf("EnsureSecret: %v", err)
		}

		if len(*requests) != 1 {
			t.Fatalf("requests = %d, want 1", len(*requests))
		}
		req := (*requests)[0]
		if req.method != http.MethodPost || req.path != "/api/v2/secrets" {
			t.Fatalf("request = %s %s, want POST /api/v2/secrets", req.method, req.path)
		}
		if req.body["name"] != "TRIVY_API_TOKEN" || req.body["value"] != "s3cret" {
			t.Fatalf("body = %v", req.body)
		}
	})

	t.Run("updates existing secret", func(t *testing.T) {
		c, requests := newTestManagementClient(t, func(w http.ResponseWriter, r recordedRequest) {
			if r.method == http.MethodPost {
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})

		if err := c.EnsureSecret(context.Background(), "TRIVY_API_TOKEN", "rotated"); err != nil {
			t.Fatalf("EnsureSecret: %v", err)
		}

		if len(*requests) != 2 {
			t.Fatalf("requests = %d, want 2", len(*requests))
		}
		req := (*requests)[1]
		if req.method != http.MethodPatch || req.path != "/api/v2/secrets/TRIVY_API_TOKEN" {
			t.Fatalf("request = %s %s, want PATCH /api/v2/secrets/TRIVY_API_TOKEN", req.method, req.path)
		}
		if req.body["value"] != "rotated" {
			t.Fatalf("body = %v", req.body)
		}
	})
}
