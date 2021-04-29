package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	sellercontroller "github.com/hieronimusbudi/komodo-backend/controllers/seller_controller"
	sellerrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/seller_repository"
	sellerusecase "github.com/hieronimusbudi/komodo-backend/usecases/seller_usecase"
)

func sellerRoutes(app *fiber.App, conn *sql.DB) {
	r := sellerrepo.NewMysqlSellerRepository(conn)
	u := sellerusecase.NewSellerUsecase(r)
	c := sellercontroller.NewSellerController(u)

	app.Post("/sellers/register", c.Register)
	app.Post("/sellers/login", c.Login)
}
