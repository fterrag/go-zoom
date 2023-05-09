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

	cacher := NewDefault()
	err := cacher.Lock(context.Background(), 1*time.Second)
	assert.NoError(err)

	err = cacher.Unlock(context.Background())
	assert.NoError(err)
}

func TestDefault_Get(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = "foo"
	cacher.expiresAt = time.Now().Add(time.Minute * 1)

	token, err := cacher.Get(context.Background())

	assert.Equal(cacher.token, token)
	assert.NoError(err)
}

func TestDefault_Get_ErrTokenNotExist(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = ""
	cacher.expiresAt = time.Now().Add(time.Minute * 1)

	token, err := cacher.Get(context.Background())

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenNotExist))
}

func TestDefault_Get_ErrTokenExpired(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = "foo"
	cacher.expiresAt = time.Now().Add(-time.Minute * 1)

	token, err := cacher.Get(context.Background())

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenExpired))
}

func TestDefault_Set_Clear(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()

	token := "foo"
	expiresAt := time.Now().Add(time.Minute * 1)

	err := cacher.Set(context.Background(), token, expiresAt)

	assert.Equal(token, cacher.token)
	assert.True(expiresAt.Equal(cacher.expiresAt))
	assert.NoError(err)

	err = cacher.Clear(context.Background())

	assert.Equal("", cacher.token)
	assert.Zero(cacher.expiresAt)
	assert.NoError(err)
}
