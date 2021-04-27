package entity

type Order struct {
	ID         int64
	Buyer      Buyer
	Seller     Seller
	Items      []Product
	Quantity   int64
	Price      int64
	TotalPrice int64
	Status     int64
}
