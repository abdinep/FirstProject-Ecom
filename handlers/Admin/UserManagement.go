package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// ====================== Display all user details in admin panel ====================================

// @Summary List User
// @Description  This will list all the users
// @Tags Admin-UserManagement
// @Accept json
// @Produce json
// @Success 200 {json} json	"Listed all users"
// @Failure 401 {json} json "Failed to list users"
// @Router /admin/usermanagement [get]
func List_user(c *gin.Context) {
	var list []models.User
	initializers.DB.Find(&list)
	c.JSON(200,gin.H{
		"data":list,
	})
}
// @Summary Edit User
// @Description This will edit  a particular user's information
// @Tags Admin-UserManagement
// @Accept json
// @Produce json
// @Success 200 {json} json	"User updated successfully"
// @Failure 401 {json} json "failed to update user"
// @Router /admin/usermanagement/{ID} [patch]
func Edit_User(c *gin.Context) {
	var edit models.User
	user := c.Param("ID")
	result := initializers.DB.First(&edit, user)
	fmt.Println("(===============", edit, "===========)(", user, "===================)")
	if result.Error != nil {
		c.JSON(500, "User not found")
	} else {

		err := c.ShouldBindJSON(&edit)
		if err != nil {
			c.JSON(500, "Failed to bind json")
		}
		error := initializers.DB.Save(&edit)
		if error.Error != nil {
			c.JSON(500, "failed to update user")
		} else {
			c.JSON(200, "User updated successfully")
		}
	}

}

// =========================== User Block/Unblock in admin panel ===========================
// @Summary Block/Unblock User
// @Description This will block  or unblock a particular user from the system
// @Tags Admin-UserManagement
// @Accept json
// @Produce json
// @Success 200 {json} json	"User Status Updated"
// @Failure 401 {json} json "failed to update user status"
// @Router /admin/usermanagement/block/{ID} [patch]
func Status(c *gin.Context) {
	var check models.User
	user := c.Param("ID")
	initializers.DB.First(&check, user)
	if check.Status == "Active" {
		initializers.DB.Model(&check).Update("status", "Blocked")
		c.JSON(200, "User Blocked")
	} else {
		initializers.DB.Model(&check).Update("status", "Active")
		c.JSON(200, "User Unblocked")
	}

}
