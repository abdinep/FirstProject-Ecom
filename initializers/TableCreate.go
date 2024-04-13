package initializers

import (
	"ecom/models"
	"log"
)
//======================================= creating tables on DB 
func TableCreate() {
	err := DB.AutoMigrate(&models.User{}, &models.Admin{}, &models.Product{}, &models.Otp{},&models.Category{},
		&models.Rating{},&models.Review{},&models.Address{},&models.Cart{},&models.Coupon{},&models.Order{},&models.OrderItem{},
	    &models.Payment{},&models.Wallet{},&models.Wishlist{},&models.Offer{})
	if err != nil {
		log.Fatal("Failed to Automigrate", err)
	}
}
