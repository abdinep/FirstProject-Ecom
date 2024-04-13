package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddWishlist(c *gin.Context) {
	var wishlist models.Wishlist

	productid := c.Param("ID")
	userid := c.GetUint("userID")
	if err := initializers.DB.Where("user_id = ? AND product_id = ?", userid, productid).First(&wishlist); err.Error != nil {
		wishlist.ProductID, _ = strconv.Atoi(productid)
		wishlist.UserID = int(userid)
		if err := initializers.DB.Create(&wishlist); err.Error != nil {
			c.JSON(500, gin.H{"Error": "Cant add to wishlist"})
			fmt.Println("Cant add to wishlist=====>", err.Error)
		} else {
			c.JSON(200, gin.H{"Message": "Product added to wishlist"})
		}
	} else {
		c.JSON(500, gin.H{"Error": "Product already added to wishlist"})
		fmt.Println("Product already added to wishlist======>", err.Error)
	}
}
func ViewWishlist(c *gin.Context) {
	var wishllist []models.Wishlist

	userid := c.GetUint("userID")
	if err := initializers.DB.Joins("Product").Find(&wishllist).Where("user_id = ?", userid); err.Error != nil {
		c.JSON(500, gin.H{"Error": "No products in Wishlist"})
		fmt.Println("No products in Wishlist=====>", err.Error)
	} else {
		c.JSON(200,gin.H{
			"data":wishllist,
		})
	}
}
func DeleteWishlist(c *gin.Context) {
	var wishlist models.Wishlist
	productid := c.Param("ID")
	userid := c.GetUint("userID")

	if err := initializers.DB.Where("user_id = ? AND product_id = ?",userid,productid).First(&wishlist); err.Error != nil{
		c.JSON(500,gin.H{"Error": "Invalid product"})
		fmt.Println("Invalid product=====>",err.Error)
	}else{
		if err := initializers.DB.Delete(&wishlist); err.Error != nil{
			c.JSON(500,gin.H{"Error":"Cant Delete the Product"})
			fmt.Println("Cant Delete the Product=======>",err.Error)
			return
		}
		c.JSON(200,gin.H{"Message":"Product removed from wishlist"})
	}
}
