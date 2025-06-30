package dependencytrack

import (
	"context"
	"net/http"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
)

func (c *dependencyTrackClient) ConfigPropertyAggregate(ctx context.Context, property client.ConfigProperty) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		_, resp, err := c.client.ConfigPropertyAPI.UpdateConfigProperty(tokenCtx).ConfigProperty(property).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return nil // Not found, return empty property
			}
			return convertError(err, "ConfigPropertyAggregate", resp)
		}
		return nil
	})
}

func (c *dependencyTrackClient) GetConfigProperties(ctx context.Context) ([]client.ConfigProperty, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]client.ConfigProperty, error) {
		res, resp, err := c.client.ConfigPropertyAPI.GetConfigProperties(tokenCtx).Execute()
		if err != nil {
			return nil, convertError(err, "GetConfigProperties", resp)
		}
		return res, nil
	})
}
