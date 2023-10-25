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

func TestClient_CreateTeam(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/team":
					assert.Equal(t, "PUT", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"uuid":"1","name":"test"}`))),
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
	team, err := c.CreateTeam(context.Background(), "test", []Permission{})
	assert.NoError(t, err)
	assert.NotEmpty(t, team)
	assert.Equal(t, "test", team.Name)
	assert.Equal(t, "1", team.Uuid)
}

func TestClient_DeleteTeam(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				if req.URL.Path == "/api/v1/team" && req.Method == "GET" {
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"uuid":"1","name":"test"}]`))),
					}
				} else if req.URL.Path == "/api/v1/team" && req.Method == "DELETE" {
					assert.Equal(t, "DELETE", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusNoContent,
						Body:       io.NopCloser(bytes.NewReader([]byte(``))),
					}
				} else {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
					}
				}
			}),
	)
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err := c.DeleteTeam(context.Background(), "1")
	assert.NoError(t, err)
}
