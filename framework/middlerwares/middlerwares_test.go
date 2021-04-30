package middlerwares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/config"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	"github.com/hieronimusbudi/komodo-backend/framework/middlerwares"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	app                  *fiber.App
	sellerToken          string
	buyerToken           string
	mockJWTPayloadSeller helpers.UserJWTPayload
	mockJWTPayloadBuyer  helpers.UserJWTPayload
}

func (suite *TestSuite) SetupTest() {
	var err error
	suite.app = fiber.New()

	suite.mockJWTPayloadSeller = helpers.UserJWTPayload{
		ID:    1,
		Email: "buyer1@mail.com",
		Name:  "buyer",
		Type:  helpers.SELLER_TYPE,
	}

	suite.sellerToken, err = helpers.GenerateToken(&suite.mockJWTPayloadSeller, []byte(config.JWT_SECRET))
	suite.NoError(err)

	suite.mockJWTPayloadSeller = helpers.UserJWTPayload{
		ID:    1,
		Email: "buyer1@mail.com",
		Name:  "buyer",
		Type:  helpers.SELLER_TYPE,
	}

	suite.buyerToken, err = helpers.GenerateToken(&suite.mockJWTPayloadBuyer, []byte(config.JWT_SECRET))
	suite.NoError(err)
}

func TestMiddlewares(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestValidateRequest() {
	authToken := fmt.Sprintf("Bearer %s", suite.buyerToken)
	suite.app.Get("/test",
		func(c *fiber.Ctx) error {
			// setup header
			c.Request().Header.Add(fiber.HeaderAuthorization, authToken)
			return c.Next()
		},
		middlerwares.ValidateRequest,
		func(c *fiber.Ctx) error {
			_, ok := c.Context().UserValue("tokenClaims").(jwt.MapClaims)

			// token claims exists, set by middlerwares.ValidateRequest
			suite.Equal(true, ok)
			return nil
		},
	)

	resp, err := suite.app.Test(httptest.NewRequest(http.MethodGet, "http://example.com/test", nil))
	suite.NoError(err)
	suite.Equal(fiber.StatusOK, resp.StatusCode)
}

func (suite *TestSuite) TestValidateRequestError() {
	invalidToken := fmt.Sprintf("Bearer x%sxinvalid", suite.buyerToken)
	suite.app.Get("/test",
		func(c *fiber.Ctx) error {
			// setup header
			c.Request().Header.Add(fiber.HeaderAuthorization, invalidToken)
			return c.Next()
		},
		middlerwares.ValidateRequest,
	)

	resp, err := suite.app.Test(httptest.NewRequest(http.MethodGet, "http://example.com/test", nil))
	suite.NoError(err)
	suite.Equal(fiber.StatusUnauthorized, resp.StatusCode)
}

func (suite *TestSuite) TestValidateRequestAndBuyerTypeChecker() {
	authToken := fmt.Sprintf("Bearer %s", suite.buyerToken)
	suite.app.Get("/test",
		func(c *fiber.Ctx) error {
			// setup header
			c.Request().Header.Add(fiber.HeaderAuthorization, authToken)
			return c.Next()
		},
		middlerwares.ValidateRequest,
		middlerwares.BuyerTypeChecker,
		func(c *fiber.Ctx) error {
			_, ok := c.Context().UserValue("tokenClaims").(jwt.MapClaims)

			// token claims exists, set by middlerwares.ValidateRequest
			suite.Equal(true, ok)
			return nil
		},
	)

	resp, err := suite.app.Test(httptest.NewRequest(http.MethodGet, "http://example.com/test", nil))
	suite.NoError(err)
	suite.Equal(fiber.StatusOK, resp.StatusCode)
}

func (suite *TestSuite) TestValidateRequestAndBuyerTypeCheckerError() {
	// pass seller token not buyer token
	authToken := fmt.Sprintf("Bearer %s", suite.sellerToken)
	suite.app.Get("/test",
		func(c *fiber.Ctx) error {
			// setup header
			c.Request().Header.Add(fiber.HeaderAuthorization, authToken)
			return c.Next()
		},
		middlerwares.ValidateRequest,
		middlerwares.BuyerTypeChecker,
	)

	resp, err := suite.app.Test(httptest.NewRequest(http.MethodGet, "http://example.com/test", nil))
	suite.NoError(err)
	suite.Equal(fiber.StatusUnauthorized, resp.StatusCode)
}

func (suite *TestSuite) TestValidateRequestAndSellerTypeChecker() {
	authToken := fmt.Sprintf("Bearer %s", suite.sellerToken)
	suite.app.Get("/test",
		func(c *fiber.Ctx) error {
			// setup header
			c.Request().Header.Add(fiber.HeaderAuthorization, authToken)
			return c.Next()
		},
		middlerwares.ValidateRequest,
		middlerwares.SellerTypeChecker,
		func(c *fiber.Ctx) error {
			_, ok := c.Context().UserValue("tokenClaims").(jwt.MapClaims)

			// token claims exists, set by middlerwares.ValidateRequest
			suite.Equal(true, ok)
			return nil
		},
	)

	resp, err := suite.app.Test(httptest.NewRequest(http.MethodGet, "http://example.com/test", nil))
	suite.NoError(err)
	suite.Equal(fiber.StatusOK, resp.StatusCode)
}

func (suite *TestSuite) TestValidateRequestAndSellerTypeCheckerError() {
	// pass buyer token not seller token
	authToken := fmt.Sprintf("Bearer %s", suite.buyerToken)
	suite.app.Get("/test",
		func(c *fiber.Ctx) error {
			// setup header
			c.Request().Header.Add(fiber.HeaderAuthorization, authToken)
			return c.Next()
		},
		middlerwares.ValidateRequest,
		middlerwares.SellerTypeChecker,
	)

	resp, err := suite.app.Test(httptest.NewRequest(http.MethodGet, "http://example.com/test", nil))
	suite.NoError(err)
	suite.Equal(fiber.StatusUnauthorized, resp.StatusCode)
}
