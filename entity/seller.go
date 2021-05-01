package entity

import resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"

type Seller struct {
	ID            int64
	Email         string
	Name          string
	Password      string
	PickUpAddress string
}

type SellerDTORequest struct {
	Email         string `json:"email" validate:"required,email"`
	Name          string `json:"name" validate:"required"`
	Password      string `json:"password" validate:"required"`
	PickUpAddress string `json:"pickupAddress" validate:"gte=0,lte=511"`
}

type SellerDTOLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SellerDTOResponse struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	PickUpAddress string `json:"pickupAddress"`
}

type SellerUseCase interface {
	Register(seller *Seller) resterrors.RestErr
	Login(seller *Seller) (Seller, resterrors.RestErr)
}

type SellerRepository interface {
	GetAll() ([]Seller, resterrors.RestErr)
	GetByID(seller *Seller) resterrors.RestErr
	Update(seller *Seller) resterrors.RestErr
	Store(seller *Seller) resterrors.RestErr
	Delete(seller *Seller) resterrors.RestErr
	GetByEmail(seller *Seller) (Seller, resterrors.RestErr)
}
