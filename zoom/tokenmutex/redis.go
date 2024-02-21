package tokenmutex

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

const redisDefaultKey = "zoom_access_token"
const redisLockKey = "zoom_access_token_lock"

type Redis struct {
	client *redis.Client
	locker *redislock.Client
	lock   *redislock.Lock
	key    string
}

func NewRedis(client *redis.Client, key string) *Redis {
	if client == nil {
		panic("client is nil")
	}

	r := &Redis{
		client: client,
		locker: redislock.New(client),
		key:    key,
	}

	if len(r.key) == 0 {
		r.key = redisDefaultKey
	}

	return r
}

func (r *Redis) Lock(ctx context.Context) error {
	lock, err := r.locker.Obtain(ctx, redisLockKey, 30*time.Second, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(500*time.Millisecond), 60),
	})
	if err != nil {
		return fmt.Errorf("obtaining lock: %w", err)
	}

	r.lock = lock

	return nil
}

func (r *Redis) Unlock(ctx context.Context) error {
	if r.lock == nil {
		return nil
	}

	err := r.lock.Release(ctx)
	if err != nil {
		return fmt.Errorf("releasing lock: %w", err)
	}

	return nil
}

func (r *Redis) Get(context.Context) (string, error) {
	val, err := r.client.Get(context.Background(), r.key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenNotExist
		}

		return "", fmt.Errorf("getting key: %w", err)
	}

	return val, nil
}

func (r *Redis) Set(ctx context.Context, token string, expiresAt time.Time) error {
	err := r.client.SetEx(ctx, r.key, token, time.Second*time.Duration(expiresAt.Unix()-time.Now().Unix())).Err()
	if err != nil {
		return fmt.Errorf("setting key: %w", err)
	}

	return nil
}

func (r *Redis) Clear(ctx context.Context) error {
	_, err := r.client.Del(ctx, r.key).Result()
	if err != nil {
		return fmt.Errorf("deleting key: %w", err)
	}

	return nil
}
