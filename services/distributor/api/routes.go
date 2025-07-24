package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"ivpn.net/auth/services/distributor/config"
	"ivpn.net/auth/services/distributor/middleware/auth"
)

func (h *Handler) SetupRoutes(cfg config.APIConfig) {
	h.Server.Use(helmet.New())
	h.Server.Use(healthcheck.New())
	h.Server.Use(auth.NewPSKCORS(cfg))
	h.Server.Use(auth.NewPSK(cfg))

	h.Server.Get("/v1/manifest", limiter.New(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}
