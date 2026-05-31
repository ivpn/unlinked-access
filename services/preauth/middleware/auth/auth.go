package auth

import (
	"crypto/subtle"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewIPFilter(allowedIPs []string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if slices.Contains(allowedIPs, "*") {
			return c.Next()
		}

		clientIP := c.IP()
		if clientIP != "" && slices.Contains(allowedIPs, clientIP) {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusForbidden)
	}
}

func NewPSK(psk string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if subtle.ConstantTimeCompare([]byte(GetToken(c)), []byte(psk)) == 1 {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUnauthorized)
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
