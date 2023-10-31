package client

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *client) GetFindings(ctx context.Context, projectUuid string) ([]*Finding, error) {
	res, err := c.get(ctx, c.baseUrl+"/api/v1/finding/project/"+projectUuid, c.authSource)
	if err != nil {
		if IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("get finding: %w", err)
	}

	var findings []*Finding
	if err = json.Unmarshal(res, &findings); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return findings, nil
}

func (c *client) TriggerAnalysis(ctx context.Context, projectUuid string) error {
	_, err := c.post(ctx, c.baseUrl+"/api/v1/finding/project/"+projectUuid+"/analyze", c.authSource, nil)
	if err != nil {
		if IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("trigger analysis: %w", err)
	}
	return nil
}
