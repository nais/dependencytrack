package dependencytrack

import (
	"context"
	"github.com/nais/dependencytrack/internal/observability"
	"github.com/nais/dependencytrack/pkg/client"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Client struct {
	client.Client
}

func NewClient(baseUrl, username, password string) *Client {
	return &Client{
		Client: client.New(baseUrl, username, password),
	}
}

func (c *Client) UpdateTotalProjects(ctx context.Context) error {
	projects, err := c.GetProjects(ctx)
	if err != nil {
		return err
	}
	for _, project := range projects {
		clusters := getTagsWithPrefix(project.Tags, client.EnvironmentTagPrefix.String())
		if len(clusters) == 0 {
			log.Warnf("project %s has no cluster tag", project.Name)
			continue
		}
		teams := getTagsWithPrefix(project.Tags, client.TeamTagPrefix.String())
		if len(teams) == 0 {
			log.Warnf("project %s has no team tag", project.Name)
			continue
		}

		workloads := getTagsWithPrefix(project.Tags, client.WorkloadTagPrefix.String())
		if len(workloads) == 0 {
			log.Warnf("project %s has no workload tag", project.Name)
			continue
		}

		for _, cluster := range clusters {
			for _, team := range teams {
				for _, workload := range workloads {
					// platform projects are not interesting for this metric
					if strings.Contains(project.Name, "nais-io") {
						continue
					}

					// only count projects that are relevant for the cluster and team
					if !validWorkload(workload, cluster, team) {
						continue
					}

					// only count projects that have metrics a.k.a. sbom exists
					if !hasSbom(project) {
						log.Warnf("project %s cluster: %s uuid: %s has no metrics", project.Name, cluster, project.Uuid)
						continue
					}

					// app or job
					workloadType := getWorkloadType(workload)
					if workloadType == "" {
						log.Warnf("project %s has no workload type", project.Name)
						continue
					}
					observability.DependencytrackTotalProjects.WithLabelValues(cluster, team, workloadType).Inc()
				}
			}
		}
	}

	return nil
}

func hasSbom(project *client.Project) bool {
	return project.Metrics != nil && project.Metrics.Components > 0
}

func getWorkloadType(workload string) string {
	parts := strings.Split(workload, "|")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

func validWorkload(workload, cluster, team string) bool {
	return strings.Contains(workload, cluster) && strings.Contains(workload, team)
}

func getTagsWithPrefix(tags []client.Tag, prefix string) []string {
	var result []string
	for _, tag := range tags {
		if strings.HasPrefix(tag.Name, prefix) {
			t, _ := strings.CutPrefix(tag.Name, prefix)
			result = append(result, t)
		}
	}
	return result
}
