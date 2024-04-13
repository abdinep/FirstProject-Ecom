package models

import "time"

type Wallet struct {
	ID         uint
	User       User
	UserID     uint
	Balance    int
	Created_at time.Time
	Updated_at time.Time
}
