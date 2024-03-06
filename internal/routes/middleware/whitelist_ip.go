package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func ValidateWhitelistIP() fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.Next()
	}
}

func GetIP(c *fiber.Ctx) string {
	ip := c.Get("Cf-Connecting-Ip")
	if ip != "" {
		return ip
	}
	return c.IP()
}

func truncatePort(ip string) string {
	if p := strings.Split(ip, ":"); len(p) != 1 { //0.0.0.0
		ip = p[0]
	}
	return ip
}
