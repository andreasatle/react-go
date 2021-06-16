package database

import (
	"log"

	"github.com/andreasatle/react-go/go-auth/routes/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	//connection, err := gorm.Open(mysql.Open("root:rootroot@/go_auth"), &gorm.Config{})
	connection, err := gorm.Open(mysql.Open("root@/go_auth"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
	}

	DB = connection
	connection.AutoMigrate(models.User{}, models.PasswordReset{})
}
