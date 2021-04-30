package helpers

import (
	"errors"
	"strings"
	"time"
)

func JwtFromHeader(auth string, authScheme string) (string, error) {
	l := len(authScheme)
	if len(auth) > l+1 && strings.EqualFold(auth[:l], authScheme) {
		return auth[l+1:], nil
	}
	return "", errors.New("missing or malformed JWT")
}

func GetTimeFromUint8(t []uint8) (time.Time, error) {
	pt, err := time.Parse("2006-01-02 15:04:05", string(t))
	return pt, err
}

func GetStringTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetTimeNow() (time.Time, error) {
	tn, err := time.Parse("2006-01-02 15:04:05", GetStringTimeNow())
	return tn, err
}
