package tokenmutex

import "errors"

var ErrTokenNotExist = errors.New("token does not exist")
var ErrTokenExpired = errors.New("token expired")
