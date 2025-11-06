package main

import (
	"backendgo/database"
	"backendgo/route"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	_ "backendgo/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title BackendGo API
// @version 1.0
// @description API untuk mengelola data backend menggunakan Fiber dan MongoDB.
// @contact.name Sherly Tanti Virginia
// @contact.url https://github.com/scherlyz
// @contact.email scheerly.tnv@gmail.com
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan format: Bearer {token}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.ConnectDB()
	database.ConnectMongoDB()
	defer database.DB.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		},
	})

	
	app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:3000, http://127.0.0.1:3000",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowCredentials: true,
}))


	// Swagger docs
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Semua route
	route.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
