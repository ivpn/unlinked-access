package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/middleware/auth"
)

func (h *Handler) SetupRoutesAdd(cfg config.APIConfig) {
	h.Server.Get("/debug-ip", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ip":              c.IP(),
			"x-forwarded-for": c.Get("X-Forwarded-For"),
			"x-real-ip":       c.Get("X-Real-IP"),
		})
	})

	h.Server.Use(helmet.New())
	h.Server.Use(healthcheck.New())

	add := h.Server.Group("/v1/preauth/add")
	add.Use(auth.NewIPFilter(cfg.ApiAllowIPs))
	add.Use(auth.NewPSK(cfg.AddPSK))
	add.Post("", h.AddPreAuth)
}

func (h *Handler) SetupRoutesGet(cfg config.APIConfig) {
	h.Server.Get("/debug-ip", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ip":              c.IP(),
			"x-forwarded-for": c.Get("X-Forwarded-For"),
			"x-real-ip":       c.Get("X-Real-IP"),
		})
	})

	h.Server.Use(helmet.New())
	h.Server.Use(healthcheck.New())

	get := h.Server.Group("/v1/preauth/get")
	get.Use(auth.NewIPFilter(cfg.ApiAllowIPs))
	get.Use(auth.NewPSK(cfg.GetPSK))
	get.Get("/:id", h.GetPreAuth)
}
