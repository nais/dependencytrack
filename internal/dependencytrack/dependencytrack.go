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

type Workload struct {
	Cluster           string
	WorkloadNamespace string
	WorkloadType      string
	Workload          string
}

func (c *Client) UpdateTotalProjects(ctx context.Context, tenant string, imagesIgnore string) error {
	projects, err := c.getProjects(ctx)
	if err != nil {
		return err
	}

	observability.WorkloadRiskscore.Reset()
	observability.WorkloadCritical.Reset()
	for _, project := range projects {

		// TODO: find a better way imagesIgnore for now
		if IsInImagesIgnoreList(project.Name, imagesIgnore) {
			continue
		}

		workloads := getWorkloads(project)
		for _, w := range workloads {
			m := map[string]string{
				"cluster":            tenant + "-" + w.Cluster,
				"workload_namespace": w.WorkloadNamespace,
				"project":            project.Name,
				"workload_type":      w.WorkloadType,
				"workload":           w.Workload,
				"has_sbom":           strconv.FormatBool(hasSbom(project)),
			}
			riskScore := 0.0
			critical := 0.0
			if hasSbom(project) {
				riskScore = project.Metrics.InheritedRiskScore
				critical = float64(project.Metrics.Critical)
			}

			observability.WorkloadRiskscore.With(m).Set(riskScore)
			observability.WorkloadCritical.With(m).Set(critical)
		}
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
	return project.Metrics != nil
}

func getWorkloads(p *client.Project) []*Workload {
	tags := tagsWithPrefix(p.Tags, client.WorkloadTagPrefix.String())
	workloads := make([]*Workload, 0)
	for _, tag := range tags {
		parts := strings.Split(tag, "|")
		if len(parts) < 4 {
			log.Warnf("workload tag %s is invalid", tag)
			continue
		}
		workloads = append(workloads, &Workload{
			Cluster:           parts[0],
			WorkloadNamespace: parts[1],
			WorkloadType:      parts[2],
			Workload:          parts[3],
		})
	}

	return workloads
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
