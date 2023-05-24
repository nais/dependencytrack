package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

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
	b, err := c.get(ctx, c.baseUrl+"/api/v1/project", c.authSource)
	if err != nil {
		return nil, fmt.Errorf("get projects: %w", err)
	}
	projects := make([]*Project, 0)
	if err = json.Unmarshal(b, &projects); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return projects, nil
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

func (c *client) CreateChildProject(ctx context.Context, parentUuid, name, version, group string, tags []string) (*Project, error) {
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
		Classifier: "CONTAINER",
		Version:    version,
		Group:      group,
		Tags:       t,
		Parent:     parentUuid,
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
