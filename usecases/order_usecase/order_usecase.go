package orderusecase

import (
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
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

func (u *orderUsecase) Store(order *entity.Order) resterrors.RestErr {
	tn, err := helpers.GetTimeNow()
	if err != nil {
		rErr := resterrors.NewInternalServerError("error when trying to save data", err)
		return rErr
	}
	order.OrderDate = tn

	repoErr := u.orderRepo.Store(order)
	if repoErr != nil {
		return repoErr
	}
	return nil
}

func (u *orderUsecase) GetByUserID(userID int64, userType helpers.UserTypeEnum) ([]entity.Order, resterrors.RestErr) {
	var orders []entity.Order
	var err resterrors.RestErr

	if userType == helpers.BUYER_TYPE {
		orders, err = u.orderRepo.GetByBuyerID(userID)
	} else if userType == helpers.SELLER_TYPE {
		orders, err = u.orderRepo.GetBySellerID(userID)
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (u *orderUsecase) AcceptOrder(order *entity.Order) (entity.Order, resterrors.RestErr) {
	repoRes, err := u.orderRepo.GetByID(order)
	if err != nil {
		return repoRes, err
	}

	repoRes.Status = entity.ACCEPTED
	updateErr := u.orderRepo.Update(&repoRes)
	if updateErr != nil {
		return repoRes, updateErr
	}

	return repoRes, nil
}
