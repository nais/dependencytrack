package dependencytrack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
)

type Project struct {
	Active                 bool           `json:"active"`
	Author                 string         `json:"author"`
	Classifier             string         `json:"classifier"`
	Group                  string         `json:"group"`
	Name                   string         `json:"name"`
	LastBomImportFormat    string         `json:"lastBomImportFormat,omitempty"`
	LastBomImport          int64          `json:"lastBomImport,omitempty"`
	LastInheritedRiskScore float64        `json:"lastInheritedRiskScore,omitempty"`
	Publisher              string         `json:"publisher"`
	Tags                   []Tag          `json:"tags"`
	Uuid                   string         `json:"uuid"`
	Version                string         `json:"version"`
	Metrics                *ProjectMetric `json:"metrics,omitempty"`
}

type ProjectMetric struct {
	Critical             int32   `json:"critical"`
	High                 int32   `json:"high"`
	Medium               int32   `json:"medium"`
	Low                  int32   `json:"low"`
	Unassigned           int32   `json:"unassigned"`
	Vulnerabilities      int64   `json:"vulnerabilities"`
	VulnerableComponents int32   `json:"vulnerableComponents"`
	Components           int32   `json:"components"`
	Suppressed           int32   `json:"suppressed"`
	FindingsTotal        int32   `json:"findingsTotal"`
	FindingsAudited      int32   `json:"findingsAudited"`
	FindingsUnaudited    int32   `json:"findingsUnaudited"`
	InheritedRiskScore   float64 `json:"inheritedRiskScore"`
	FirstOccurrence      int64   `json:"firstOccurrence"`
	LastOccurrence       int64   `json:"lastOccurrence"`
}

type Tag struct {
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag `json:"tags"`
}

func (c *dependencyTrackClient) GetProject(ctx context.Context, name, version string) (*Project, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*Project, error) {
		project, resp, err := c.client.ProjectAPI.GetProjectByNameAndVersion(tokenCtx).
			Name(name).
			Version(version).
			Execute()

		if err != nil && resp != nil && resp.Body != nil && resp.StatusCode == 404 {
			body, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read response body: %w", readErr)
			}

			if strings.Contains(string(body), "The project could not be found") {
				return nil, nil
			}

			return nil, convertError(err, "GetProject", resp)
		}

		return parseProject(project), err
	})
}

func (c *dependencyTrackClient) GetProjects(ctx context.Context, limit, offset int32) ([]Project, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]Project, error) {
		return c.paginateProjects(tokenCtx, limit, offset, func(ctx context.Context, offset int32) ([]Project, error) {
			pageNumber := (offset / limit) + 1
			// convert to string from int32
			projects, resp, err := c.client.ProjectAPI.GetProjects(ctx).
				PageSize(strconv.Itoa(int(limit))).
				PageNumber(strconv.Itoa(int(pageNumber))).
				Execute()
			pp := make([]Project, 0)
			if len(projects) > 0 {
				pp = make([]Project, len(projects))
				for i, project := range projects {
					pp[i] = *parseProject(&project)
				}
			}
			return pp, convertError(err, "GetProjects", resp)
		})
	})
}

func (c *dependencyTrackClient) CreateProjectWithSbom(ctx context.Context, sbom *in_toto.CycloneDXStatement, imageName, imageTag string) (string, error) {
	p, err := c.GetProject(ctx, imageName, imageTag)
	if err != nil {
		return "", fmt.Errorf("failed to lookup project: %w", err)
	}

	if p == nil {
		p, err = c.CreateProject(ctx, imageName, imageTag, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create project: %w", err)
		}
		if p == nil {
			return "", fmt.Errorf("created project is unexpectedly nil")
		}
	}

	if p.Uuid == "" {
		return "", fmt.Errorf("project UUID is empty")
	}

	if err = c.uploadSbom(ctx, p.Uuid, sbom); err != nil {
		var clientErr *ClientError
		if errors.As(err, &clientErr) {
			if deleteErr := c.DeleteProject(ctx, p.Uuid); deleteErr != nil {
				return "", fmt.Errorf("upload failed: %w (also failed to delete project: %v)", err, deleteErr)
			}
		}
		return "", fmt.Errorf("failed to upload SBOM: %w", err)
	}

	return p.Uuid, nil
}

