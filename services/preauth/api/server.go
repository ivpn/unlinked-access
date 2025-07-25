package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/auth/services/preauth/config"
)

type Service interface {
	AddPreAuth(string) error
	GetPreAuth(string) error
}

type Handler struct {
	Cfg     config.APIConfig
	Server  *fiber.App
	Service Service
}

func Start(cfg config.APIConfig, service Service) error {
	log.Printf("preauth server starting on :%s", cfg.Port)

	app := fiber.New()

	h := &Handler{
		Cfg:     cfg,
		Server:  app,
		Service: service,
	}

	h.SetupRoutes(h.Cfg)

	return h.Server.Listen(":" + h.Cfg.Port)
}

func (h *Handler) AddPreAuth(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) GetPreAuth(c *fiber.Ctx) error {
	return nil
}
