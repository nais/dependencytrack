package dependencytrack

import (
	"context"
)

// GetEcosystems retrieves all ecosystems from the Dependency-Track server.
func (c *dependencyTrackClient) GetEcosystems(ctx context.Context) ([]string, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]string, error) {
		res, resp, err := c.client.IntegrationAPI.GetAllEcosystems(tokenCtx).Execute()
		if err != nil {
			return nil, convertError(err, "GetEcosystems", resp)
		}
		return res, nil
	})
}
