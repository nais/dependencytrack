package dependencytrack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nais/dependencytrack/pkg/dependencytrack/auth"
	"github.com/nais/dependencytrack/pkg/httpclient"
	log "github.com/sirupsen/logrus"
)

const EmailPostfix = "@nais.io"

type Client interface {
	ChangeAdminPassword(ctx context.Context, oldPassword string, newPassword string) error
	CreateAdminUsers(ctx context.Context, users *AdminUsers, teamUuid string) error
	RemoveAdminUsers(ctx context.Context, users *AdminUsers) error
	CreateTeam(ctx context.Context, teamName string, permissions []Permission) (*Team, error)
	GetTeams(ctx context.Context) ([]Team, error)
	GetTeamUuid(ctx context.Context, team string) (string, error)
	CreateOidcUser(ctx context.Context, email string) error
	CreateManagedUser(ctx context.Context, username, password string) error
	AddToTeam(ctx context.Context, username string, uuid string) error
	DeleteTeam(ctx context.Context, uuid string) error
	DeleteUserMembership(ctx context.Context, uuid string, username string) error
	auth.Auth
}

type client struct {
	baseUrl    string
	authSource auth.Auth
	httpClient *httpclient.HttpClient
	log        *log.Entry
}

func (c *client) Headers(ctx context.Context) (http.Header, error) {
	return c.authSource.Headers(ctx)
}

func NewClient(baseUrl string, username string, password string, c *httpclient.HttpClient, log *log.Entry) Client {
	if c == nil {
		c = httpclient.New(http.DefaultClient, log)
	}
	source := auth.NewUsernamePasswordSource(baseUrl+"/user/login", username, password, c, log)

	return &client{
		baseUrl:    baseUrl,
		authSource: source,
		httpClient: c,
		log:        log,
	}
}

func (c *client) ChangeAdminPassword(ctx context.Context, oldPassword string, newPassword string) error {
	data := url.Values{
		"username":        {"admin"},
		"password":        {oldPassword},
		"newPassword":     {newPassword},
		"confirmPassword": {newPassword},
	}

	_, err := c.httpClient.SendRequest(ctx, http.MethodPost, c.baseUrl+"/user/forceChangePassword", map[string][]string{
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
		err := c.DeleteUser(ctx, user.Username)
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

func (c *client) GetTeams(ctx context.Context) ([]Team, error) {
	b, err := c.Get(ctx, c.baseUrl+"/team", c.authSource)

	if err != nil {
		return nil, err
	}

	var teams []Team
	if err := json.Unmarshal(b, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

func (c *client) GetTeamUuid(ctx context.Context, team string) (string, error) {
	teams, err := c.GetTeams(ctx)
	if err != nil {
		return "", nil
	}

	var uuid string
	for _, t := range teams {
		if t.Name == team {
			uuid = t.Uuid
		}
	}
	if uuid == "" {
		return "", fmt.Errorf("team %s not found", team)
	}
	return uuid, nil
}

func (c *client) CreateTeam(ctx context.Context, teamName string, permissions []Permission) (*Team, error) {
	body, _ := json.Marshal(&Team{
		Name: teamName,
	})

	t := &Team{}
	b, err := c.Put(ctx, c.baseUrl+"/team", c.authSource, body)

	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, t); err != nil {
		return nil, err
	}

	for _, p := range permissions {
		if _, err := c.Post(ctx, c.baseUrl+"/permission/"+string(p)+"/team/"+t.Uuid, c.authSource, nil); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (c *client) CreateManagedUser(ctx context.Context, username, password string) error {
	user := &User{
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

	_, err = c.Put(ctx, c.baseUrl+"/user/managed", c.authSource, body)
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

	_, err = c.Put(ctx, c.baseUrl+"/user/oidc", c.authSource, body)

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

func (c *client) DeleteUser(ctx context.Context, username string) error {
	user := &User{
		Username: username,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshalling user: %w", err)
	}

	_, err = c.Delete(ctx, c.baseUrl+"/user/managed", c.authSource, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) AddToTeam(ctx context.Context, username string, uuid string) error {
	_, err := c.Post(ctx, c.baseUrl+"/user/"+username+"/membership", c.authSource, []byte(`{"uuid": "`+uuid+`"}`))

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

	_, err = c.Delete(ctx, c.baseUrl+"/team", c.authSource, body)
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

func (c *client) DeleteUserMembership(ctx context.Context, uuid string, username string) error {
	body, err := json.Marshal(map[string]string{
		"uuid": uuid,
	})
	if err != nil {
		log.Errorf("marshalling body: %v", err)
		return err
	}

	_, err = c.Delete(ctx, c.baseUrl+"/user/"+username+"/membership", c.authSource, body)
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

func (c *client) Get(ctx context.Context, url string, authSource auth.Auth) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, http.MethodGet, url, headers, nil)
}

func (c *client) Post(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, "POST", url, headers, body)
}

func (c *client) Put(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, "PUT", url, headers, body)
}

func (c *client) Delete(ctx context.Context, url string, authSource auth.Auth, body []byte) ([]byte, error) {
	headers, err := authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = []string{"application/json"}
	headers["Accept"] = []string{"application/json"}
	return c.httpClient.SendRequest(ctx, "DELETE", url, headers, body)
}
