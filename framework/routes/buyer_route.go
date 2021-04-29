package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	buyercontroller "github.com/hieronimusbudi/komodo-backend/controllers/buyer_controller"
	buyerrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/buyer_repository"
	buyerusecase "github.com/hieronimusbudi/komodo-backend/usecases/buyer_usecase"
)

func buyerRoutes(app *fiber.App, conn *sql.DB) {
	r := buyerrepo.NewMysqlBuyerRepository(conn)
	u := buyerusecase.NewBuyerUsecase(r)
	c := buyercontroller.NewBuyerController(u)

	app.Post("/buyers/register", c.Register)
	app.Post("/buyers/login", c.Login)
}
