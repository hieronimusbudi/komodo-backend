package entity

import (
	"time"

	"github.com/hieronimusbudi/komodo-backend/framework/helper"
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
	Buyer                      BuyerResponse            `json:"buyer"`
	Seller                     SellerResponse           `json:"seller"`
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
	Product  ProductResponse `json:"product"`
	Quantity int64           `json:"quantity"`
	Price    float64         `json:"price"`
}

type OrderUseCase interface {
	Store(order *Order) error
	GetByUserID(userID int64, userType helper.UserTypeEnum) ([]Order, error)
	AcceptOrder(order *Order) error
}

type OrderRepository interface {
	GetAll() ([]Order, error)
	GetByBuyerID(buyerID int64) ([]Order, error)
	GetBySellerID(buyerID int64) ([]Order, error)
	GetByID(order *Order) error
	Update(order *Order) error
	Store(order *Order) error
	Delete(order *Order) error
}
