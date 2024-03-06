package service

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/repository/commondb"
)

type PingService struct {
	commonRepo *commondb.SecureRepository
}

func NewPingService() *PingService {
	return &PingService{}
}
func (sv *PingService) Ping(c *fiber.Ctx) error {
	return nil
}
