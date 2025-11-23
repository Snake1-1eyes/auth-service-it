package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

func (r *SessionRepository) CreateSession(ctx context.Context, sessionID string, userID int64, duration time.Duration) error {
	err := r.client.Set(ctx, sessionID, userID, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (int64, error) {
	val, err := r.client.Get(ctx, sessionID).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get session: %w", err)
	}
	return val, nil
}
