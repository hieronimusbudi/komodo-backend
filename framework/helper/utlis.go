package helper

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ResponseSuccess struct {
	Data interface{} `json:"data"`
}

func JwtFromFiberHeader(c *fiber.Ctx, header string, authScheme string) (string, error) {
	auth := c.Get(header)
	l := len(authScheme)
	if len(auth) > l+1 && strings.EqualFold(auth[:l], authScheme) {
		return auth[l+1:], nil
	}
	return "", errors.New("missing or malformed JWT")
}
