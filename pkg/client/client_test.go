package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/nais/dependencytrack/pkg/client/test"
	"github.com/stretchr/testify/assert"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestClient_CreateProject(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/project":
					assert.Equal(t, http.MethodPut, req.Method)
					var project Project
					b, err := io.ReadAll(req.Body)
					err = json.Unmarshal(b, &project)
					if err != nil {
						assert.Error(t, err, "unmarshalling request body")
					}
					assert.NotEmpty(t, project.Name)
					assert.NotEmpty(t, project.Publisher)
					assert.NotEmpty(t, project.Version)
					assert.NotEmpty(t, project.Group)
					assert.NotEmpty(t, project.Tags)
					assert.True(t, project.Active)
					assert.Equal(t, "APPLICATION", project.Classifier)
					assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					project.Uuid = "1234"
					b, err = json.Marshal(project)
					if err != nil {
						assert.Error(t, err, "marshalling response body")
					}
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
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	project, err := c.CreateProject(context.Background(), "Team", "version1", "group1", []string{"tag1", "tag2"})
	assert.NoError(t, err)
	assert.Equal(t, "1234", project.Uuid)
}

func TestUploadSbom(t *testing.T) {
	att, err := os.ReadFile("testdata/attestation.json")
	assert.NoError(t, err)

	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/project/lookup":
					assert.Equal(t, http.MethodGet, req.Method)
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
				case "/api/v1/bom":
					assert.Equal(t, http.MethodPut, req.Method)
					b, err := io.ReadAll(req.Body)
					var p BomSubmitRequest
					err = json.Unmarshal(b, &p)
					if err != nil {
						assert.Error(t, err)
					}
					assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
					assert.NotEmpty(t, p.ProjectName)
					assert.NotEmpty(t, p.ProjectVersion)
					assert.Equal(t, p.AutoCreate, true)
					assert.NotEmpty(t, p.Bom)
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte{})),
					}
				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
					}
				}
			},
		),
	)

	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err = c.UploadProject(context.Background(), "project1", "1.0.1", att)
	assert.NoError(t, err)
}

func TestClient_UpdateProjectInfo(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/project/1234":
					assert.Equal(t, http.MethodPatch, req.Method)
					var tag Tags
					b, err := io.ReadAll(req.Body)
					err = json.Unmarshal(b, &tag)
					if err != nil {
						assert.Error(t, err, "unmarshalling request body")
					}
					assert.NotEmpty(t, tag.Tags)
					assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte{})),
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
	err := c.UpdateProjectInfo(context.Background(), "1234", "version1", "group1", []string{"tag1", "tag2"})
	assert.NoError(t, err)
}

func TestClient_DeleteProject(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/project/1234":
					assert.Equal(t, http.MethodDelete, req.Method)
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte("1234"))),
					}
				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
					}
				}
			}))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err := c.DeleteProject(context.Background(), "1234")
	assert.NoError(t, err)
}

func TestClient_DeleteProjects(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/project":
					assert.Equal(t, http.MethodGet, req.Method)
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"name":"project1","description":"project1","uuid":"1","team":{"name":"team1","description":"team1","uuid":"1"}}]`))),
					}
				case "/api/v1/project/1":
					assert.Equal(t, http.MethodDelete, req.Method)
					return &http.Response{
						StatusCode: http.StatusOK,
					}
				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
					}
				}
			}))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err := c.DeleteProjects(context.Background(), "project1")
	assert.NoError(t, err)
}

func TestClient_GetProject(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/project/lookup":
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"name":"project1","version":"1.0.1"}`))),
					}
				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
					}
				}
			}))

	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	p, err := c.GetProject(context.Background(), "project1", "1.0.1")
	assert.NoError(t, err)
	assert.Equal(t, p.Name, "project1")
	assert.Equal(t, p.Version, "1.0.1")
}

func TestIsHttpStatusCodes(t *testing.T) {
	t.Run("is not found", func(t *testing.T) {
		c := test.NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}
		})
		client := New("http://localhost", "", "", WithHttpClient(c))
		_, err := client.GetProject(context.Background(), "", "1.0.0")
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
		switch req.URL.Path {
		case "/api/v1/user/login":
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
		case "/api/v1/team":
			assert.Equal(t, "GET", req.Method)
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
		default:
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
