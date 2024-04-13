package controllers

import (
	"ecom/initializers"
	"ecom/middleware"
	"ecom/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	// "github.com/dgrijalva/jwt-go"
	// "main/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// var Signup models.User
var Otp string
var Roleuser = "User"

type userlogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
// ============================== User Authentication =============================================

// @Summary Login as an User
// @Description Logs in the User with email and password, returns a token to be used for authentication
// @Tags User-Auth
// @Accept json
// @Produce  json
// @Param credentials body userlogin true "Login Data"
// @Success 200 {json} json "Login Successful".
// @Failure 401 {json} json "Unauthorized."
// @Router /user/signin [post]
func Userlogin(c *gin.Context) {
	var form userlogin
	var table models.User
	var token string
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(501, "failed to bind json")
	}
	initializers.DB.First(&table, "email=?", strings.ToLower(form.Email))
	fmt.Println("(=======================", table, ")(====================", form.Email, "==============)")

	err = bcrypt.CompareHashAndPassword([]byte(table.Password), []byte(form.Password))
	if err != nil {
		c.JSON(501, "invalid user name or password")
	} else {
		if table.Status == "Active" {
			fmt.Println("id======>", table.ID)
			token = middleware.GenerateJwt(c, form.Email, Roleuser, table.ID)
			fmt.Println("token----->", token)
			c.SetCookie("jwtTokenUser", token, int((time.Hour * 5).Seconds()), "/", "localhost", false, true)
			c.JSON(200, gin.H{
				"Message": "Welcome to Home page",
				"Token":   token,
			})
		} else {
			c.JSON(200, "User Blocked")
		}
	}
}

// =============================== END ===============================================
// @Summary Logout as an User
// @Description Logs out as an user, clear the token used for authentication
// @Tags User-Auth
// @Accept json
// @Produce  json
// @Success 200 {json} json "Logout Successful".
// @Failure 401 {json} json "Unauthorized."
// @Router /user/logout [get]
func User_Logout(c *gin.Context) {

	c.SetCookie("jwtTokenUser", "", -1, "", "", false, false)
	c.JSON(200, gin.H{"message": "Logout succesful"})
}

// ========================= Sending OTP by clicking Signup =========================
type signupdata struct {
	Name     string
	Email    string
	Mobile   string
	Password string
	Gender   string
}

// @Summary User Signup
// @Description Collecting the New User data and Send  an OTP to verify Email ID
// @Tags User-Auth
// @Accept json
// @Produce  json
// @Param request body signupdata true "SignUp Request Body"
// @Success 200 {json} json " Successfuly sent OTP".
// @Failure 401 {json} json "Failed to sent OTP."
// @Router /user/registration [post]
func Usersignup(c *gin.Context) {
	var check models.Otp
	var Signup signupdata
	er := c.ShouldBindJSON(&Signup)
	if er != nil {
		c.JSON(501, "failed to bind json")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(Signup.Password), 10)
	if err != nil {
		c.JSON(501, "Failed to hash password")
	}
	Signup.Password = string(hash)
	Otp = GenerateOtp()
	check.Otp = Otp
	err = SendOtp(Signup.Email, Otp)
	if err != nil {
		c.JSON(501, "Failed to sent otp")
	}
	result := initializers.DB.First(&check, "email=?", Signup.Email)
	if result.Error != nil {

		check = models.Otp{
			Email:     Signup.Email,
			Otp:       Otp,
			Expire_at: time.Now().Add(60 * time.Second),
		}

		initializers.DB.Create(&check)
	} else {
		initializers.DB.Model(&check).Where("email=?", Signup.Email).Updates(models.Otp{
			Otp:       Otp,
			Expire_at: time.Now().Add(60 * time.Second),
		})
	}
	// initializers.DB.Delete(&check)
	c.JSON(200, "OTP sent to your mail: "+Otp)

}

//================================== END ======================================

// ========================== OTP validation and Signup =================================
type otpvalidation struct {
	Otp string `json:"otp"`
}

// @Summary OTP Check
// @Description Validating OTP
// @Tags User-Auth
// @Accept json
// @Produce  json
// @Param request body otpvalidation true "Otp check"
// @Success 200 {json} json " Successfully signed up"
// @Failure 401 {json} json "Failed to Signup"
// @Router /user/registration/otp [post]
func Otpcheck(c *gin.Context) {
	var Signup signupdata
	var check models.Otp
	var userotp otpvalidation
	var existinigOtp models.Otp
	var wallet models.Wallet
	var userid models.User
	c.ShouldBindJSON(&userotp)
	initializers.DB.First(&check, "email=?", Signup.Email)
	fmt.Println("=======(", check.Otp, ")=========(", userotp.Otp, ")=========", "(", Signup.Email, ")=========")
	value := initializers.DB.Where("otp=? AND expire_at > ?", userotp.Otp, time.Now()).First(&existinigOtp)
	if value.Error != nil {
		c.JSON(501, "Incorrect OTP or OTP expired")
	} else {
		result := initializers.DB.Create(&Signup)
		if result.Error != nil {
			c.JSON(501, "User already exist")
			return
		} else {
			initializers.DB.First(&userid, "email = ?", Signup.Email)
			wallet.Created_at = time.Now()
			wallet.UserID = userid.ID
			if err := initializers.DB.Create(&wallet); err.Error != nil {
				c.JSON(500, "Failed to create wallet")
				fmt.Println("Failed to create wallet====>", err.Error)
				return
			}
			c.JSON(200, "Successfully signed up")
		}
	}
	Signup = signupdata{}
}

// @Summary Resend OTP 
// @Description  This API is used for resending the OTP
// @Tags User-Auth
// @Accept json
// @Produce  json
// @Success 200 {json} json " Resent OTP"
// @Failure 401 {json} json "Failed to send OTP"
// @Router /user/registration/resendotp [post]
func Resend_Otp(c *gin.Context) {
	var Signup signupdata
	var check models.Otp
	Otp = GenerateOtp()
	err := SendOtp(Signup.Email, Otp)
	if err != nil {
		c.JSON(501, "Failed to sent otp")
	} else {

		result := initializers.DB.First(&check, "email=?", Signup.Email)
		if result.Error != nil {

			check = models.Otp{
				Email:     Signup.Email,
				Otp:       Otp,
				Expire_at: time.Now().Add(15 * time.Second),
			}

			result := initializers.DB.Create(&check)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, "Failed to save OTP")
			}
		} else {
			err := initializers.DB.Model(&check).Where("email=?", Signup.Email).Updates(models.Otp{
				Otp:       Otp,
				Expire_at: time.Now().Add(15 * time.Second),
			})
			if err.Error != nil {
				c.JSON(http.StatusBadRequest, "Failed to update data")
			}
		}
		c.JSON(200, "OTP sent to your mail: "+Otp)
	}

}

//============================= END =====================================
