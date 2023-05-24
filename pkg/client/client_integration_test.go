//go:build integration_test

package client

import (
	"context"
	"fmt"
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

	cwp := New(baseUrl, "admin", "test")

	err := cwp.ChangeAdminPassword(context.Background(), "admin", "test")
	if err != nil {
		log.Fatalf("Could not change admin password: %s", err)
	}
	team, err := cwp.GetTeam(context.Background(), "Administrators")
	if err != nil {
		log.Fatalf("Could not get team: %s", err)
	}
	_, err = cwp.GenerateApiKey(context.Background(), team.Uuid)
	if err != nil {
		log.Fatalf("Could not generate api key: %s", err)
	}

	c = New(baseUrl, "admin", "test", WithApiKeySource("Administrators"))
	code := m.Run()

	cleanup()

	os.Exit(code)
}

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	team, e := c.GetTeam(ctx, "Administrators")
	if e != nil {
		log.Fatalf("get team uuid: %v", e)
	}

	t.Run("Get version", func(t *testing.T) {
		version, err := c.Version(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, version)
	})

	t.Run("Create admin users", func(t *testing.T) {
		a := &AdminUsers{Users: []AdminUser{{"u1", "p1"}}}
		err := c.CreateAdminUsers(ctx, a, team.Uuid)
		assert.NoError(t, err)
		err = c.RemoveAdminUsers(ctx, a)
		assert.NoError(t, err)
	})

	t.Run("Create, get and delete team", func(t *testing.T) {
		team, err := c.CreateTeam(ctx, "test-team", []Permission{ViewPortfolioPermission})
		assert.NoError(t, err)
		assert.Equal(t, "test-team", team.Name)

		// dependencytrack allows creating teams with the same name
		team, err = c.CreateTeam(ctx, "test-team", []Permission{ViewPortfolioPermission})
		assert.NoError(t, err)
		assert.Equal(t, "test-team", team.Name)

		// dependencytrack comes with 3 teams by default, so we expect 5 teams after creating 2 new teams
		teams, err := c.GetTeams(ctx)
		assert.NoError(t, err)
		assert.Len(t, teams, 5)

		err = c.DeleteTeam(ctx, team.Uuid)
		assert.NoError(t, err)
	})

	t.Run("Create, get and delete project", func(t *testing.T) {
		project, err := c.CreateProject(context.Background(), "projectname", "version1", "group1", []string{"tag1", "tag2"})
		assert.NoError(t, err)
		assert.Equal(t, "projectname", project.Name)
		assert.Equal(t, "version1", project.Version)
		assert.Equal(t, "group1", project.Group)
		assert.Equal(t, "tag1", project.Tags[0].Name)
		assert.Equal(t, "tag2", project.Tags[1].Name)

		p, err := c.GetProject(ctx, "projectname", "version1")
		if err != nil {
			log.Fatalf("get project: %v", err)
		}
		assert.Equal(t, "projectname", p.Name)
		assert.Equal(t, "tag1", project.Tags[0].Name)
		assert.Equal(t, "tag2", project.Tags[1].Name)

		projects, err := c.GetProjectsByTag(ctx, "tag1")
		assert.NoError(t, err)
		assert.Equal(t, "projectname", projects[0].Name)

		projects, err = c.GetProjects(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, projects)

		err = c.DeleteProject(ctx, project.Uuid)
		assert.NoError(t, err)
	})

	t.Run("Create  and delete projects", func(t *testing.T) {
		project, err := c.CreateProject(context.Background(), "projectname", "version1", "group1", []string{"tag1", "tag2"})
		assert.NoError(t, err)

		err = c.DeleteProjects(ctx, project.Name)
		assert.NoError(t, err)
	})

	t.Run("CreateManagedUser", func(t *testing.T) {
		err := c.CreateManagedUser(ctx, "testuser", "testpassword")
		assert.NoError(t, err)
	})

	t.Run("AddToTeam and delete membership", func(t *testing.T) {
		err := c.AddToTeam(ctx, "testuser", team.Uuid)
		assert.NoError(t, err)
		err = c.DeleteUserMembership(ctx, team.Uuid, "testuser")
		assert.NoError(t, err)
	})

	t.Run("DeleteManagedUser", func(t *testing.T) {
		err := c.DeleteManagedUser(ctx, "testuser")
		assert.NoError(t, err)
	})

	t.Run("CreateOidcUser", func(t *testing.T) {
		err := c.CreateOidcUser(ctx, "email@nav.no")
		assert.NoError(t, err)
	})

	t.Run("GetOidcUsers", func(t *testing.T) {
		users, err := c.GetOidcUsers(ctx)
		fmt.Printf("*** users: %+v\n", users)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "email@nav.no", users[0].Username)
	})

	t.Run("DeleteOidcUser", func(t *testing.T) {
		err := c.DeleteOidcUser(ctx, "email@nav.no")
		assert.NoError(t, err)
	})

	t.Run("PortfolioRefresh", func(t *testing.T) {
		err := c.PortfolioRefresh(ctx)
		assert.NoError(t, err)
	})

	t.Run("UploadProjectBom", func(t *testing.T) {
		err := c.UploadProject(ctx, "projectname", "version1", []byte("test"))
		assert.NoError(t, err)
	})
}
