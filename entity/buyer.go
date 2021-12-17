package entity

import (
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
)

type Buyer struct {
	ID             int64
	Email          string
	Name           string
	Password       string
	SendingAddress string
}

type BuyerDTORequest struct {
	Email          string `json:"email" validate:"required,email"`
	Name           string `json:"name" validate:"required"`
	Password       string `json:"password" validate:"required"`
	SendingAddress string `json:"sendingAddress" validate:"gte=0,lte=511"`
}

// LoginRequest represents login request
//
// swagger:model loginRequest
type BuyerDTOLogin struct {
	// api key as user identification
	//
	// example: a639d2829e9e11eb9ed1bf478603cc01
	// required: true
	Email string `json:"email" validate:"required,email"`
	// api key as user identification
	//
	// example: a639d2829e9e11eb9ed1bf478603cc01
	// required: true
	Password string `json:"password" validate:"required"`
}

type BuyerDTOResponse struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	SendingAddress string `json:"sendingAddress"`
}

type BuyerUseCase interface {
	Register(buyer *Buyer) resterrors.RestErr
	Login(buyer *Buyer) (Buyer, resterrors.RestErr)
}

type BuyerRepository interface {
	GetAll() ([]Buyer, resterrors.RestErr)
	GetByID(buyer *Buyer) resterrors.RestErr
	Update(buyer *Buyer) resterrors.RestErr
	Store(buyer *Buyer) resterrors.RestErr
	Delete(buyer *Buyer) resterrors.RestErr
	GetByEmail(buyer *Buyer) (Buyer, resterrors.RestErr)
}
