package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// =============================== Rating ==============================================
type productRating struct {
	// Users     int `json:"rating_user"`
	ProductId int `json:"rating_product"`
	Value     int `json:"rating_value"`
}

// @Summary Store or update rating for a product
// @Description Store or update rating for a product based on user input
// @Tags User-Rating
// @Accept json
// @Produce json
// @Param request body productRating true "Rating data"
// @Success 200 {json} json "Thanks for rating"
// @Failure 400 {json} json "Failed to bind data or failed to update rating"
// @Router /products/rating [post]
func RatingStrore(c *gin.Context) {
	var userrate productRating
	var dbrate models.Rating
	// ID := c.Param("ID")
	if err := c.ShouldBindJSON(&userrate); err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to bind data",
			"status": 401,
		})
	}
	result := initializers.DB.First(&dbrate, "product_id=?", userrate.ProductId)
	dbrate.ProductId = userrate.ProductId
	dbrate.Value = userrate.Value
	dbrate.Users = 1
	if result.Error != nil {
		if err := initializers.DB.Create(&dbrate).Error; err != nil {
			c.JSON(401, gin.H{
				"error":  "failed to store",
				"status": 401,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Thanks for rating",
				"status":  200,
			})
		}
	} else {
		err := initializers.DB.Model(&dbrate).Where("product_id=?", userrate.ProductId).Updates(models.Rating{
			Users: dbrate.Users + 1,
			Value: dbrate.Value + userrate.Value,
		})
		if err.Error != nil {
			c.JSON(401, gin.H{
				"error":  "failed to update data",
				"status": 401,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Thanks for rating",
				"status":  200,
			})
		}
	}
	dbrate = models.Rating{}

}
func Ratingcalc(id string, c *gin.Context) float64 {
	var ratinguser models.Rating
	if err := initializers.DB.First(&ratinguser, "product_id=?", id); err.Error != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error":"failed to fetch product data",
		// })
	} else {
		averageratio := float64(ratinguser.Value) / float64(ratinguser.Users)
		ratinguser = models.Rating{}
		result := fmt.Sprintf("%.1f", averageratio)
		averageratio, _ = strconv.ParseFloat(result, 64)
		return averageratio
	}
	return 0
}

// ========================================= END ==================================================
