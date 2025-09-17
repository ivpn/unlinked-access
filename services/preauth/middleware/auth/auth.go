package auth

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewPSK(psk string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if GetToken(c) != psk {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func NewCORS(allowRemoteOrigins string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     allowRemoteOrigins,
		AllowMethods:     fiber.MethodGet,
		AllowCredentials: true,
	})
}

func NewIPFilter(allowedIPs []string) fiber.Handler {
	allowed := make(map[string]struct{})
	for _, ip := range allowedIPs {
		allowed[ip] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		log.Println("Client IP:", clientIP)
		log.Println("Allowed IPs:", allowedIPs)

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
