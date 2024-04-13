package initializers

import (
	"log"

	"github.com/joho/godotenv"
)
//================================ Load values from .env ============================
func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
//================================== END ==============================================