package entity

type Buyer struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Password       string `json:"password"`
	SendingAddress string `json:"sendingAddress"`
}

type BuyerResponse struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	SendingAddress string `json:"sendingAddress"`
}

type BuyerLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BuyerUseCase interface {
	Register(buyer *Buyer) error
	Login(buyer *Buyer) (Buyer, error)
}

type BuyerRepository interface {
	GetAll() ([]Buyer, error)
	GetByID(buyer *Buyer) error
	Update(buyer *Buyer) error
	Store(buyer *Buyer) error
	Delete(buyer *Buyer) error
	GetByEmail(buyer *Buyer) (Buyer, error)
}
