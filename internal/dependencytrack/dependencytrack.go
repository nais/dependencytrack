package dependencytrack

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/nais/dependencytrack/internal/observability"
	"github.com/nais/dependencytrack/pkg/client"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
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
		Cache:  cache.New(20*time.Minute, 3*time.Minute),
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

func (c *Client) UpdateTotalProjects(ctx context.Context, tenant string, imagesIgnore string) error {
	projects, err := c.getProjects(ctx)
	if err != nil {
		return err
	}

	observability.WorkloadRegistered.Reset()
	observability.WorkloadRiskscore.Reset()
	observability.WorkloadCritical.Reset()
	teamWorkloads := map[string]bool{}
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

		// TODO: find a better way imagesIgnore for now
		if IsInImagesIgnoreList(project.Name, imagesIgnore) {
			continue
		}

		for _, cluster := range clusters {
			for _, team := range teams {
				for _, workload := range workloads {
					w := workloadName(workload)
					if w == "" {
						log.Warnf("project %s has no workload name", project.Name)
						continue
					}
					sbom := strconv.FormatBool(hasSbom(project))
					tenantCluster := tenant + "-" + cluster
					riskScore := 0.0
					critical := 0.0
					if hasSbom(project) {
						riskScore = project.Metrics.InheritedRiskScore
						critical = float64(project.Metrics.Critical)
					}
					teamWorkloads[tenantCluster+":"+team+":"+workloadName(workload)] = hasSbom(project)
					observability.WorkloadRiskscore.WithLabelValues(tenantCluster, team, sbom, project.Name, project.Version).Set(riskScore)
					observability.WorkloadCritical.WithLabelValues(tenantCluster, team, sbom, project.Name, project.Version).Set(critical)
				}
			}
		}
	}

	for teamWorkload, sbom := range teamWorkloads {
		parts := strings.Split(teamWorkload, ":")
		tenantCluster := parts[0]
		team := parts[1]
		workload := parts[2]
		observability.WorkloadRegistered.WithLabelValues(tenantCluster, team, workload, strconv.FormatBool(sbom)).Set(1)
	}

	return nil
}

func IsInImagesIgnoreList(image string, imagesIgnore string) bool {
	if imagesIgnore == "" {
		return false
	}
	images := strings.Split(imagesIgnore, ",")
	for _, i := range images {
		if strings.Contains(image, i) {
			return true
		}
	}
	return false
}

func hasSbom(project *client.Project) bool {
	return project.Metrics != nil && project.Metrics.Components > 0
}

func workloadName(workload string) string {
	parts := strings.Split(workload, "|")
	if len(parts) < 4 {
		return ""
	}
	return parts[3]
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
