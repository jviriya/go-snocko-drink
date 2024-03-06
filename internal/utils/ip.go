package utils

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetIP(ctx interface{}) string {
	var ip string
	switch c := ctx.(type) {
	case *fiber.Ctx:
		ip = c.Get("Cf-Connecting-Ip")
		if ip != "" {
			return ip
		}
		return c.IP()
	}
	return ""
}

func truncatePort(ip string) string {
	if p := strings.Split(ip, ":"); len(p) != 1 { //0.0.0.0
		ip = p[0]
	}
	return ip
}
