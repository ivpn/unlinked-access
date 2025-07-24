package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/auth/services/distributor/config"
	"ivpn.net/auth/services/distributor/model"
)

type Service interface {
	GetManifest() (model.Manifest, error)
}

type Handler struct {
	Cfg     config.APIConfig
	Server  *fiber.App
	Service Service
}

func Start(cfg config.APIConfig, service Service) error {
	log.Printf("distributor server starting on :%s", cfg.Port)

	app := fiber.New()

	h := &Handler{
		Cfg:     cfg,
		Server:  app,
		Service: service,
	}

	h.SetupRoutes(h.Cfg)

	return h.Server.Listen(":" + h.Cfg.Port)
}

func (h *Handler) GetManifest(c *fiber.Ctx) error {
	manifest, err := h.Service.GetManifest()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(manifest)
}
