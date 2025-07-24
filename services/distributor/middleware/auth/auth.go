package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"ivpn.net/auth/services/distributor/config"
)

func NewPSK(cfg config.APIConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if GetToken(c) != cfg.PSK {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func NewPSKCORS(cfg config.APIConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.PSKAllowOrigin,
		AllowMethods:     fiber.MethodGet,
		AllowCredentials: true,
	})
}

func GetToken(c *fiber.Ctx) string {
	var token string
	authorization := c.Get("Authorization")

	if after, ok := strings.CutPrefix(authorization, "Bearer "); ok {
		token = after
	}

	return token
}
