package main

import (
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/lib"
	"golang-fiber-prisma/middleware"
	"log"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var prisma *db.PrismaClient

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func fiberConfig() fiber.Config {
	return fiber.Config{
		Prefork:           true,
		IdleTimeout:       10 * time.Second,
		EnablePrintRoutes: true,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	prisma = db.NewClient()
	err := lib.ConnectToDatabase(prisma)
	if err != nil {
		log.Fatal(err)
	}
	defer func(prisma *db.PrismaClient) {
		err := lib.DisconnectFromDatabase(prisma)
		if err != nil {
			log.Fatal(err)
		}
	}(prisma)

	app := fiber.New(fiberConfig())
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	app.Use(middleware.RateLimiter(60, 30))
	app.Use(middleware.Cache(5))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	if err := app.Listen(":" + port); err != nil {
		log.Panic(err)
	}
}