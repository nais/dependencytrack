package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nais/dependencytrack/pkg/httpclient"
	log "github.com/sirupsen/logrus"
)

func (c *client) GenerateApiKey(ctx context.Context, uuid string) (string, error) {
	res, err := c.put(ctx, fmt.Sprintf("%s/api/v1/team/%s/key", c.baseUrl, uuid), c.authSource, nil)
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
