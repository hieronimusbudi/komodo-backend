package sellercontroller_test

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	sellercontroller "github.com/hieronimusbudi/komodo-backend/controllers/seller_controller"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type TestSuite struct {
	suite.Suite
	mockSellerUCase    *mocks.SellerUseCase
	mockSeller         entity.Seller
	mockSellerDTOReq   entity.SellerDTORequest
	mockSellerLoginReq entity.SellerDTOLogin
	app                *fiber.App
	validate           *validator.Validate
}

// for each test
func (suite *TestSuite) SetupTest() {
	suite.mockSellerUCase = new(mocks.SellerUseCase)
	suite.app = fiber.New()
	suite.validate = validator.New()

	suite.mockSeller = entity.Seller{
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		PickUpAddress: "sending address",
	}

	suite.mockSellerDTOReq = entity.SellerDTORequest{
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "123456",
		PickUpAddress: "sending address",
	}

	suite.mockSellerLoginReq = entity.SellerDTOLogin{
		Email:    "seller1@mail.com",
		Password: "123456",
	}
}

func TestSellerController(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestRegister() {
	suite.mockSellerUCase.On("Register", mock.AnythingOfType("*entity.Seller")).Return(nil).Once()

	j, err := json.Marshal(suite.mockSellerDTOReq)
	suite.NoError(err)

	// setup fiber ctx
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
	ctx.Request().SetBody(j)
	defer suite.app.ReleaseCtx(ctx)

	handler := sellercontroller.NewSellerController(suite.mockSellerUCase, suite.validate)

	hErr := handler.Register(ctx)
	suite.NoError(hErr)
}

func (suite *TestSuite) TestLogin() {
	suite.mockSellerUCase.On("Login", mock.AnythingOfType("*entity.Seller")).Return(suite.mockSeller, nil).Once()

	j, err := json.Marshal(suite.mockSellerLoginReq)
	suite.NoError(err)

	// setup fiber ctx
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
	ctx.Request().SetBody(j)
	defer suite.app.ReleaseCtx(ctx)

	handler := sellercontroller.NewSellerController(suite.mockSellerUCase, suite.validate)

	hErr := handler.Login(ctx)
	suite.NoError(hErr)
}
