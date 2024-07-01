package dependencytrack

import (
	"context"
	"github.com/nais/dependencytrack/internal/observability"
	"github.com/nais/dependencytrack/pkg/client"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	ProjectsCacheKey = "projects"
)

type Client struct {
	client.Client
	Cache *cache.Cache
}

func NewClient(baseUrl, username, password string) *Client {
	return &Client{
		Client: client.New(baseUrl, username, password),
		Cache:  cache.New(10*time.Minute, 5*time.Minute),
	}
}

func (c *Client) getProjects(ctx context.Context) ([]*client.Project, error) {
	if v, ok := c.Cache.Get(ProjectsCacheKey); ok {
		return v.([]*client.Project), nil
	}

	projects, err := c.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	c.Cache.Set(ProjectsCacheKey, projects, cache.DefaultExpiration)
	return projects, nil
}

func (c *Client) UpdateTotalProjects(ctx context.Context) error {
	projects, err := c.getProjects(ctx)
	if err != nil {
		return err
	}
	for _, project := range projects {
		clusters := tagsWithPrefix(project.Tags, client.EnvironmentTagPrefix.String())
		if len(clusters) == 0 {
			log.Warnf("project %s has no cluster tag", project.Name)
			continue
		}
		teams := tagsWithPrefix(project.Tags, client.TeamTagPrefix.String())
		if len(teams) == 0 {
			log.Warnf("project %s has no team tag", project.Name)
			continue
		}

		workloads := tagsWithPrefix(project.Tags, client.WorkloadTagPrefix.String())
		if len(workloads) == 0 {
			log.Warnf("project %s has no workload tag", project.Name)
			continue
		}

		for _, cluster := range clusters {
			for _, team := range teams {
				for _, workload := range workloads {
					// only count projects that are relevant for the cluster and team
					if !validWorkload(workload, cluster, team) {
						continue
					}

					workloadType := workloadType(workload)
					if workloadType == "" {
						log.Warnf("project %s has no workload type", project.Name)
						continue
					}
					sbom := strconv.FormatBool(hasSbom(project))
					if strings.Contains(project.Name, "nais-io") {
						name := platformName(project.Name)
						observability.DependencytrackTotalPlatformProjects.WithLabelValues(cluster, team, workloadType, name, sbom).Inc()
					} else {
						observability.DependencytrackTotalProjects.WithLabelValues(cluster, team, workloadType, sbom).Inc()
					}
				}
			}
		}
	}

	return nil
}

func hasSbom(project *client.Project) bool {
	return project.Metrics != nil && project.Metrics.Components > 0
}

func platformName(workload string) string {
	return workload[strings.LastIndex(workload, "/")+1:]
}

func workloadType(workload string) string {
	parts := strings.Split(workload, "|")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

func validWorkload(workload, cluster, team string) bool {
	return strings.Contains(workload, cluster) && strings.Contains(workload, team)
}

func tagsWithPrefix(tags []client.Tag, prefix string) []string {
	var result []string
	for _, tag := range tags {
		if strings.HasPrefix(tag.Name, prefix) {
			t, _ := strings.CutPrefix(tag.Name, prefix)
			result = append(result, t)
		}
	}
	return result
}
