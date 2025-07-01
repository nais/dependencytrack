package dependencytrack

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nais/dependencytrack/pkg/dependencytrack/auth"
	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
	log "github.com/sirupsen/logrus"
)

type AdminUser struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (c *dependencyTrackClient) AddToTeam(ctx context.Context, username, uuid string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		_, resp, err := c.client.UserAPI.AddTeamToUser(tokenCtx, username).
			IdentifiableObject(client.IdentifiableObject{Uuid: &uuid}).
			Execute()
		if err != nil {
			switch resp.StatusCode {
			case http.StatusNotFound:
				log.Infof("team %s does not exist", uuid)
				return nil
			case http.StatusNotModified:
				log.Infof("username %s is a part of the team", username)
				return nil
			default:
				return convertError(err, "AddToTeam", resp)
			}
		}
		return nil
	})
}

func (c *dependencyTrackClient) ChangeAdminPassword(ctx context.Context, oldPassword, newPassword string) error {
	_, err := c.Login(ctx, "admin", oldPassword)
	var authErr auth.ClientError
	if errors.As(err, &authErr) {
		if authErr.ForcePasswordChange {
			resp, err := c.client.UserAPI.ForceChangePassword(ctx).
				Username("admin").
				Password(oldPassword).
				NewPassword(newPassword).
				ConfirmPassword(newPassword).
				Execute()

			if err != nil {
				return convertError(err, "ChangeAdminPassword", resp)
			}
		}
	}

	return nil
}

func (c *dependencyTrackClient) CreateAdminUser(ctx context.Context, username, password string, teamUuid string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		user := client.ManagedUser{
			Username:            stringPtr(username),
			NewPassword:         stringPtr(password),
			ConfirmPassword:     stringPtr(password),
			Fullname:            stringPtr(username),
			Email:               stringPtr(username + EmailPostfix),
			Suspended:           boolPtr(false),
			ForcePasswordChange: boolPtr(false),
			NonExpiryPassword:   boolPtr(true),
		}

		mu, resp, err := c.client.UserAPI.CreateManagedUser(tokenCtx).ManagedUser(user).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode != http.StatusConflict {
				return convertError(err, "CreateAdminUser", resp)
			} else {
				log.Infof("admin user %s already exists", username)
				return nil // User already exists, nothing to do
			}
		}

		_, resp, err = c.client.UserAPI.AddTeamToUser(tokenCtx, *mu.Username).IdentifiableObject(
			client.IdentifiableObject{
				Uuid: stringPtr(teamUuid),
			}).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode != http.StatusNotModified {
				return convertError(err, "AddTeamToUser", resp)
			} else {
				log.Infof("admin user %s already in team, nothing to do", username)
				return nil // User already in team, nothing to do
			}
		}
		c.log.Infof("created admin user %s", username)
		return nil
	})
}

func (c *dependencyTrackClient) CreateAdminUsers(ctx context.Context, users []*AdminUser, teamUuid string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		for _, user := range users {
			if err := c.CreateAdminUser(tokenCtx, user.Username, user.Password, teamUuid); err != nil {
				return fmt.Errorf("failed to create admin user %s: %w", user.Username, err)
			}
			log.Infof("created admin user %s", user.Username)
		}
		return nil
	})
}

func (c *dependencyTrackClient) RemoveAdminUser(ctx context.Context, username string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		resp, err := c.client.UserAPI.DeleteManagedUser(tokenCtx).ManagedUser(client.ManagedUser{
			Username: stringPtr(username),
		}).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode != http.StatusNotFound {
				return convertError(err, "RemoveAdminUser", resp)
			} else {
				log.Infof("admin user %s does not exist, nothing to remove", username)
				return nil // User does not exist, nothing to do
			}
		}
		return nil
	})
}

func (c *dependencyTrackClient) RemoveAdminUsers(ctx context.Context, users []*AdminUser) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		for _, user := range users {
			err := c.RemoveAdminUser(tokenCtx, user.Username)
			if err != nil {
				return fmt.Errorf("failed to remove admin user %s: %w", user.Username, err)
			} else {
				log.Infof("removed admin user %s", user.Username)
			}
		}
		return nil
	})
}

func (c *dependencyTrackClient) GetOidcUsers(ctx context.Context) ([]*User, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]*User, error) {
		users, resp, err := c.client.UserAPI.GetOidcUsers(tokenCtx).Execute()
		if err != nil {
			return nil, convertError(err, "GetOidcUsers", resp)
		}

		uu := make([]*User, len(users))
		for i, u := range users {
			uu[i] = parseUser(u)
		}
		return uu, nil
	})
}

func (c *dependencyTrackClient) GetOidcUser(ctx context.Context, email string) (*User, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*User, error) {
		users, resp, err := c.client.UserAPI.GetOidcUsers(tokenCtx).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return nil, fmt.Errorf("oidc user %s not found", email)
			}
			return nil, convertError(err, "GetOidcUser", resp)
		}

		var user *User
		for _, u := range users {
			if u.Username != nil && *u.Username == email {
				user = parseUser(u)
				break
			}
		}
		return user, nil
	})
}

func (c *dependencyTrackClient) CreateOidcUser(ctx context.Context, email string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		user := client.OidcUser{
			Username: stringPtr(email),
			Email:    stringPtr(email),
		}

		_, resp, err := c.client.UserAPI.CreateOidcUser(tokenCtx).OidcUser(user).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode != http.StatusConflict {
				return convertError(err, "CreateOidcUser", resp)
			} else {
				log.Infof("oidc user %s already exists", email)
				return nil // User already exists, nothing to do
			}
		}
		log.Infof("created oidc user %s", email)
		return nil
	})
}

func (c *dependencyTrackClient) DeleteManagedUser(ctx context.Context, username string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		resp, err := c.client.UserAPI.DeleteManagedUser(tokenCtx).ManagedUser(client.ManagedUser{
			Username: stringPtr(username),
		}).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode != http.StatusNotFound {
				return convertError(err, "DeleteManagedUser", resp)
			} else {
				log.Infof("managed user %s does not exist, nothing to remove", username)
				return nil // User does not exist, nothing to do
			}
		}
		return nil
	})
}

func (c *dependencyTrackClient) DeleteOidcUser(ctx context.Context, email string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		resp, err := c.client.UserAPI.DeleteOidcUser(tokenCtx).OidcUser(client.OidcUser{
			Username: stringPtr(email),
		}).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode != http.StatusNotFound {
				return convertError(err, "DeleteOidcUser", resp)
			} else {
				log.Infof("oidc user %s does not exist, nothing to remove", email)
				return nil // User does not exist, nothing to do
			}
		}
		return nil
	})
}

func (c *dependencyTrackClient) DeleteUserMembership(ctx context.Context, uuid, username string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		_, resp, err := c.client.UserAPI.RemoveTeamFromUser(tokenCtx, username).
			IdentifiableObject(client.IdentifiableObject{
				Uuid: stringPtr(uuid),
			}).
			Execute()
		if err != nil {
			switch resp.StatusCode {
			case http.StatusNotFound:
				log.Infof("user %s does not exist", username)
				return nil
			case http.StatusNotModified:
				log.Infof("user %s is not a part of the team", username)
				return nil
			default:
				return convertError(err, "DeleteUserMembership", resp)
			}
		}
		return nil
	})
}

func parseUser(users client.OidcUser) *User {
	if users.Username == nil {
		return &User{}
	}
	return &User{
		Username: SafeString(users.Username),
		Email:    SafeString(users.Email),
	}
}
