package orderusecase_test

import (
	"testing"

	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
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

	time, err := helpers.GetTimeNow()
	assert.NoError(t, err)
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

func TestByUserID(t *testing.T) {
	mockOrderRepo := new(mocks.OrderRepository)

	mockBuyer1 := entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		SendingAddress: "sending address",
	}

	mockBuyer2 := entity.Buyer{
		ID:             2,
		Email:          "buyer2@mail.com",
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

	mockSeller2 := entity.Seller{
		ID:            2,
		Email:         "seller2@mail.com",
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

	time, err := helpers.GetTimeNow()
	assert.NoError(t, err)
	mockOrderForBuyer := entity.Order{
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

	mockOrderForSeller := entity.Order{
		ID:                         1,
		Buyer:                      mockBuyer2,
		Seller:                     mockSeller2,
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

	t.Run("success get order by buyer", func(t *testing.T) {
		mockOrdersForBuyer := []entity.Order{mockOrderForBuyer}
		mockOrderRepo.On("GetByBuyerID", mock.AnythingOfType("int64")).Return(mockOrdersForBuyer, nil).Once()

		u := orderusecase.NewOrderUsecase(mockOrderRepo)
		uRes, err := u.GetByUserID(mockBuyer1.ID, helpers.BUYER_TYPE)

		assert.NoError(t, err)
		assert.Equal(t, len(uRes), len(mockOrdersForBuyer))
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("success get order by seller", func(t *testing.T) {
		mockOrdersForSeller := []entity.Order{mockOrderForSeller}
		mockOrderRepo.On("GetBySellerID", mock.AnythingOfType("int64")).Return(mockOrdersForSeller, nil).Once()

		u := orderusecase.NewOrderUsecase(mockOrderRepo)
		uRes, err := u.GetByUserID(mockSeller2.ID, helpers.SELLER_TYPE)

		assert.NoError(t, err)
		assert.Equal(t, len(uRes), len(mockOrdersForSeller))
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestAcceptOrder(t *testing.T) {
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

	time, err := helpers.GetTimeNow()
	assert.NoError(t, err)
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
		mockOrderRepo.On("GetByID", mock.AnythingOfType("*entity.Order")).Return(mockOrder1, nil).Once()
		mockOrderRepo.On("Update", mock.AnythingOfType("*entity.Order")).Return(nil).Once()

		u := orderusecase.NewOrderUsecase(mockOrderRepo)
		uRes, err := u.AcceptOrder(&tmpMockOrder)

		assert.NoError(t, err)
		assert.Equal(t, uRes.Status, entity.ACCEPTED)
		mockOrderRepo.AssertExpectations(t)
	})
}
