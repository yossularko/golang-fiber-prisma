package main

import (
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/inits"
	"golang-fiber-prisma/lib"
	"golang-fiber-prisma/middleware"
	"golang-fiber-prisma/routes"
	"log"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	inits.LoadEnv()
	inits.PrismaInit()
	inits.ValidateInit()
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

	defer func(prisma *db.PrismaClient) {
		err := lib.DisconnectFromDatabase(prisma)
		if err != nil {
			log.Fatal(err)
		}
	}(inits.Prisma)

	app := fiber.New(fiberConfig())
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	app.Use(middleware.RateLimiter(60, 30))
	app.Use(middleware.Cache(5))

	routes.SetupRoutes(app)

	if err := app.Listen(":" + port); err != nil {
		log.Panic(err)
	}
}
