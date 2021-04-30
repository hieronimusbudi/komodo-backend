package middlerwares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/config"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
)

func ValidateRequest(c *fiber.Ctx) error {
	// Get token from header
	auth := c.Get(fiber.HeaderAuthorization)
	token, err := helpers.JwtFromHeader(auth, "Bearer")
	if err != nil {
		rErr := resterrors.NewUnauthorizedError(err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// Validate token
	tokenClaims, err := helpers.ValidateToken(token, []byte(config.JWT_SECRET))
	if err != nil {
		rErr := resterrors.NewUnauthorizedError(err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	c.Context().SetUserValue("tokenClaims", tokenClaims)
	return c.Next()
}
