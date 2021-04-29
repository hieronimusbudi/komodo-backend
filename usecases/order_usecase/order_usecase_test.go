package orderusecase_test

import (
	"testing"
	"time"

	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	orderusecase "github.com/hieronimusbudi/komodo-backend/usecases/order_usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStore(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)

	mockBuyer1 := entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		SendingAddress: "sending address",
	}

	mockSeller1 := entity.Seller{
		ID:            1,
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		PickUpAddress: "pickup address",
	}

	mockProduct1 := entity.Product{
		ID:          1,
		Name:        "product1",
		Description: "desc",
		Price:       decimal.NewFromFloat(181818.11),
		Seller:      mockSeller1,
	}

	mockOrderDetail1 := entity.OrderDetail{
		ID:       1,
		Product:  mockProduct1,
		Quantity: 10,
	}

	time, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	mockOrder1 := entity.Order{
		ID:                         1,
		Buyer:                      mockBuyer1,
		Seller:                     mockSeller1,
		DeliverySourceAddress:      "pickup address",
		DeliveryDestinationAddress: "sending address",
		TotalQuantity:              10,
		TotalPrice:                 decimal.NewFromFloat(181818.11),
		Status:                     entity.PENDING,
		OrderDate:                  time,
		Items: []entity.OrderDetail{
			mockOrderDetail1,
		},
	}

	t.Run("success", func(t *testing.T) {
		tmpMockOrder := mockOrder1
		mockOrderRepo.On("Store", mock.AnythingOfType("*entity.Order")).Return(nil).Once()

		u := orderusecase.NewOrderUsecase(mockOrderRepo)
		err := u.Store(&tmpMockOrder)

		assert.NoError(t, err)
		assert.Equal(t, mockOrder1.Status, tmpMockOrder.Status)
		mockOrderRepo.AssertExpectations(t)
	})
}
