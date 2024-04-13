package controllers

import (
	"context"
	"ecom/initializers"
	"ecom/middleware"
	"ecom/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func OauthSetup() *oauth2.Config {
	clientid := os.Getenv("CLIENTID")
	clientsecret := os.Getenv("CLIENTSECRET")

	conf := &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientsecret,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}

func Googlelogin(c *gin.Context) {
	var googleConfig *oauth2.Config
	googleConfig = OauthSetup()
	url := googleConfig.AuthCodeURL("state")
	c.Redirect(http.StatusFound, url)
	// fmt.Println("url-------->", url)
}

func GoogleCallback(c *gin.Context) {
	var userdetails models.User
	var jwttoken string
	code := c.Request.URL.Query().Get("code")
	googleConfig := OauthSetup()
	fmt.Println("code---------->", code)
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Please try again-------->")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	client := googleConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Fatalf("GET /userinfo error %v", resp)
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var user GUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Couldn't decode response from Google: %v", err)
		return
	}
	guser := models.User{
		Email: user.Email,
		Name:  user.Name,
	}
	if user.VerifiedEmail {
		if err := initializers.DB.Where("email = ?", user.Email).First(&userdetails); err.Error != nil {
			if err := initializers.DB.Create(&guser); err.Error != nil {
				fmt.Println("-----failed to signup using OAuth-----", err.Error)
			} else {
				fmt.Println("------user data updated to the database------")
				initializers.DB.First(&guser, "email = ?", guser.Email)
			}
		}
		// } else {
		// 	fmt.Println("------welcome back to home------")
		// 	c.Redirect(303, "/")
		// }
		jwttoken = middleware.GenerateJwt(c, guser.Email, Roleuser, guser.ID)
		fmt.Println("jwttoken------->", jwttoken)
		c.SetCookie("jwtTokenUser", jwttoken, int((time.Hour * 5).Seconds()), "/", "localhost", false, true)
		c.Redirect(303, "/")
	}
	fmt.Println("user ------------", user)
}
