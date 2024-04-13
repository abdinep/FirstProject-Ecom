package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string `gorm:"not null" json:"adminName"`
	Email    string `gorm:"unique,not null" json:"adminMail"`
	Password string `gorm:"not null" json:"adminPassword"`
}
