package tokenmutex

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefault_Lock_Unlock(t *testing.T) {
	assert := assert.New(t)

	mutex := NewDefault()
	err := mutex.Lock(context.Background(), 1*time.Second)
	assert.NoError(err)

	err = mutex.Unlock(context.Background())
	assert.NoError(err)
}

func TestDefault_Get(t *testing.T) {
	assert := assert.New(t)

	mutex := NewDefault()
	mutex.token = "foo"
	mutex.expiresAt = time.Now().Add(time.Minute * 1)

	token, err := mutex.Get(context.Background())

	assert.Equal(mutex.token, token)
	assert.NoError(err)
}

func TestDefault_Get_ErrTokenNotExist(t *testing.T) {
	assert := assert.New(t)

	mutex := NewDefault()
	mutex.token = ""
	mutex.expiresAt = time.Now().Add(time.Minute * 1)

	token, err := mutex.Get(context.Background())

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenNotExist))
}

func TestDefault_Get_ErrTokenExpired(t *testing.T) {
	assert := assert.New(t)

	mutex := NewDefault()
	mutex.token = "foo"
	mutex.expiresAt = time.Now().Add(-time.Minute * 1)

	token, err := mutex.Get(context.Background())

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenExpired))
}

func TestDefault_Set_Clear(t *testing.T) {
	assert := assert.New(t)

	mutex := NewDefault()

	token := "foo"
	expiresAt := time.Now().Add(time.Minute * 1)

	err := mutex.Set(context.Background(), token, expiresAt)

	assert.Equal(token, mutex.token)
	assert.True(expiresAt.Equal(mutex.expiresAt))
	assert.NoError(err)

	err = mutex.Clear(context.Background())

	assert.Equal("", mutex.token)
	assert.Zero(mutex.expiresAt)
	assert.NoError(err)
}
