package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

func (c *client) GenerateApiKey(ctx context.Context, uuid string) (string, error) {
	res, err := c.put(ctx, fmt.Sprintf("%s/team/%s/key", c.baseUrl, uuid), c.authSource, nil)
	if err != nil {
		return "", err
	}
	var apiKey ApiKey
	err = json.Unmarshal(res, &apiKey)
	if err != nil {
		return "", err
	}
	return apiKey.Key, nil
}

func (c *client) ChangeAdminPassword(ctx context.Context, oldPassword, newPassword string) error {
	data := url.Values{
		"username":        {"admin"},
		"password":        {oldPassword},
		"newPassword":     {newPassword},
		"confirmPassword": {newPassword},
	}

	_, err := c.httpClient.SendRequest(ctx, http.MethodPost, c.baseUrl+"/api/v1/user/forceChangePassword", map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"text/plain"},
	}, []byte(data.Encode()))
	if err != nil {
		return err
	}
	return nil
}

func (c *client) CreateAdminUsers(ctx context.Context, users *AdminUsers, teamUuid string) error {
	for _, user := range users.Users {
		err := c.CreateManagedUser(ctx, user.Username, user.Password)
		if err != nil {
			e, ok := err.(*httpclient.RequestError)
			if !ok {
				return err
			} else {
				if e.StatusCode != http.StatusConflict {
					return err
				} else {
					log.Infof("admin user %s already exists", user.Username)
				}
			}
		}
		err = c.AddToTeam(ctx, user.Username, teamUuid)
		if err != nil {
			return err
		}
		c.log.Infof("created admin user %s", user.Username)
	}
	return nil
}

func (c *client) RemoveAdminUsers(ctx context.Context, users *AdminUsers) error {
	for _, user := range users.Users {
		err := c.DeleteManagedUser(ctx, user.Username)
		if err != nil {
			e, ok := err.(*httpclient.RequestError)
			if !ok {
				return err
			} else {
				if e.StatusCode != http.StatusNotFound {
					return err
				} else {
					c.log.Infof("user %s does not exist, nothing to remove", user.Username)
				}
			}
		}
		c.log.Infof("removed user %s", user.Username)
	}
	return nil
}

