package main

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Setup the logging, for if program crashes
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := gorm.Open(mysql.Open("root:rootroot@/go_basics"), &gorm.Config{})
	if err != nil {
		log.Printf("Could not connect to database: %v\n", err)
	}
	db.Migrator().DropTable(&User{})
	db.Migrator().CreateTable(&User{})
	user := User{
		FirstName: newNullString("John", true),
		LastName:  newNullString("Doe", true),
		Email:     "john.doe@email.com",
	}
	result := db.Create(&user)
	if result.Error != nil {
		log.Printf("Error creating the user: %v\n", result.Error)
	}
	fmt.Println("Created:", result.Error, result.RowsAffected, user)
	user.FirstName = newNullString("Johnny", true)

	result = db.Where("id", user.ID).Updates(&user)
	if result.Error != nil {
		log.Printf("Error updating the user: %v\n", result.Error)
	}
	fmt.Println("Updated:", result.Error, result.RowsAffected, user)

	result = db.Where("id", user.ID).Delete(&user)
	if result.Error != nil {
		log.Printf("Error deleting the user: %v\n", result.Error)
	}
	fmt.Println("Deleted:", result.Error, result.RowsAffected, user)

	users := []User{
		User{
			Email: "foo.bar@fmail.com",
		},
		User{
			FirstName: newNullString("Sven", true),
			Email:     "Sven.Ek@fmail.com",
		},
		User{
			LastName: newNullString("Anka", true),
			Email:    "Arne.Anka@fmail.com",
		},
		User{
			FirstName: newNullString("Arne", true),
			LastName:  newNullString("Panka", true),
			Email:     "Arne.Panka@fmail.com",
		},
		User{
			FirstName: newNullString("Arne", true),
			LastName:  newNullString("Panka", true),
			Email:     "Arne.Panka@fmail.com",
		},
	}

	for _, user := range users {
		fmt.Println(user)
		res := db.Create(&user)
		if res.Error != nil {
			log.Printf("Error creating user: %v\n", res.Error)
		}
	}

	//result = db.Where("last_name", "Ek").First(&users)
	//fmt.Println(result, users)
}

type User struct {
	gorm.Model
	FirstName sql.NullString `gorm:"type:VARCHAR(64); null"`
	LastName  sql.NullString `gorm:"size:64; default:'Svensson'"`
	Email     string         `gorm:"unique; not null"`
}

func newNullString(str string, valid bool) sql.NullString {
	return sql.NullString{
		String: str,
		Valid:  valid,
	}
}
