package dependencytrack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"net/http"

	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
)

func (c *dependencyTrackClient) GetProject(ctx context.Context, name, version string) (*client.Project, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*client.Project, error) {
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

		return project, err
	})
}

func (c *dependencyTrackClient) GetProjects(ctx context.Context, limit, offset int32) ([]client.Project, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]client.Project, error) {
		return c.paginateProjects(tokenCtx, limit, offset, func(ctx context.Context, offset int32) ([]client.Project, error) {
			pageNumber := (offset / limit) + 1
			// convert to string from int32
			projects, resp, err := c.client.ProjectAPI.GetProjects(ctx).
				PageSize(strconv.Itoa(int(limit))).
				PageNumber(strconv.Itoa(int(pageNumber))).
				Execute()
			return projects, convertError(err, "GetProjects", resp)
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

func (c *dependencyTrackClient) CreateProject(ctx context.Context, name, version string, tags []client.Tag) (*client.Project, error) {
	c.log.Debugf("creating project: %s", name+":"+version)
	p, err := withAuthContextValue(c, ctx, func(tokenCtx context.Context) (*client.Project, error) {
		active := true
		classifier := "APPLICATION"

		req := c.client.ProjectAPI.CreateProject(tokenCtx).Project(client.Project{
			Name:       &name,
			Active:     &active,
			Classifier: &classifier,
			Version:    &version,
			Tags:       tags,
			Parent:     nil,
		})

		project, resp, err := req.Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusConflict {
				return nil, convertError(fmt.Errorf("project %s:%s already exists", name, version), "CreateProject", resp)
			}
			return nil, convertError(err, "CreateProject", resp)
		}

		return project, nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}
