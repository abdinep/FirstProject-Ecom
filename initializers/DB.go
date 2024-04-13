package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
//=========================== Connecting to DB ==================================
var DB *gorm.DB

func Dbinit() {
	var err error

	dsn := os.Getenv("DSN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB")
	}
	fmt.Println("============================ CONNECTED TO DB =====================================")
}
//========================== END ============================================
