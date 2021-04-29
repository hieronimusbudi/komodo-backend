package sellercontroller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/config"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helper"
)

type SellerController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type sellerController struct {
	sellerUseCase entity.SellerUseCase
}

func NewSellerController(sellerUseCase entity.SellerUseCase) SellerController {
	return &sellerController{
		sellerUseCase: sellerUseCase,
	}
}

func (s *sellerController) Register(c *fiber.Ctx) error {
	seller := new(entity.Seller)
	if err := c.BodyParser(seller); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	err := s.sellerUseCase.Register(seller)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(helper.ResponseSuccess{
		Data: seller,
	})
}

func (s *sellerController) Login(c *fiber.Ctx) error {
	loginReq := new(entity.SellerLoginRequest)
	if err := c.BodyParser(loginReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	seller := entity.Seller{}
	seller.Email = loginReq.Email
	seller.Password = loginReq.Password
	seller, err := s.sellerUseCase.Login(&seller)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(err.Error())
	}

	// create jwt token
	jwtUserType := helper.SELLER_TYPE
	token, tokenErr := helper.GenerateToken(&helper.UserJWTPayload{
		ID:    seller.ID,
		Email: seller.Email,
		Name:  seller.Name,
		Type:  jwtUserType,
	}, []byte(config.JWT_SECRET))
	if tokenErr != nil {
		return c.Status(http.StatusUnauthorized).JSON(tokenErr.Error())
	}

	res := helper.JWTResponse{
		Data: entity.SellerResponse{
			ID:            seller.ID,
			Email:         seller.Email,
			Name:          seller.Name,
			PickUpAddress: seller.PickUpAddress,
		},
		Type:  jwtUserType,
		Token: token,
	}
	return c.Status(http.StatusCreated).JSON(res)
}
