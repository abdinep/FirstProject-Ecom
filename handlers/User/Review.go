package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ================================== Review ======================================================
type productReview struct{
    UserId    int `json:"review_user"`
    ProductId uint `json:"review_product"`
    Review    string `json:"review"`
    Time      string `json:"time"`
}
// @Summary Store or update review for a product
// @Description Store or update review for a product based on user input
// @Tags User-Review
// @Accept json
// @Produce json
// @Param request body productReview true "Review data"
// @Success 200 {json} JSON "Thanks for your review"
// @Failure 401 {json} JSON "Failed to bind data or failed to update review"
// @Router /products/review [post]
func ReviewStore(c *gin.Context) {
	var reviewstore productReview
	if err := c.ShouldBindJSON(&reviewstore); err != nil {
		c.JSON(401,gin.H{
			"error": "failed to bind data",
			"status":401,
		})
	} else {
		reviewstore.Time = time.Now().Format("2006-01-02")
		reviewstore.UserId = c.GetInt("userID")
		if err := initializers.DB.Create(&reviewstore); err.Error != nil {
			c.JSON(401,gin.H{
				"error": "failed to store review",
				"status":401,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":"Thank you for your valuable feedback",
				"status":200,
			})
		}
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
