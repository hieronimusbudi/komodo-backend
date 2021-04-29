package productcontroller

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helper"
	"github.com/shopspring/decimal"
)

type ProductController interface {
	Store(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type productController struct {
	productUsecase entity.ProductUseCase
}

func NewProductController(productUsecase entity.ProductUseCase) ProductController {
	return &productController{
		productUsecase: productUsecase,
	}
}

func (p *productController) Store(c *fiber.Ctx) error {
	// parse product from request body
	productReq := new(entity.ProductRequest)
	if err := c.BodyParser(productReq); err != nil {
		log.Printf("%v \n", err)
		return c.Status(http.StatusBadRequest).JSON(err.Error())
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
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// transform Product to ProductResponse
	fP, _ := product.Price.Float64()
	res := entity.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       fP,
		Seller: entity.SellerResponse{
			ID:            product.Seller.ID,
			Email:         product.Seller.Email,
			Name:          product.Seller.Name,
			PickUpAddress: product.Seller.PickUpAddress,
		},
	}

	return c.Status(http.StatusCreated).JSON(helper.ResponseSuccess{
		Data: res,
	})
}

func (p *productController) GetAll(c *fiber.Ctx) error {
	products, err := p.productUsecase.GetAll()
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(err.Error())
	}

	var productRes entity.ProductResponse
	var res []entity.ProductResponse

	for _, product := range products {
		fP, _ := product.Price.Float64()
		productRes = entity.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       fP,
			Seller: entity.SellerResponse{
				ID:            product.Seller.ID,
				Email:         product.Seller.Email,
				Name:          product.Seller.Name,
				PickUpAddress: product.Seller.PickUpAddress,
			},
		}

		res = append(res, productRes)
	}

	return c.Status(http.StatusOK).JSON(helper.ResponseSuccess{
		Data: res,
	})
}
