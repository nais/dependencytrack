//go:build integration_test

package client

import (
	"context"
	"os"
	"testing"

	"github.com/nais/dependencytrack/pkg/client/test"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var c Client

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	baseUrl, cleanup := test.DependencyTrackPool()

	c = New(baseUrl, "admin", "test")
	err := c.ChangeAdminPassword(context.Background(), "admin", "test")
	if err != nil {
		log.Fatalf("Could not change admin password: %s", err)
	}

	code := m.Run()

	cleanup()

	os.Exit(code)
}

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	team, err := c.GetTeam(ctx, "Administrators")
	if err != nil {
		log.Fatalf("get team uuid: %v", err)
	}

	t.Run("Create admin users", func(t *testing.T) {
		a := &AdminUsers{Users: []AdminUser{{"u1", "p1"}}}
		err := c.CreateAdminUsers(ctx, a, team.Uuid)
		assert.NoError(t, err)
		// TODO: should we implement GetAdminUser and check that the users are created?
	})

	t.Run("Create team", func(t *testing.T) {
		team, err := c.CreateTeam(ctx, "test-team", []Permission{ViewPortfolioPermission})
		assert.NoError(t, err)
		assert.Equal(t, "test-team", team.Name)
	})

	// dependencytrack allows creating teams with the same name
	t.Run("Create team with existing name", func(t *testing.T) {
		team, err := c.CreateTeam(ctx, "test-team", []Permission{ViewPortfolioPermission})
		assert.NoError(t, err)
		assert.Equal(t, "test-team", team.Name)
	})

	// dependencytrack comes with 3 teams by default, so we expect 5 teams after creating 2 new teams
	t.Run("Get teams with new team", func(t *testing.T) {
		teams, err := c.GetTeams(ctx)
		assert.NoError(t, err)
		assert.Len(t, teams, 5)
	})
}
