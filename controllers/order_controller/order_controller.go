package ordercontroller

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
	"github.com/shopspring/decimal"
)

type OrderController interface {
	Store(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	AcceptOrder(c *fiber.Ctx) error
}

type orderController struct {
	orderUsecase entity.OrderUseCase
}

// NewOrderController will create a object with OrderController interface representation
func NewOrderController(orderUsecase entity.OrderUseCase) OrderController {
	return &orderController{
		orderUsecase: orderUsecase,
	}
}

func (oc *orderController) Store(c *fiber.Ctx) error {
	// parse order from request body
	oDTOReq := new(entity.OrderDTORequest)
	if err := c.BodyParser(oDTOReq); err != nil {
		rErr := resterrors.NewRestError("unprocessable entity", http.StatusUnprocessableEntity, err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// transform OrderDTORequest to Order
	dTP := decimal.NewFromFloat(oDTOReq.TotalPrice)
	order := entity.Order{
		Buyer:                      entity.Buyer{ID: oDTOReq.BuyerID},
		Seller:                     entity.Seller{ID: oDTOReq.SellerID},
		DeliverySourceAddress:      oDTOReq.DeliverySourceAddress,
		DeliveryDestinationAddress: oDTOReq.DeliveryDestinationAddress,
		TotalQuantity:              oDTOReq.TotalQuantity,
		TotalPrice:                 dTP,
		Status:                     entity.PENDING,
	}

	for _, od := range oDTOReq.Items {
		order.Items = append(order.Items, entity.OrderDetail{
			Product:  entity.Product{ID: od.ProductId},
			Quantity: od.Quantity,
		})
	}

	// store Order
	err := oc.orderUsecase.Store(&order)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	// transform Order to OrderDTOResponse
	fTP, _ := order.TotalPrice.Float64()
	res := entity.OrderDTOResponse{
		ID:                         order.ID,
		Buyer:                      entity.BuyerDTOResponse{ID: order.Buyer.ID},
		Seller:                     entity.SellerDTOResponse{ID: order.Seller.ID},
		DeliverySourceAddress:      order.DeliverySourceAddress,
		DeliveryDestinationAddress: order.DeliveryDestinationAddress,
		TotalQuantity:              order.TotalQuantity,
		TotalPrice:                 fTP,
		Status:                     order.Status,
		OrderDate:                  order.OrderDate,
	}

	for _, od := range order.Items {
		res.Items = append(res.Items, entity.OrderDetailDTOResponse{
			Product:  entity.ProductDTOResponse{ID: od.Product.ID},
			Quantity: od.Quantity,
		})
	}

	return c.Status(http.StatusCreated).JSON(helpers.SuccessResponse{
		Data: res,
	})
}

func (oc *orderController) GetByUserID(c *fiber.Ctx) error {
	// take token from user value context
	tokenClaims, ok := c.Context().UserValue("tokenClaims").(jwt.MapClaims)
	if !ok {
		rErr := resterrors.NewUnauthorizedError("token claims not exists")
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// get id & user type
	userID := int64(tokenClaims["id"].(float64))
	userType := helpers.UserTypeEnum(tokenClaims["type"].(float64))

	// get orders by buyer/seller id
	orders, err := oc.orderUsecase.GetByUserID(userID, userType)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.ErrorResponse())
	}

	// transform Order to OrderDTOResponse
	var orderRes entity.OrderDTOResponse
	var res []entity.OrderDTOResponse
	for _, order := range orders {
		fTP, _ := order.TotalPrice.Float64()
		orderRes = entity.OrderDTOResponse{
			ID: order.ID,
			Buyer: entity.BuyerDTOResponse{
				ID: order.Buyer.ID,
			},
			Seller: entity.SellerDTOResponse{
				ID: order.Seller.ID,
			},
			DeliverySourceAddress:      order.DeliverySourceAddress,
			DeliveryDestinationAddress: order.DeliveryDestinationAddress,
			TotalQuantity:              order.TotalQuantity,
			TotalPrice:                 fTP,
			Status:                     order.Status,
			OrderDate:                  order.OrderDate,
		}

		odRow := entity.OrderDetailDTOResponse{}
		orderRes.Items = []entity.OrderDetailDTOResponse{}
		for _, od := range order.Items {
			odRow.Product.ID = od.Product.ID
			odRow.Quantity = od.Quantity

			fP, _ := od.Product.Price.Float64()
			odRow.Price = fP

			orderRes.Items = append(orderRes.Items, odRow)
		}

		res = append(res, orderRes)
	}

	return c.Status(http.StatusOK).JSON(helpers.SuccessResponse{
		Data: res,
	})
}

func (oc *orderController) AcceptOrder(c *fiber.Ctx) error {
	// extract params
	orderId, idErr := c.ParamsInt("id")
	if idErr != nil {
		rErr := resterrors.NewBadRequestError(idErr.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// accept order
	order := new(entity.Order)
	order.ID = int64(orderId)
	uOrderRes, err := oc.orderUsecase.AcceptOrder(order)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	// transform Order to OrderDTOResponse
	fTP, _ := uOrderRes.TotalPrice.Float64()
	res := entity.OrderDTOResponse{
		ID:                         uOrderRes.ID,
		Buyer:                      entity.BuyerDTOResponse{ID: uOrderRes.Buyer.ID},
		Seller:                     entity.SellerDTOResponse{ID: uOrderRes.Seller.ID},
		DeliverySourceAddress:      uOrderRes.DeliverySourceAddress,
		DeliveryDestinationAddress: uOrderRes.DeliveryDestinationAddress,
		TotalQuantity:              uOrderRes.TotalQuantity,
		TotalPrice:                 fTP,
		Status:                     uOrderRes.Status,
		OrderDate:                  uOrderRes.OrderDate,
	}

	for _, od := range uOrderRes.Items {
		res.Items = append(res.Items, entity.OrderDetailDTOResponse{
			Product:  entity.ProductDTOResponse{ID: od.Product.ID},
			Quantity: od.Quantity,
		})
	}

	return c.Status(http.StatusCreated).JSON(helpers.SuccessResponse{
		Data: res,
	})
}
