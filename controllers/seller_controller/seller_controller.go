package sellercontroller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/config"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
)

type SellerController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type sellerController struct {
	sellerUseCase entity.SellerUseCase
	validate      *validator.Validate
}

func NewSellerController(u entity.SellerUseCase, v *validator.Validate) SellerController {
	return &sellerController{
		sellerUseCase: u,
		validate:      v,
	}
}

func (sctr *sellerController) Register(c *fiber.Ctx) error {
	sellerReq := new(entity.SellerDTORequest)
	if err := c.BodyParser(sellerReq); err != nil {
		rErr := resterrors.NewRestError("unprocessable entity", http.StatusUnprocessableEntity, err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// validate request
	vErr := sctr.validate.Struct(sellerReq)
	if vErr != nil {
		message, _ := helpers.CreateValidationMessage(vErr)
		rErr := resterrors.NewBadRequestError(message)
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	seller := entity.Seller{
		Email:         sellerReq.Email,
		Name:          sellerReq.Name,
		Password:      sellerReq.Password,
		PickUpAddress: sellerReq.PickUpAddress,
	}
	err := sctr.sellerUseCase.Register(&seller)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	sellerRes := entity.SellerDTOResponse{
		ID:            seller.ID,
		Email:         seller.Email,
		Name:          seller.Name,
		PickUpAddress: seller.PickUpAddress,
	}
	return c.Status(http.StatusCreated).JSON(helpers.SuccessResponse{
		Data: sellerRes,
	})
}

func (sctr *sellerController) Login(c *fiber.Ctx) error {
	loginReq := new(entity.SellerDTOLogin)
	if err := c.BodyParser(loginReq); err != nil {
		rErr := resterrors.NewRestError("unprocessable entity", http.StatusUnprocessableEntity, err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// validate request
	vErr := sctr.validate.Struct(loginReq)
	if vErr != nil {
		message, _ := helpers.CreateValidationMessage(vErr)
		rErr := resterrors.NewBadRequestError(message)
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	seller := entity.Seller{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	}
	seller, err := sctr.sellerUseCase.Login(&seller)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	// create jwt token
	jwtUserType := helpers.SELLER_TYPE
	token, tokenErr := helpers.GenerateToken(&helpers.UserJWTPayload{
		ID:    seller.ID,
		Email: seller.Email,
		Name:  seller.Name,
		Type:  jwtUserType,
	}, []byte(config.JWT_SECRET))
	if tokenErr != nil {
		rErr := resterrors.NewInternalServerError("generate token error ", tokenErr)
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	res := helpers.JWTResponse{
		Data: entity.SellerDTOResponse{
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
