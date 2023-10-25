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

func TestClient_GetOidcUsers(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/user/oidc":
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"email":"test@email.com","username":"test"}]`))),
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
	users, err := c.GetOidcUsers(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, "test", users[0].Username)
	assert.Equal(t, "test@email.com", users[0].Email)
}

func TestClient_CreateOidcUser(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/user/oidc":
					assert.Equal(t, "PUT", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(nil))),
					}
				}
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
				}
			},
		),
	)
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err := c.CreateOidcUser(context.Background(), "test@email.com")
	assert.NoError(t, err)
}

func TestClient_CreateAdminUsers(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/user/admin/membership":
					assert.Equal(t, "POST", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"uuid":"1","name":"test"}`))),
					}

				case "/api/v1/user/managed":
					assert.Equal(t, "PUT", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"email":"test@email.com","username":"test"}]`))),
					}
				}
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
				}
			}),
	)
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	admins := AdminUsers{
		Users: []AdminUser{
			{
				Username: "admin",
				Password: "admin",
			},
		},
	}
	err := c.CreateAdminUsers(context.Background(), &admins, "1")
	assert.NoError(t, err)
}

func TestClient_RemoveAdminUsers(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/user/managed":
					assert.Equal(t, "DELETE", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusNoContent,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"email":"test@email.com","username":"test"}]`))),
					}
				}
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
				}
			},
		),
	)
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	admins := AdminUsers{
		Users: []AdminUser{
			{
				Username: "test",
				Password: "test",
			},
		},
	}
	err := c.RemoveAdminUsers(context.Background(), &admins)
	assert.NoError(t, err)
}

func TestClient_CreateManagedUser(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/user/managed":
					assert.Equal(t, "PUT", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"username":"test","password":"test"}`))),
					}
				}
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
				}
			},
		),
	)
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err := c.CreateManagedUser(context.Background(), "test", "test")
	assert.NoError(t, err)
}
