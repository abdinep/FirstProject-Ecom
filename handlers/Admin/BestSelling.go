package handlers

import (
	"ecom/initializers"
	"ecom/models"

	"github.com/gin-gonic/gin"
)
// @Summary  Get best sold product and category
// @Description Get best sold product and category
// @Tags Admin- BestSell
// @Accept json
// @Produce  json
// @Param sort query string true "Sort by 'product' or 'category'"
// @Success 200 {json} json	"best sell"
// @Failure 401 {json} json "not found"
// @Router /admin/bestsell [get]
func BestSelling(c *gin.Context) {
	sortBy := c.Query("sort")
	// query := initializers.DB
	switch sortBy {
	case "product":
		var mostOrderedProduct []models.Product
		query := `SELECT p.id,p.product_name, p.price, COUNT(oi.order_quantity) AS quantity
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		GROUP BY p.product_name, p.price, p.id
		ORDER BY quantity DESC
		LIMIT 10;`
		if err := initializers.DB.Raw(query).Scan(&mostOrderedProduct); err.Error != nil {
			c.JSON(400, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(200, gin.H{
			"mostOrderedProduct":mostOrderedProduct,
			"status":200,
		})
	case "category":
		var mostOrderedCategory []models.Category
		// c.JSON(200, mostOrderedCategory)
		query := `SELECT c.id,c.name, COUNT(oi.order_quantity) as quantity
				FROM order_items oi
				JOIN products p ON oi.product_id = p.id
				JOIN categories c ON c.id = p.category_id
				GROUP BY c.name,c.id
				ORDER BY quantity DESC
				LIMIT 10;`
		if err := initializers.DB.Raw(query).Scan(&mostOrderedCategory); err.Error != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{
			"mostOrderedCategory":mostOrderedCategory,
			"status":200,
		})
	}
}
