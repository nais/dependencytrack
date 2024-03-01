package client

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *client) GetCurrentProjectMetric(ctx context.Context, projectUuid string) (*ProjectMetric, error) {
	res, err := c.get(ctx, c.baseUrl+"/api/v1/metrics/project/"+projectUuid+"/current", c.authSource)
	if err != nil {
		if IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("get metric: %w", err)
	}

	var metrics *ProjectMetric
	if err = json.Unmarshal(res, &metrics); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return metrics, nil
}

func (c *client) GetProjectMetricsByDate(ctx context.Context, projectUuid, date string) ([]*ProjectMetric, error) {
	res, err := c.get(ctx, c.baseUrl+"/api/v1/metrics/project/"+projectUuid+"/since/"+date, c.authSource)
	if err != nil {
		if IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("get metric: %w", err)
	}

	var metrics []*ProjectMetric
	if err = json.Unmarshal(res, &metrics); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return metrics, nil
}
