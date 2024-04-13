package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	PaymentID     string `json:"paymentid"`
	OrderID       string `json:"orderid"`
	Receipt       int
	PaymentStatus string
	PaymentAmount int
}
