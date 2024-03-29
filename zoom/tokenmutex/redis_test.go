package tokenmutex

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedis_Lock_Unlock(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	mutex := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), "")

	err = mutex.Unlock(context.Background())
	assert.NoError(err)

	err = mutex.Lock(context.Background())
	assert.NoError(err)

	err = mutex.Unlock(context.Background())
	assert.NoError(err)

	err = mutex.Unlock(context.Background())
	assert.ErrorIs(err, redislock.ErrLockNotHeld)
}

func TestRedis_Get(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	expectedToken := "foo"

	s.Set(redisDefaultKey, expectedToken)
	s.SetTTL(redisDefaultKey, time.Minute*1)

	mutex := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), "")

	token, err := mutex.Get(context.Background())

	assert.Equal(expectedToken, token)
	assert.NoError(err)
}

func TestRedis_Get_ErrTokenNotExist(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	mutex := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), "")

	token, err := mutex.Get(context.Background())

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenNotExist))
}

func TestRedis_Set_Clear(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	mutex := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), "")

	expectedToken := "foo"
	err = mutex.Set(context.Background(), expectedToken, time.Now().Add(time.Minute*1))

	assert.NoError(err)

	token, _ := s.Get(redisDefaultKey)
	ttl := s.TTL(redisDefaultKey)

	assert.Equal(expectedToken, token)
	assert.True(time.Now().Add(time.Second * ttl).After(time.Now()))

	err = mutex.Clear(context.Background())
	assert.NoError(err)
	token, err = s.Get(redisDefaultKey)
	assert.Equal("", token)
	assert.Error(err)
	assert.False(s.Exists(redisDefaultKey))
}
