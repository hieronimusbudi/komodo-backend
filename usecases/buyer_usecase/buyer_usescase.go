package buyerusecase

import (
	"github.com/hieronimusbudi/komodo-backend/entity"
	"golang.org/x/crypto/bcrypt"
)

type buyerUsecase struct {
	buyerRepo entity.BuyerRepository
}

// NewBuyerUsecase will create a object with entity.BuyerUseCase interface representation
func NewBuyerUsecase(buyerRepo entity.BuyerRepository) entity.BuyerUseCase {
	return &buyerUsecase{
		buyerRepo: buyerRepo,
	}
}

func (b *buyerUsecase) Register(buyer *entity.Buyer) error {
	// encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(buyer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	buyer.Password = string(hashedPassword)

	repoErr := b.buyerRepo.Store(buyer)
	if repoErr != nil {
		return repoErr
	}
	return nil
}

func (b *buyerUsecase) Login(buyer *entity.Buyer) (entity.Buyer, error) {
	oriPass := buyer.Password
	repoRes, err := b.buyerRepo.GetByEmail(buyer)
	if err != nil {
		return *buyer, err
	}

	// compare password
	cprErr := bcrypt.CompareHashAndPassword([]byte(repoRes.Password), []byte(oriPass))
	if cprErr != nil {
		return *buyer, cprErr
	}

	return repoRes, nil
}
