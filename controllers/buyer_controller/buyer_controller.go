package buyercontroller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/config"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
)

type BuyerController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type buyerController struct {
	buyerUsecase entity.BuyerUseCase
	validate     *validator.Validate
}

// NewBuyerController will create a object with BuyerController interface representation
func NewBuyerController(u entity.BuyerUseCase, v *validator.Validate) BuyerController {
	return &buyerController{
		buyerUsecase: u,
		validate:     v,
	}
}

func (bctr *buyerController) Register(c *fiber.Ctx) error {
	buyerReq := new(entity.BuyerDTORequest)
	if err := c.BodyParser(buyerReq); err != nil {
		rErr := resterrors.NewRestError("unprocessable entity", http.StatusUnprocessableEntity, err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// validate request
	vErr := bctr.validate.Struct(buyerReq)
	if vErr != nil {
		message, _ := helpers.CreateValidationMessage(vErr)
		rErr := resterrors.NewBadRequestError(message)
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	buyer := entity.Buyer{
		Email:          buyerReq.Email,
		Name:           buyerReq.Name,
		Password:       buyerReq.Password,
		SendingAddress: buyerReq.SendingAddress,
	}
	err := bctr.buyerUsecase.Register(&buyer)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	buyerRes := entity.BuyerDTOResponse{
		ID:             buyer.ID,
		Email:          buyer.Email,
		Name:           buyer.Name,
		SendingAddress: buyer.SendingAddress,
	}
	return c.Status(http.StatusCreated).JSON(helpers.SuccessResponse{
		Data: buyerRes,
	})
}

func (bctr *buyerController) Login(c *fiber.Ctx) error {
	loginReq := new(entity.BuyerDTOLogin)
	if err := c.BodyParser(loginReq); err != nil {
		rErr := resterrors.NewRestError("unprocessable entity", http.StatusUnprocessableEntity, err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// validate request
	vErr := bctr.validate.Struct(loginReq)
	if vErr != nil {
		message, _ := helpers.CreateValidationMessage(vErr)
		rErr := resterrors.NewBadRequestError(message)
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	buyer := entity.Buyer{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	}
	buyer, err := bctr.buyerUsecase.Login(&buyer)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	// create jwt token
	jwtUserType := helpers.BUYER_TYPE
	token, tokenErr := helpers.GenerateToken(&helpers.UserJWTPayload{
		ID:    buyer.ID,
		Email: buyer.Email,
		Name:  buyer.Name,
		Type:  jwtUserType,
	}, []byte(config.JWT_SECRET))
	if tokenErr != nil {
		rErr := resterrors.NewInternalServerError("generate token error ", tokenErr)
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	res := helpers.JWTResponse{
		Data: entity.BuyerDTOResponse{
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
