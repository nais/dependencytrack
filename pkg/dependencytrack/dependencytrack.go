package dependencytrack

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nais/dependencytrack/pkg/dependencytrack/auth"
	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
	"github.com/sirupsen/logrus"
)

const EmailPostfix = "@nais.io"

var _ Client = &dependencyTrackClient{}

type Client interface {
	CreateProject(ctx context.Context, imageName, imageTag string, tags []string) (*Project, error)
	CreateProjectWithSbom(ctx context.Context, imageName, imageTag string, sbom []byte) (*UploadSbomResponse, error)
	DeleteProject(ctx context.Context, uuid string) error
	GetAnalysisTrailForImage(ctx context.Context, projectId, componentId, vulnerabilityId string) (*Analysis, error)
	GetFindings(ctx context.Context, uuid string, suppressed bool, filterSource ...string) ([]*Vulnerability, error)
	GetProject(ctx context.Context, name, version string) (*Project, error)
	GetProjects(ctx context.Context, limit, offset int32) ([]Project, error)
	IsTaskInProgress(ctx context.Context, uuid string) (bool, error)
	TriggerAnalysis(ctx context.Context, uuid string) error
	TriggerAnalysisToken(ctx context.Context, uuid string) (string, error)
	UpdateFinding(ctx context.Context, request AnalysisRequest) error
}

type ManagementClient interface {
	AddToTeam(ctx context.Context, username, uuid string) error
	AllMetricsRefresh(ctx context.Context) error
	ChangeAdminPassword(ctx context.Context, oldPassword, newPassword auth.Password) error
	ConfigPropertyAggregate(ctx context.Context, property ConfigProperty) error
	CreateAdminUser(ctx context.Context, username string, password auth.Password, teamUuid string) error
	CreateAdminUsers(ctx context.Context, users []*AdminUser, teamUuid string) error
	CreateOidcUser(ctx context.Context, email string) error
	CreateTeam(ctx context.Context, teamName string, permissions []Permission) (*Team, error)
	DeleteManagedUser(ctx context.Context, username string) error
	DeleteOidcUser(ctx context.Context, username string) error
	DeleteTeam(ctx context.Context, uuid string) error
	DeleteUserMembership(ctx context.Context, teamUuid, username string) error
	GenerateApiKey(ctx context.Context, uuid string) (string, error)
	GetConfigProperties(ctx context.Context) ([]ConfigProperty, error)
	GetEcosystems(ctx context.Context) ([]string, error)
	GetOidcUser(ctx context.Context, username string) (*User, error)
	GetOidcUsers(ctx context.Context) ([]*User, error)
	GetTeam(ctx context.Context, team string) (*Team, error)
	GetTeams(ctx context.Context) ([]*Team, error)
	ProjectMetricsRefresh(ctx context.Context, uuid string) error
	RemoveAdminUser(ctx context.Context, username string) error
	RemoveAdminUsers(ctx context.Context, users []*AdminUser) error
	Version(ctx context.Context) (string, error)
}

type ClientError struct {
	error
}

type ServerError struct {
	error
}

type dependencyTrackClient struct {
	client *client.APIClient
	auth   auth.Auth
	log    logrus.FieldLogger
}

type managementClient struct {
	client *client.APIClient
	auth   auth.Auth
	log    logrus.FieldLogger
}

type Options struct {
	Client *http.Client
}

type Option = func(*Options)

func NewClient(url string, username auth.Username, password auth.Password, log *logrus.Entry, options ...Option) (Client, error) {
	if url == "" {
		return nil, fmt.Errorf("NewClient: URL cannot be empty")
	}
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}

	clientConfig := setupConfig(url, opts)
	apiClient := client.NewAPIClient(clientConfig)
	return &dependencyTrackClient{
		client: apiClient,
		auth:   auth.NewUsernamePasswordSource(username, password, apiClient, log),
		log:    log,
	}, nil
}

func NewManagementClient(url string, username auth.Username, password auth.Password, log *logrus.Entry, options ...Option) (ManagementClient, error) {
	if url == "" {
		return nil, fmt.Errorf("NewClient: URL cannot be empty")
	}
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}

	clientConfig := setupConfig(url, opts)
	apiClient := client.NewAPIClient(clientConfig)
	return &managementClient{
		client: apiClient,
		auth:   auth.NewUsernamePasswordSource(username, password, apiClient, log),
		log:    log,
	}, nil
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(opts *Options) {
		opts.Client = httpClient
	}
}

func (c *managementClient) Version(ctx context.Context) (string, error) {
	about, resp, err := c.client.VersionAPI.GetVersion(ctx).Execute()
	if err != nil {
		if resp != nil {
			return "", convertError(err, "GetVersion", resp)
		}
		return "", fmt.Errorf("GetVersion failed: %w", err)
	}
	return about.GetVersion(), nil
}

func setupConfig(rawURL string, options *Options) *client.Configuration {
	cfg := client.NewConfiguration()
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		cfg.Scheme = "https"
	} else {
		cfg.Scheme = parsedURL.Scheme
		if cfg.Scheme == "" {
			cfg.Scheme = "https"
		}
	}

	cfg.Servers = client.ServerConfigurations{{URL: rawURL}}
	cfg.HTTPClient = options.Client

	return cfg
}

func (c *dependencyTrackClient) paginateProjects(ctx context.Context, limit, offset int32, callFunc func(ctx context.Context, offset int32) ([]Project, error)) ([]Project, error) {
	var allProjects []Project

	for {
		projects, err := callFunc(ctx, offset)
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)

		if len(projects) < int(limit) {
			break
		}

		offset += limit
	}

	return allProjects, nil
}

func (c *dependencyTrackClient) withAuthContext(ctx context.Context, fn func(ctx context.Context) error) error {
	tokenCtx, err := c.auth.AuthContext(ctx)
	if err != nil {
		return fmt.Errorf("auth error: %w", err)
	}
	return fn(tokenCtx)
}

func (c *dependencyTrackClient) AuthContext(ctx context.Context) (context.Context, error) {
	return c.auth.AuthContext(ctx)
}

func (c *managementClient) Login(ctx context.Context, username, password string) (string, error) {
	return c.auth.Login(ctx, username, password)
}

func (c *managementClient) withAuthContext(ctx context.Context, fn func(ctx context.Context) error) error {
	tokenCtx, err := c.auth.AuthContext(ctx)
	if err != nil {
		return fmt.Errorf("auth error: %w", err)
	}
	return fn(tokenCtx)
}

func (c *managementClient) AuthContext(ctx context.Context) (context.Context, error) {
	return c.auth.AuthContext(ctx)
}

func convertError(err error, msg string, resp *http.Response) error {
	switch {
	case resp == nil:
		return fmt.Errorf("%s, err=%w, statuscode=unknown, body=unknown", msg, err)
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		return ClientError{fmt.Errorf("%s, err=%w, statuscode=%d, body=%s", msg, err, resp.StatusCode, parseErrorResponseBody(resp))}
	case resp.StatusCode >= 500:
		return ServerError{fmt.Errorf("%s, err=%w, statuscode=%d, body=%s", msg, err, resp.StatusCode, parseErrorResponseBody(resp))}
	default:
		return nil
	}
}

func parseErrorResponseBody(resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return "no response body"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("failed to read response body: %v", err)
	}
	return string(body)
}

func withAuthContextValue[T any](auth auth.Auth, ctx context.Context, fn func(ctx context.Context) (T, error)) (T, error) {
	tokenCtx, err := auth.AuthContext(ctx)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("auth error: %w", err)
	}
	return fn(tokenCtx)
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func boolPtr(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func SafeBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
