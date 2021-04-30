package sellerusecase

import (
	"github.com/hieronimusbudi/komodo-backend/entity"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
	"golang.org/x/crypto/bcrypt"
)

type sellerUsecase struct {
	sellerRepo entity.SellerRepository
}

// NewSellerUsecase will create a object with entity.NewSellerUsecase interface representation
func NewSellerUsecase(sellerRepo entity.SellerRepository) entity.SellerUseCase {
	return &sellerUsecase{
		sellerRepo: sellerRepo,
	}
}

func (s *sellerUsecase) Register(buyer *entity.Seller) resterrors.RestErr {
	// encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(buyer.Password), bcrypt.DefaultCost)
	if err != nil {
		return resterrors.NewInternalServerError(err.Error(), err)
	}

	buyer.Password = string(hashedPassword)

	repoErr := s.sellerRepo.Store(buyer)
	if repoErr != nil {
		return repoErr
	}
	return nil
}

func (s *sellerUsecase) Login(seller *entity.Seller) (entity.Seller, resterrors.RestErr) {
	oriPass := seller.Password
	repoRes, err := s.sellerRepo.GetByEmail(seller)
	if err != nil {
		return *seller, err
	}

	// compare hashed password and requested password
	cprErr := bcrypt.CompareHashAndPassword([]byte(repoRes.Password), []byte(oriPass))
	if cprErr != nil {
		return *seller, resterrors.NewInternalServerError(cprErr.Error(), cprErr)
	}

	return repoRes, nil
}
