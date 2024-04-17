package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

// @Summary Add a new product
// @Description Add a new product with images and other details
// @Tags Admin-ProductManagement
// @Accept multipart/form-data
// @Param prodName formData string true "Product name"
// @Param price formData integer true "Product price"
// @Param size formData string true "Product size"
// @Param quantity formData integer true "Product quantity"
// @Param description formData string true "Product description"
// @Param category formData  int true "Category ID of the product"
// @Param images formData []file true "Product images"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/products [post]
func Add_Product(c *gin.Context) {
	var cat_id models.Category

	description := c.PostForm("description")
	price, _ := strconv.Atoi(c.PostForm("price"))
	prodName := c.PostForm("productName")
	quantity, _ := strconv.Atoi(c.PostForm("quantity"))
	size, _ := strconv.Atoi(c.PostForm("size"))
	fmt.Println("============", product, "==================")
	result := initializers.DB.First(&cat_id, c.PostForm("category"))
	fmt.Println("============", cat_id, "=====================")
	if result.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Category not found",
			"status": 401,
		})
		return
	}
	files := c.Request.MultipartForm.File["images"]
	if len(files) < 3 {
		c.JSON(401, gin.H{
			"error":  "atleast 3 images required",
			"status": 401,
		})
		return
	}
	var images []string
	for _, image := range files {
		path := filepath.Join("./assets", image.Filename)
		if err := c.SaveUploadedFile(image, path); err != nil {
			c.JSON(401, gin.H{
				"error":  "Failed to upload images",
				"status": 401,
			})
			return
		}
		images = append(images, path)
	}
	datas := models.Product{
		Product_Name: prodName,
		Price:        price,
		Description:  description,
		Category_id:  cat_id.ID,
		ImagePath:    pq.StringArray(images),
		Quantity:     quantity,
		Size:         size,
	}
	if err := initializers.DB.Create(&datas); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "failed to upload product",
			"status": 401,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Product added successfully",
		"status":  200,
	})
}

// ==================================== END =========================================
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
	var listproduct []gin.H
	//  var checkcategory models.Categories
	initializers.DB.Find(&View)
	for _, v := range View {
		listproduct = append(listproduct, gin.H{
			"Product_Name": v.Product_Name,
			"Price":        v.Price,
			"Quantity":     v.Quantity,
			"Size":         v.Size,
		})
	}
	c.JSON(200, gin.H{
		"data": listproduct,
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
		c.JSON(401, gin.H{
			"error":  "Product not found",
			"status": 401,
		})
		return
	} else {
		err := c.ShouldBindJSON(&edit)
		if err != nil {
			c.JSON(401, gin.H{
				"error":  "failed to bind json",
				"status": 401,
			})
			return
		}
		save := initializers.DB.Save(&edit)
		if save.Error != nil {
			c.JSON(401, gin.H{
				"error":  "Failed to update Product",
				"status": 401,
			})
			return
		} else {
			c.JSON(200, gin.H{
				"message": "Product updated successfully",
				"status":  200,
			})
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
	if err := initializers.DB.First(&delete, product); err.Error != nil {
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
