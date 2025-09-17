package api

import (
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/middleware/auth"
)

func (h *Handler) SetupRoutes(cfg config.APIConfig) {
	h.Server.Use(helmet.New())
	h.Server.Use(healthcheck.New())

	get := h.Server.Group("/v1/preauth/get")
	get.Use(auth.NewCORS(cfg.AllowRemoteOrigins))
	get.Use(auth.NewIPFilter(cfg.AllowedRemoteIPs))
	get.Use(auth.NewPSK(cfg.PSK))
	get.Get("/:id", h.GetPreAuth)

	add := h.Server.Group("/v1/preauth/add")
	add.Use(auth.NewIPFilter(cfg.AllowedLocalIPs))
	add.Use(auth.NewPSK(cfg.PSK))
	add.Post("", h.AddPreAuth)
}
