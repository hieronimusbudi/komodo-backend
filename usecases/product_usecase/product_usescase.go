package productusecase

import (
	"github.com/hieronimusbudi/komodo-backend/entity"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
)

type productUsecase struct {
	productRepo entity.ProductRepository
}

// NewProductUsecase will create a object with entity.ProductUseCase interface representation
func NewProductUsecase(productRepo entity.ProductRepository) entity.ProductUseCase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (p *productUsecase) Store(product *entity.Product) resterrors.RestErr {
	repoErr := p.productRepo.Store(product)
	if repoErr != nil {
		return repoErr
	}
	return nil
}

func (p *productUsecase) GetAll() ([]entity.Product, resterrors.RestErr) {
	products, err := p.productRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return products, nil
}
