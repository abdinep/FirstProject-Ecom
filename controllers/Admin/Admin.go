package controllers

import (
	"ecom/initializers"
	"ecom/middleware"
	"ecom/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// =============================== Admin login & logout ========================
var Adminrole = "Admin"

type adminlogin struct {
	Email    string `json:"adminMail"`
	Password string `json:"adminPassword"`
}

// @Summary Login as an admin
// @Description Logs in the Admin with email and password, returns a token to be used for authentication
// @Tags Admin-Auth
// @Accept json
// @Produce  json
// @Param credentials body adminlogin true "Login Data"
// @Success 200 {json} json "Login Successful".
// @Failure 401 {json} json "Unauthorized."
// @Router /admin/login [post]
func Login(c *gin.Context) {
	var log adminlogin
	var admin models.Admin
	err := c.ShouldBindJSON(&log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	initializers.DB.First(&admin)

	if log.Email == admin.Email && log.Password == admin.Password {
		adminID := admin.ID
		fmt.Println("==========>", admin.Email, admin.Password, adminID, "<=============")
		token := middleware.GenerateJwt(c, log.Email, Adminrole, adminID)
		c.SetCookie("jwtTokenAdmin", token, int((time.Hour * 1).Seconds()), "/", "abdin.online", false, false)
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
		})
	} else {
		c.JSON(401, gin.H{
			"message": "invalid password or Username",
		})
	}
}

// @Summary Logout as an admin
// @Description Logs out as an admin, clear the token used for authentication
// @Tags Admin-Auth
// @Accept json
// @Produce  json
// @Success 200 {json} json "Logout Successful".
// @Failure 401 {json} json "Unauthorized."
// @Router /admin/logout [get]
func Admin_Logout(c *gin.Context) {

	c.SetCookie("jwtTokenAdmin", "", -1, "/", "abdin.online", false, false)
	c.JSON(200, gin.H{"message": "Logout succesful"})
}

//============================= END =======================================
// @Summary Admin Landing Page
// @Description Admin Landing Page
// @Tags Admin-Dashboard
// @Accept json
// @Produce  json
// @Success 200 {json} json "Data llisted".
// @Failure 401 {json} json "Unauthorized."
// @Router /admin/landingPage [get]
func AdminLandingPage(c *gin.Context){
	var orders []models.OrderItem
	count :=0
	flag:=0
	var total float64
	initializers.DB.Preload("Order.User").Find(&orders)

	for _,v := range orders{
		if v.Orderstatus =="Order cancelled"{
			count++
		}else{
			flag++
		}
		total +=v.Subtotal
	}
	c.JSON(200,gin.H{
		"totalSale":total,
		"totalOrder":flag,
		"totalCancelled":count,
		"status":200,
	})
}
