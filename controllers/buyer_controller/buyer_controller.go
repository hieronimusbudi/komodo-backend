package buyercontroller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/config"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helper"
)

type BuyerController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type buyerController struct {
	buyerUsecase entity.BuyerUseCase
}

func NewBuyerController(buyerUsecase entity.BuyerUseCase) BuyerController {
	return &buyerController{
		buyerUsecase: buyerUsecase,
	}
}

func (b *buyerController) Register(c *fiber.Ctx) error {
	buyer := new(entity.Buyer)
	if err := c.BodyParser(buyer); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	err := b.buyerUsecase.Register(buyer)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	return c.Status(http.StatusCreated).JSON(helper.ResponseSuccess{
		Data: buyer,
	})
}

func (b *buyerController) Login(c *fiber.Ctx) error {
	loginReq := new(entity.BuyerLoginRequest)
	if err := c.BodyParser(loginReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	buyer := entity.Buyer{}
	buyer.Email = loginReq.Email
	buyer.Password = loginReq.Password
	buyer, err := b.buyerUsecase.Login(&buyer)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(err.Error())
	}

	// create jwt token
	jwtUserType := helper.BUYER_TYPE
	token, tokenErr := helper.GenerateToken(&helper.UserJWTPayload{
		ID:    buyer.ID,
		Email: buyer.Email,
		Name:  buyer.Name,
		Type:  jwtUserType,
	}, []byte(config.JWT_SECRET))
	if tokenErr != nil {
		return c.Status(http.StatusUnauthorized).JSON(tokenErr.Error())
	}

	res := helper.JWTResponse{
		Data: entity.BuyerResponse{
			ID:             buyer.ID,
			Email:          buyer.Email,
			Name:           buyer.Name,
			SendingAddress: buyer.SendingAddress,
		},
		Type:  jwtUserType,
		Token: token,
	}
	return c.Status(http.StatusCreated).JSON(res)
}
