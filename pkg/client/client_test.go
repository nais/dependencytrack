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

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestIsHttpStatusCodes(t *testing.T) {
	t.Run("is not found", func(t *testing.T) {
		c := test.NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}
		})
		cc := New("http://localhost", "", "", WithHttpClient(c)).(*client)

		_, err := cc.get(context.Background(), "http://localhost", cc.authSource)
		assert.True(t, IsNotFound(err))
	})

	t.Run("is not modified", func(t *testing.T) {
		c := test.NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusNotModified,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}
		})
		client := New("http://localhost", "", "", WithHttpClient(c))
		_, err := client.GetProject(context.Background(), "", "1.0.0")
		assert.True(t, IsNotModified(err))
	})

	t.Run("is conflict", func(t *testing.T) {
		c := test.NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusConflict,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}
		})
		client := New("http://localhost", "", "", WithHttpClient(c))
		_, err := client.GetProject(context.Background(), "", "1.0.0")
		assert.True(t, IsConflict(err))
	})

	t.Run("is unauthorized", func(t *testing.T) {
		c := test.NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}
		})
		client := New("http://localhost", "", "", WithHttpClient(c))
		_, err := client.GetProject(context.Background(), "", "1.0.0")
		assert.True(t, IsUnauthorized(err))
	})
}

func authenticate(t *testing.T) func(req *http.Request) *http.Response {
	return func(req *http.Request) *http.Response {
		if req.URL.Path == "/api/v1/user/login" {
			assert.Equal(t, "POST", req.Method)
			assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
			assert.Equal(t, "text/plain", req.Header.Get("Accept"))
			b, err := io.ReadAll(req.Body)
			assert.NoError(t, err)

			assert.Contains(t, string(b), "username=admin")
			assert.Contains(t, string(b), "password=admin")

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(token))),
			}
		} else if req.URL.Path == "/api/v1/team" && req.Method == "GET" {
			assert.Equal(t, "Bearer "+token, req.Header.Get("Authorization"))
			teams := []Team{
				{
					Name: "Administrators",
					Uuid: "1234",
					ApiKeys: []ApiKey{
						{
							Key: "key",
						},
					},
				},
			}
			b, err := json.Marshal(teams)
			assert.NoError(t, err)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(b)),
			}
		} else {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
			}
		}
	}
}

func routes(fn ...test.RoundTripFunc) test.RoundTripFunc {
	return func(req *http.Request) *http.Response {
		for _, f := range fn {
			resp := f(req)
			if resp.StatusCode != http.StatusNotFound {
				return resp
			}
		}
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
		}
	}
}
