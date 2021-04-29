package entity

type Seller struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	PickUpAddress string `json:"pickupAddress"`
}

type SellerResponse struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	PickUpAddress string `json:"pickupAddress"`
}

type SellerLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SellerUseCase interface {
	Register(seller *Seller) error
	Login(seller *Seller) (Seller, error)
}

type SellerRepository interface {
	GetAll() ([]Seller, error)
	GetByID(seller *Seller) error
	Update(seller *Seller) error
	Store(seller *Seller) error
	Delete(seller *Seller) error
	GetByEmail(seller *Seller) (Seller, error)
}
