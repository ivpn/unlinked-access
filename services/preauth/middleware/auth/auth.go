package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewPSK(psk string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if GetToken(c) != psk {
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
