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
		cluster := getTag(project, "cluster:")
		if cluster == "" {
			log.Errorf("missing cluster tag for project %s, skipping", project.Name)
			continue
		}
		observability.DependencytrackTotalProjects.WithLabelValues(cluster).Inc()
	}

	return nil
}

// TODO: fix to get tag with prefix, e.g. cluster:
func getTag(p *client.Project, t string) string {
	for _, tag := range p.Tags {
		if strings.HasPrefix(tag.Name, t) {
			s, _ := strings.CutPrefix(tag.Name, t)
			return s
		}
	}
	return ""
}
