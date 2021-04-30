package buyercontroller_test

import (
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v2"
	buyercontroller "github.com/hieronimusbudi/komodo-backend/controllers/buyer_controller"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type TestSuite struct {
	suite.Suite
	mockBuyerUCase    *mocks.BuyerUseCase
	mockBuyer         entity.Buyer
	mockBuyerDTOReq   entity.BuyerDTORequest
	mockBuyerLoginReq entity.BuyerDTOLogin
	app               *fiber.App
}

// for each test
func (suite *TestSuite) SetupTest() {
	suite.mockBuyerUCase = new(mocks.BuyerUseCase)
	suite.app = fiber.New()

	suite.mockBuyer = entity.Buyer{
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		SendingAddress: "sending address",
	}

	suite.mockBuyerDTOReq = entity.BuyerDTORequest{
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "123456",
		SendingAddress: "sending address",
	}

	suite.mockBuyerLoginReq = entity.BuyerDTOLogin{
		Email:    "buyer1@mail.com",
		Password: "123456",
	}
}

func TestBuyerController(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestRegister() {
	suite.mockBuyerUCase.On("Register", mock.AnythingOfType("*entity.Buyer")).Return(nil).Once()

	j, err := json.Marshal(suite.mockBuyerDTOReq)
	suite.NoError(err)

	// setup fiber ctx
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
	ctx.Request().SetBody(j)
	defer suite.app.ReleaseCtx(ctx)

	handler := buyercontroller.NewBuyerController(suite.mockBuyerUCase)

	hErr := handler.Register(ctx)
	suite.NoError(hErr)
}

func (suite *TestSuite) TestLogin() {
	suite.mockBuyerUCase.On("Login", mock.AnythingOfType("*entity.Buyer")).Return(suite.mockBuyer, nil).Once()

	j, err := json.Marshal(suite.mockBuyerLoginReq)
	suite.NoError(err)

	// setup fiber ctx
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
	ctx.Request().SetBody(j)
	defer suite.app.ReleaseCtx(ctx)

	handler := buyercontroller.NewBuyerController(suite.mockBuyerUCase)

	hErr := handler.Login(ctx)
	suite.NoError(hErr)
}
