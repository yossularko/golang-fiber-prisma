package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func Cache(durationInSecond int) fiber.Handler {
	return cache.New(cache.Config{
		Expiration: time.Duration(durationInSecond) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.OriginalURL()
		},
	})
}
