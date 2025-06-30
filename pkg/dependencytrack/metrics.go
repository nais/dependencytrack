package dependencytrack

import (
	"context"
)

// ProjectMetricsRefresh refreshes the metrics for a specific project identified by its UUID.
func (c *dependencyTrackClient) ProjectMetricsRefresh(ctx context.Context, uuid string) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		resp, err := c.client.MetricsAPI.RefreshProjectMetrics(tokenCtx, uuid).Execute()
		if err != nil {
			return convertError(err, "PortfolioRefresh", resp)
		}
		return nil
	})
}

// AllMetricsRefresh refreshes all portfolio metrics.
func (c *dependencyTrackClient) AllMetricsRefresh(ctx context.Context) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		resp, err := c.client.MetricsAPI.RefreshPortfolioMetrics(tokenCtx).Execute()
		if err != nil {
			return convertError(err, "PortfolioRefresh", resp)
		}
		return nil
	})
}
