package zoom

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestMeetingSDKJWT(t *testing.T) {
	assert := assert.New(t)

	sdkKey := "foobar"
	sdkSecret := "bazcat"
	meetingNumber := int64(123)
	role := 1
	expiration := 1 * time.Hour

	tokenStr, err := MeetingSDKJWT(sdkKey, sdkSecret, meetingNumber, role, expiration)
	assert.NoError(err)

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(sdkSecret), nil
	}, jwt.WithJSONNumber())
	assert.NoError(err)

	assert.Equal(sdkKey, token.Claims.(jwt.MapClaims)["appKey"])
	assert.Equal(sdkKey, token.Claims.(jwt.MapClaims)["sdkKey"])

	mn, err := token.Claims.(jwt.MapClaims)["mn"].(json.Number).Int64()
	assert.NoError(err)

	claimRole, err := token.Claims.(jwt.MapClaims)["role"].(json.Number).Int64()
	assert.NoError(err)

	iat, err := token.Claims.(jwt.MapClaims)["iat"].(json.Number).Int64()
	assert.NoError(err)

	exp, err := token.Claims.(jwt.MapClaims)["exp"].(json.Number).Int64()
	assert.NoError(err)

	tokenExp, err := token.Claims.(jwt.MapClaims)["tokenExp"].(json.Number).Int64()
	assert.NoError(err)

	assert.Equal(meetingNumber, mn)
	assert.Equal(int64(role), claimRole)
	assert.True(iat > 0)
	assert.True(exp > iat)
	assert.True(tokenExp > iat)
}
