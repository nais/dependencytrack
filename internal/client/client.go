package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	baseUrl     string
	username    string
	password    string
	accessToken string
}

type User struct {
	Username            string `json:"username"`
	NewPassword         string `json:"newPassword"`
	ConfirmPassword     string `json:"confirmPassword"`
	Fullname            string `json:"fullname"`
	Email               string `json:"email"`
	Suspended           bool   `json:"suspended"`
	ForcePasswordChange bool   `json:"forcePasswordChange"`
	NonExpiryPassword   bool   `json:"nonExpiryPassword"`
}

type Team struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

func New(baseUrl string, username string, password string) *Client {
	return &Client{
		baseUrl:  baseUrl,
		username: username,
		password: password,
	}
}

func (c *Client) ChangeAdminPassword(ctx context.Context, oldPassword, newPassword string) error {
	_, err := sendRequest(ctx, http.MethodPost, c.baseUrl+"/api/v1/user/forceChangePassword", map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"text/plain"},
	}, []byte("username=admin&password="+oldPassword+"&newPassword="+newPassword+"&confirmPassword="+newPassword))

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) NewManagedUser(ctx context.Context, username, password string) error {
	user := &User{
		Username:            username,
		NewPassword:         password,
		ConfirmPassword:     password,
		Fullname:            username,
		Email:               username + "@nais.io",
		Suspended:           false,
		ForcePasswordChange: false,
		NonExpiryPassword:   true,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshalling user: %w", err)
	}

	token, err := c.Token(ctx)
	if err != nil {
		return fmt.Errorf("getting Token: %w", err)
	}

	_, err = sendRequest(ctx, http.MethodPut, c.baseUrl+"/api/v1/user/managed", map[string][]string{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Authorization": {"Bearer " + token},
	}, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) AddToTeam(ctx context.Context, username string, team string) error {
	token, err := c.Token(ctx)
	if err != nil {
		return fmt.Errorf("getting Token: %w", err)
	}

	resp, err := sendRequest(ctx, http.MethodGet, c.baseUrl+"/api/v1/team", map[string][]string{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Authorization": {"Bearer " + token},
	}, nil)
	if err != nil {
		return fmt.Errorf("getting teams: %w", err)
	}

	var teams []Team
	if err = json.Unmarshal(resp, &teams); err != nil {
		return fmt.Errorf("unmarshalling teams: %w", err)
	}

	var uuid string
	for _, t := range teams {
		if t.Name == team {
			uuid = t.Uuid
		}
	}
	if uuid == "" {
		return fmt.Errorf("team %s not found", team)
	}

	_, err = sendRequest(ctx, http.MethodPost, c.baseUrl+"/api/v1/user/"+username+"/membership", map[string][]string{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Authorization": {"Bearer " + token},
	}, []byte(`{"uuid": "`+uuid+`"}`))

	if err != nil {
		return fmt.Errorf("adding user to team: %w", err)
	}
	return nil
}

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func fail(status int, err error) *RequestError {
	return &RequestError{
		StatusCode: status,
		Err:        err,
	}
}

func sendRequest(ctx context.Context, httpMethod string, url string, headers map[string][]string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	log.Debugf("Response: %s", resp.Status)

	if resp.StatusCode > 299 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body: %w", err)
		}
		return nil, fail(resp.StatusCode, fmt.Errorf("%s\n", string(b)))
	}
	resBody, err := io.ReadAll(resp.Body)
	return resBody, err
}
