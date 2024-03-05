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

func Test_client_AddProjectToNotification(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/notification/rule/1234/project/1234":
					assert.Equal(t, "POST", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"uuid": "1234", "name": "test"}`))),
					}
				default:
					return defaultResponseCallback(req)
				}
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	notice, err := c.AddProjectToNotification(context.Background(), "1234", "1234")
	assert.NoError(t, err)
	assert.NotNil(t, notice)
	assert.Equal(t, "1234", notice.Uuid)
}

func Test_client_CreateNotification(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/notification/rule":
					assert.Equal(t, "PUT", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"uuid": "1234", "name": "test notification", "scope": "PORTFOLIO", "notificationLevel": "INFORMATIONAL", "publisher": {"name": "test"}}`))),
					}
				default:
					return defaultResponseCallback(req)
				}
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	notice, err := c.CreateNotification(context.Background(),
		"test notification",
		PortfolioNotificationScope,
		InformationalNotificationLevel,
		Publisher{
			Name: "test",
		})
	assert.NoError(t, err)
	assert.NotNil(t, notice)
	assert.Equal(t, "1234", notice.Uuid)
	assert.Equal(t, "test notification", notice.Name)
	assert.Equal(t, PortfolioNotificationScope, notice.Scope)
	assert.Equal(t, InformationalNotificationLevel, notice.Level)
}

func Test_client_DeleteNotification(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/notification/rule":
					assert.Equal(t, "DELETE", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"uuid": "1234", "name": "test notification", "scope": "PORTFOLIO", "notificationLevel": "INFORMATIONAL", "publisher": {"name": "test"	}}`))),
					}
				default:
					return defaultResponseCallback(req)
				}
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	err := c.DeleteNotification(context.Background(), &NotificationRule{
		Uuid:  "1234",
		Name:  "test notification",
		Scope: PortfolioNotificationScope,
		Level: InformationalNotificationLevel,
		Publisher: Publisher{
			Name: "test",
		},
	})
	assert.NoError(t, err)
}

func Test_client_GetNotifications(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/notification/rule":
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"uuid": "1234", "name": "test notification", "scope": "PORTFOLIO", "notificationLevel": "INFORMATIONAL", "publisher": {"name": "test"}}]`))),
					}
				default:
					return defaultResponseCallback(req)
				}
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	nots, err := c.GetNotifications(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, nots)
	assert.Equal(t, 1, len(nots))
}

func Test_client_UpdateNotification(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/notification/rule":
					assert.Equal(t, "POST", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"uuid": "1234", "name": "test notification","scope": "PORTFOLIO", "notificationLevel": "INFORMATIONAL", "publisher": {"name": "test"	}, "logSuccessfulPublish": false, "enabled": false, "notifyChildren": false, "notifyOn": ["NEW_VULNERABLE_DEPENDENCY"]}`))),
					}
				default:
					return defaultResponseCallback(req)
				}
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	not, err := c.UpdateNotification(context.Background(), &NotificationRule{
		Uuid:  "1234",
		Name:  "test notification",
		Scope: PortfolioNotificationScope,
	})

	assert.NoError(t, err)
	assert.NotNil(t, not)
	assert.Equal(t, "1234", not.Uuid)
	assert.Equal(t, false, not.LogSuccessfulPublish)
	assert.Equal(t, false, not.Enabled)
	assert.Equal(t, false, not.NotifyChildren)
	assert.Equal(t, 1, len(not.NotifyOn))
	assert.Equal(t, NewVulnerableDependencyNotify, not.NotifyOn[0])
}

func Test_client_GetPublishers(t *testing.T) {
	httpClient := test.NewTestClient(
		routes(
			authenticate(t),
			func(req *http.Request) *http.Response {
				switch req.URL.Path {
				case "/api/v1/notification/publisher":
					assert.Equal(t, "GET", req.Method)
					assert.Equal(t, "key", req.Header.Get("X-Api-Key"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`[{"name": "slack", "description": "", "publisherClass": "slack", "defaultPublisher": false, "uuid": "123"}]`))),
					}
				default:
					return defaultResponseCallback(req)
				}
			},
		))
	c := New("http://localhost", "admin", "admin", WithHttpClient(httpClient), WithApiKeySource("Administrators"))
	pubs, err := c.GetPublishers(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, pubs)
	assert.Len(t, pubs, 1)
}

func defaultResponseCallback(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewReader([]byte("Not found: " + req.URL.Path))),
	}
}
