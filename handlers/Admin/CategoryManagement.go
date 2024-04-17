package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

//======================= Category Adding to the DB ================================
type  addCategoryForm struct {
	Name        string `json:"categoryName"`
    Description string `json:"catDescription"`
}
// @Summary  Add Category
// @Description Admin  can add a new category
// @Tags Admin-CategoryManagement
// @Accept json
// @Produce  json
// @Param addCategoryForm string true "New Category  Info"
// @Success 200 {json} json	"Added Category"
// @Failure 401 {json} json "Failed to add Category"
// @Router /admin/category/addcategory [post]
func Category(c *gin.Context) {
	var cat addCategoryForm
	c.ShouldBindJSON(&cat)
	upload := initializers.DB.Create(&cat)
	if upload.Error != nil {
		c.JSON(401, gin.H{
			"error":"failed to upload category",
			"status":401,
		})
	} else {
		c.JSON(200, gin.H{
			"message":"New Category added",
			"status":200,
		})
	}
}

//==================================== END ===========================================
// @Summary  List Category
// @Description Admin  can list categories
// @Tags Admin-CategoryManagement
// @Accept json
// @Produce  json
// @Success 200 {json} json	"Added Category"
// @Failure 401 {json} json "Failed to add Category"
// @Router /admin/category [get]
func View_Category(c *gin.Context) {
	var View []models.Category
	//  var checkcategory models.Categories
	initializers.DB.Find(&View)
	c.JSON(200,gin.H{
		"data":View,
		"status":200,
	})
}

// ==================================== Editing Category ===========================================
// @Summary Edit Category 
// @Description Admin  can Edit category
// @Tags Admin-CategoryManagement
// @Accept json
// @Produce  json
// @Param  ID path string true "Category ID"
// @Success 200 {json} json	"Edited Category"
// @Failure 401 {json} json "Category not Found or Failed to Edit Category"
// @Router /admin/category/edit/{ID} [post]
func Edit_Category(c *gin.Context) {
	var edit models.Category
	id := c.Param("ID")
	result := initializers.DB.First(&edit, id)
	fmt.Println("(===============", edit, "===========)(", id, "===================)")
	if result.Error != nil {
		c.JSON(401, gin.H{
			"error":"Category not found",
			"status":401,
		})
	} else {
		err := c.ShouldBindJSON(&edit)
		if err != nil {
			c.JSON(401, gin.H{
				"error":"failed to bind json",
				"status":401,
			})
		}
		save := initializers.DB.Save(&edit)
		if save.Error != nil {
			c.JSON(401, gin.H{
				"error":"Failed to update Category",
				"status":401,
			})
		} else {
			c.JSON(200, gin.H{
				"message":"Category updated successfully",
				"status":200,
			})
		}
	}
}

//==================================== END ===========================================
//==================================== Deleting Categories ===========================================
// @Summary Delete Category 
// @Description Admin  can Delete category by selecting the category
// @Tags Admin-CategoryManagement
// @Accept json
// @Produce  json
// @Param  ID path string true "Category ID"
// @Success 200 {json} json	"Deleted Category"
// @Failure 401 {json} json "Category not Found or Failed to Delete Category"
// @Router /admin/category/delete/{ID} [delete]
func Delete_Category(c *gin.Context) {
	var delete models.Category	
	cat := c.Param("ID")
	err := initializers.DB.First(&delete, cat)
	if err.Error != nil {
		c.JSON(401,gin.H{
			"error":"Category cant be deleted",
			"status":401,
		} )
	} else {
		initializers.DB.Delete(&delete)
		c.JSON(200, gin.H{
			"message":"Category Deleted",
			"status":200,
		})
	}
}

//==================================== END ===========================================
