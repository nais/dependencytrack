package dependencytrack

import (
	"context"
	"net/http"

	"github.com/nais/dependencytrack/pkg/dependencytrack/client"
)

type ConfigProperty struct {
	GroupName     *string `json:"groupName,omitempty"`
	PropertyName  *string `json:"propertyName,omitempty"`
	PropertyValue *string `json:"propertyValue,omitempty"`
	PropertyType  string  `json:"propertyType"`
	Description   *string `json:"description,omitempty"`
}

func (c *dependencyTrackClient) ConfigPropertyAggregate(ctx context.Context, property ConfigProperty) error {
	return c.withAuthContext(ctx, func(tokenCtx context.Context) error {
		_, resp, err := c.client.ConfigPropertyAPI.UpdateConfigProperty(tokenCtx).ConfigProperty(client.ConfigProperty{
			GroupName:     property.GroupName,
			PropertyName:  property.PropertyName,
			PropertyValue: property.PropertyValue,
			PropertyType:  property.PropertyType,
			Description:   property.Description,
		}).Execute()
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return nil // Not found, return empty property
			}
			return convertError(err, "ConfigPropertyAggregate", resp)
		}
		return nil
	})
}

func (c *dependencyTrackClient) GetConfigProperties(ctx context.Context) ([]ConfigProperty, error) {
	return withAuthContextValue(c, ctx, func(tokenCtx context.Context) ([]ConfigProperty, error) {
		res, resp, err := c.client.ConfigPropertyAPI.GetConfigProperties(tokenCtx).Execute()
		if err != nil {
			return nil, convertError(err, "GetConfigProperties", resp)
		}

		configProperties := make([]ConfigProperty, len(res))
		for i, property := range res {
			configProperties[i] = parseConfigProperty(property)
		}
		return configProperties, nil
	})
}

func parseConfigProperty(res client.ConfigProperty) ConfigProperty {
	return ConfigProperty{
		GroupName:     res.GroupName,
		PropertyName:  res.PropertyName,
		PropertyValue: res.PropertyValue,
		PropertyType:  res.PropertyType,
		Description:   res.Description,
	}
}
