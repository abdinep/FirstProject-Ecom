package main

import (
	"ecom/initializers"
	"ecom/routers"
	"os"

	_ "ecom/docs"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title swagger documentation for the application
// @version 1.0
// @description ecom application
// @host localhost:8080
// @BasePath /
func init() {
	initializers.LoadEnvVariables()
	initializers.Dbinit()
	initializers.TableCreate()
	// controllers.GenerateOTP()
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user := server.Group("/")
	routers.UserGroup(user)

	admin := server.Group("/admin")
	routers.AdminGroup(admin)

	server.Run(os.Getenv("PORT"))

}
