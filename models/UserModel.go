package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"userName"`
	Email    string `gorm:"unique,not null" json:"userEmail"`
	Mobile   string `gorm:"not null" json:"Mob"`
	Password string `gorm:"not null" json:"userPassword"`
	Gender   string `gorm:"check: gender IN ('male', 'female','')" json:"gender"`
	Status   string `gorm:"default:Active" json:"status"`
}
