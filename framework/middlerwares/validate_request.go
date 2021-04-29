package middlerwares

import (
	"errors"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/config"
	"github.com/hieronimusbudi/komodo-backend/framework/helper"
)

func ValidateRequest(c *fiber.Ctx) error {
	// Get token from header
	token, restJwtErr := helper.JwtFromFiberHeader(c, fiber.HeaderAuthorization, "Bearer")
	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(restJwtErr.Error())
	}

	// Validate token
	tokenClaims, tokenErr := helper.ValidateToken(token, []byte(config.JWT_SECRET))
	log.Printf("ValidateRequest %v \n", tokenClaims)
	if tokenErr != nil {
		tokenErr := errors.New("unauthorized")
		return c.Status(http.StatusUnauthorized).JSON(tokenErr.Error())
	}

	c.Context().SetUserValue("tokenClaims", tokenClaims)
	return c.Next()
}
