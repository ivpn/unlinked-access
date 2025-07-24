package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/auth/services/distributor/config"
)

type Handler struct {
	Cfg    config.APIConfig
	Server *fiber.App
}

func Start(cfg config.APIConfig) error {
	log.Printf("distributor server starting on :%s", cfg.Port)

	app := fiber.New()

	h := &Handler{
		Cfg:    cfg,
		Server: app,
	}

	h.SetupRoutes(h.Cfg)

	return h.Server.Listen(":" + h.Cfg.Port)
}
