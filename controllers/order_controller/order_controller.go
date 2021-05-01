package ordercontroller

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
)

type OrderController interface {
	Store(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	AcceptOrder(c *fiber.Ctx) error
}

type orderController struct {
	orderUsecase entity.OrderUseCase
	validate     *validator.Validate
}

// NewOrderController will create a object with OrderController interface representation
func NewOrderController(o entity.OrderUseCase, v *validator.Validate) OrderController {
	return &orderController{
		orderUsecase: o,
		validate:     v,
	}
}

func (octr *orderController) Store(c *fiber.Ctx) error {
	// parse order from request body
	oDTOReq := new(entity.OrderDTORequest)
	if err := c.BodyParser(oDTOReq); err != nil {
		rErr := resterrors.NewRestError("unprocessable entity", http.StatusUnprocessableEntity, err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// validate request
	vErr := octr.validate.Struct(oDTOReq)
	if vErr != nil {
		message, _ := helpers.CreateValidationMessage(vErr)
		rErr := resterrors.NewBadRequestError(message)
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	order := entity.Order{
		Buyer:                      entity.Buyer{ID: oDTOReq.BuyerID},
		Seller:                     entity.Seller{ID: oDTOReq.SellerID},
		DeliverySourceAddress:      oDTOReq.DeliverySourceAddress,
		DeliveryDestinationAddress: oDTOReq.DeliveryDestinationAddress,
		Status:                     entity.PENDING,
	}

	for _, od := range oDTOReq.Items {
		order.Items = append(order.Items, entity.OrderDetail{
			Product: entity.Product{
				ID: od.ProductId,
			},
			Quantity: od.Quantity,
		})
	}

	// store Order
	err := octr.orderUsecase.Store(&order)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	// transform Order to OrderDTOResponse
	fTP, _ := order.TotalPrice.Float64()
	res := entity.OrderDTOResponse{
		ID:                         order.ID,
		BuyerID:                    order.Buyer.ID,
		SellerID:                   order.Seller.ID,
		DeliverySourceAddress:      order.DeliverySourceAddress,
		DeliveryDestinationAddress: order.DeliveryDestinationAddress,
		TotalQuantity:              order.TotalQuantity,
		TotalPrice:                 fTP,
		Status:                     order.Status,
		OrderDate:                  order.OrderDate,
	}

	for _, od := range order.Items {
		fP, _ := od.Product.Price.Float64()
		res.Items = append(res.Items, entity.OrderDetailDTOResponse{
			Product: entity.ProductDTOResponse{
				ID:          od.Product.ID,
				Name:        od.Product.Name,
				Description: od.Product.Description,
				Price:       fP,
				SellerID:    od.Product.Seller.ID,
			},
			Price:    fP,
			Quantity: od.Quantity,
		})
	}

	return c.Status(http.StatusCreated).JSON(helpers.SuccessResponse{
		Data: res,
	})
}

func (octr *orderController) GetByUserID(c *fiber.Ctx) error {
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
	orders, err := octr.orderUsecase.GetByUserID(userID, userType)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.ErrorResponse())
	}

	// transform Order to OrderDTOResponse
	var orderRes entity.OrderDTOResponse
	var res []entity.OrderDTOResponse
	for _, order := range orders {
		fTP, _ := order.TotalPrice.Float64()
		orderRes = entity.OrderDTOResponse{
			ID:                         order.ID,
			BuyerID:                    order.Buyer.ID,
			SellerID:                   order.Seller.ID,
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
			odRow.Quantity = od.Quantity
			fP, _ := od.Product.Price.Float64()

			odRow.Price = fP
			odRow.Product.ID = od.Product.ID
			odRow.Product.Price = fP
			odRow.Product.Description = od.Product.Description
			odRow.Product.Name = od.Product.Name
			odRow.Product.SellerID = od.Product.Seller.ID

			orderRes.Items = append(orderRes.Items, odRow)
		}

		res = append(res, orderRes)
	}

	return c.Status(http.StatusOK).JSON(helpers.SuccessResponse{
		Data: res,
	})
}

func (octr *orderController) AcceptOrder(c *fiber.Ctx) error {
	// extract params
	orderId, idErr := c.ParamsInt("id")
	if idErr != nil {
		rErr := resterrors.NewBadRequestError(idErr.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// accept order
	order := new(entity.Order)
	order.ID = int64(orderId)
	uOrderRes, err := octr.orderUsecase.AcceptOrder(order)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	// transform Order to OrderDTOSimpleResponse
	fTP, _ := uOrderRes.TotalPrice.Float64()
	res := entity.OrderDTOSimpleResponse{
		ID:                         uOrderRes.ID,
		BuyerID:                    order.Buyer.ID,
		SellerID:                   order.Seller.ID,
		DeliverySourceAddress:      uOrderRes.DeliverySourceAddress,
		DeliveryDestinationAddress: uOrderRes.DeliveryDestinationAddress,
		TotalQuantity:              uOrderRes.TotalQuantity,
		TotalPrice:                 fTP,
		Status:                     uOrderRes.Status,
		OrderDate:                  uOrderRes.OrderDate,
	}

	return c.Status(http.StatusOK).JSON(helpers.SuccessResponse{
		Data: res,
	})
}
