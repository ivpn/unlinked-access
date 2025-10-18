package http

import (
	"errors"
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

func (h Http) PostSession(session model.Session) error {
	req := fiber.Post(h.Cfg.SessionURL)
	req.Set("Content-Type", "application/json")
	req.Set("Accept", "application/json")
	req.Set("Authorization", "Bearer "+h.Cfg.SessionPSK)
	req.Body([]byte(`{"id": "` + session.ID + `", "token": "` + session.Token + `", "preauth_id": "` + session.PreAuthID + `"}`))

	status, res, err := req.Bytes()
	if err != nil {
		log.Printf("Error calling session webhook: %v", err)
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