func (c *client) GetOidcUsers(ctx context.Context) ([]User, error) {
	b, err := c.get(ctx, c.baseUrl+"/api/v1/user/oidc", c.authSource)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(b, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *client) GetTeams(ctx context.Context) ([]Team, error) {
	b, err := c.get(ctx, c.baseUrl+"/api/v1/team", c.authSource)
	if err != nil {
		return nil, err
	}

	var teams []Team
	if err := json.Unmarshal(b, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

func (c *client) GetTeam(ctx context.Context, team string) (*Team, error) {
	teams, err := c.GetTeams(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range teams {
		if t.Name == team {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("team %s not found", team)
}

func (c *client) CreateTeam(ctx context.Context, teamName string, permissions []Permission) (*Team, error) {
	body, _ := json.Marshal(&Team{
		Name: teamName,
	})

	t := &Team{}
	b, err := c.put(ctx, c.baseUrl+"/api/v1/team", c.authSource, body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, t); err != nil {
		return nil, err
	}

	for _, p := range permissions {
		if _, err := c.post(ctx, c.baseUrl+"/api/v1/permission/"+string(p)+"/team/"+t.Uuid, c.authSource, nil); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (c *client) CreateManagedUser(ctx context.Context, username, password string) error {
	user := &NewUser{
		Username:            username,
		NewPassword:         password,
		ConfirmPassword:     password,
		Fullname:            username,
		Email:               username + EmailPostfix,
		Suspended:           false,
		ForcePasswordChange: false,
		NonExpiryPassword:   true,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshalling user: %w", err)
	}

	_, err = c.put(ctx, c.baseUrl+"/api/v1/user/managed", c.authSource, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) CreateOidcUser(ctx context.Context, email string) error {
	body, err := json.Marshal(map[string]string{
		"username": email,
		"email":    email,
	})
	if err != nil {
		return err
	}

	_, err = c.put(ctx, c.baseUrl+"/api/v1/user/oidc", c.authSource, body)

	if err != nil {
		e, ok := err.(*httpclient.RequestError)
		if !ok {
			return fmt.Errorf("create user: %w", err)
		}
		if e.StatusCode == http.StatusConflict {
			log.Infof("user %s already exists", email)
			return nil
		}
		return fmt.Errorf("create user: %w", err)

	}
	return nil
}

func (c *client) DeleteManagedUser(ctx context.Context, username string) error {
	user := &User{
		Username: username,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshalling user: %w", err)
	}

	_, err = c.delete(ctx, c.baseUrl+"/api/v1/user/managed", c.authSource, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) DeleteOidcUser(ctx context.Context, username string) error {
	user := &User{
		Username: username,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshalling user: %w", err)
	}

	_, err = c.delete(ctx, c.baseUrl+"/api/v1/user/oidc", c.authSource, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) AddToTeam(ctx context.Context, username, uuid string) error {
	_, err := c.post(ctx, c.baseUrl+"/api/v1/user/"+username+"/membership", c.authSource, []byte(`{"uuid": "`+uuid+`"}`))
	if err != nil {
		e, ok := err.(*httpclient.RequestError)
		if !ok {
			return fmt.Errorf("adding user to team: %w", err)
		}
		if e.StatusCode == http.StatusNotModified {
			log.Infof("user %s already in team", username)
			return nil
		}
		return fmt.Errorf("adding user to team: %w", err)
	}
	return nil
}

func (c *client) DeleteTeam(ctx context.Context, uuid string) error {
	body, err := json.Marshal(map[string]string{
		"uuid": uuid,
	})
	if err != nil {
		return err
	}

	_, err = c.delete(ctx, c.baseUrl+"/api/v1/team", c.authSource, body)
	if err != nil {
		e, ok := err.(*httpclient.RequestError)
		if !ok {
			return fmt.Errorf("deleting team: %w", err)
		}
		if e.StatusCode == http.StatusNotFound {
			log.Infof("team %s does not exist", uuid)
			return nil
		}
		return fmt.Errorf("deleting team: %w", err)
	}
	return nil
}

func (c *client) DeleteUserMembership(ctx context.Context, uuid, username string) error {
	body, err := json.Marshal(map[string]string{
		"uuid": uuid,
	})
	if err != nil {
		log.Errorf("marshalling body: %v", err)
		return err
	}

	_, err = c.delete(ctx, c.baseUrl+"/api/v1/user/"+username+"/membership", c.authSource, body)
	if err != nil {
		e, ok := err.(*httpclient.RequestError)
		if !ok {
			return fmt.Errorf("deleting user membership: %w", err)
		}
		if e.StatusCode == http.StatusNotFound {
			log.Infof("user %s does not exist", username)
			return nil
		}
		return fmt.Errorf("deleting user membership: %w", err)

	}
	return nil
}

func (c *client) UploadProject(ctx context.Context, name, version string, bom []byte) error {
	c.log.WithFields(log.Fields{
		"name":    name,
		"version": version,
	}).Info("uploading sbom")

	body, err := json.Marshal(&BomSubmitRequest{
		ProjectName:    name,
		ProjectVersion: version,
		AutoCreate:     true,
		Bom:            base64.StdEncoding.EncodeToString(bom),
	})
	if err != nil {
		return fmt.Errorf("marshalling bom submit request: %w", err)
	}

	_, err = c.put(ctx, c.baseUrl+"/api/v1/bom", c.authSource, body)
	if err != nil {
		return fmt.Errorf("uploading bom: %w", err)
	}

	c.log.Info("sbom uploaded")
	return nil
}

func (c *client) GetProject(ctx context.Context, name, version string) (*Project, error) {
	res, err := c.get(ctx, c.baseUrl+"/api/v1/project/lookup?name="+name+"&version="+version, c.authSource)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	var project Project
	if err = json.Unmarshal(res, &project); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return &project, nil
}

func (c *client) CreateProject(ctx context.Context, name, version, group string, tags []string) (*Project, error) {
	c.log.WithFields(log.Fields{
		"group": group,
		"tags":  tags,
	}).Debug("creating project")

	t := make([]Tag, 0)
	for _, tag := range tags {
		t = append(t, Tag{
			Name: tag,
		})
	}

	pp := Project{
		Name:       name,
		Publisher:  group,
		Active:     true,
		Classifier: "APPLICATION",
		Version:    version,
		Group:      group,
		Tags:       t,
	}

	body, err := json.Marshal(pp)

	p, err := c.put(ctx, c.baseUrl+"/api/v1/project", c.authSource, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var project Project
	if err = json.Unmarshal(p, &project); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return &project, nil
}

func (c *client) UpdateProjectInfo(ctx context.Context, uuid, version, group string, tags []string) error {
	c.log.WithFields(log.Fields{
		"uuid":  uuid,
		"group": group,
		"tags":  tags,
	}).Debug("adding additional info to project")

	t := make([]Tag, 0)
	for _, tag := range tags {
		t = append(t, Tag{
			Name: tag,
		})
	}

	body, err := json.Marshal(Project{
		Publisher:  group,
		Active:     true,
		Classifier: "APPLICATION",
		Version:    version,
		Group:      group,
		Tags:       t,
	})

	_, err = c.patch(ctx, c.baseUrl+"/api/v1/project/"+uuid, c.authSource, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	c.log.Info("additional info added to project")

	return nil
}

func (c *client) DeleteProjects(ctx context.Context, name string) error {
	b, err := c.get(ctx, c.baseUrl+"/api/v1/project"+"?name="+name+"&excludeInactive=false", c.authSource)
	if err != nil {
		return fmt.Errorf("getting projects: %w", err)
	}
	var projects []Project
	if err = json.Unmarshal(b, &projects); err != nil {
		return fmt.Errorf("unmarshalling response body: %w", err)
	}

	for _, project := range projects {
		_, err := c.delete(ctx, c.baseUrl+"/api/v1/project/"+project.Uuid, c.authSource, nil)
		if err != nil {
			return fmt.Errorf("deleting project: %w", err)
		}
	}
	return nil
}

func (c *client) DeleteProject(ctx context.Context, uuid string) error {
	_, err := c.delete(ctx, c.baseUrl+"/api/v1/project/"+uuid, c.authSource, nil)
	if err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}
	return nil
}

func (c *client) PortfolioRefresh(ctx context.Context) error {
	_, err := c.get(ctx, c.baseUrl+"/api/v1/metrics/portfolio/refresh", c.authSource)
	if err != nil {
		return fmt.Errorf("refreshing portfolio: %w", err)
	}
	return nil
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
