package models

import "time"

type Otp struct {
	Email     string    `gorm:"not null" json:"otpEmail"`
	Otp       string    `json:"otp"`
	Expire_at time.Time `gorm:"type:timestamp;not null"`
}
