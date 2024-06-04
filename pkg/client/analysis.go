package client

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *client) RecordAnalysis(ctx context.Context, analysisRequest *AnalysisRequest) error {
	body, err := json.Marshal(analysisRequest)
	if err != nil {
		return fmt.Errorf("marshalling bom submit request: %w", err)
	}
	_, err = c.put(ctx, c.baseUrl+"/api/v1/analysis", c.authSource, body)
	if err != nil {
		if IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("trigger analysis: %w", err)
	}
	return nil
}

func (c *client) GetAnalysisTrail(ctx context.Context, projectUuid, componentUuid, vulnerabilityUuid string) (*Analysis, error) {
	var trail *Analysis
	res, err := c.get(ctx, c.baseUrl+"/api/v1/analysis?project="+projectUuid+"&component="+componentUuid+"&vulnerability="+vulnerabilityUuid, c.authSource)
	if err != nil {
		if IsNotFound(err) {
			return trail, nil
		}
		return nil, fmt.Errorf("get analysis: %w", err)
	}
	if err := json.Unmarshal(res, &trail); err != nil {
		return nil, fmt.Errorf("unmarshal analysis: %w", err)
	}
	return trail, nil
}
