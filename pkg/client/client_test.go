package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/nais/dependencytrack/pkg/client/test"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestUploadSbom(t *testing.T) {
	att, err := os.ReadFile("testdata/attestation.json")
	assert.NoError(t, err)
	var a *in_toto.CycloneDXStatement
	err = json.Unmarshal(att, &a)
	assert.NoError(t, err)

	httpClient := test.NewTestClient(func(req *http.Request) *http.Response {
		switch req.Method {
		case http.MethodGet:
			switch req.URL.Path {
			case "/api/v1/project/lookup":
				p, err := json.Marshal(Project{
					Name:    "project1",
					Uuid:    "1234",
					Version: "1.0.1",
				})
				assert.NoError(t, err)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(p)),
				}
			case "/api/v1/team":
				assert.Equal(t, req.Method, http.MethodGet)
				assert.Equal(t, req.Header.Get("Accept"), "application/json")
				assert.Equal(t, req.Header.Get("Authorization"), "Bearer "+token)
				tt, err := json.Marshal([]Team{
					{
						Name: "Administrators",
						Uuid: "1234",
						ApiKeys: []ApiKey{
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
			}
		case http.MethodPut:
			err = requestIsValid(t, req, http.MethodPut, "/api/v1/bom")
			assert.NoError(t, err)
		case http.MethodPost:
			switch req.URL.Path {
			case "/api/v1/user/login":
				assert.Equal(t, req.Method, http.MethodPost)
				assert.Equal(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
				assert.Equal(t, req.Header.Get("Accept"), "text/plain")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(token))),
				}
			}
			err = requestIsValid(t, req, http.MethodPost, "/api/v1/user/login")
		case http.MethodPatch:
			err = requestIsValid(t, req, http.MethodPatch, "/api/v1/project/1234")
			assert.NoError(t, err)
		default:
			assert.Fail(t, "unexpected method")
		}
		assert.NoError(t, err)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
	})

	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err = c.UploadProject(context.Background(), "project1", "1.0.1", a)
	assert.NoError(t, err)
}

func TestClient_DeleteProject(t *testing.T) {
	client := test.NewTestClient(func(req *http.Request) *http.Response {
		switch req.URL.Path {
		case "/api/v1/user/login":
			assert.Equal(t, req.Method, http.MethodPost)
			assert.Equal(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
			assert.Equal(t, req.Header.Get("Accept"), "text/plain")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"))),
			}
		case "/api/v1/team":
			assert.Equal(t, req.Method, http.MethodGet)
			assert.Equal(t, req.Header.Get("Accept"), "application/json")
			assert.Equal(t, req.Header.Get("Authorization"), "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`[{"name":"team1","description":"team1","id":1, "apiKeys": [{"apiKey":"key"}]}]`))),
			}
		case "/api/v1/project/1234":
			assert.Equal(t, req.Method, http.MethodDelete)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte("1234"))),
			}
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}
		}
	})
	c := New("http://localhost", "admin", "admin", WithHttpClient(client), WithApiKeySource("team1"))
	err := c.DeleteProject(context.Background(), "1234")
	assert.NoError(t, err)
}

func TestClient_DeleteProjects(t *testing.T) {
	client := test.NewTestClient(func(req *http.Request) *http.Response {
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
			assert.Equal(t, req.Header.Get("Accept"), "application/json")
			assert.Equal(t, req.Header.Get("Authorization"), "Bearer "+token)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`[{"name":"team1","description":"team1","id":1, "apiKeys": [{"apiKey":"key"}]}]`))),
			}
		case "/api/v1/project":
			assert.Equal(t, req.Method, http.MethodGet)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`[{"name":"project1","description":"project1","id":1,"team":{"name":"team1","description":"team1","id":1}}]`))),
			}
		case "/api/v1/project/":
			assert.Equal(t, req.Method, http.MethodDelete)
			return &http.Response{
				StatusCode: http.StatusOK,
			}
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}
		}
	})
	c := New("http://localhost", "admin", "admin", WithHttpClient(client), WithApiKeySource("team1"))
	err := c.DeleteProjects(context.Background(), "project1")
	assert.NoError(t, err)
}

func TestClient_GetProject(t *testing.T) {
	client := test.NewTestClient(func(req *http.Request) *http.Response {
		switch req.URL.Path {
		case "/api/v1/user/login":
			assert.Equal(t, "POST", req.Method)
			assert.Equal(t, req.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
			assert.Equal(t, req.Header.Get("Accept"), "text/plain")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(token))),
			}
		case "/api/v1/project/lookup":
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, req.Header.Get("X-Api-Key"), "key")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"name":"project1","version":"1.0.1"}`))),
			}
		case "/api/v1/team":
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, req.Header.Get("Authorization"), "Bearer "+token)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`[{"name":"Administrators","uuid":"1234","apiKeys":[{"key":"key"}]}]`))),
			}
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}
		}
	})

	c := New("http://localhost", "admin", "admin", WithHttpClient(client), WithApiKeySource("Administrators"))
	p, err := c.GetProject(context.Background(), "project1", "1.0.1")
	assert.NoError(t, err)
	assert.Equal(t, p.Name, "project1")
	assert.Equal(t, p.Version, "1.0.1")
}

func requestIsValid(t *testing.T, r *http.Request, expectedMethod, expectedURL string) error {
	switch expectedMethod {
	case http.MethodPost:
		if expectedURL == "api/v1/user/login" {
			b, err := io.ReadAll(r.Body)
			var a map[string]string
			err = json.Unmarshal(b, &a)
			if err != nil {
				assert.Error(t, err)
			}

			assert.NotEmpty(t, a["username"], "username is empty")
			assert.NotEmpty(t, a["password"], "password is empty")
			assert.Equal(t, expectedMethod, r.Method, "request method is not POST")
			assert.Equal(t, expectedURL, r.URL.Host, "request URL is not /api/v1/user/login")
			assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		}

	case http.MethodGet:
		assert.Equal(t, expectedMethod, r.Method)
		assert.Equal(t, expectedURL, r.URL.Path)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		if expectedURL == "api/v1/user/login" {
			assert.NotEmpty(t, r.Header.Get("X-Api-Key"), "X-Api-Key header is empty")
		} else {
			assert.NotEmpty(t, r.Header.Get("Authorization"), "Authorization header is empty")
		}
	case http.MethodPut:
		b, err := io.ReadAll(r.Body)
		var p BomSubmitRequest
		err = json.Unmarshal(b, &p)
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, expectedMethod, r.Method)
		assert.Equal(t, expectedURL, r.URL.Path)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.NotEmpty(t, p.ProjectName)
		assert.NotEmpty(t, p.ProjectVersion)
		assert.Equal(t, p.AutoCreate, true)
		assert.NotEmpty(t, p.Bom)
	case http.MethodPatch:
		var tag Tags
		b, err := io.ReadAll(r.Body)
		err = json.Unmarshal(b, &tag)
		if err != nil {
			assert.Error(t, err, "unmarshalling request body")
		}
		assert.Equal(t, expectedMethod, r.Method)
		assert.Equal(t, expectedURL, r.URL.Path)
		assert.NotEmpty(t, tag.Tags)
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.NotEmpty(t, r.Header.Get("X-Api-Key"), "Authorization header is empty")
	}
	return nil
}
