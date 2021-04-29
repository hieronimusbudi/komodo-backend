package ordercontroller

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helper"
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

func NewOrderController(orderUsecase entity.OrderUseCase) OrderController {
	return &orderController{
		orderUsecase: orderUsecase,
	}
}

func (oc *orderController) Store(c *fiber.Ctx) error {
	// parse order from request body
	oDTOReq := new(entity.OrderDTORequest)
	if err := c.BodyParser(oDTOReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
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
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// transform Order to OrderDTOResponse
	fTP, _ := order.TotalPrice.Float64()
	res := entity.OrderDTOResponse{
		ID:                         order.ID,
		Buyer:                      entity.BuyerResponse{ID: order.Buyer.ID},
		Seller:                     entity.SellerResponse{ID: order.Seller.ID},
		DeliverySourceAddress:      order.DeliverySourceAddress,
		DeliveryDestinationAddress: order.DeliveryDestinationAddress,
		TotalQuantity:              order.TotalQuantity,
		TotalPrice:                 fTP,
		Status:                     order.Status,
		OrderDate:                  order.OrderDate,
	}

	for _, od := range order.Items {
		res.Items = append(res.Items, entity.OrderDetailDTOResponse{
			Product:  entity.ProductResponse{ID: od.Product.ID},
			Quantity: od.Quantity,
		})
	}

	return c.Status(http.StatusCreated).JSON(helper.ResponseSuccess{
		Data: res,
	})
}

func (oc *orderController) GetByUserID(c *fiber.Ctx) error {
	// take token from user value context
	tokenClaims, ok := c.Context().UserValue("tokenClaims").(jwt.MapClaims)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(errors.New("token claims not exists"))
	}

	// get id & user type
	userID := int64(tokenClaims["id"].(float64))
	userType := helper.UserTypeEnum(tokenClaims["type"].(float64))

	// get orders by buyer/seller id
	orders, err := oc.orderUsecase.GetByUserID(userID, userType)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	// transform Order to OrderDTOResponse
	var orderRes entity.OrderDTOResponse
	var res []entity.OrderDTOResponse
	for _, order := range orders {
		fTP, _ := order.TotalPrice.Float64()
		orderRes = entity.OrderDTOResponse{
			ID: order.ID,
			Buyer: entity.BuyerResponse{
				ID: order.Buyer.ID,
			},
			Seller: entity.SellerResponse{
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

	return c.Status(http.StatusOK).JSON(helper.ResponseSuccess{
		Data: res,
	})
}

func (oc *orderController) AcceptOrder(c *fiber.Ctx) error {
	// extract params
	orderId, idErr := c.ParamsInt("id")
	if idErr != nil {
		return c.Status(http.StatusBadRequest).JSON(idErr.Error())
	}

	// accept order
	order := new(entity.Order)
	order.ID = int64(orderId)
	err := oc.orderUsecase.AcceptOrder(order)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// transform Order to OrderDTOResponse
	fTP, _ := order.TotalPrice.Float64()
	res := entity.OrderDTOResponse{
		ID:                         order.ID,
		Buyer:                      entity.BuyerResponse{ID: order.Buyer.ID},
		Seller:                     entity.SellerResponse{ID: order.Seller.ID},
		DeliverySourceAddress:      order.DeliverySourceAddress,
		DeliveryDestinationAddress: order.DeliveryDestinationAddress,
		TotalQuantity:              order.TotalQuantity,
		TotalPrice:                 fTP,
		Status:                     order.Status,
		OrderDate:                  order.OrderDate,
	}

	for _, od := range order.Items {
		res.Items = append(res.Items, entity.OrderDetailDTOResponse{
			Product:  entity.ProductResponse{ID: od.Product.ID},
			Quantity: od.Quantity,
		})
	}

	return c.Status(http.StatusCreated).JSON(helper.ResponseSuccess{
		Data: res,
	})
}
