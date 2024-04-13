package controllers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var pass models.User
var forgotOTP string
var flag = false

func ForgotPassword_OTP(c *gin.Context) {
	var otp models.Otp
	if err := c.ShouldBindJSON(&pass); err != nil {
		c.JSON(500, "Failed to bind data")
		fmt.Println("Failed to bind data=====>", err)
	}
	if err := initializers.DB.First(&otp, "email = ?", pass.Email); err.Error != nil {
		c.JSON(500, "User Not found, Please Signup First")
		fmt.Println("User Not found, Please Signup First=====>", err.Error)
	} else {
		forgotOTP = GenerateOtp()
		if err := SendOtp(pass.Email, forgotOTP); err != nil {
			c.JSON(500, "Failed to send OTP")
			fmt.Println("Failed to send OTP=====>", err)
		}
		if err := initializers.DB.Model(&otp).Where("email = ?", pass.Email).Updates(models.Otp{
			Otp:       forgotOTP,
			Expire_at: time.Now().Add(15 * time.Second),
		}); err.Error != nil {
			c.JSON(500, "Failed to update data")
			fmt.Println("Failed to update data=====>", err.Error)
		}
		c.JSON(200, "OTP sent to your mail:"+forgotOTP)
	}
}
func Forgot_Pass_OTP_Check(c *gin.Context) {
	var check models.Otp
	var userotp models.Otp
	var existingotp models.Otp
	c.ShouldBindJSON(&userotp)
	if err := initializers.DB.First(&check, "email = ?", pass.Email); err.Error != nil {
		c.JSON(500, "Failed to fetch data")
		fmt.Println("Failed to fetch data=====>", err.Error)
		return
	}
	if err := initializers.DB.Where("otp = ? AND expire_at > ?", userotp.Otp, time.Now()).First(&existingotp); err.Error != nil {
		c.JSON(500, "Incorrect OTP or OTP expired")
		fmt.Println("Incorrect OTP or OTP expired", err.Error)
	} else {
		c.JSON(200, "Verified Email")
		flag = true
	}
}

func ForgotPassword_Change(c *gin.Context) {

	if flag {

		if err := c.ShouldBindJSON(&pass); err != nil {
			c.JSON(500, "Failed to bind data")
			fmt.Println("Failed to bind data", err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(pass.Password), 10)
		if err != nil {
			c.JSON(501, "Failed to hash password")
			fmt.Println("Failed to hash password", err)
			return
		}
		pass.Password = string(hash)
		if err := initializers.DB.Model(&pass).Where("email=?", pass.Email).Updates(models.User{Password: pass.Password}); err.Error != nil {
			c.JSON(500, "User not found Please signup first")
			fmt.Println("User not found Please signup first", err.Error)
		} else {
			c.JSON(200, "Password changed")
		}
	} else {
		c.JSON(http.StatusUnauthorized, "Please verify your Email")
	}
	flag = false
}
