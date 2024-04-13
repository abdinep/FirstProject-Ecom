package models

type Cart struct {
	Product    Product
	User       User
	ID         uint
	Product_Id int  `json:"product_id"`
	User_id    int  `json:"user_id"`
	Quantity   uint `json:"quantity"`
}
