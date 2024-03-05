package client

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *client) GetNotifications(ctx context.Context) ([]NotificationRule, error) {
	var notificationRules []NotificationRule
	res, err := c.get(ctx, c.baseUrl+"/api/v1/notification/rule", c.authSource)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, &notificationRules); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return notificationRules, nil
}

func (c *client) CreateNotification(ctx context.Context, name string, scope NotificationScope, level NotificationLevel) (*NotificationRule, error) {
	notificationRule := &NotificationRule{
		Name:  name,
		Scope: scope,
		Level: level,
	}

	body, err := json.Marshal(notificationRule)
	if err != nil {
		return nil, fmt.Errorf("marshalling request body: %w", err)
	}
	rule, err := c.put(ctx, c.baseUrl+"/api/v1/notification/rule", c.authSource, body)
	if err != nil {
		return nil, fmt.Errorf("creating notification rule: %w", err)
	}
	var notify NotificationRule
	if err = json.Unmarshal(rule, &notify); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return &notify, err
}

func (c *client) DeleteNotification(ctx context.Context, notificationRule *NotificationRule) error {
	body, err := json.Marshal(notificationRule)
	if err != nil {
		return fmt.Errorf("marshalling request body: %w", err)
	}
	_, err = c.delete(ctx, c.baseUrl+"/api/v1/notification/rule", c.authSource, body)
	if err != nil {
		return fmt.Errorf("deleting notification rule: %w", err)
	}
	return err
}

func (c *client) UpdateNotification(ctx context.Context, notificationRule *NotificationRule) (*NotificationRule, error) {
	body, err := json.Marshal(notificationRule)
	if err != nil {
		return nil, fmt.Errorf("marshalling request body: %w", err)
	}

	r, err := c.post(ctx, c.baseUrl+"/api/v1/notification/rule", c.authSource, body)
	if err != nil {
		return nil, fmt.Errorf("updating notification rule: %w", err)
	}

	var notify NotificationRule
	if err = json.Unmarshal(r, &notify); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return &notify, err
}

func (c *client) AddProjectToNotification(ctx context.Context, notificationUuid, projectUuid string) (*NotificationRule, error) {
	update, err := c.post(ctx, c.baseUrl+"/api/v1/notification/rule/"+notificationUuid+"/project/"+projectUuid, c.authSource, nil)
	if err != nil {
		return nil, fmt.Errorf("adding project to notification: %w", err)
	}
	var notify NotificationRule
	if err = json.Unmarshal(update, &notify); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return &notify, err
}

func (c *client) GetPublishers(ctx context.Context) ([]Publisher, error) {
	var publishers []Publisher
	res, err := c.get(ctx, c.baseUrl+"/api/v1/notification/publisher", c.authSource)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, &publishers); err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}
	return publishers, nil
}
