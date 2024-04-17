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

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// var Signup models.User
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
			c.SetCookie("jwtTokenUser", token, int((time.Hour * 5).Seconds()), "/", "abdin.online", false, false)
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
	Name     string `json:"userName"`
	Email    string `json:"userEmail"`
	Mobile   string `json:"Mob"`
	Password string `json:"userPassword"`
	Gender   string `json:"gender"`
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
	var Otp string
	er := c.ShouldBindJSON(&Signup)
	if er != nil {
		c.JSON(401, gin.H{
			"error":  "failed to bind json",
			"status": 401,
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(Signup.Password), 10)
	if err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to hash password",
			"status": 401,
		})
	}
	Signup.Password = string(hash)
	Otp = GenerateOtp()
	check.Otp = Otp
	fmt.Println("otp------->", Otp, "email----->", Signup.Email)
	err = SendOtp(Signup.Email, Otp)
	if err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to sent otp",
			"status": 401,
		})
		return
	}
	fmt.Println("email-------------->", Signup.Email)
	result := initializers.DB.First(&check, "email=?", Signup.Email)
	if result.Error != nil {
		check = models.Otp{
			Email:     Signup.Email,
			Otp:       Otp,
			Expire_at: time.Now().Add(2 * time.Minute),
		}
		fmt.Println("check--------->", check.Email, check.Otp)
		initializers.DB.Create(&check)
	} else {
		initializers.DB.Model(&check).Where("email=?", Signup.Email).Updates(models.Otp{
			Otp:       Otp,
			Expire_at: time.Now().Add(60 * time.Second),
		})
	}
	userDetails := map[string]interface{}{
		"name":     Signup.Name,
		"email":    Signup.Email,
		"password": Signup.Password,
		"phone":    Signup.Mobile,
		"gender":   Signup.Gender,
	}
	session := sessions.Default(c)
	session.Set("user"+Signup.Email, userDetails)
	session.Save()
	c.SetCookie("sessionID", "user"+Signup.Email, int((time.Hour * 5).Seconds()), "/", "abdin.online", false, false)
	c.JSON(200, gin.H{
		"message": "OTP sent to your mail: " + Otp,
		"status":  200,
	})

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
// @Router /user/registration/otp [post]Second
func Otpcheck(c *gin.Context) {
	var signupData models.User
	// var check models.Otp
	var userotp otpvalidation
	var existinigOtp models.Otp
	var wallet models.Wallet
	var userid models.User
	c.ShouldBindJSON(&userotp)
	cookie, _ := c.Cookie("sessionID")
	session := sessions.Default(c)
	userData := session.Get(cookie)
	// initializers.DB.First(&check, "email=?", Signup.Email)
	// fmt.Println("=======(", check.Otp, ")=========(", userotp.Otp, ")=========", "(", Signup.Email, ")=========")
	value := initializers.DB.Where("otp=? AND expire_at > ?", userotp.Otp, time.Now()).First(&existinigOtp)
	if value.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Incorrect OTP or OTP expired",
			"status": 401,
		})
		return
	}
	userMap := userData.(map[string]interface{})

	phn := fmt.Sprintf("%v", userMap["phone"])
	signupData = models.User{
		Name:     userMap["name"].(string),
		Email:    userMap["email"].(string),
		Password: userMap["password"].(string),
		Mobile:   phn,
		Gender:   userMap["gender"].(string),
	}
	fmt.Println("======>signupdata", signupData, "<========end")
	result := initializers.DB.Create(&signupData)
	if result.Error != nil {
		c.JSON(401, gin.H{
			"error":  "User already exist",
			"status": 401,
		})
		return
	} else {
		initializers.DB.First(&userid, "email = ?", signupData.Email)
		wallet.Created_at = time.Now()
		wallet.UserID = userid.ID
		if err := initializers.DB.Create(&wallet); err.Error != nil {
			c.JSON(500, "Failed to create wallet")
			fmt.Println("Failed to create wallet====>", err.Error)
			return
		}

		session.Delete(cookie)
		session.Save()
		c.SetCookie("sessionID", "", -1, "", "", false, false)
		c.JSON(200, "Successfully signed up")
	}
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
	// var Signup signupdata
	var Otp string
	var check models.Otp
	cookie, _ := c.Cookie("sessionID")
	session := sessions.Default(c)
	userData := session.Get(cookie)
	if userData == nil {
		c.JSON(401, gin.H{
			"error":  "User data not found in session",
			"status": 401,
		})
		return
	}
	userMap, ok := userData.(map[string]interface{})
	if !ok {
		c.JSON(401, gin.H{
			"error":  "Failed to convert user data to map",
			"status": 401,
		})
		return
	}
	// email := userMap["email"].(string)
	Otp = GenerateOtp()
	err := SendOtp(userMap["email"].(string), Otp)
	if err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to sent otp",
			"status": 401,
		})
	} else {

		result := initializers.DB.First(&check, "email=?", userMap["email"].(string))
		if result.Error != nil {

			check = models.Otp{
				Email:     userMap["email"].(string),
				Otp:       Otp,
				Expire_at: time.Now().Add(2 * time.Minute),
			}

			result := initializers.DB.Create(&check)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, "Failed to save OTP")
			}
		} else {
			err := initializers.DB.Model(&check).Where("email=?", userMap["email"].(string)).Updates(models.Otp{
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
