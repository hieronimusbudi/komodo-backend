package sellerusecase_test

import (
	"testing"

	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	sellerusecase "github.com/hieronimusbudi/komodo-backend/usecases/seller_usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	mockSellerRepo := new(mocks.SellerRepository)
	mockSeller := entity.Seller{
		ID:            1,
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "123456",
		PickUpAddress: "pickup address",
	}

	mockSellerEmpty := entity.Seller{}

	mockSellerExist := entity.Seller{
		ID:            1,
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		PickUpAddress: "pickup address",
	}

	t.Run("success", func(t *testing.T) {
		tmpMockSeller := mockSeller

		mockSellerRepo.On("GetByEmail", mock.AnythingOfType("*entity.Seller")).Return(mockSellerEmpty, nil).Once()
		mockSellerRepo.On("Store", mock.AnythingOfType("*entity.Seller")).Return(nil).Once()

		u := sellerusecase.NewSellerUsecase(mockSellerRepo)
		err := u.Register(&tmpMockSeller)

		assert.NoError(t, err)
		assert.Equal(t, mockSeller.Email, tmpMockSeller.Email)
		mockSellerRepo.AssertExpectations(t)
	})

	t.Run("user is already exist", func(t *testing.T) {
		tmpMockSeller := mockSeller

		mockSellerRepo.On("GetByEmail", mock.AnythingOfType("*entity.Seller")).Return(mockSellerExist, nil).Once()
		u := sellerusecase.NewSellerUsecase(mockSellerRepo)
		err := u.Register(&tmpMockSeller)

		assert.Error(t, err)
		mockSellerRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockSellerRepo := new(mocks.SellerRepository)

	mockSellerRepoResponse := entity.Seller{
		ID:            1,
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		PickUpAddress: "pickup address",
	}

	mockSeller := entity.Seller{
		ID:            0,
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "12345",
		PickUpAddress: "pickup address",
	}

	t.Run("success", func(t *testing.T) {
		mockSellerRepo.On("GetByEmail", mock.AnythingOfType("*entity.Seller")).Return(mockSellerRepoResponse, nil).Once()

		u := sellerusecase.NewSellerUsecase(mockSellerRepo)
		uRes, err := u.Login(&mockSeller)

		assert.NoError(t, err)
		assert.Equal(t, mockSellerRepoResponse.ID, uRes.ID)
		mockSellerRepo.AssertExpectations(t)
	})
}
