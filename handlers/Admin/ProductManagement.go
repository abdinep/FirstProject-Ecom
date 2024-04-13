package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var product ProductHandler

type ProductHandler struct {
	Product_Name string `json:"prodName"`
	Price        int    `json:"price"`
	Quantity     int    `json:"quantity"`
	Size         int    `json:"size"`
	Category_id  uint   `json:"category"`
	Description  string `json:"description"`
	ImagePath1   string
	ImagePath2   string
	ImagePath3   string
}

// @Summary Add Product
// @Description Collecting the Product Details
// @Tags Admin-ProductManagement
// @Accept json
// @Produce json
// @Param product body ProductHandler true "Collecting Product Details"
// @Success 200 {json} json	"Upload Product Images"
// @Failure 401 {json} json "Failed to collect  data, please check your inputs again."
// @Router /admin/products [get]
func Add_Product(c *gin.Context) {
	var cat_id models.Category
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(501, "failed to bind json")
	}
	fmt.Println("============", product, "==================")
	result := initializers.DB.First(&cat_id, product.Category_id)
	fmt.Println("============", cat_id, "=====================")
	if result.Error != nil {
		c.JSON(404, "Category not found")
	} else {

		c.JSON(200, "Upload Product Images ")
	}
}

// ==================================== END =========================================
// ================================= Upload Product Image ===========================
// @Summary Add Product Images
// @Description Collecting all Images of the Product
// @Tags Admin-ProductManagement
// @Accept mulipart/form-data
// @Produce json
// @Param images formData file true "Product images to upload"
// @Success 200 {json} json	"Product added successfully"
// @Failure 401 {json} json "Product already exist"
// @Router /admin/products [post]
func ProductImage(c *gin.Context) {
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to fetch images")
	}
	files := file.File["images"]
	var imagepaths []string

	for _, val := range files {
		filepath := "./images/" + val.Filename
		if err = c.SaveUploadedFile(val, filepath); err != nil {
			c.JSON(http.StatusBadRequest, "Failed to save images")
			return
		}
		imagepaths = append(imagepaths, filepath)
	}
	product.ImagePath1 = imagepaths[0]
	product.ImagePath2 = imagepaths[1]
	product.ImagePath3 = imagepaths[2]
	upload := initializers.DB.Create(&product)
	if upload.Error != nil {
		c.JSON(501, "Product already exist")
		return
	} else {
		c.JSON(200, "Product added successfully")
	}
	product = ProductHandler{}

}

// ====================== Showing all products in admin page ==========================

// @Summary List Products
// @Description Listing All Products
// @Tags Admin-ProductManagement
// @Accept mulipart/form-data
// @Produce json
// @Success 200 {json} json	"Listed All Products"
// @Failure 401 {json} json "Failed To Fetch Products"
// @Router /admin/productlist [get]
func View_Product(c *gin.Context) {
	var View []models.Product
	//  var checkcategory models.Categories
	initializers.DB.Where("deleted_at IS NULL").Preload("Category").Find(&View)
	c.JSON(200, gin.H{
		"data": View,
	})
}

// =========================== Editing product detailes in admin panel =========================
type UpdateProduct struct {
	Product_Name string `json:"prodName"`
	Price        int    `json:"price"`
	Quantity     int    `json:"quantity"`
	Size         int    `json:"size"`
	Category_id  uint   `json:"category"`
	Description  string `json:"description"`
	ImagePath1   string
	ImagePath2   string
	ImagePath3   string
}

// @Summary Edit Product Details
// @Description Editing all  details of a particular product
// @Tags Admin-ProductManagement
// @Accept json
// @Produce json
// @Param id path int true "Prouct ID"
// @Param  data body UpdateProduct true "Edit Product Data"
// @Success 200 {json} json	"Product updated successfully"
// @Failure 401 {json} json "Failed to update Product"
// @Router /admin/products/{ID} [patch]
func Edit_Product(c *gin.Context) {
	var edit UpdateProduct
	product := c.Param("ID")
	result := initializers.DB.First(&edit, product)
	fmt.Println("(===============", edit, "===========)(", product, "===================)")
	if result.Error != nil {
		c.JSON(501, "Product not found")
		return
	} else {
		err := c.ShouldBindJSON(&edit)
		if err != nil {
			c.JSON(501, "failed to bind json")
			return
		}
		save := initializers.DB.Save(&edit)
		if save.Error != nil {
			c.JSON(501, "Failed to update Product")
			return
		} else {
			c.JSON(200, "Product updated successfully")
		}
	}
}

//========================= END ================================

// =================== Deleting products in admin panel =================================

// @Summary Delete Products
// @Description Deleting  a specific product from the database
// @Tags Admin-ProductManagement
// @Accept json
// @Produce json
// @Param id path int true "Prouct ID"
// @Success 200 {json} json	"Product Deleted successfully"
// @Failure 401 {json} json "Product cant be deleted"
// @Router /admin/products/{ID} [delete]
func Delete_Product(c *gin.Context) {
	var delete models.Product
	product := c.Param("ID")
	if err := initializers.DB.First(&delete, product);
	err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Product cant be deleted",
			"status": 401,
		})
	} else {
		initializers.DB.Delete(&delete)
		c.JSON(200, gin.H{
			"message": "Product Deleted",
			"status":  200,
		})
	}
}
