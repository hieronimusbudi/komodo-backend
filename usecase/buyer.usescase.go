package usecase

import (
	"time"

	"github.com/hieronimusbudi/komodo-backend/entity"
)

type buyerUsecase struct {
	buyerRepo      entity.BuyerRepository
	contextTimeout time.Duration
}

func NewBuyerUsecase(buyerRepo entity.BuyerRepository, timeout time.Duration) entity.BuyerUseCase {
	return &buyerUsecase{
		buyerRepo:      buyerRepo,
		contextTimeout: timeout,
	}
}

func (b *buyerUsecase) Register(buyer *entity.Buyer) error {
	err := b.buyerRepo.Store(buyer)
	if err != nil {
		return err
	}
	return nil
}

func (b *buyerUsecase) Login(logReq *entity.LoginRequest) error {
	buyer := entity.Buyer{}
	buyer.Email = logReq.Email
	buyer.Password = logReq.Password

	err := b.buyerRepo.GetByEmailAndPassword(&buyer)
	if err != nil {
		return err
	}
	return nil
}
