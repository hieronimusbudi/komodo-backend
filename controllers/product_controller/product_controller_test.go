package productcontroller_test

import (
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v2"
	productcontroller "github.com/hieronimusbudi/komodo-backend/controllers/product_controller"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type TestSuite struct {
	suite.Suite
	mockProductUCase  *mocks.ProductUseCase
	mockProduct       entity.Product
	mockProductDTOReq entity.ProductDTORequest
	mockSeller        entity.Seller
	app               *fiber.App
}

// for each test
func (suite *TestSuite) SetupTest() {
	suite.mockProductUCase = new(mocks.ProductUseCase)
	suite.app = fiber.New()

	suite.mockSeller = entity.Seller{
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		PickUpAddress: "sending address",
	}

	suite.mockProduct = entity.Product{
		ID:          1,
		Name:        "product1",
		Description: "desc",
		Price:       decimal.NewFromFloat(181818.11),
		Seller:      suite.mockSeller,
	}

	suite.mockProductDTOReq = entity.ProductDTORequest{
		Name:        "product1",
		Description: "desc",
		Price:       181818.11,
		SellerID:    1,
	}
}

func TestProductController(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestStore() {
	suite.mockProductUCase.On("Store", mock.AnythingOfType("*entity.Product")).Return(nil).Once()

	j, err := json.Marshal(suite.mockProductDTOReq)
	suite.NoError(err)

	// setup fiber ctx
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
	ctx.Request().SetBody(j)
	defer suite.app.ReleaseCtx(ctx)

	handler := productcontroller.NewProductController(suite.mockProductUCase)

	hErr := handler.Store(ctx)
	suite.NoError(hErr)
}

func (suite *TestSuite) TestGetAll() {
	suite.mockProductUCase.On("GetAll").Return([]entity.Product{suite.mockProduct}, nil).Once()

	// setup fiber ctx
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
	defer suite.app.ReleaseCtx(ctx)

	handler := productcontroller.NewProductController(suite.mockProductUCase)

	hErr := handler.GetAll(ctx)
	suite.NoError(hErr)
}
