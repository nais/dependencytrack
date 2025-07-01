package dependencytrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
	log "github.com/sirupsen/logrus"
)

type Permission string

const (
	AccessManagementPermission        = Permission("ACCESS_MANAGEMENT")
	PolicyManagementPermission        = Permission("POLICY_MANAGEMENT")
	PolicyViolationAnalysisPermission = Permission("POLICY_VIOLATION_ANALYSIS")
	SystemConfigurationPermission     = Permission("SYSTEM_CONFIGURATION")
	ViewPolicyViolationPermission     = Permission("VIEW_POLICY_VIOLATION")
	ViewPortfolioPermission           = Permission("VIEW_PORTFOLIO")
	ViewVulnerabilityPermission       = Permission("VIEW_VULNERABILITY")
)

type Team struct {
	Uuid      string   `json:"uuid,omitempty"`
	Name      string   `json:"name,omitempty"`
	OidcUsers []User   `json:"oidcUsers,omitempty"`
	ApiKeys   []ApiKey `json:"apiKeys,omitempty"`
}

type ApiKey struct {
	Key string `json:"key,omitempty"`
}

func (c *dependencyTrackClient) GetTeam(ctx context.Context, team string) (*Team, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*Team, error) {
		teams, resp, err := c.client.TeamAPI.GetTeams(tokenCtx).Execute()
		if err != nil {
			return nil, convertError(err, "GetTeam", resp)
		}

		for _, t := range teams {
			if *t.Name == team {
				teamDetails, err := parseTeam(&t)
				if err != nil {
					return nil, fmt.Errorf("failed to parse team %s: %w", team, err)
				}
				return teamDetails, nil
			}
		}
		return nil, fmt.Errorf("team %s not found", team)
	})
}

func (c *dependencyTrackClient) GetTeams(ctx context.Context) ([]*Team, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]*Team, error) {
		res, resp, err := c.client.TeamAPI.GetTeams(tokenCtx).Execute()
		if err != nil {
			return nil, convertError(err, "GetTeams", resp)
		}

		teams := make([]*Team, 0, len(res))
		for _, team := range res {
			teamDetails, err := parseTeam(&team)
			if err != nil {
				return nil, fmt.Errorf("failed to parse team %s: %w", *team.Name, err)
			}
			teams = append(teams, teamDetails)
		}
		return teams, nil
	})
}

func (c *dependencyTrackClient) CreateTeam(ctx context.Context, teamName string, permissions []Permission) (*Team, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*Team, error) {
		team, resp, err := c.client.TeamAPI.CreateTeam(tokenCtx).Team(client.Team{
			Name: &teamName,
		}).Execute()
		if err != nil {
			return nil, convertError(err, "CreateTeam", resp)
		}

		for _, p := range permissions {
			if _, resp, err = c.client.PermissionAPI.AddPermissionToTeam(tokenCtx, team.Uuid, string(p)).Execute(); err != nil {
				return nil, convertError(err, "AddPermissionToTeam", resp)
			}
		}
		teamDetails, err := parseTeam(team)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created team %s: %w", teamName, err)
		}
		log.Infof("created team %s with UUID %s", teamDetails.Name, teamDetails.Uuid)
		return teamDetails, nil
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

func parseTeam(team *client.Team) (*Team, error) {
	if team == nil {
		return nil, fmt.Errorf("team is nil")
	}

	parsed := &Team{
		Uuid:      team.Uuid,
		Name:      SafeString(team.Name),
		OidcUsers: make([]User, 0, len(team.OidcUsers)),
		ApiKeys:   make([]ApiKey, 0, len(team.ApiKeys)),
	}

	for _, user := range team.OidcUsers {
		parsed.OidcUsers = append(parsed.OidcUsers, User{
			Username: SafeString(user.Username),
			Email:    SafeString(user.Email),
		})
	}

	for _, key := range team.ApiKeys {
		parsed.ApiKeys = append(parsed.ApiKeys, ApiKey{
			Key: SafeString(key.Key),
		})
	}

	return parsed, nil
}
