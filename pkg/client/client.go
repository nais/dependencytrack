package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nais/dependencytrack/pkg/client/auth"
	"github.com/nais/dependencytrack/pkg/httpclient"
	log "github.com/sirupsen/logrus"
)

const EmailPostfix = "@nais.io"

type Client interface {
	AddToTeam(ctx context.Context, username, uuid string) error
	ChangeAdminPassword(ctx context.Context, oldPassword, newPassword string) error
	ConfigPropertyAggregate(ctx context.Context, properties []ConfigProperty) ([]ConfigProperty, error)
	CreateAdminUsers(ctx context.Context, users *AdminUsers, teamUuid string) error
	CreateChildProject(ctx context.Context, project *Project, name, version, group, classifier string, tags []string) (*Project, error)
	CreateManagedUser(ctx context.Context, username, password string) error
	CreateOidcUser(ctx context.Context, email string) error
	CreateProject(ctx context.Context, name, version, group string, tags []string) (*Project, error)
	CreateTeam(ctx context.Context, teamName string, permissions []Permission) (*Team, error)
	DeleteManagedUser(ctx context.Context, username string) error
	DeleteOidcUser(ctx context.Context, username string) error
	DeleteProject(ctx context.Context, uuid string) error
	DeleteProjects(ctx context.Context, name string) error
	DeleteTeam(ctx context.Context, uuid string) error
	DeleteUserMembership(ctx context.Context, uuid, username string) error
	GenerateApiKey(ctx context.Context, uuid string) (string, error)
	GetAnalysisTrail(ctx context.Context, projectUuid, componentUuid, vulnerabilityUuid string) ([]*Analysis, error)
	GetConfigProperties(ctx context.Context) ([]ConfigProperty, error)
	GetCurrentProjectMetric(ctx context.Context, projectUuid string) (*ProjectMetric, error)
	GetEcosystems(ctx context.Context) ([]string, error)
	GetFindings(ctx context.Context, projectUuid string) ([]*Finding, error)
	GetOidcUsers(ctx context.Context) ([]User, error)
	GetProject(ctx context.Context, name, version string) (*Project, error)
	GetProjectMetricsByDate(ctx context.Context, projectUuid, date string) ([]*ProjectMetric, error)
	GetProjects(ctx context.Context) ([]*Project, error)
	GetProjectById(ctx context.Context, uuid string) (*Project, error)
	GetProjectsByTag(ctx context.Context, tag string) ([]*Project, error)
	GetProjectsByPrefixedTag(ctx context.Context, prefix TagPrefix, tag string) ([]*Project, error)
	GetTeam(ctx context.Context, team string) (*Team, error)
	GetTeams(ctx context.Context) ([]Team, error)
	PortfolioRefresh(ctx context.Context) error
	RecordAnalysis(ctx context.Context, analysis *AnalysisRequest) error
	RemoveAdminUsers(ctx context.Context, users *AdminUsers) error
	TriggerAnalysis(ctx context.Context, projectUuid string) error
	UpdateProject(ctx context.Context, uuid, name, version, group string, tags []string) (*Project, error)
	UpdateProjectInfo(ctx context.Context, uuid, version, group string, tags []string) error
	UploadProject(ctx context.Context, name, version, parentUuid string, autoCreate bool, bom []byte) error
	Version(ctx context.Context) (string, error)
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

func (c *client) Version(ctx context.Context) (string, error) {
	res, err := c.httpClient.SendRequest(ctx, http.MethodGet, c.baseUrl+"/api/version", map[string][]string{}, nil)
	if err != nil {
		return "", err
	}
	return string(res.Body), nil
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
	resp, err := c.httpClient.SendRequest(ctx, http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *client) getAllFrom(ctx context.Context, url string, authSource auth.Auth, offset, size int) ([]*httpclient.Response, error) {
	headers, err := c.authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}

	headers["Accept"] = []string{"application/json"}

	u := url
	if strings.Contains(url, "?") {
		u += "&"
	} else {
		u += "?"
	}
	u += "offset=" + strconv.Itoa(offset)

	resp, err := c.httpClient.SendRequest(ctx, http.MethodGet, u, headers, nil)
	if err != nil {
		return nil, err
	}

	result := []*httpclient.Response{resp}

	header := resp.Headers.Get("X-Total-Count")
	if header != "" {
		total, err := strconv.Atoi(header)
		if err != nil {
			return nil, fmt.Errorf("parsing X-Total-Count header: %w", err)
		}

		if total > offset+size {
			nextOffset := offset + size
			nextPage, err := c.getAllFrom(ctx, url, authSource, nextOffset, size)
			if err != nil {
				return nil, err
			}

			result = append(result, nextPage...)
		}
	}

	return result, nil
}

func (c *client) post(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	resp, err := c.httpClient.SendRequest(ctx, http.MethodPost, url, headers, body)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *client) put(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	resp, err := c.httpClient.SendRequest(ctx, http.MethodPut, url, headers, body)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *client) patch(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	resp, err := c.httpClient.SendRequest(ctx, http.MethodPatch, url, headers, body)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *client) delete(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	resp, err := c.httpClient.SendRequest(ctx, http.MethodDelete, url, headers, body)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
