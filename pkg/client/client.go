package client

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/nais/dependencytrack/pkg/client/auth"
	"github.com/nais/dependencytrack/pkg/httpclient"
	log "github.com/sirupsen/logrus"
)

const EmailPostfix = "@nais.io"

type Client interface {
	ChangeAdminPassword(ctx context.Context, oldPassword, newPassword string) error
	CreateAdminUsers(ctx context.Context, users *AdminUsers, teamUuid string) error
	RemoveAdminUsers(ctx context.Context, users *AdminUsers) error
	CreateTeam(ctx context.Context, teamName string, permissions []Permission) (*Team, error)
	GetTeams(ctx context.Context) ([]Team, error)
	GetTeam(ctx context.Context, team string) (*Team, error)
	CreateOidcUser(ctx context.Context, email string) error
	CreateManagedUser(ctx context.Context, username, password string) error
	DeleteManagedUser(ctx context.Context, username string) error
	DeleteOidcUser(ctx context.Context, username string) error
	AddToTeam(ctx context.Context, username, uuid string) error
	DeleteTeam(ctx context.Context, uuid string) error
	DeleteUserMembership(ctx context.Context, uuid, username string) error
	GetOidcUsers(ctx context.Context) ([]User, error)
	UploadProject(ctx context.Context, name, version string, bom []byte) error
	GetProject(ctx context.Context, name, version string) (*Project, error)
	CreateProject(ctx context.Context, name, version, group string, tags []string) (*Project, error)
	UpdateProjectInfo(ctx context.Context, uuid, version, group string, tags []string) error
	GenerateApiKey(ctx context.Context, uuid string) (string, error)
	DeleteProjects(ctx context.Context, name string) error
	DeleteProject(ctx context.Context, uuid string) error
	PortfolioRefresh(ctx context.Context) error
	auth.Auth
}

var _ Client = &client{}

type client struct {
	baseUrl    string
	authSource auth.Auth
	httpClient *httpclient.HttpClient
	log        *log.Entry
}

type Options struct {
	authSource       auth.Auth
	log              *log.Entry
	client           *http.Client
	team             string
	responseCallback func(res *http.Response, err error)
}

type Option = func(c *Options)

func New(baseUrl, username, password string, opts ...Option) Client {
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	if o.client == nil {
		o.client = http.DefaultClient
	}
	if o.log == nil {
		o.log = log.NewEntry(log.New())
	}
	httpClient := httpclient.New(
		httpclient.WithClient(o.client),
		httpclient.WithLogger(o.log),
		httpclient.WithResponseCallback(o.responseCallback),
	)
	if o.team != "" {
		u := auth.NewUsernamePasswordSource(baseUrl, username, password, httpClient, o.log)
		o.authSource = auth.NewApiKeySource(baseUrl, o.team, u, httpClient, o.log)
	}
	if o.authSource == nil {
		o.authSource = auth.NewUsernamePasswordSource(baseUrl, username, password, httpClient, o.log)
	}
	return &client{
		baseUrl:    baseUrl,
		authSource: o.authSource,
		log:        o.log,
		httpClient: httpClient,
	}
}

func WithApiKeySource(team string) Option {
	return func(o *Options) {
		o.team = team
	}
}

func WithHttpClient(client *http.Client) Option {
	return func(o *Options) {
		o.client = client
	}
}

func WithLogger(log *log.Entry) Option {
	return func(o *Options) {
		o.log = log
	}
}

func WithAuthSource(authSource auth.Auth) Option {
	return func(o *Options) {
		o.authSource = authSource
	}
}

func WithResponseCallback(callback func(res *http.Response, err error)) Option {
	return func(o *Options) {
		o.responseCallback = callback
	}
}

func (c *client) Headers(ctx context.Context) (http.Header, error) {
	return c.authSource.Headers(ctx)
}

func IsNotFound(err error) bool {
	return hasStatusCode(err, http.StatusNotFound)
}

func IsNotModified(err error) bool {
	return hasStatusCode(err, http.StatusNotModified)
}

func IsUnauthorized(err error) bool {
	return hasStatusCode(err, http.StatusUnauthorized)
}

func IsConflict(err error) bool {
	return hasStatusCode(err, http.StatusConflict)
}

func hasStatusCode(err error, statusCode int) bool {
	r, ok := err.(*httpclient.RequestError)
	if ok {
		if r.StatusCode == statusCode {
			return true
		}
	}
	wrapped := errors.Unwrap(err)
	if wrapped != nil {
		return hasStatusCode(wrapped, statusCode)
	}
	return false
}

func (c *client) get(ctx context.Context, url string, authSource auth.Auth) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, http.MethodGet, url, headers, nil)
}

func (c *client) post(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, http.MethodPost, url, headers, body)
}

func (c *client) put(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, http.MethodPut, url, headers, body)
}

func (c *client) patch(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, http.MethodPatch, url, headers, body)
}

func (c *client) delete(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, http.MethodDelete, url, headers, body)
}
