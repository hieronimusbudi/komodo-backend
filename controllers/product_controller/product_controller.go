package productcontroller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
	"github.com/shopspring/decimal"
)

type ProductController interface {
	Store(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type productController struct {
	productUsecase entity.ProductUseCase
}

// NewProductController will create a object with ProductController interface representation
func NewProductController(productUsecase entity.ProductUseCase) ProductController {
	return &productController{
		productUsecase: productUsecase,
	}
}

func (p *productController) Store(c *fiber.Ctx) error {
	// parse product from request body
	productReq := new(entity.ProductDTORequest)
	if err := c.BodyParser(productReq); err != nil {
		rErr := resterrors.NewRestError("unprocessable entity", http.StatusUnprocessableEntity, err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	// store Product
	dP := decimal.NewFromFloat(productReq.Price)
	product := entity.Product{
		Name:        productReq.Name,
		Description: productReq.Description,
		Price:       dP,
		Seller:      entity.Seller{ID: productReq.SellerID},
	}
	err := p.productUsecase.Store(&product)
	if err != nil {
		return c.Status(err.Status()).JSON(err.ErrorResponse())
	}

	// transform Product to ProductResponse
	fP, _ := product.Price.Float64()
	res := entity.ProductDTOResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       fP,
		Seller: entity.SellerDTOResponse{
			ID:            product.Seller.ID,
			Email:         product.Seller.Email,
			Name:          product.Seller.Name,
			PickUpAddress: product.Seller.PickUpAddress,
		},
	}

	return c.Status(http.StatusCreated).JSON(helpers.SuccessResponse{
		Data: res,
	})
}

func (p *productController) GetAll(c *fiber.Ctx) error {
	products, err := p.productUsecase.GetAll()
	if err != nil {
		rErr := resterrors.NewRestError("unprocessable entity", http.StatusUnprocessableEntity, err.Error())
		return c.Status(rErr.Status()).JSON(rErr.ErrorResponse())
	}

	var productRes entity.ProductDTOResponse
	var res []entity.ProductDTOResponse

	for _, product := range products {
		fP, _ := product.Price.Float64()
		productRes = entity.ProductDTOResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       fP,
			Seller: entity.SellerDTOResponse{
				ID:            product.Seller.ID,
				Email:         product.Seller.Email,
				Name:          product.Seller.Name,
				PickUpAddress: product.Seller.PickUpAddress,
			},
		}

		res = append(res, productRes)
	}

	return c.Status(http.StatusOK).JSON(helpers.SuccessResponse{
		Data: res,
	})
}
