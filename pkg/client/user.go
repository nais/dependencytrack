package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nais/dependencytrack/pkg/httpclient"
	log "github.com/sirupsen/logrus"
)

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
		} else {
			c.log.Infof("removed user %s", user.Username)
		}
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
