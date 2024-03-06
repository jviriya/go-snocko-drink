package service

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/pkg/email"
	"go-pentor-bank/internal/pkg/jwtManager"
	"go-pentor-bank/internal/repository/commondb"
)

type PingService struct {
	commonRepo *commondb.SecureRepository
	Jwtmanager *jwtManager.JWTManager
}

func NewPingService(rp *commondb.SecureRepository, jwtmanager *jwtManager.JWTManager) *PingService {
	return &PingService{
		commonRepo: rp,
		Jwtmanager: jwtmanager,
	}
}

func (sv *PingService) Ping(c *fiber.Ctx) error {

	return nil
}

func (sv *PingService) TestSendEmail(c *fiber.Ctx) error {
	log := clog.GetLog()
	sendTo := c.Query("email")
	err := email.SendHTMLEmail("test body", "test", []string{sendTo}, nil)
	if err != nil {
		log.Error().Err(err).Msg("TestSendEmail  got err")
		return err
	}
	return nil
}
