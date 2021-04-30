package entity

import (
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       decimal.Decimal
	Seller      Seller
}

type ProductDTORequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SellerID    int64   `json:"sellerId"`
}

type ProductDTOResponse struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Seller      SellerDTOResponse `json:"seller"`
}

type ProductUseCase interface {
	Store(product *Product) resterrors.RestErr
	GetAll() ([]Product, resterrors.RestErr)
}

type ProductRepository interface {
	GetAll() ([]Product, resterrors.RestErr)
	GetByID(product *Product) resterrors.RestErr
	Update(product *Product) resterrors.RestErr
	Store(product *Product) resterrors.RestErr
	Delete(product *Product) resterrors.RestErr
}
