package tokenmutex

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/eleanorhealth/go-common/pkg/errs"
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

func (r *Redis) Lock(ctx context.Context, d time.Duration) error {
	lock, err := r.locker.Obtain(ctx, redisLockKey, d, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(500*time.Millisecond), 6),
	})
	if err != nil {
		return errs.Wrap(err, "obtaining lock")
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
		return errs.Wrap(err, "releasing lock")
	}

	return nil
}

func (r *Redis) Get(context.Context) (string, error) {
	val, err := r.client.Get(context.Background(), r.key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenNotExist
		}

		return "", errs.Wrap(err, "getting key")
	}

	return val, nil
}

func (r *Redis) Set(ctx context.Context, token string, expiresAt time.Time) error {
	err := r.client.SetEx(ctx, r.key, token, time.Second*time.Duration(expiresAt.Unix()-time.Now().Unix())).Err()
	if err != nil {
		return errs.Wrap(err, "setting key")
	}

	return nil
}

func (r *Redis) Clear(ctx context.Context) error {
	_, err := r.client.Del(ctx, r.key).Result()
	if err != nil {
		return errs.Wrap(err, "deleting key")
	}

	return nil
}
