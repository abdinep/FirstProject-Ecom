package models

type OrderItem struct {
	ID            uint
	Order         Order
	OrderID       uint `json:"orderid"`
	Product       Product
	ProductID     int `json:"product_id"`
	OrderQuantity int
	Subtotal      float64 `json:"subtotal"`
	Orderstatus   string  `json:"status"`
}
