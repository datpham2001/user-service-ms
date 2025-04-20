package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type AuthEventType string

const (
	UserResetPasswordEvent AuthEventType = "user_reset_password"
)

type UserEvent struct {
	UserID    string        `json:"user_id"`
	EventType AuthEventType `json:"event_type"`
	Timestamp time.Time     `json:"timestamp"`
	Data      any           `json:"data"`
}

func (s *AuthService) publishUserEvent(ctx context.Context, event *UserEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return s.rabbitMQ.Publish(ctx, "", eventJSON)
}
