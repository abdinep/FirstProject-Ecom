package models

type Review struct {
	ID        uint
	UserId    int `json:"review_user"`
	User      User
	ProductId uint `json:"review_product"`
	Product   Product
	Review    string `json:"review"`
	Time      string
}
