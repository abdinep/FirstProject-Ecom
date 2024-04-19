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
	var productData []gin.H
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
		for _,value := range mostOrderedProduct {
			productData = append(productData, gin.H{
				"id":value.ID,
				"product_name":value.Product_Name,
				"price":value.Price,
				"quantity":value.Quantity,
				"size":value.Size,
				"desciption":value.Description,
				"offer":value.Offer,
				"category":value.Category.Name,
				"status":200,
			})
		}
		c.JSON(200, gin.H{
			"mostOrderedProduct":productData,
			"status":200,
		})
	case "category":
		var mostOrderedCategory []models.Category
		var categoryData []gin.H
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
		for _,value := range mostOrderedCategory {
			categoryData = append(categoryData, gin.H{
				"id":value.ID,
				"name":value.Name,
				"description":value.Description,
				"status":200,
			})
		}
		c.JSON(200, gin.H{
			"mostOrderedCategory":categoryData,
			"status":200,
		})
	}
}
