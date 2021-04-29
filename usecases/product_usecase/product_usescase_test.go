package productusecase_test

import (
	"testing"

	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/entity/mocks"
	productusecase "github.com/hieronimusbudi/komodo-backend/usecases/product_usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStore(t *testing.T) {
	mockProductRepo := new(mocks.ProductRepository)
	mockSeller := entity.Seller{
		ID:            1,
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "123456",
		PickUpAddress: "pickup address",
	}

	mockProduct := entity.Product{
		ID:          0,
		Name:        "product1",
		Description: "desc",
		Price:       decimal.NewFromFloat(181818.11),
		Seller:      mockSeller,
	}

	t.Run("success", func(t *testing.T) {
		tmpMockProduct := mockProduct
		mockProductRepo.On("Store", mock.AnythingOfType("*entity.Product")).Return(nil).Once()

		u := productusecase.NewProductUsecase(mockProductRepo)
		err := u.Store(&tmpMockProduct)

		assert.NoError(t, err)
		assert.Equal(t, mockProduct.Name, tmpMockProduct.Name)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestGetAll(t *testing.T) {
	mockProductRepo := new(mocks.ProductRepository)
	mockProducts := []entity.Product{}

	mockSeller := entity.Seller{
		ID:            1,
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      "123456",
		PickUpAddress: "pickup address",
	}

	mockProduct1 := entity.Product{
		ID:          1,
		Name:        "product1",
		Description: "desc",
		Price:       decimal.NewFromFloat(181818.11),
		Seller:      mockSeller,
	}

	mockProduct2 := entity.Product{
		ID:          2,
		Name:        "product1",
		Description: "desc",
		Price:       decimal.NewFromFloat(181818.11),
		Seller:      mockSeller,
	}

	mockProducts = append(mockProducts, mockProduct1, mockProduct2)

	t.Run("success", func(t *testing.T) {
		mockProductRepo.On("GetAll").Return(mockProducts, nil).Once()

		u := productusecase.NewProductUsecase(mockProductRepo)
		uRes, err := u.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, len(uRes), len(mockProducts))
		mockProductRepo.AssertExpectations(t)
	})
}
