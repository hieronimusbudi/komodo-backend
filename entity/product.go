package entity

import "github.com/shopspring/decimal"

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       decimal.Decimal
	Seller      Seller
}

type ProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SellerID    int64   `json:"sellerId"`
}

type ProductResponse struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Seller      SellerResponse `json:"seller"`
}

type ProductUseCase interface {
	Store(product *Product) error
	GetAll() ([]Product, error)
}

type ProductRepository interface {
	GetAll() ([]Product, error)
	GetByID(product *Product) error
	Update(product *Product) error
	Store(product *Product) error
	Delete(product *Product) error
}
