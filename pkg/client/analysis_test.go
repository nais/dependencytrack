package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/nais/dependencytrack/pkg/client/test"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetAnalysisTrail(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/analysis":
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					analysis := []Analysis{
						{
							AnalysisState:         "GOOD",
							AnalysisJustification: "yes",
							AnalysisResponse:      "yes",
							AnalysisDetails:       "yes",
							AnalysisComments:      nil,
							IsSuppressed:          true,
						},
					}
					b, err := json.Marshal(analysis)
					assert.NoError(t, err)
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(b)),
					}
				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
					}
				}
			}),
	)
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	analysis, err := c.GetAnalysisTrail(context.Background(), "uuid", "0", "10")
	assert.NoError(t, err)
	assert.NotEmpty(t, analysis)
	assert.Equal(t, "GOOD", analysis[0].AnalysisState)
	assert.Equal(t, "yes", analysis[0].AnalysisJustification)
	assert.Equal(t, "yes", analysis[0].AnalysisResponse)
	assert.Equal(t, "yes", analysis[0].AnalysisDetails)
}

func TestClient_RecordAnalysis(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/analysis":
					assert.Equal(t, "PUT", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
					}
				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
					}
				}
			}),
	)
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err := c.RecordAnalysis(context.Background(), &AnalysisRequest{
		Project:               "uuid",
		Component:             "0",
		Vulnerability:         "10",
		AnalysisState:         "GOOD",
		AnalysisJustification: "yes",
		AnalysisResponse:      "yes",
		AnalysisDetails:       "yes",
		Comment:               "yes",
		IsSuppressed:          true,
	})
	assert.NoError(t, err)
}
