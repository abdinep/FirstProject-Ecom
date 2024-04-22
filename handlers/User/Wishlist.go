package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Add product to wishlist
// @Description Adds a product to the user's wishlist.
// @Tags User-Wishlist
// @Accept json
// @Produce json
// @Param ID path string true "Product ID"
// @Success 200 {json} JSON "Product added to wishlist"
// @Failure 401 {json} JSON "Cant add to wishlist"
// @Failure 401 {json} JSON "Product already added to wishlist"
// @Router /user/wishlist/{ID} [post]
func AddWishlist(c *gin.Context) {
	var wishlist models.Wishlist
	productid := c.Param("ID")
	userid := c.GetUint("userID")
	if err := initializers.DB.Where("user_id = ? AND product_id = ?", userid, productid).First(&wishlist); err.Error != nil {
		wishlist.ProductID, _ = strconv.Atoi(productid)
		wishlist.UserID = int(userid)
		if err := initializers.DB.Create(&wishlist); err.Error != nil {
			c.JSON(401, gin.H{
				"error":  "Cant add to wishlist",
				"status": 401,
			})
			fmt.Println("Cant add to wishlist=====>", err.Error)
		} else {
			c.JSON(200, gin.H{
				"message": "Product added to wishlist",
				"status":  200,
			})
		}
	} else {
		c.JSON(401, gin.H{
			"error":  "Product already added to wishlist",
			"status": 401,
		})
		fmt.Println("Product already added to wishlist======>", err.Error)
	}
}

// @Summary List wishlist
// @Description Lists the products in the user's wishlist.
// @Tags User-Wishlist
// @Accept json
// @Produce json
// @Success 200 {json} JSON "Listed Products from the wishlist"
// @Failure 401 {json} JSON "Product not found"
// @Router /user/wishlist [get]
func ViewWishlist(c *gin.Context) {
	var wishllist []models.Wishlist
	var listWishlist []gin.H
	userid := c.GetUint("userID")
	if err := initializers.DB.Joins("Product").Where("user_id = ?", userid).Find(&wishllist); err.Error != nil {
		c.JSON(500, gin.H{"Error": "No products in Wishlist"})
		fmt.Println("No products in Wishlist=====>", err.Error)
	} else {
		for _, v := range wishllist {
			listWishlist = append(listWishlist, gin.H{
				"wishlist_id":   v.ID,
				"product_id":    v.ProductID,
				"product_name":  v.Product.Product_Name,
				"product_price": v.Product.Price,
			})
		}
		c.JSON(200, gin.H{
			"data":   listWishlist,
			"status": 200,
		})
	}
}
// @Summary Remove product from wishlist
// @Description Removes a product from the user's wishlist.
// @Tags User-Wishlist
// @Accept json
// @Produce json
// @Param ID path string true "Product ID"
// @Success 200 {json} JSON "Product removed from wishlist"
// @Failure 401 {json} JSON "Cant Delete the Product",
// @Router /user/wishlist/{ID} [delete]
func DeleteWishlist(c *gin.Context) {
	var wishlist models.Wishlist
	productid := c.Param("ID")
	userid := c.GetUint("userID")

	if err := initializers.DB.Where("user_id = ? AND product_id = ?", userid, productid).First(&wishlist); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Product not Found",
			"status": 401,
		})
		fmt.Println("Invalid product=====>", err.Error)
	} else {
		if err := initializers.DB.Delete(&wishlist); err.Error != nil {
			c.JSON(401, gin.H{
				"error":  "Cant Delete the Product",
				"status": 401,
			})
			fmt.Println("Cant Delete the Product=======>", err.Error)
			return
		}
		c.JSON(200, gin.H{
			"message": "Product removed from wishlist",
			"status":  200,
		})
	}
}
