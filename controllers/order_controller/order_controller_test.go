package ordercontroller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	ordercontroller "github.com/hieronimusbudi/komodo-backend/controllers/order_controller"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type TestSuite struct {
	suite.Suite
	mockOrderUCase            *mocks.OrderUseCase
	mockOrder                 entity.Order
	mockOrderDTOReq           entity.OrderDTORequest
	mockOrderDetail           entity.OrderDetail
	mockOrderDetailDTORequest entity.OrderDetailDTORequest
	mockBuyer                 entity.Buyer
	mockSeller                entity.Seller
	mockProduct               entity.Product
	app                       *fiber.App
	validate                  *validator.Validate
}

// for each test
func (suite *TestSuite) SetupTest() {
	suite.mockOrderUCase = new(mocks.OrderUseCase)
	suite.app = fiber.New()
	suite.validate = validator.New()

	suite.mockBuyer = entity.Buyer{
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		SendingAddress: "sending address",
	}

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

	suite.mockOrderDetail = entity.OrderDetail{
		ID:       1,
		Product:  suite.mockProduct,
		Quantity: 10,
	}

	suite.mockOrderDetailDTORequest = entity.OrderDetailDTORequest{
		ProductId: suite.mockProduct.ID,
		Quantity:  10,
	}

	t, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	suite.mockOrder = entity.Order{
		ID:                         1,
		Buyer:                      suite.mockBuyer,
		Seller:                     suite.mockSeller,
		DeliverySourceAddress:      "pickup address",
		DeliveryDestinationAddress: "sending address",
		TotalQuantity:              10,
		TotalPrice:                 decimal.NewFromFloat(181818.11),
		Status:                     entity.PENDING,
		OrderDate:                  t,
		Items: []entity.OrderDetail{
			suite.mockOrderDetail,
		},
	}

	suite.mockOrderDTOReq = entity.OrderDTORequest{
		BuyerID:                    suite.mockBuyer.ID,
		SellerID:                   suite.mockSeller.ID,
		DeliverySourceAddress:      "pickup address",
		DeliveryDestinationAddress: "sending address",
		Items: []entity.OrderDetailDTORequest{
			suite.mockOrderDetailDTORequest,
		},
	}
}

func TestOrderController(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestStore() {
	suite.mockOrderUCase.On("Store", mock.AnythingOfType("*entity.Order")).Return(nil).Once()

	j, err := json.Marshal(suite.mockOrderDTOReq)
	suite.NoError(err)

	// setup fiber ctx
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
	ctx.Request().SetBody(j)
	defer suite.app.ReleaseCtx(ctx)

	handler := ordercontroller.NewOrderController(suite.mockOrderUCase, suite.validate)

	hErr := handler.Store(ctx)
	suite.NoError(hErr)
}

func (suite *TestSuite) TestStoreError() {
	suite.mockOrderDTOReq.BuyerID = 0
	suite.mockOrderDTOReq.SellerID = 0

	j, err := json.Marshal(suite.mockOrderDTOReq)
	suite.NoError(err)

	handler := ordercontroller.NewOrderController(suite.mockOrderUCase, suite.validate)
	suite.app.Post("/orders",
		func(c *fiber.Ctx) error {
			c.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
			return c.Next()
		},
		handler.Store,
	)

	resp, err := suite.app.Test(
		httptest.NewRequest(
			http.MethodPost,
			"http://example.com/orders",
			strings.NewReader(string(j)),
		))
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *TestSuite) TestGetByUserID() {
	suite.mockOrderUCase.On("GetByUserID", suite.mockOrder.Buyer.ID, helpers.BUYER_TYPE).Return([]entity.Order{suite.mockOrder}, nil).Once()

	// setup fiber ctx
	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().Header.SetContentType(fiber.MIMEApplicationJSON)
	defer suite.app.ReleaseCtx(ctx)

	handler := ordercontroller.NewOrderController(suite.mockOrderUCase, suite.validate)

	hErr := handler.GetByUserID(ctx)
	suite.NoError(hErr)
}

func (suite *TestSuite) TestAcceptOrder() {
	suite.mockOrderUCase.On("AcceptOrder", mock.AnythingOfType("*entity.Order")).Return(suite.mockOrder, nil).Once()

	j, err := json.Marshal(suite.mockOrderDTOReq)
	suite.NoError(err)

	handler := ordercontroller.NewOrderController(suite.mockOrderUCase, suite.validate)
	suite.app.Put("/orders/:id/accept", func(c *fiber.Ctx) error {
		hErr := handler.AcceptOrder(c)
		suite.NoError(hErr)

		return nil
	})

	_, err = suite.app.Test(
		httptest.NewRequest(
			http.MethodPut,
			// embed id in url
			fmt.Sprintf("/orders/%d/accept",
				suite.mockOrder.ID),
			strings.NewReader(string(j))))
	suite.NoError(err)
}
