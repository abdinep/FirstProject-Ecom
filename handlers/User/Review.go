package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ================================== Review ======================================================
// @Summary Store or update review for a product
// @Description Store or update review for a product based on user input
// @Tags User-Review
// @Accept json
// @Produce json
// @Param review formData string true "Review data"
// @Param product_id path int true "Product ID"
// @Success 200 {json} JSON "Thanks for your review"
// @Failure 401 {json} JSON "Failed to bind data or failed to update review"
// @Router /products/review/{ID} [post]
func ReviewStore(c *gin.Context) {
	var reviewstore models.Review
	userid := c.GetInt("userID")
	productid, _ := strconv.Atoi(c.Param("ID"))
	reviewstore.Time = time.Now().Format("2006-01-02")
	reviewstore.UserId = userid
	reviewstore.Review = c.Request.FormValue("review")
	reviewstore.ProductId = uint(productid)
	fmt.Println("-------->", productid, reviewstore.Time, reviewstore.UserId, reviewstore.Review)
	if err := initializers.DB.Create(&reviewstore); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "failed to store review",
			"status": 401,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Thank you for your valuable feedback",
			"status":  200,
		})
	}
}
func ReviewView(id string, c *gin.Context) []models.Review {
	var reviewView []models.Review
	if err := initializers.DB.Joins("User").Joins("Product.Category").Find(&reviewView).Where("product_id=?", id); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to fetch reviews",
		})
		return []models.Review{}
	}
	return reviewView

}

// ========================================= END ==================================================
