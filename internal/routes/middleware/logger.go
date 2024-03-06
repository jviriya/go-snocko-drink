package middleware

import (
	"crypto/rand"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/common"
)

// LoggerWithRequestMeta is a middleware that inject request information into a logger
func LoggerWithRequestMeta() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Return new handler

		//rqID := c.Get("request-id", Generator())
		rqID := Generator()
		appID := "go-pentor-bank"
		projectID := "pbank"

		logger2 := clog.WithField(map[string]interface{}{
			"path":       c.Path(),
			"request-id": rqID,
			"app-id":     appID,
			"project-id": projectID,
		})

		c.Locals(clog.CLoggerKey, logger2)

		c.Set(common.RequestIDContext, rqID)
		c.Locals(common.RequestIDContext, rqID)
		c.Locals("app-id", appID)
		c.Locals("project-id", projectID)

		return c.Next()
	}
}

func Generator() string {
	//timeNow := time.Now().Unix()
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("ffffffff")
		//return fmt.Sprintf("%v", timeNow)
	}
	return fmt.Sprintf("%x", b)
}
