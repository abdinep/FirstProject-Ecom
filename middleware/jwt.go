package middleware

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var SecretKey = []byte("secretkey")
var UserDetails models.User

type customClaims struct {
	UserID   uint   `json:"userid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateJwt(c *gin.Context, mail string, role string, userid uint) string {
	fmt.Println("mail===>", mail, "role===>", role, "userid===>", userid)
	tokenkey, err := CreateToken(mail, role, userid)
	if err != nil {
		fmt.Println("failed to create token")
		return ""
	}
	fmt.Println("session check======>", tokenkey)
	return tokenkey

}
func CreateToken(mail string, role string, userid uint) (string, error) {
	claims := customClaims{
		Email:  mail,
		Role:   role,
		UserID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenkey, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenkey, nil
}
func JwtMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		UserDetails = models.User{}
		tokenString,_ := c.Cookie("jwtToken"+role)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login First"})
			c.Abort()
			return
		}
		claim := &customClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			fmt.Println("invalid token=========>", err)
			c.Abort()
			return
		}
		if claim.Role == "User" {
			fmt.Println("email =====>", claim.Email)
			if err := initializers.DB.First(&UserDetails, "email=?", claim.Email); err.Error != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to fetch user details"})
				c.Abort()
				return
			}
		}
		if claim.Role != role {
			c.JSON(403, gin.H{"error": "You have No Access"})
			c.Abort()
			return
		}
		c.Set("userID", claim.UserID)
		c.Next()
	}
}
