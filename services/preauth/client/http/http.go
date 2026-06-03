package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/model"
)

type Http struct {
	Cfg config.APIConfig
}

func New(cfg config.APIConfig) *Http {
	return &Http{
		Cfg: cfg,
	}
}

func (h Http) PostSession(session model.Session, url string, psk string) error {
	body, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	req := fiber.Post(url)
	req.Set("Content-Type", "application/json")
	req.Set("Accept", "application/json")
	req.Set("Authorization", "Bearer "+psk)
	req.Body(body)

	status, res, errs := req.Bytes()
	if len(errs) > 0 {
		log.Printf("Error calling session webhook: %v", errs)
		return errors.New("error calling session webhook")
	}

	if status != http.StatusOK {
		// Log response for debugging
		log.Printf("Error calling session webhook, status: %d", status)
		log.Printf("Session webhook response: %s", string(res))
		return errors.New("error response from session webhook")
	}

	return nil
}
