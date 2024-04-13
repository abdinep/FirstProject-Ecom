package models

import "time"

type Order struct {
	ID             uint
	User           User
	UserID         int
	Address        Address
	AddressID      int    `json:"address_id"`
	CouponCode     string `json:"coupon_code"`
	OrderPrice     int
	PaymentMethod  string `json:"payment_method"`
	DeliveryCharge int    `json:"delivery_charge"`
	OrderDate      time.Time
	UpdateDate     time.Time
}
