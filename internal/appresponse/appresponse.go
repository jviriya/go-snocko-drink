package appresponse

import (
	"github.com/gofiber/fiber/v2"
	"math"
)

// JSONResponse is a function to return response in JSON-format
func JSONResponse(c *fiber.Ctx, status int, v IResponse) error {
	v.ErrorCode = v.ErrorCode.WithLocale(c)
	if c.Locals("request-id") != nil {
		v.TraceID = c.Locals("request-id").(string)
	}
	return c.Status(status).JSON(v)
}

func GetTotalPages(totalRecord, pageLimit int64) int64 {
	return int64(math.Ceil(float64(totalRecord) / float64(pageLimit)))
}
