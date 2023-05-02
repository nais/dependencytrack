package auth

import (
	"context"
	"encoding/json"
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

var _ Auth = &apikeySource{}

type apikeySource struct {
	baseUrl    string
	source     Auth
	team       string
	log        *log.Entry
	httpClient *httpclient.HttpClient
	apiKey     string
}

type team struct {
	Uuid    string   `json:"uuid"`
	Name    string   `json:"name"`
	Apikeys []apiKey `json:"apikeys"`
}

type apiKey struct {
	Key string `json:"key"`
}

func NewUsernamePasswordSource(loginUrl LoginUrl, username Username, password Password, c *httpclient.HttpClient, log *log.Entry) Auth {
	if c == nil {
		c = httpclient.New(httpclient.WithLogger(log))
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

func NewApiKeySource(baseUrl string, team string, source Auth, c *httpclient.HttpClient, log *log.Entry) Auth {
	if c == nil {
		c = httpclient.New(httpclient.WithLogger(log))
	}
	return &apikeySource{
		baseUrl:    baseUrl,
		source:     source,
		team:       team,
		log:        log,
		httpClient: c,
	}
}

func (c *apikeySource) Headers(ctx context.Context) (http.Header, error) {
	key, err := c.refreshApiKey(ctx)
	if err != nil {
		return nil, err
	}
	return map[string][]string{"X-Api-Key": {key}}, nil
}

func (c *apikeySource) refreshApiKey(ctx context.Context) (string, error) {
	if c.apiKey == "" {
		c.log.Debug("Fetching apiKey")
		key, err := c.getApiKey(ctx)
		if err != nil {
			return "", err
		}
		c.apiKey = key
	}
	return c.apiKey, nil
}

func (c *apikeySource) getApiKey(ctx context.Context) (string, error) {
	headers, err := c.source.Headers(ctx)
	if err != nil {
		return "", err
	}

	headers.Set("Accept", "application/json")
	responseTeams, err := c.httpClient.SendRequest(ctx, http.MethodGet, c.baseUrl+"/team", headers, nil)

	if err != nil {
		return "", err
	}

	var teams []team
	err = json.Unmarshal(responseTeams, &teams)
	if err != nil {
		return "", err
	}

	for _, t := range teams {
		if t.Name == c.team {
			if len(t.Apikeys) == 0 {
				return "", fmt.Errorf("no apikeys found for team %s", c.team)
			}
			return t.Apikeys[0].Key, nil
		}
	}

	return "", fmt.Errorf("no team found with name %s", c.team)
}
