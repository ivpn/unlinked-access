package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
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

func NewIPFilter(cfg config.APIConfig) fiber.Handler {
	allowed := make(map[string]struct{})
	for _, ip := range cfg.AllowedIPs {
		allowed[ip] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		clientIP := c.IP()
		if _, ok := allowed[clientIP]; !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func GetToken(c *fiber.Ctx) string {
	var token string
	authorization := c.Get("Authorization")

	if after, ok := strings.CutPrefix(authorization, "Bearer "); ok {
		token = after
	}

	return token
}
