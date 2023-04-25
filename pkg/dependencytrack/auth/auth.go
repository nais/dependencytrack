package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/nais/dependencytrack/pkg/httpclient"
	log "github.com/sirupsen/logrus"
)

type Auth interface {
	Headers(ctx context.Context) (http.Header, error)
}

type LoginUrl = string
type Username = string
type Password = string

type usernamePasswordSource struct {
	loginUrl    string
	username    string
	password    string
	accessToken string
	httpClient  *httpclient.HttpClient
	log         *log.Entry
}

func NewUsernamePasswordSource(loginUrl LoginUrl, username Username, password Password, c *httpclient.HttpClient, log *log.Entry) Auth {
	if c == nil {
		c = httpclient.New(http.DefaultClient, log)
	}
	return &usernamePasswordSource{
		loginUrl:   loginUrl,
		username:   username,
		password:   password,
		httpClient: c,
		log:        log,
	}
}

func (c *usernamePasswordSource) Headers(ctx context.Context) (http.Header, error) {
	expired, err := c.isExpired()
	if err != nil {
		return nil, err
	}
	if c.accessToken == "" || expired {
		c.log.Debugf("accessToken expired, getting new one")
		t, err := c.login(ctx)
		if err != nil {
			return nil, err
		}
		c.accessToken = t
	}
	return map[string][]string{"Authorization": {"Bearer " + c.accessToken}}, nil
}

func (c *usernamePasswordSource) isExpired() (bool, error) {
	if c.accessToken == "" {
		return true, nil
	}
	parseOpts := []jwt.ParseOption{
		jwt.WithVerify(false),
	}
	token, err := jwt.ParseString(c.accessToken, parseOpts...)
	if err != nil {
		return true, fmt.Errorf("parsing accessToken: %w", err)
	}
	if token.Expiration().Before(time.Now().Add(-1 * time.Minute)) {
		return true, nil
	}
	return false, nil
}

func (c *usernamePasswordSource) login(ctx context.Context) (string, error) {
	data := url.Values{
		"username": {c.username},
		"password": {c.password},
	}

	token, err := c.httpClient.SendRequest(ctx, "POST", c.loginUrl, map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"text/plain"},
	}, []byte(data.Encode()))

	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	_, err = c.parseToken(string(token))
	if err != nil {
		return "", fmt.Errorf("could not parse token from body after login request: %w, response body: %s", err, token)
	}

	return string(token), nil
}

func (c *usernamePasswordSource) parseToken(token string) (jwt.Token, error) {
	parseOpts := []jwt.ParseOption{
		jwt.WithVerify(false),
	}
	return jwt.ParseString(token, parseOpts...)
}
