package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
	"github.com/sirupsen/logrus"
)

var (
	_ Auth = &usernamePasswordSource{}
)

type Auth interface {
	AuthContext(ctx context.Context) (context.Context, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type (
	Username = string
	Password = string
	Team     = string
)

type usernamePasswordSource struct {
	username    Username
	password    Password
	accessToken string
	lock        sync.Mutex
	client      *client.APIClient
	log         *logrus.Entry
}

func NewUsernamePasswordSource(username Username, password Password, c *client.APIClient, log *logrus.Entry) Auth {
	return &usernamePasswordSource{
		username: username,
		password: password,
		client:   c,
		log:      log,
	}
}

func newUsernamePasswordSourceWithToken(username Username, password Password, token string, c *client.APIClient, log *logrus.Entry) *usernamePasswordSource {
	return &usernamePasswordSource{
		username:    username,
		password:    password,
		accessToken: token,
		client:      c,
		log:         log,
	}
}

func (ups *usernamePasswordSource) AuthContext(ctx context.Context) (context.Context, error) {
	ups.lock.Lock()
	defer ups.lock.Unlock()

	t, expired, err := ups.checkAccessToken()
	if err != nil {
		return nil, err
	}

	if expired {
		t, err = ups.Login(ctx, ups.username, ups.password)
		if err != nil {
			return nil, err
		}
		ups.accessToken = t
	}

	return context.WithValue(ctx, client.ContextAccessToken, ups.accessToken), nil
}

func (ups *usernamePasswordSource) checkAccessToken() (string, bool, error) {
	if ups.accessToken == "" {
		return "", true, nil
	}

	token, err := jwt.ParseString(ups.accessToken, jwt.WithVerify(false))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired()) {
			ups.log.Debug("access token is expired, will re-login")
			return "", true, nil
		}
		return "", true, fmt.Errorf("checkAccessToken: failed to parse access token: %w", err)
	}
	expired := token.Expiration().Before(time.Now().Add(1 * time.Minute))
	return ups.accessToken, expired, nil
}

func (ups *usernamePasswordSource) Login(ctx context.Context, username Username, password Password) (string, error) {
	res, resp, err := ups.client.UserAPI.ValidateCredentials(ctx).
		Username(username).
		Password(password).
		Execute()
	if err != nil {
		return "", convertError(err, resp)
	}

	_, err = jwt.ParseString(res, jwt.WithVerify(false))
	if err != nil {
		return "", convertError(err, resp)
	}

	ups.log.Info("Successfully logged in with username and password")
	return res, nil
}

type ClientError struct {
	StatusCode          int
	ForcePasswordChange bool
	Body                string
	error
}

type ServerError struct {
	StatusCode int
	Body       string
	error
}

func convertError(err error, resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	switch {
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		if resp.StatusCode == http.StatusUnauthorized && body != nil && string(body) == "FORCE_PASSWORD_CHANGE" {
			return ClientError{
				StatusCode:          resp.StatusCode,
				ForcePasswordChange: true,
				Body:                string(body),
				error:               err,
			}
		}
		return ClientError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
			error:      err,
		}
	case resp.StatusCode >= 500:
		return ServerError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
			error:      err,
		}
	default:
		return nil
	}
}
