package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type ProjectMetric struct {
	Critical             int   `json:"critical"`
	High                 int   `json:"high"`
	Medium               int   `json:"medium"`
	Low                  int   `json:"low"`
	Unassigned           int   `json:"unassigned"`
	Vulnerabilities      int   `json:"vulnerabilities"`
	VulnerableComponents int   `json:"vulnerableComponents"`
	Components           int   `json:"components"`
	Suppressed           int   `json:"suppressed"`
	FindingsTotal        int   `json:"findingsTotal"`
	FindingsAudited      int   `json:"findingsAudited"`
	FindingsUnaudited    int   `json:"findingsUnaudited"`
	InheritedRiskScore   int   `json:"inheritedRiskScore"`
	FirstOccurrence      int64 `json:"firstOccurrence"`
	LastOccurrence       int64 `json:"lastOccurrence"`
}

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
