package dependencytrack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
)

// Extension points in Dependency-Track v5. Vulnerability data sources and
// analyzers are extensions configured through the v2 API instead of v4
// config properties.
const (
	ExtensionPointVulnDataSource = "vuln-data-source"
	ExtensionPointVulnAnalyzer   = "vuln-analyzer"
)

// ExtensionConfig is the runtime configuration of a v5 extension. Its shape
// is defined by the extension's config schema. Fields marked x-secret-ref in
// the schema hold the name of a managed secret, never the secret value.
type ExtensionConfig = map[string]any

type extensionConfigBody struct {
	Config ExtensionConfig `json:"config"`
}

func (c *managementClient) GetExtensionConfig(ctx context.Context, extensionPoint, extension string) (ExtensionConfig, error) {
	return withAuthContextValue(c.auth, ctx, func(tokenCtx context.Context) (ExtensionConfig, error) {
		resp, err := c.doV2(tokenCtx, http.MethodGet, extensionConfigPath(extensionPoint, extension), nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, convertError(fmt.Errorf("unexpected status"), "GetExtensionConfig", resp)
		}

		var body extensionConfigBody
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("GetExtensionConfig: decode response: %w", err)
		}
		return body.Config, nil
	})
}

// UpdateExtensionConfig replaces the configuration of an extension. It
// reports whether the configuration changed; the server answers 304 when the
// supplied configuration equals the stored one.
func (c *managementClient) UpdateExtensionConfig(ctx context.Context, extensionPoint, extension string, config ExtensionConfig) (bool, error) {
	return withAuthContextValue(c.auth, ctx, func(tokenCtx context.Context) (bool, error) {
		resp, err := c.doV2(tokenCtx, http.MethodPut, extensionConfigPath(extensionPoint, extension), extensionConfigBody{Config: config})
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusNoContent:
			return true, nil
		case http.StatusNotModified:
			return false, nil
		default:
			return false, convertError(fmt.Errorf("unexpected status"), "UpdateExtensionConfig", resp)
		}
	})
}

// EnsureSecret creates a managed secret, or updates its value when a secret
// with the same name already exists.
func (c *managementClient) EnsureSecret(ctx context.Context, name, value string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		resp, err := c.doV2(tokenCtx, http.MethodPost, "/secrets", map[string]string{"name": name, "value": value})
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusCreated, http.StatusOK, http.StatusNoContent:
			return nil
		case http.StatusConflict:
		default:
			return convertError(fmt.Errorf("unexpected status"), "CreateSecret", resp)
		}

		resp, err = c.doV2(tokenCtx, http.MethodPatch, "/secrets/"+url.PathEscape(name), map[string]string{"value": value})
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return convertError(fmt.Errorf("unexpected status"), "UpdateSecret", resp)
		}
		return nil
	})
}

func extensionConfigPath(extensionPoint, extension string) string {
	return fmt.Sprintf("/extension-points/%s/extensions/%s/config", url.PathEscape(extensionPoint), url.PathEscape(extension))
}

// doV2 sends a request to the v2 API. The generated client only covers the v1
// API; v2 lives next to it under the same base URL, e.g. <base>/api/v2.
func (c *managementClient) doV2(ctx context.Context, method, path string, body any) (*http.Response, error) {
	cfg := c.client.GetConfig()
	if len(cfg.Servers) == 0 {
		return nil, fmt.Errorf("%s %s: no server configured", method, path)
	}
	endpoint := strings.TrimSuffix(cfg.Servers[0].URL, "/") + "/v2" + path

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("%s %s: encode request: %w", method, path, err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return nil, fmt.Errorf("%s %s: create request: %w", method, path, err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token, ok := ctx.Value(client.ContextAccessToken).(string); ok && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", method, path, err)
	}
	return resp, nil
}
