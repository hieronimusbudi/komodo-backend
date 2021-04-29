package orderusecase

import (
	"log"
	"time"

	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helper"
)

type orderUsecase struct {
	orderRepo entity.OrderRepository
}

// NewOrderUsecase will create a object with entity.OrderUseCase interface representation
func NewOrderUsecase(orderRepo entity.OrderRepository) entity.OrderUseCase {
	return &orderUsecase{
		orderRepo: orderRepo,
	}
}

func (u *orderUsecase) Store(order *entity.Order) error {
	order.OrderDate, _ = time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	repoErr := u.orderRepo.Store(order)
	if repoErr != nil {
		return repoErr
	}
	return nil
}

func (u *orderUsecase) GetByUserID(userID int64, userType helper.UserTypeEnum) ([]entity.Order, error) {
	var orders []entity.Order
	var err error

	if userType == helper.BUYER_TYPE {
		orders, err = u.orderRepo.GetByBuyerID(userID)
	} else if userType == helper.SELLER_TYPE {
		orders, err = u.orderRepo.GetBySellerID(userID)
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (u *orderUsecase) AcceptOrder(order *entity.Order) error {
	err := u.orderRepo.GetByID(order)
	if err != nil {
		log.Println(1, err)
		return err
	}

	order.Status = entity.ACCEPTED
	updateErr := u.orderRepo.Update(order)
	if updateErr != nil {
		log.Println(2, err)
		return updateErr
	}

	return nil
}
