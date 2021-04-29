package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	productcontroller "github.com/hieronimusbudi/komodo-backend/controllers/product_controller"
	"github.com/hieronimusbudi/komodo-backend/framework/middlerwares"
	productrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/product_repository"
	productusecase "github.com/hieronimusbudi/komodo-backend/usecases/product_usecase"
)

func productRoutes(app *fiber.App, conn *sql.DB) {
	r := productrepo.NewMysqlProductRepository(conn)
	u := productusecase.NewProductUsecase(r)
	c := productcontroller.NewProductController(u)

	app.Get("/products", c.GetAll)
	app.Post("/products", middlerwares.ValidateRequest, c.Store)
}
