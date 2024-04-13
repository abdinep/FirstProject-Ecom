package models

type Wishlist struct {
	ID        uint
	Product   Product
	ProductID int
	User      User
	UserID    int
}

