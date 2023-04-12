package client

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	log "github.com/sirupsen/logrus"
)

func (c *Client) Token(ctx context.Context) (string, error) {
	if c.accessToken == "" || c.isExpired() {
		log.Debugf("accessToken expired, getting new one")
		token, err := c.login(ctx)
		if err != nil {
			return "", err
		}
		c.accessToken = token
	}
	return c.accessToken, nil
}

func (c *Client) login(ctx context.Context) (string, error) {
	token, err := sendRequest(ctx, "POST", c.baseUrl+"/api/v1/user/login", map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"text/plain"},
	}, []byte("username="+c.username+"&password="+c.password))

	if err != nil {
		return "", err
	}
	return string(token), nil
}

func (c *Client) isExpired() bool {
	if c.accessToken == "" {
		return true
	}
	parseOpts := []jwt.ParseOption{
		jwt.WithVerify(false),
	}
	token, err := jwt.ParseString(c.accessToken, parseOpts...)
	if err != nil {
		log.Errorf("error parsing accessToken: %v", err)
		return true
	}
	if token.Expiration().Before(time.Now().Add(-1 * time.Minute)) {
		return true
	}
	return false
}
