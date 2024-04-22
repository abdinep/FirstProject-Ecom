package handlers

import (
	"ecom/initializers"
	"ecom/middleware"
	"ecom/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Get user details
// @Description Get user details
// @Tags User-Profile
// @Produce json
// @Success 200 {json} json "listed Profile Details"
// @Failure 401 {json} json "Failed to list Profile Details"
// @Router /user/profile [get]
func User_Details(c *gin.Context) {
	var details models.User
	var detailsList []gin.H
	// id := c.Param("ID")
	if err := initializers.DB.First(&details, middleware.UserDetails.ID); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch data",
			"status": 401,
		})
		fmt.Println("Error", err.Error)
	} else {
		detailsList = append(detailsList, gin.H{
			"user_name":   details.Name,
			"user_email":  details.Email,
			"user_mobile": details.Mobile,
			"user_gender": details.Gender,
			"user_status": details.Status,
		})
		c.JSON(200, gin.H{
			"data":   detailsList,
			"status": 200,
		})
	}
}

type profileEdit struct {
	Name     string `json:"userName"`
	Email    string `json:"userEmail"`
	Mobile   string `json:"Mob"`
	Password string `json:"userPassword"`
	Gender   string `json:"gender"`
}

// @Summary Edit user details
// @Description Edit user details
// @Tags User-Profile
// @Accept json
// @Produce json
// @Param  data body profileEdit true "Edit User Profile"
// @Success 200 {json} json "Edited User Details"
// @Failure 401 {json} json "Failed to edit User Details"
// @Router /user/profile [patch]
func Edit_Profile(c *gin.Context) {
	var edit models.User
	var editdetails profileEdit
	id := c.GetUint("userID")
	if err := initializers.DB.First(&edit, id); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch data from DB",
			"status": 401,
		})
		fmt.Println("Failed to fetch data from DB=====>", err.Error)
	} else {
		if err := c.ShouldBindJSON(&editdetails); err != nil {
			c.JSON(401, gin.H{
				"error":  "failed to bind profile details",
				"status": 401,
			})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(editdetails.Password), 10)
		if err != nil {
			c.JSON(401, gin.H{
				"error":  "Failed to hash password",
				"status": 401,
			})
		}
		edit.Email = editdetails.Email
		edit.Name = editdetails.Name
		edit.Mobile = editdetails.Mobile
		edit.Password = string(hash)
		edit.Gender = editdetails.Gender
		if err := initializers.DB.Save(&edit); err.Error != nil {
			c.JSON(401, gin.H{
				"error":  "failed to edit details",
				"status": 401,
			})
			fmt.Println("failed to edit details", err.Error)
		} else {
			c.JSON(200, gin.H{
				"message": "Updated Profile details",
				"status":  200,
			})
		}
	}

}
