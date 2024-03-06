package binding

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"

	"os"
)

func BodyParserAndValidate(ctx *fiber.Ctx, req interface{}) error {

	log := clog.GetContextLog(ctx.Context())

	if err := ctx.BodyParser(&req); err != nil {
		log.Error().Err(err).Send()
		return err
	}

	if _, err := config.ValidateStruct(req); err != nil {
		log.Error().Err(err).Send()
		return err
	}

	return nil
}

func PrintJson(data interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		panic(err)
	}
}
