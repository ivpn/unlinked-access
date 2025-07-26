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
	h.Server.Use(auth.NewCORS(cfg))
	h.Server.Use(auth.NewPSK(cfg))

	h.Server.Get("/v1/preauth/:id", h.GetPreAuth)
	h.Server.Post("/v1/preauth", h.AddPreAuth)
}
