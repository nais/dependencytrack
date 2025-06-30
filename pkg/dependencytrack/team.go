package dependencytrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
	log "github.com/sirupsen/logrus"
)

func (c *dependencyTrackClient) GetTeam(ctx context.Context, team string) (*client.Team, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*client.Team, error) {
		teams, resp, err := c.client.TeamAPI.GetTeams(tokenCtx).Execute()
		if err != nil {
			return nil, convertError(err, "GetTeam", resp)
		}

		for _, t := range teams {
			if *t.Name == team {
				return &t, nil
			}
		}
		return nil, fmt.Errorf("team %s not found", team)
	})
}

func (c *dependencyTrackClient) GetTeams(ctx context.Context) ([]client.Team, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]client.Team, error) {
		res, resp, err := c.client.TeamAPI.GetTeams(tokenCtx).Execute()
		if err != nil {
			return nil, convertError(err, "GetTeams", resp)
		}
		return res, nil
	})
}

func (c *dependencyTrackClient) CreateTeam(ctx context.Context, teamName string, permissions []client.Permission) (*client.Team, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*client.Team, error) {
		team, resp, err := c.client.TeamAPI.CreateTeam(tokenCtx).Team(client.Team{
			Name: &teamName,
		}).Execute()
		if err != nil {
			return nil, convertError(err, "CreateTeam", resp)
		}

		for _, p := range permissions {
			if _, resp, err := c.client.PermissionAPI.AddPermissionToTeam(tokenCtx, team.Uuid, p.GetName()).Execute(); err != nil {
				return nil, convertError(err, "AddPermissionToTeam", resp)
			}
		}

		return team, nil
	})
}

func (c *dependencyTrackClient) DeleteTeam(ctx context.Context, uuid string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		resp, err := c.client.TeamAPI.DeleteTeam(tokenCtx).Team(client.Team{
			Uuid: uuid,
		}).Execute()
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				log.Infof("team %s does not exist", uuid)
				return nil // Team does not exist, return nil
			}
			return convertError(err, "DeleteTeam", resp)
		}
		return nil
	})
}

func (c *dependencyTrackClient) GenerateApiKey(ctx context.Context, uuid string) (string, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (string, error) {
		res, resp, err := c.client.TeamAPI.GenerateApiKey(tokenCtx, uuid).Execute()
		if err != nil {
			return "", convertError(err, "GenerateApiKey", resp)
		}
		return *res.Key, nil
	})
}
