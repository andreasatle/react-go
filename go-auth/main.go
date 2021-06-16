package main

import (
	"log"

	"github.com/andreasatle/react-go/go-auth/database"
	"github.com/andreasatle/react-go/go-auth/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Setup the logging, for if program crashes
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Failed listening to port: %v\n", err)
	}
}
