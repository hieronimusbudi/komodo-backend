package entity

type Buyer struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Password       string `json:"password"`
	SendingAddress string `json:"sendingAddress"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BuyerUseCase interface {
	Register(buyer *Buyer) error
	Login(logReq *LoginRequest) error
}

type BuyerRepository interface {
	Init()
	GetAll() ([]Buyer, error)
	GetByID(buyer *Buyer) error
	Update(buyer *Buyer) error
	Store(buyer *Buyer) error
	Delete(buyer *Buyer) error
	GetByEmailAndPassword(buyer *Buyer) error
}
