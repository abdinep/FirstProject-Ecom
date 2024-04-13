package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	Users     int `json:"rating_user"`
	ProductId int `gorm:"unique" json:"rating_product"`
	Product   Product
	Value     int `json:"rating_value"`
}
