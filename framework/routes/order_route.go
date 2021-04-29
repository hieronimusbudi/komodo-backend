package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	ordercontroller "github.com/hieronimusbudi/komodo-backend/controllers/order_controller"
	"github.com/hieronimusbudi/komodo-backend/framework/middlerwares"
	orderrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/order_repository"
	orderusecase "github.com/hieronimusbudi/komodo-backend/usecases/order_usecase"
)

func orderRoutes(app *fiber.App, conn *sql.DB) {
	r := orderrepo.NewMysqlOrderRepository(conn)
	u := orderusecase.NewOrderUsecase(r)
	c := ordercontroller.NewOrderController(u)

	app.Get("/orders/by-id", middlerwares.ValidateRequest, c.GetByUserID)
	app.Post("/orders", middlerwares.ValidateRequest, middlerwares.BuyerTypeChecker, c.Store)
	app.Put("/orders/:id/accept", middlerwares.ValidateRequest, middlerwares.SellerTypeChecker, c.AcceptOrder)
}
