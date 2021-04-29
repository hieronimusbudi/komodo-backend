package middlerwares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/framework/helper"
)

func userTypeChecker(c *fiber.Ctx, as helper.UserTypeEnum) error {
	// take token from user value context
	tokenClaims, ok := c.Context().UserValue("tokenClaims").(jwt.MapClaims)
	if !ok {
		return errors.New("token claims not exists")
	}

	// get name & user type
	userName := string(tokenClaims["name"].(string))
	userType := helper.UserTypeEnum(tokenClaims["type"].(float64))

	if userType != as {
		return fmt.Errorf("%s is not authorized", userName)
	}

	return nil
}

func BuyerTypeChecker(c *fiber.Ctx) error {
	err := userTypeChecker(c, helper.BUYER_TYPE)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(err.Error())
	}

	return c.Next()
}

func SellerTypeChecker(c *fiber.Ctx) error {
	err := userTypeChecker(c, helper.SELLER_TYPE)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(err.Error())
	}

	return c.Next()
}
