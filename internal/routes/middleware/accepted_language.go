package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/common"
)

var ()

func AcceptLanguage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.Get("Accept-Language")
		if lang == "" {
			lang = "en"
		}
		c.Locals(common.LocaleContext, lang)
		return c.Next()
	}
}
