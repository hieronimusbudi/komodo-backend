package buyerusecase_test

import (
	"testing"

	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	buyerusecase "github.com/hieronimusbudi/komodo-backend/usecases/buyer_usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	mockBuyerRepo := new(mocks.BuyerRepository)
	mockBuyer := entity.Buyer{
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "123456",
		SendingAddress: "sending address",
	}

	mockBuyerEmpty := entity.Buyer{}

	mockBuyerExist := entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		SendingAddress: "sending address",
	}

	t.Run("success", func(t *testing.T) {
		tmpMockBuyer := mockBuyer

		mockBuyerRepo.On("GetByEmail", mock.AnythingOfType("*entity.Buyer")).Return(mockBuyerEmpty, nil).Once()
		mockBuyerRepo.On("Store", mock.AnythingOfType("*entity.Buyer")).Return(nil).Once()

		u := buyerusecase.NewBuyerUsecase(mockBuyerRepo)
		err := u.Register(&tmpMockBuyer)

		assert.NoError(t, err)
		assert.Equal(t, mockBuyer.Email, tmpMockBuyer.Email)
		mockBuyerRepo.AssertExpectations(t)
	})

	t.Run("user is already exist", func(t *testing.T) {
		tmpMockBuyer := mockBuyer

		mockBuyerRepo.On("GetByEmail", mock.AnythingOfType("*entity.Buyer")).Return(mockBuyerExist, nil).Once()
		u := buyerusecase.NewBuyerUsecase(mockBuyerRepo)
		err := u.Register(&tmpMockBuyer)

		assert.Error(t, err)
		mockBuyerRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockBuyerRepo := new(mocks.BuyerRepository)

	mockBuyerRepoResponse := entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "$2a$10$634oWhFDuTohq7suxGn5TuRQ8BGmWu9wFfiHZelLwfqSgWk/45vzu",
		SendingAddress: "sending address",
	}

	mockBuyer := entity.Buyer{
		ID:             0,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       "12345",
		SendingAddress: "sending address",
	}

	t.Run("success", func(t *testing.T) {
		mockBuyerRepo.On("GetByEmail", mock.AnythingOfType("*entity.Buyer")).Return(mockBuyerRepoResponse, nil).Once()

		u := buyerusecase.NewBuyerUsecase(mockBuyerRepo)
		uRes, err := u.Login(&mockBuyer)

		assert.NoError(t, err)
		assert.Equal(t, mockBuyerRepoResponse.ID, uRes.ID)
		mockBuyerRepo.AssertExpectations(t)
	})
}
