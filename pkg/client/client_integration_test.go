//go:build integration_test

package client

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

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
	baseUrl, cleanup := test.DependencyTrackPool("4.9.0")

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

	t.Run("Create, get, update and delete project", func(t *testing.T) {
		project, err := c.CreateProject(context.Background(), "projectname", "version1", "group1", []string{"tag1", "tag2", "team:app:container"})
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

		p, err = c.UpdateProject(ctx, p.Uuid, p.Name, "newVersion", p.Group, []string{"tag1", "tag2", "team:app:container"})
		assert.NoError(t, err)
		assert.Equal(t, "newVersion", p.Version)

		projects, err := c.GetProjectsByTag(ctx, "team:app:container")
		assert.NoError(t, err)
		assert.Equal(t, "projectname", projects[0].Name)

		projects, err = c.GetProjects(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, projects)

		err = c.DeleteProject(ctx, project.Uuid)
		assert.NoError(t, err)
	})

	t.Run("delete project that does not exist", func(t *testing.T) {
		err := c.DeleteProjects(ctx, "non-existing-project")
		assert.NoError(t, err)
	})

	t.Run("Create  and delete projects", func(t *testing.T) {
		project, err := c.CreateProject(context.Background(), "projectname", "version1", "group1", []string{"tag1", "tag2"})
		assert.NoError(t, err)

		err = c.DeleteProjects(ctx, project.Name)
		assert.NoError(t, err)
	})

	t.Run("Get projects", func(t *testing.T) {
		projects1, err := c.GetProjects(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 0, len(projects1))
		total := 201

		for i := 0; i < total; i++ {
			_, err := c.CreateProject(context.Background(), "projectname"+strconv.Itoa(i), "version1", "group1", []string{"tag1", "tag2"})
			assert.NoError(t, err)
		}

		// tests pagination as max limit per page is 100
		projects2, err := c.GetProjects(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, total, len(projects2))
		names := make([]string, len(projects2))
		for _, p := range projects2 {
			names = append(names, p.Name)
		}
		for i := 0; i < total; i++ {
			assert.Contains(t, names, "projectname"+strconv.Itoa(i))
		}
	})

	t.Run("CreateManagedUser", func(t *testing.T) {
		err := c.CreateManagedUser(ctx, "testuser", "testpassword")
		assert.NoError(t, err)
	})

	t.Run("AddToTeam, add existing member and delete membership", func(t *testing.T) {
		err := c.AddToTeam(ctx, "testuser", team.Uuid)
		assert.NoError(t, err)
		err = c.AddToTeam(ctx, "testuser", team.Uuid)
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

	t.Run("GetConfigProperties", func(t *testing.T) {
		props, err := c.GetConfigProperties(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, props)
	})

	t.Run("ConfigPropertyAggregate", func(t *testing.T) {
		props, err := c.ConfigPropertyAggregate(ctx, []ConfigProperty{
			{
				GroupName:     "vuln-source",
				PropertyName:  "google.osv.enabled",
				PropertyValue: "false",
				PropertyType:  "STRING",
				Description:   "List of enabled ecosystems to mirror OSV",
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, props)
	})

	t.Run("GetEcosystems", func(t *testing.T) {
		ecos, err := c.GetEcosystems(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, ecos)
	})

	t.Run("GetFindings", func(t *testing.T) {
		b, err := os.ReadFile("testdata/attestation.json")
		assert.NoError(t, err)
		p, err := c.CreateProject(context.Background(), "project-with-findings", "version1", "group1", []string{"tag1", "tag2", "team:app:container"})
		assert.NoError(t, err)
		err = c.UploadProject(ctx, "project-with-findings", "version1", b)
		assert.NoError(t, err)

		err = c.TriggerAnalysis(ctx, p.Uuid)
		assert.NoError(t, err)
		time.Sleep(2 * time.Second)

		_, err = c.GetFindings(ctx, p.Uuid)
		assert.NoError(t, err)
	})
}
