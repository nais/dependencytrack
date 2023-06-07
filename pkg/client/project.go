package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (c *client) UploadProject(ctx context.Context, name, version string, bom []byte) error {
	c.log.WithFields(log.Fields{
		"name":    name,
		"version": version,
	}).Info("uploading sbom")

	body, err := json.Marshal(&BomSubmitRequest{
		ProjectName:    name,
		ProjectVersion: version,
		AutoCreate:     true,
		Bom:            base64.StdEncoding.EncodeToString(bom),
	})
	if err != nil {
		return fmt.Errorf("marshalling bom submit request: %w", err)
	}

	_, err = c.put(ctx, c.baseUrl+"/api/v1/bom", c.authSource, body)
	if err != nil {
		return fmt.Errorf("uploading bom: %w", err)
	}

	c.log.Info("sbom uploaded")
	return nil
}

func (c *client) GetProjects(ctx context.Context) ([]*Project, error) {
	headers, err := c.authSource.Headers(ctx)
	if err != nil {
		return nil, err
	}

	headers["Accept"] = []string{"application/json"}
	resp, err := c.httpClient.SendRequest(ctx, http.MethodGet, c.baseUrl+"/api/v1/project?offset=0", headers, nil)
	if err != nil {
		return nil, err
	}

	allProjects := make([]*Project, 0)
	if err = json.Unmarshal(resp.Body, &allProjects); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	header := resp.Headers.Get("X-Total-Count")
	if header == "" {
		return nil, fmt.Errorf("missing X-Total-Count header")
	}

	numProjects, err := strconv.Atoi(header)
	if err != nil {
		return nil, fmt.Errorf("parsing X-Total-Count header: %w", err)
	}
	remainingPages := numProjects - 1/100

	if remainingPages > 0 {
		offset := 0
		for i := 0; i < remainingPages; i++ {
			offset += 100
			projects := make([]*Project, 0)
			b, err := c.get(ctx, c.baseUrl+"/api/v1/project?offset="+strconv.Itoa(offset), c.authSource)
			if err != nil {
				return nil, fmt.Errorf("get projects: %w", err)
			}
			if err = json.Unmarshal(b, &projects); err != nil {
				return nil, fmt.Errorf("unmarshalling response body: %w", err)
			}
			allProjects = append(allProjects, projects...)
		}
	}

	return allProjects, nil
}

func (c *client) GetProject(ctx context.Context, name, version string) (*Project, error) {
	res, err := c.get(ctx, c.baseUrl+"/api/v1/project/lookup?name="+name+"&version="+version, c.authSource)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	var project Project
	if err = json.Unmarshal(res, &project); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return &project, nil
}

func (c *client) GetProjectsByTag(ctx context.Context, tag string) ([]*Project, error) {
	res, err := c.get(ctx, c.baseUrl+"/api/v1/project/tag/"+tag, c.authSource)
	if err != nil {
		return nil, fmt.Errorf("get projects by tag: %w", err)
	}

	var projects []*Project
	if err = json.Unmarshal(res, &projects); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return projects, nil
}

func (c *client) CreateProject(ctx context.Context, name, version, group string, tags []string) (*Project, error) {
	c.log.WithFields(log.Fields{
		"group": group,
		"tags":  tags,
	}).Debug("creating project")

	t := make([]Tag, 0)
	for _, tag := range tags {
		t = append(t, Tag{
			Name: tag,
		})
	}

	pp := Project{
		Name:       name,
		Publisher:  group,
		Active:     true,
		Classifier: "APPLICATION",
		Version:    version,
		Group:      group,
		Tags:       t,
		Parent:     nil,
	}

	body, err := json.Marshal(pp)

	p, err := c.put(ctx, c.baseUrl+"/api/v1/project", c.authSource, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var project Project
	if err = json.Unmarshal(p, &project); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return &project, nil
}

func (c *client) UpdateProject(ctx context.Context, uuid, name, version, group string, tags []string) (*Project, error) {
	c.log.WithFields(log.Fields{
		"group": group,
		"tags":  tags,
	}).Debug("updating project")

	t := make([]Tag, 0)
	for _, tag := range tags {
		t = append(t, Tag{
			Name: tag,
		})
	}

	pp := Project{
		Uuid:       uuid,
		Name:       name,
		Publisher:  group,
		Active:     true,
		Classifier: "APPLICATION",
		Version:    version,
		Group:      group,
		Tags:       t,
		Parent:     nil,
	}

	body, err := json.Marshal(pp)

	p, err := c.post(ctx, c.baseUrl+"/api/v1/project", c.authSource, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var project Project
	if err = json.Unmarshal(p, &project); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return &project, nil
}

func (c *client) CreateChildProject(ctx context.Context, parent *Project, name, version, group string, tags []string) (*Project, error) {
	c.log.WithFields(log.Fields{
		"group": group,
		"tags":  tags,
	}).Debug("creating child project")

	t := make([]Tag, 0)
	for _, tag := range tags {
		t = append(t, Tag{
			Name: tag,
		})
	}

	pp := Project{
		Name:       name,
		Publisher:  group,
		Active:     true,
		Classifier: "APPLICATION",
		Version:    version,
		Group:      group,
		Tags:       t,
		Parent:     parent,
	}

	body, err := json.Marshal(pp)

	p, err := c.put(ctx, c.baseUrl+"/api/v1/project", c.authSource, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var project Project
	if err = json.Unmarshal(p, &project); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return &project, nil
}

func (c *client) UpdateProjectInfo(ctx context.Context, uuid, version, group string, tags []string) error {
	c.log.WithFields(log.Fields{
		"uuid":  uuid,
		"group": group,
		"tags":  tags,
	}).Debug("adding additional info to project")

	t := make([]Tag, 0)
	for _, tag := range tags {
		t = append(t, Tag{
			Name: tag,
		})
	}

	body, err := json.Marshal(Project{
		Publisher:  group,
		Active:     true,
		Classifier: "APPLICATION",
		Version:    version,
		Group:      group,
		Tags:       t,
	})

	_, err = c.patch(ctx, c.baseUrl+"/api/v1/project/"+uuid, c.authSource, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	c.log.Info("additional info added to project")

	return nil
}

func (c *client) DeleteProjects(ctx context.Context, name string) error {
	b, err := c.get(ctx, c.baseUrl+"/api/v1/project"+"?name="+name+"&excludeInactive=false", c.authSource)
	if err != nil {
		return fmt.Errorf("getting projects: %w", err)
	}
	var projects []Project
	if err = json.Unmarshal(b, &projects); err != nil {
		return fmt.Errorf("unmarshalling response body: %w", err)
	}

	for _, project := range projects {
		_, err := c.delete(ctx, c.baseUrl+"/api/v1/project/"+project.Uuid, c.authSource, nil)
		if err != nil {
			return fmt.Errorf("deleting project: %w", err)
		}
	}
	return nil
}

func (c *client) DeleteProject(ctx context.Context, uuid string) error {
	_, err := c.delete(ctx, c.baseUrl+"/api/v1/project/"+uuid, c.authSource, nil)
	if err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}
	return nil
}

func (c *client) PortfolioRefresh(ctx context.Context) error {
	_, err := c.get(ctx, c.baseUrl+"/api/v1/metrics/portfolio/refresh", c.authSource)
	if err != nil {
		return fmt.Errorf("refreshing portfolio: %w", err)
	}
	return nil
}
