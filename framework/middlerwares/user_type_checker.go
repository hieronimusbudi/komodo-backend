package middlerwares

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
)

func userTypeChecker(c *fiber.Ctx, as helpers.UserTypeEnum) resterrors.RestErr {
	// take token from user value context
	tokenClaims, ok := c.Context().UserValue("tokenClaims").(jwt.MapClaims)
	if !ok {
		return resterrors.NewUnauthorizedError("token claims not exists")
	}

	// get name & user type
	userName := string(tokenClaims["name"].(string))
	userType := helpers.UserTypeEnum(tokenClaims["type"].(float64))

	if userType != as {
		return resterrors.NewUnauthorizedError(fmt.Sprintf("%s is not authorized", userName))
	}

	return nil
}

func BuyerTypeChecker(c *fiber.Ctx) error {
	err := userTypeChecker(c, helpers.BUYER_TYPE)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	return c.Next()
}

func SellerTypeChecker(c *fiber.Ctx) error {
	err := userTypeChecker(c, helpers.SELLER_TYPE)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	return c.Next()
}
