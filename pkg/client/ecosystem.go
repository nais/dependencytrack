package client

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *client) GetEcosystems(ctx context.Context) ([]string, error) {
	var ecosystems []string
	res, err := c.get(ctx, c.baseUrl+"/api/v1/integration/osv/ecosystem/inactive", c.authSource)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, &ecosystems); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return ecosystems, nil
}