func (c *dependencyTrackClient) uploadSbom(ctx context.Context, projectId string, sbom *in_toto.CycloneDXStatement) error {
	b, err := json.Marshal(sbom.Predicate)
	if err != nil {
		return fmt.Errorf("failed to marshal sbom: %w", err)
	}

	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		req := c.client.BomAPI.UploadBom(tokenCtx).Bom(string(b)).Project(projectId).AutoCreate(false)
		_, resp, err := req.Execute()
		if err != nil {
			return convertError(err, "uploadSbom", resp)
		}
		return nil
	})
}

func (c *dependencyTrackClient) DeleteProject(ctx context.Context, uuid string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		resp, err := c.client.ProjectAPI.DeleteProject(tokenCtx, uuid).Execute()
		if err != nil {
			return convertError(err, "DeleteProject", resp)
		}
		return nil
	})
}

func (c *dependencyTrackClient) CreateProject(ctx context.Context, name, version string, tags []string) (*Project, error) {
	c.log.Debugf("creating project: %s", name+":"+version)
	p, err := withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*Project, error) {
		active := true
		classifier := "APPLICATION"

		var clientTags []client.Tag
		if tags != nil {
			clientTags = make([]client.Tag, len(tags))
			for i, tag := range tags {
				clientTags[i] = client.Tag{Name: &tag}
			}
		} else {
			clientTags = nil
		}

		req := c.client.ProjectAPI.CreateProject(tokenCtx).Project(client.Project{
			Name:       &name,
			Active:     &active,
			Classifier: &classifier,
			Version:    &version,
			Tags:       clientTags,
			Parent:     nil,
		})

		project, resp, err := req.Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusConflict {
				return nil, convertError(fmt.Errorf("project %s:%s already exists", name, version), "CreateProject", resp)
			}
			return nil, convertError(err, "CreateProject", resp)
		}

		return parseProject(project), nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func parseProject(project *client.Project) *Project {
	if project == nil {
		return nil
	}

	var tags []Tag
	if project.Tags != nil {
		tags = make([]Tag, len(project.Tags))
		for i, tag := range project.Tags {
			tags[i] = parseTag(&tag)
		}
	}

	return &Project{
		Active:                 SafeBool(project.Active),
		Author:                 SafeString(project.Author),
		Classifier:             SafeString(project.Classifier),
		Group:                  SafeString(project.Group),
		Name:                   SafeString(project.Name),
		LastBomImportFormat:    SafeString(project.LastBomImportFormat),
		LastBomImport:          project.LastBomImport,
		LastInheritedRiskScore: safeFloat64(project.LastInheritedRiskScore),
		Publisher:              SafeString(project.Publisher),
		Tags:                   tags,
		Uuid:                   project.Uuid,
		Version:                SafeString(project.Version),
		Metrics:                parseMetrics(project.Metrics),
	}

}

func parseMetrics(metrics *client.ProjectMetrics) *ProjectMetric {
	if metrics == nil {
		return nil
	}

	return &ProjectMetric{
		Critical:             metrics.Critical,
		High:                 metrics.High,
		Medium:               metrics.Medium,
		Low:                  metrics.Low,
		Unassigned:           safeInt32(metrics.Unassigned),
		Vulnerabilities:      safeInt64(metrics.Vulnerabilities),
		VulnerableComponents: safeInt32(metrics.VulnerableComponents),
		Components:           safeInt32(metrics.Components),
		Suppressed:           safeInt32(metrics.Suppressed),
		FindingsTotal:        safeInt32(metrics.FindingsTotal),
		FindingsAudited:      safeInt32(metrics.FindingsAudited),
		FindingsUnaudited:    safeInt32(metrics.FindingsUnaudited),
		InheritedRiskScore:   safeFloat64(metrics.InheritedRiskScore),
		FirstOccurrence:      metrics.FirstOccurrence,
		LastOccurrence:       metrics.LastOccurrence,
	}
}

func parseTag(tags *client.Tag) Tag {
	if tags == nil {
		return Tag{}
	}
	return Tag{
		Name: SafeString(tags.Name),
	}
}

func safeInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}

func safeInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func safeFloat64(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}
