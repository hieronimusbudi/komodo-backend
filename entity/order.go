package entity

import (
	"time"

	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
	"github.com/shopspring/decimal"
)

type OrderStatusEnum int

const (
	PENDING OrderStatusEnum = iota
	ACCEPTED
)

type Order struct {
	ID                         int64
	Buyer                      Buyer
	Seller                     Seller
	DeliverySourceAddress      string
	DeliveryDestinationAddress string
	TotalQuantity              int64
	TotalPrice                 decimal.Decimal
	Status                     OrderStatusEnum
	OrderDate                  time.Time
	Items                      []OrderDetail
}

type OrderDetail struct {
	ID       int64
	Product  Product
	Quantity int64
}

type OrderDTORequest struct {
	BuyerID                    int64                   `json:"buyerId"`
	SellerID                   int64                   `json:"sellerId"`
	DeliverySourceAddress      string                  `json:"deliverySourceAddress"`
	DeliveryDestinationAddress string                  `json:"deliveryDestinationAddress"`
	TotalQuantity              int64                   `json:"totalQuantity"`
	TotalPrice                 float64                 `json:"totalPrice"`
	Items                      []OrderDetailDTORequest `json:"items"`
}

type OrderDTOResponse struct {
	ID                         int64                    `json:"id"`
	Buyer                      BuyerDTOResponse         `json:"buyer"`
	Seller                     SellerDTOResponse        `json:"seller"`
	DeliverySourceAddress      string                   `json:"deliverySourceAddress"`
	DeliveryDestinationAddress string                   `json:"deliveryDestinationAddress"`
	TotalQuantity              int64                    `json:"totalQuantity"`
	TotalPrice                 float64                  `json:"totalPrice"`
	Status                     OrderStatusEnum          `json:"status"`
	OrderDate                  time.Time                `json:"orderDate"`
	Items                      []OrderDetailDTOResponse `json:"items"`
}

type OrderDetailDTORequest struct {
	ProductId int64 `json:"productId"`
	Quantity  int64 `json:"quantity"`
}

type OrderDetailDTOResponse struct {
	Product  ProductDTOResponse `json:"product"`
	Quantity int64              `json:"quantity"`
	Price    float64            `json:"price"`
}

type OrderUseCase interface {
	Store(order *Order) resterrors.RestErr
	GetByUserID(userID int64, userType helpers.UserTypeEnum) ([]Order, resterrors.RestErr)
	AcceptOrder(order *Order) (Order, resterrors.RestErr)
}

type OrderRepository interface {
	GetAll() ([]Order, resterrors.RestErr)
	GetByBuyerID(buyerID int64) ([]Order, resterrors.RestErr)
	GetBySellerID(buyerID int64) ([]Order, resterrors.RestErr)
	GetByID(order *Order) (Order, resterrors.RestErr)
	Update(order *Order) resterrors.RestErr
	Store(order *Order) resterrors.RestErr
	Delete(order *Order) resterrors.RestErr
}
