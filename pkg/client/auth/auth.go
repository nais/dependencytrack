package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/nais/dependencytrack/pkg/httpclient"
	"github.com/sirupsen/logrus"
)

type Auth interface {
	Headers(ctx context.Context) (http.Header, error)
}

type (
	BaseUrl  = string
	Username = string
	Password = string
)

type usernamePasswordSource struct {
	baseUrl     string
	username    string
	password    string
	accessToken string
	httpClient  *httpclient.HttpClient
	log         *logrus.Entry
}

type staticApiKeySource struct {
	headerName string
	apiKey     string
}

type apiKeySource struct {
	baseUrl    string
	source     Auth
	team       string
	log        *logrus.Entry
	httpClient *httpclient.HttpClient
	apiKey     string
}

type team struct {
	Uuid    string   `json:"uuid"`
	Name    string   `json:"name"`
	ApiKeys []apiKey `json:"apiKeys"`
}

type apiKey struct {
	Key string `json:"key"`
}

var (
	_ Auth = &usernamePasswordSource{}
	_ Auth = &staticApiKeySource{}
	_ Auth = &apiKeySource{}
)

func NewUsernamePasswordSource(baseUrl BaseUrl, username Username, password Password, c *httpclient.HttpClient, log *logrus.Entry) Auth {
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	if log == nil {
		log = logrus.NewEntry(logrus.StandardLogger())
	}
	if c == nil {
		c = httpclient.New(httpclient.WithLogger(log))
	}
	return &usernamePasswordSource{
		baseUrl:    baseUrl,
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

	res, err := c.httpClient.SendRequest(ctx, "POST", c.baseUrl+"/api/v1/user/login", map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"text/plain"},
	}, []byte(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	_, err = c.parseToken(string(res.Body))
	if err != nil {
		return "", fmt.Errorf("could not parse token from body after login request: %w, response body: %s", err, res.Body)
	}

	return string(res.Body), nil
}

func (c *usernamePasswordSource) parseToken(token string) (jwt.Token, error) {
	parseOpts := []jwt.ParseOption{
		jwt.WithVerify(false),
	}
	return jwt.ParseString(token, parseOpts...)
}

func NewStaticApiKeySource(headerName string, apiKey string) Auth {
	return &staticApiKeySource{
		headerName: headerName,
		apiKey:     apiKey,
	}
}

func (s *staticApiKeySource) Headers(_ context.Context) (http.Header, error) {
	return map[string][]string{s.headerName: {s.apiKey}}, nil
}

func NewApiKeySource(baseUrl string, team string, source Auth, c *httpclient.HttpClient, log *logrus.Entry) Auth {
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	if log == nil {
		log = logrus.NewEntry(logrus.StandardLogger())
	}
	if c == nil {
		c = httpclient.New(httpclient.WithLogger(log))
	}
	return &apiKeySource{
		baseUrl:    baseUrl,
		source:     source,
		team:       team,
		log:        log,
		httpClient: c,
	}
}

func (c *apiKeySource) Headers(ctx context.Context) (http.Header, error) {
	key, err := c.refreshApiKey(ctx)
	if err != nil {
		return nil, err
	}
	return map[string][]string{"X-Api-Key": {key}}, nil
}

func (c *apiKeySource) refreshApiKey(ctx context.Context) (string, error) {
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

func (c *apiKeySource) getApiKey(ctx context.Context) (string, error) {
	headers, err := c.source.Headers(ctx)
	if err != nil {
		return "", err
	}

	headers.Set("Accept", "application/json")
	resp, err := c.httpClient.SendRequest(ctx, http.MethodGet, c.baseUrl+"/api/v1/team", headers, nil)
	if err != nil {
		return "", err
	}

	var teams []team
	err = json.Unmarshal(resp.Body, &teams)
	if err != nil {
		return "", err
	}

	for _, t := range teams {
		if t.Name == c.team {
			if len(t.ApiKeys) == 0 {
				return "", fmt.Errorf("no apikeys found for team %s", c.team)
			}
			return t.ApiKeys[0].Key, nil
		}
	}

	return "", fmt.Errorf("no team found with name %s", c.team)
}
