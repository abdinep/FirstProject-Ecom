package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"strings"

	"github.com/gin-gonic/gin"
)
// @Summary Advanced Search
// @Description User can search products  by name, category or price range. It returns a list of product objects that match
// @Tags User-Search
// @Accept json
// @Produce  json
// @Param query query string false "Search query"
// @Param sort query string false "Sorting criteria: price_low_to_high, price_high_to_low, new_arrivals, category, a_to_z, z_to_a, popularity"
// @Success 200 {json} json	"List of products"
// @Failure 401 {json} json "failed to fetch products"
// @Router /user/search [get]
func SearchProduct(c *gin.Context) {
	var listproduct []gin.H
	var listitems []gin.H
	searchQuery := c.Query("query")
	sortBy := strings.ToLower(c.DefaultQuery("sort", "a_to_z"))

	query := initializers.DB
	if searchQuery != "" {
		query = query.Where("product_name ILIKE ?", "%"+searchQuery+"%")
	}

	switch sortBy {
	case "price_low_to_high":
		query = query.Order("price asc")
	case "price_high_to_low":
		query = query.Order("price desc")
	case "new_arrivals":
		query = query.Order("created_at desc")
	case "category":
		query = query.Order("category_id")
	case "a_to_z":
		query = query.Order("product_name asc")
	case "z_to_a":
		query = query.Order("product_name desc")
	case "popularity":
		var products []models.Product
		query := `SELECT * FROM products
                JOIN (
                    SELECT
						product_id,
                        SUM(order_quantity) as total_quantity
                    FROM
                        orders
                    GROUP BY
                        product_id
                    ORDER BY
                        total_quantity DESC
                    LIMIT 10
                ) AS subq ON products.id = subq.product_id
                WHERE
                    products.deleted_at IS NULL
                ORDER BY
                    subq.total_quantity DESC`
		initializers.DB.Raw(query).Scan(&products)

		for _, v := range products {
			listproduct = append(listproduct, gin.H{
				"Name":  v.Product_Name,
				"Price": v.Price,
				"ID":    v.ID,
			})
		}
		c.JSON(200,gin.H{
			"data": listproduct,
			"status":200,
		})
		return
	default:
		query = query.Order("product_name asc")
	}
	var items []models.Product
	query.Joins("Category").Find(&items)

	for _, v := range items {
		listitems = append(listitems,gin.H{
			"Name":     v.Product_Name,
			"Price":    v.Price,
			"Category": v.Category.Name,
			"ID":       v.ID,
		})
	}
	c.JSON(200,gin.H{
		"data":   listitems,
		"status": 200,
	})
}
