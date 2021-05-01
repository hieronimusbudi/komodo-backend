package sellerusecase

import (
	"fmt"

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

func (s *sellerUsecase) Register(seller *entity.Seller) resterrors.RestErr {
	// check existing user
	ss := new(entity.Seller)
	ss.Email = seller.Email
	repoRes, err := s.sellerRepo.GetByEmail(ss)
	if err != nil {
		if err.Causes() != "sql: no rows in result set" {
			return err
		}
	} else if repoRes.Email == seller.Email {
		return resterrors.NewBadRequestError(fmt.Sprintf("user with email %s is already exist", repoRes.Email))
	}

	// encrypt password
	hashedPassword, bErr := bcrypt.GenerateFromPassword([]byte(seller.Password), bcrypt.DefaultCost)
	if bErr != nil {
		return resterrors.NewInternalServerError(bErr.Error(), bErr)
	}

	seller.Password = string(hashedPassword)

	repoErr := s.sellerRepo.Store(seller)
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
