package sellerusecase

import (
	"github.com/hieronimusbudi/komodo-backend/entity"
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

func (s *sellerUsecase) Register(buyer *entity.Seller) error {
	// encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(buyer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	buyer.Password = string(hashedPassword)

	repoErr := s.sellerRepo.Store(buyer)
	if repoErr != nil {
		return repoErr
	}
	return nil
}

func (s *sellerUsecase) Login(seller *entity.Seller) (entity.Seller, error) {
	oriPass := seller.Password
	repoRes, err := s.sellerRepo.GetByEmail(seller)
	if err != nil {
		return *seller, err
	}

	// compare hashed password and requested password
	cprErr := bcrypt.CompareHashAndPassword([]byte(repoRes.Password), []byte(oriPass))
	if cprErr != nil {
		return *seller, cprErr
	}

	return repoRes, nil
}
