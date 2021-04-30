package routes

import (
	"github.com/gofiber/fiber/v2"
	productcontroller "github.com/hieronimusbudi/komodo-backend/controllers/product_controller"
	"github.com/hieronimusbudi/komodo-backend/dependencies"
	"github.com/hieronimusbudi/komodo-backend/framework/middlerwares"
	productrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/product_repository"
	productusecase "github.com/hieronimusbudi/komodo-backend/usecases/product_usecase"
)

// productRoutes used to define route and inject dependencies to repository, usecase and controller
func productRoutes(app *fiber.App, d *dependencies.Dependencies) {
	// inject connection to repository
	r := productrepo.NewMysqlProductRepository(d.Conn)
	// inject repository to usecase
	u := productusecase.NewProductUsecase(r)
	// inject usecase to controller
	c := productcontroller.NewProductController(u)

	app.Get("/products", c.GetAll)
	app.Post("/products", middlerwares.ValidateRequest, c.Store)
}
