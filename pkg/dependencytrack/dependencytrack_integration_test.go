//go:build integration_test
// +build integration_test

package dependencytrack_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	uuid "github.com/google/uuid"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/nais/dependencytrack/pkg/dependencytrack"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var c dependencytrack.Client

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	baseUrl, cleanup := DependencyTrackPool("4.13.2")
	cwp, err := dependencytrack.NewClient(baseUrl, "admin", "test", log.WithField("test", "client_integration_test"))
	if err != nil {
		log.Fatalf("Failed to create DependencyTrack client: %v", err)
	}

	err = cwp.ChangeAdminPassword(context.Background(), "admin", "test")
	if err != nil {
		log.Fatalf("Could not change admin password: %s", err)
	}

	c = cwp
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
	assert.NotEmpty(t, team.Uuid, "Team UUID should not be empty")

	t.Run("Get version", func(t *testing.T) {
		version, err := c.Version(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, version)
	})

	t.Run("Create admin users", func(t *testing.T) {
		type AdminUsers struct {
			Users []*dependencytrack.AdminUser `json:"users"`
		}
		a := &AdminUsers{Users: []*dependencytrack.AdminUser{{"u1", "p1"}}}
		err := c.CreateAdminUsers(ctx, a.Users, team.Uuid)
		assert.NoError(t, err)
		err = c.RemoveAdminUsers(ctx, a.Users)
		assert.NoError(t, err)
	})

	t.Run("Create, get and delete team", func(t *testing.T) {
		team, err := c.CreateTeam(ctx, "test-team", []dependencytrack.Permission{dependencytrack.ViewPortfolioPermission})
		assert.NoError(t, err)
		assert.Equal(t, "test-team", team.Name)

		// dependencytrack allows creating teams with the same name
		team, err = c.CreateTeam(ctx, "test-team", []dependencytrack.Permission{dependencytrack.ViewPortfolioPermission})
		assert.NoError(t, err)
		assert.Equal(t, "test-team", team.Name)

		// dependencytrack comes with 4 teams by default, so we expect 5 teams after creating 2 new teams
		teams, err := c.GetTeams(ctx)
		assert.NoError(t, err)
		for _, t := range teams {
			log.Debugf("Team: %s, UUID: %s", t.Name, t.Uuid)
		}
		// created 2 new teams, so we expect 6 teams in total
		assert.Len(t, teams, 6)

		err = c.DeleteTeam(ctx, team.Uuid)
		assert.NoError(t, err)
		teams, err = c.GetTeams(ctx)
		assert.NoError(t, err)
		assert.Len(t, teams, 5) // should be back to 5 teams after deletion
	})

	t.Run("Create, get and delete project", func(t *testing.T) {
		project, err := c.CreateProject(ctx, "projectname", "version1", []string{"tag1", "tag2", "team:app:container"})
		assert.NoError(t, err)
		assert.Equal(t, "projectname", project.Name)
		assert.Equal(t, "version1", project.Version)
		assert.Equal(t, "tag1", project.Tags[0].Name)
		assert.Equal(t, "tag2", project.Tags[1].Name)

		p, err := c.GetProject(ctx, "projectname", "version1")
		assert.NoError(t, err)
		assert.Equal(t, "projectname", p.Name)
		assert.Equal(t, "tag1", project.Tags[0].Name)
		assert.Equal(t, "tag2", project.Tags[1].Name)

		projects, err := c.GetProjects(ctx, 10, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, projects)
		assert.Len(t, projects, 1)

		err = c.DeleteProject(ctx, project.Uuid)
		assert.NoError(t, err)
		projects, err = c.GetProjects(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Empty(t, projects) // should be empty after deletion
	})

	t.Run("delete project that does not exist", func(t *testing.T) {
		err := c.DeleteProject(ctx, "non-existing-project")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid UUID")

		uuid, err := uuid.NewUUID()
		assert.NoError(t, err)
		err = c.DeleteProject(ctx, uuid.String())
		assert.Error(t, err)
	})

	t.Run("Get projects", func(t *testing.T) {
		projects1, err := c.GetProjects(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(projects1))
		total := 201

		for i := 0; i < total; i++ {
			_, err := c.CreateProject(ctx, "projectname"+strconv.Itoa(i), "version1", []string{"tag1", "tag2"})
			assert.NoError(t, err)
		}

		// tests pagination as max limit per page is 100
		projects2, err := c.GetProjects(ctx, 10, 0)
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

	t.Run("AddToTeam, add existing member and delete membership", func(t *testing.T) {
		err := c.AddToTeam(ctx, "testuser", team.Uuid)
		assert.NoError(t, err)
		err = c.AddToTeam(ctx, "testuser", team.Uuid)
		assert.NoError(t, err)
		err = c.DeleteUserMembership(ctx, team.Uuid, "testuser")
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

	t.Run("GetOidcUsers, GetOidcUser", func(t *testing.T) {
		users, err := c.GetOidcUsers(ctx)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "email@nav.no", users[0].Username)

		user, err := c.GetOidcUser(ctx, users[0].Username)
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("DeleteOidcUser", func(t *testing.T) {
		err := c.DeleteOidcUser(ctx, "email@nav.no")
		assert.NoError(t, err)
	})

	t.Run("ProjectMetricsRefresh, AllMetricsRefresh", func(t *testing.T) {
		p, err := c.CreateProject(ctx, "projectname", "version1", nil)
		assert.NoError(t, err)
		assert.NotNil(t, p)
		err = c.ProjectMetricsRefresh(ctx, p.Uuid)
		assert.NoError(t, err)

		err = c.AllMetricsRefresh(ctx)
		assert.NoError(t, err)
	})

	t.Run("CreateProjectWithSbom", func(t *testing.T) {
		sbom, err := getSbom(t)
		pWithSbom, err := c.CreateProjectWithSbom(ctx, sbom, "project-with-sbom", "version1")
		assert.NoError(t, err)
		assert.NotEmpty(t, pWithSbom)
		p, err := c.GetProject(ctx, "project-with-sbom", "version1")
		assert.NoError(t, err)
		assert.Equal(t, "project-with-sbom", p.Name)
		assert.Equal(t, "version1", p.Version)
	})

	t.Run("GetConfigProperties", func(t *testing.T) {
		props, err := c.GetConfigProperties(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, props)
	})

	t.Run("ConfigPropertyAggregate", func(t *testing.T) {
		groupName := "vuln-source"
		propertyName := "google.osv.enabled"
		propertyValue := "false"
		description := "List of enabled ecosystems to mirror OSV"
		err := c.ConfigPropertyAggregate(ctx, dependencytrack.ConfigProperty{
			GroupName:     &groupName,
			PropertyName:  &propertyName,
			PropertyValue: &propertyValue,
			PropertyType:  "STRING",
			Description:   &description,
		})
		assert.NoError(t, err)
	})

	t.Run("GetEcosystems", func(t *testing.T) {
		ecos, err := c.GetEcosystems(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, ecos)
	})

	t.Run("CreateProjectWithSbom, GetVulnerabilities, TriggerAnalysis", func(t *testing.T) {
		sbom, err := getSbom(t)
		assert.NoError(t, err)
		p, err := c.CreateProjectWithSbom(ctx, sbom, "project-with-findings", "version1")
		assert.NoError(t, err)
		assert.NotNil(t, p)

		err = c.TriggerAnalysis(ctx, p)
		assert.NoError(t, err)
		time.Sleep(3 * time.Second) // wait for analysis to complete

		v, err := c.GetVulnerabilities(ctx, p, "", true)
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	})

	t.Run("GetAnalysisTrail", func(t *testing.T) {
		sbom, err := getSbom(t)
		assert.NoError(t, err)

		p, err := c.CreateProjectWithSbom(ctx, sbom, "project-with-findings", "version2")
		assert.NoError(t, err)
		assert.NotNil(t, p)

		findings, err := c.GetVulnerabilities(ctx, p, "", true)
		assert.NoError(t, err)
		assert.NotNil(t, findings)
	})

	t.Run("RecordAnalysis", func(t *testing.T) {
		p, err := c.GetProject(ctx, "project-with-findings", "version1")
		assert.NoError(t, err)
		assert.NotNil(t, p)

		vulnerabilities, err := c.GetVulnerabilities(ctx, p.Uuid, "", true)
		assert.NoError(t, err)

		err = c.TriggerAnalysis(ctx, p.Uuid)
		assert.NoError(t, err)
		err = c.ProjectMetricsRefresh(ctx, p.Uuid)
		assert.NoError(t, err)
		time.Sleep(2 * time.Second) // wait for analysis to complete

		if len(vulnerabilities) == 0 {
			t.Skip("no findings")
		}
		assert.NoError(t, err)
		assert.NotNil(t, vulnerabilities)
		// will never be run if there are no findings
		err = c.UpdateFinding(ctx, dependencytrack.AnalysisRequest{
			ProjectId:       p.Uuid,
			Reason:          "dont care",
			ComponentId:     vulnerabilities[0].Metadata.ComponentId,
			VulnerabilityId: vulnerabilities[0].Metadata.VulnerabilityUuid,
			State:           "NOT_AFFECTED",
			SuppressedBy:    "ME",
			Suppressed:      true,
		})
		assert.NoError(t, err)

		trail, err := c.GetAnalysisTrailForImage(ctx, p.Uuid, vulnerabilities[0].Metadata.ComponentId, vulnerabilities[0].Metadata.VulnerabilityUuid)
		assert.NoError(t, err)
		assert.Equal(t, "NOT_AFFECTED", trail.AnalysisState)
		assert.Equal(t, "NOT_SET", trail.AnalysisResponse)
		assert.Equal(t, "dont care", trail.AnalysisDetails)
		assert.Len(t, trail.AnalysisComments, 4)
		assert.Equal(t, "Analysis: NOT_SET → NOT_AFFECTED", trail.AnalysisComments[0].Comment)
		assert.Equal(t, "admin", *trail.AnalysisComments[0].Commenter)
	})
}

func getSbom(t *testing.T) (*in_toto.CycloneDXStatement, error) {
	read, err := os.ReadFile("testdata/attestation.json")
	assert.NoError(t, err)

	metadata := &in_toto.CycloneDXStatement{}
	err = json.Unmarshal(read, metadata)
	assert.NoError(t, err)

	return metadata, nil
}

func DependencyTrackPool(version string) (string, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	opts := &dockertest.RunOptions{
		Repository:   "dependencytrack/apiserver",
		Tag:          version,
		ExposedPorts: []string{"8080"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8080": {
				{HostIP: "", HostPort: "1234"},
			},
		},
	}

	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = resource.Expire(60)
	if err != nil {
		log.Fatalf("Could not set expiration: %s", err)
	}
	port := "1234" //resource.GetPort("8080")
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:"+port, nil)
		if err != nil {
			return fmt.Errorf("creating request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			return fmt.Errorf("got status code %d", resp.StatusCode)
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to dtrack: %s", err)
	}

	return "http://localhost:" + port + "/api", func() {
		// You can't defer this because os.Exit doesn't care for defer
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
