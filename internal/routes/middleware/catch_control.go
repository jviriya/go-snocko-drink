package middleware

import (
	"github.com/gofiber/fiber/v2"
)

const HTTPCacheControlHeaderKey = "Cache-Control"

func CacheControl() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(HTTPCacheControlHeaderKey, "no-cache,no-store")
		return c.Next()
	}
}
