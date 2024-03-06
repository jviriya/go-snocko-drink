package routes

import "github.com/gofiber/fiber/v2"

type Route struct {
	Name        string
	Description string
	Method      string
	Pattern     string
	Endpoint    fiber.Handler
	Middleware  []fiber.Handler
	Test        bool
}
