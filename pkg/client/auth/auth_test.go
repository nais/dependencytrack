package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/nais/dependencytrack/pkg/client/test"
	"github.com/nais/dependencytrack/pkg/httpclient"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestClient_ApiKey(t *testing.T) {
	httpClient := test.NewTestClient(func(req *http.Request) *http.Response {
		switch req.URL.Path {
		case "/api/v1/user/login":
			assert.Equal(t, req.Method, http.MethodPost)
			assert.Equal(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
			assert.Equal(t, req.Header.Get("Accept"), "text/plain")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(token))),
			}
		case "/api/v1/team":
			assert.Equal(t, req.Method, http.MethodGet)
			assert.Equal(t, req.Header.Get("Authorization"), "Bearer "+token)
			tt, err := json.Marshal([]team{
				{
					Name: "Administrators",
					Uuid: "1234",
					ApiKeys: []apiKey{
						{
							Key: "key",
						},
					},
				},
			})
			assert.NoError(t, err)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(tt)),
			}
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
			}
		}
	})

	h := httpclient.New(httpclient.WithClient(httpClient))
	u := NewUsernamePasswordSource("http://localhost", "username", "password", h, nil)
	a := NewApiKeySource("http://localhost", "Administrators", u, h, nil)
	headers, err := a.Headers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "key", headers.Get("X-Api-Key"))
}
