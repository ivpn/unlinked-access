package auth

import (
	"log"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewIPFilter(allowedIPs []string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		clientIP := c.IP()
		if slices.Contains(allowedIPs, clientIP) || c.IP() == "" {
			return c.Next()
		}

		log.Println("Unauthorized IP: ", clientIP)

		return c.SendStatus(fiber.StatusForbidden)
	}
}

func NewPSK(psk string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if GetToken(c) == psk {
			return c.Next()
		}

		log.Println("Unauthorized PSK: ", GetToken(c))

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
