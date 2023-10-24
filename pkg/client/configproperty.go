package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type ConfigProperty struct {
	GroupName     string `json:"groupName"`
	PropertyName  string `json:"propertyName"`
	PropertyValue string `json:"propertyValue"`
	PropertyType  string `json:"propertyType"`
	Description   string `json:"description"`
}

func (c *client) ConfigPropertyAggregate(ctx context.Context, properties []ConfigProperty) ([]ConfigProperty, error) {
	body, err := json.Marshal(properties)
	if err != nil {
		return nil, fmt.Errorf("marshalling request body: %w", err)
	}

	res, err := c.post(ctx, c.baseUrl+"/api/v1/configProperty/aggregate", c.authSource, body)
	if err != nil {
		return nil, fmt.Errorf("aggregating config properties: %w", err)
	}

	var props []ConfigProperty
	if err = json.Unmarshal(res, &props); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return props, nil
}

func (c *client) GetConfigProperties(ctx context.Context) ([]ConfigProperty, error) {
	res, err := c.get(ctx, c.baseUrl+"/api/v1/configProperty", c.authSource)
	if err != nil {
		return nil, fmt.Errorf("get config properties: %w", err)
	}

	var properties []ConfigProperty
	if err = json.Unmarshal(res, &properties); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return properties, nil
}
