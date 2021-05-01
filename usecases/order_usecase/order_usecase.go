package orderusecase

import (
	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
	"github.com/shopspring/decimal"
)

type orderUsecase struct {
	orderRepo   entity.OrderRepository
	productRepo entity.ProductRepository
}

// NewOrderUsecase will create a object with entity.OrderUseCase interface representation
func NewOrderUsecase(orderRepo entity.OrderRepository, productRepo entity.ProductRepository) entity.OrderUseCase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *orderUsecase) Store(order *entity.Order) resterrors.RestErr {
	tn, err := helpers.GetTimeNow()
	order.OrderDate = tn
	if err != nil {
		rErr := resterrors.NewInternalServerError("error when trying to save data", err)
		return rErr
	}

	totalPrice := decimal.NewFromFloat(0)
	totalQuantity := int64(0)
	items := []entity.OrderDetail{}
	for _, od := range order.Items {
		p, err := u.productRepo.GetByID(&od.Product)
		if err != nil {
			return err
		}

		nOd := entity.OrderDetail{
			ID: p.ID,
			Product: entity.Product{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Seller:      p.Seller,
			},
			Quantity: od.Quantity,
		}

		qD := decimal.NewFromInt(nOd.Quantity)
		totalPrice = totalPrice.Add(p.Price.Mul(qD))
		totalQuantity += nOd.Quantity
		items = append(items, nOd)
	}

	order.TotalPrice = totalPrice
	order.TotalQuantity = totalQuantity
	order.Items = append(items[:0:0], items...)

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

	var newOrder entity.Order
	var newOrders []entity.Order
	var items []entity.OrderDetail

	for _, o := range orders {
		items = []entity.OrderDetail(nil)
		newOrder = entity.Order{
			ID:                         o.ID,
			Buyer:                      o.Buyer,
			Seller:                     o.Seller,
			DeliverySourceAddress:      o.DeliverySourceAddress,
			DeliveryDestinationAddress: o.DeliveryDestinationAddress,
			TotalQuantity:              o.TotalQuantity,
			TotalPrice:                 o.TotalPrice,
			Status:                     o.Status,
			OrderDate:                  o.OrderDate,
			Items:                      items,
		}
		for _, od := range o.Items {
			p, err := u.productRepo.GetByID(&od.Product)
			if err != nil {
				return nil, err
			}

			nOd := entity.OrderDetail{
				ID: od.ID,
				Product: entity.Product{
					ID:          p.ID,
					Name:        p.Name,
					Description: p.Description,
					Price:       p.Price,
					Seller:      p.Seller,
				},
				Quantity: od.Quantity,
			}

			items = append(items, nOd)
		}

		newOrder.Items = items
		newOrders = append(newOrders, newOrder)
	}

	if err != nil {
		return nil, err
	}
	return newOrders, nil
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
