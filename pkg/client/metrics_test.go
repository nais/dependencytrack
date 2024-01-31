package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/nais/dependencytrack/pkg/client/test"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetCurrentProjectMetric(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/metrics/project/1234/current":
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"critical": 0, "high": 0, "medium": 0, "low": 0, "unassigned": 0, "vulnerabilities": 0, "components": 0, "suppressed": 0, "findingsTotal": 0, "findingsAudited": 0, "findingsUnaudited": 0, "findingsSuppressed": 0, "findingsCritical": 0, "findingsHigh": 0, "findingsMedium": 0, "findingsLow": 0, "findingsUnassigned": 0, "findingsVulnerable": 0, "findingsInconclusive": 0, "findingsTotalAudited": 0, "findingsTotalUnaudited": 0, "findingsTotalSuppressed": 0, "findingsTotalCritical": 0, "findingsTotalHigh": 0, "findingsTotalMedium": 0, "findingsTotalLow": 0, "findingsTotalUnassigned": 0, "findingsTotalVulnerable": 0, "findingsTotalInconclusive": 0}`))),
					}

				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
					}
				}
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	metric, err := c.GetCurrentProjectMetric(context.Background(), "1234")
	assert.NoError(t, err)
	assert.NotNil(t, metric)
	assert.Equal(t, 0, metric.Critical)
}

func TestClient_GetProjectMetricsByDate(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/metrics/project/1234/since/2021-01-01":
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"critical": 1, "high": 0, "medium": 0, "low": 0, "unassigned": 0, "vulnerabilities": 0, "components": 0, "suppressed": 0, "findingsTotal": 0, "findingsAudited": 0, "findingsUnaudited": 0, "findingsSuppressed": 0, "findingsCritical": 0, "findingsHigh": 0, "findingsMedium": 0, "findingsLow": 0, "findingsUnassigned": 0, "findingsVulnerable": 0, "findingsInconclusive": 0, "findingsTotalAudited": 0, "findingsTotalUnaudited": 0, "findingsTotalSuppressed": 0, "findingsTotalCritical": 0, "findingsTotalHigh": 0, "findingsTotalMedium": 0, "findingsTotalLow": 0, "findingsTotalUnassigned": 0, "findingsTotalVulnerable": 0, "findingsTotalInconclusive": 0}]`))),
					}

				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
					}
				}
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	metric, err := c.GetProjectMetricsByDate(context.Background(), "1234", "2021-01-01")
	assert.NoError(t, err)
	assert.NotNil(t, metric)
	assert.Equal(t, 1, len(metric))
	assert.Equal(t, 1, metric[0].Critical)
}
