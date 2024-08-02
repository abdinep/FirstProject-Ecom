package handlers

import (
	"ecom/initializers"
	"ecom/models"

	"github.com/gin-gonic/gin"
)

func DisplayWallet(c *gin.Context){
	var walletData []models.Wallet
	UserId := c.GetUint("userID")

	if err := initializers.DB.Where("user_id = ?",UserId).Find(&walletData); err.Error != nil {
	
	}
}