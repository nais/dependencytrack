package client

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *client) RecordAnalysis(ctx context.Context, projectUuid string) error {
	_, err := c.post(ctx, c.baseUrl+"/api/v1/finding/project/"+projectUuid+"/analyze", c.authSource, nil)
	if err != nil {
		if IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("trigger analysis: %w", err)
	}
	return nil
}

func (c *client) GetAnalysisTrail(ctx context.Context, projectUuid string) ([]*AnalysisTrail, error) {
	var trails []*AnalysisTrail
	res, err := c.get(ctx, c.baseUrl+"/api/v1/finding/project/"+projectUuid+"/analyze", c.authSource)
	if err != nil {
		if IsNotFound(err) {
			return trails, nil
		}
		return nil, fmt.Errorf("get analysis: %w", err)
	}
	if err := json.Unmarshal(res, &trails); err != nil {
		return nil, fmt.Errorf("unmarshal analysis: %w", err)
	}
	return trails, nil
}
