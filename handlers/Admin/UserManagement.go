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
	var listUser []gin.H
	initializers.DB.Find(&list)
	count := 0
	for _, user := range list {
		listUser = append(listUser, gin.H{
			"ID":        user.ID,
			"CreatedAt": user.CreatedAt,
			"userName":  user.Name,
			"userEmail": user.Email,
			"Mob":       user.Mobile,
			"gender":    user.Gender,
			"status":    user.Status,
		})
		count++
	}
	c.JSON(200, gin.H{
		"data":       listUser,
		"totalUsers": count,
		"status":     200,
	})
}

type editUser struct {
	Name   string `json:"userName"`
	Email  string `son:"userEmail"`
	Mobile string `json:"Mob"`
	Gender string `json:"gender"`
}

// @Summary Edit User
// @Description This will edit  a particular user's information
// @Tags Admin-UserManagement
// @Accept json
// @Produce json
// @Param ID path string true "User ID"
// @Param user body editUser true "User Details"
// @Success 200 {json} json	"User updated successfully"
// @Failure 401 {json} json "failed to update user"
// @Router /admin/usermanagement/{ID} [patch]
func Edit_User(c *gin.Context) {
	var edit editUser
	var userDetails models.User
	user := c.Param("ID")
	result := initializers.DB.First(&userDetails, user)
	fmt.Println("(===============", userDetails, "===========)(", user, "===================)")
	if result.Error != nil {
		c.JSON(401, gin.H{
			"error":  "User not found",
			"status": 401,
		})
		return
	}

	err := c.ShouldBindJSON(&edit)
	if err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to bind json",
			"status": 401,
		})
		return
	}
	userDetails.Email = edit.Email
	userDetails.Name = edit.Name
	userDetails.Mobile = edit.Mobile
	userDetails.Gender = edit.Gender
	error := initializers.DB.Save(&userDetails)
	if error.Error != nil {
		c.JSON(401, gin.H{
			"error":  "failed to update user",
			"status": 401,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "User updated successfully",
		"status":  200,
	})
}

// =========================== User Block/Unblock in admin panel ===========================
// @Summary Block/Unblock User
// @Description This will block  or unblock a particular user from the system
// @Tags Admin-UserManagement
// @Accept json
// @Produce json
// @Param ID path string true "User ID"
// @Success 200 {json} json	"User Status Updated"
// @Failure 401 {json} json "failed to update user status"
// @Router /admin/usermanagement/block/{ID} [patch]
func Status(c *gin.Context) {
	var check models.User
	user := c.Param("ID")
	initializers.DB.First(&check, user)
	if check.Status == "Active" {
		initializers.DB.Model(&check).Update("status", "Blocked")
		c.JSON(200, gin.H{
			"message": "User Blocked",
			"status":  200,
		})
	} else {
		initializers.DB.Model(&check).Update("status", "Active")
		c.JSON(200, gin.H{
			"message": "User Unblocked",
			"status":  200,
		})
	}

}
