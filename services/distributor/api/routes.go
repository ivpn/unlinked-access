package api

import (
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"ivpn.net/auth/services/distributor/config"
	"ivpn.net/auth/services/distributor/middleware/auth"
	"ivpn.net/auth/services/distributor/middleware/compress"
)

func (h *Handler) SetupRoutes(cfg config.APIConfig) {
	h.Server.Use(helmet.New())
	h.Server.Use(healthcheck.New())
	h.Server.Use(auth.NewIPFilter(cfg.ApiAllowIPs))
	h.Server.Use(auth.NewPSK(cfg.PSK))
	h.Server.Use(compress.New())

	h.Server.Get("/v1/manifest", h.GetManifest)
}
