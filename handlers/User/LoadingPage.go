package handlers

import (
	handlers "ecom/handlers/Admin"
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Landing Page
// @Description Load products along with their categories and images for display on a webpage
// @Tags User-LandingPage
// @Accept json
// @Produce json
// @Success 200 {json} JSON "Products loaded successfully"
// @Failure 401 {json} JSON "Failed to fetch product data"
// @Router / [get]
func ProductLoadingPage(c *gin.Context) {
	var load []models.Product
	var loads []gin.H
	if err := initializers.DB.Joins("Category").Find(&load); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to fetch product data",
			"status": 400,
		})
		return
	}
	for _, v := range load {
		loads = append(loads, gin.H{
			"product_name":     v.Product_Name,
			"product_price":    v.Price,
			"product_category": v.Category.Name,
			"product_image":    v.ImagePath1,
			"product_id":       v.ID,
		})
	}
	c.JSON(200, gin.H{
		"data":   loads,
		"status": 200,
	})
}

// @Summary Product Details
// @Description Load products along with their details and recommended products
// @Tags User-LandingPage
// @Produce json
// @Param ID path string true "Product ID"
// @Success 200 {json} JSON "Products details loaded successfully"
// @Failure 401 {json} JSON "Failed to fetch product data"
// @Router /{ID} [get]
func ProductDetails(c *gin.Context) {
	var product models.Product
	var reviews []gin.H
	var load []models.Product
	// var productdetails []gin.H
	productID := c.Param("ID")
	if err := initializers.DB.Joins("Category").First(&product, productID); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "failed to fetch product data",
			"status": 401,
		})
		return
	}
	if err := initializers.DB.Joins("Category").Where("category_id = ? AND products.id != ?", product.Category_id, productID).Find(&load); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "failed to fetch similiar product data",
			"status": 401,
		})
	}

	result := handlers.OfferCalc(int(product.ID), c)
	fmt.Println("---->", result)
	rating := Ratingcalc(productID, c)
	fmt.Println("rating--------->",rating)
	stock := ""
	if product.Quantity == 0 {
		stock = "Out of stock"
	} else {
		stock = "Item is available"
	}
	review := ReviewView(productID, c)
	for _, values := range review {
		if values.ProductId == product.ID {
			reviews = append(reviews, gin.H{
				"ID":      values.ID,
				"product": values.Product.Product_Name,
				"review":  values.Review,
				"time":    values.Time,
				"user":    values.User.Name,
			})
		}
	}

	var similiarproducts []gin.H
	var results float64

	for _, value := range load {
		similiarproducts = append(similiarproducts, gin.H{
			"message":             "similar products",
			"category_id":         value.Category_id,
			"product Name":        value.Product_Name,
			"product image":       value.ImagePath1,
			"product Price":       value.Price,
			"product size":        value.Size,
			"product discription": value.Description,
			"category":            value.Category.Name,
		})
		fmt.Println("similar------------>", similiarproducts)
		results = handlers.OfferCalc(int(value.ID), c)
		fmt.Println("---->", results)
	}
	c.JSON(200, gin.H{
		"product_name":        product.Product_Name,
		"product_quantity":    product.Quantity,
		"product_price":       product.Price,
		"product_description": product.Description,
		"product_image2":      product.ImagePath2,
		"product_image3":      product.ImagePath3,
		"product_size":        product.Size,
		"product_category":    product.Category.Name,
		"rating":            rating,
		"stock":             stock,
		"review":            reviews,
		"similiar_products": similiarproducts,
		"offer_price":       results,
		"status":            200,
	})
}

//================================ END ================================================
