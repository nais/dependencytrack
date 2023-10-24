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

func TestClient_GetConfigProperties(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/configProperty":
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"groupName":"test","propertyName":"test","propertyValue":"test","propertyType":"test","description":"test"}]`))),
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
	properties, err := c.GetConfigProperties(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, properties)
}

func TestClient_ConfigPropertyAggregate(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/configProperty/aggregate":
					assert.Equal(t, "POST", req.Method)
					var properties []ConfigProperty
					b, err := io.ReadAll(req.Body)
					assert.NoError(t, err)
					err = json.Unmarshal(b, &properties)
					if err != nil {
						assert.Error(t, err, "unmarshalling request body")
					}
					assert.Len(t, properties, 1)
					for _, property := range properties {
						assert.NotEmpty(t, property.GroupName)
						assert.NotEmpty(t, property.PropertyName)
						assert.NotEmpty(t, property.PropertyValue)
						assert.NotEmpty(t, property.PropertyType)
						assert.NotEmpty(t, property.Description)
					}
					assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
					assert.NotEmpty(t, req.Header.Get("X-Api-Key"), "Authorization header is empty")
					b, err = json.Marshal(properties)
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
	properties, err := c.ConfigPropertyAggregate(context.Background(), []ConfigProperty{
		{
			GroupName:     "test",
			PropertyName:  "test",
			PropertyValue: "test",
			PropertyType:  "test",
			Description:   "test",
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, properties)
}
