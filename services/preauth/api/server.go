package api

import (
	"context"
	"log"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gofiber/fiber/v2"
	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/model"
	"ivpn.net/auth/services/preauth/utils"
)

var (
	ErrInvalidRequest = "the request is invalid."
	AddPreAuthError   = "failed to add pre-authentication."
	GetPreAuthError   = "failed to retrieve pre-authentication."
)

type Service interface {
	AddPreAuth(context.Context, string, bool, time.Time, string) (string, error)
	GetPreAuth(context.Context, string) (model.PreAuth, error)
}

type Handler struct {
	Cfg       config.APIConfig
	Server    *fiber.App
	Service   Service
	Validator utils.Validator
}

func Start(cfg config.APIConfig, service Service) error {
	// Channel to collect errors from both servers
	errCh := make(chan error, 2)

	// Start /add server in a goroutine
	go func() {
		h := &Handler{
			Cfg:       cfg,
			Server:    fiber.New(),
			Service:   service,
			Validator: utils.NewValidator(),
		}

		h.SetupRoutesAdd(h.Cfg)

		log.Printf("preauth /add server starting on :%s", cfg.AddPort)
		errCh <- h.Server.Listen(":" + cfg.AddPort)
	}()

	// Start /get server in a goroutine
	go func() {
		h := &Handler{
			Cfg:       cfg,
			Server:    fiber.New(),
			Service:   service,
			Validator: utils.NewValidator(),
		}

		h.SetupRoutesGet(h.Cfg)

		log.Printf("preauth /get server starting on :%s", cfg.GetPort)
		errCh <- h.Server.Listen(":" + cfg.GetPort)
	}()

	// Wait for any server to return an error
	return <-errCh
}

func (h *Handler) AddPreAuth(c *fiber.Ctx) error {
	req := PreauthReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Println("failed to parse request body:", err)
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	err = h.Validator.Struct(req)
	if err != nil {
		log.Println("failed to validate request body:", err)
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	activeUntil, err := dateparse.ParseAny(req.ActiveUntil)
	if err != nil {
		log.Println("failed to parse ActiveUntil:", err)
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	sessionID, err := h.Service.AddPreAuth(c.Context(), req.AccountID, req.IsActive, activeUntil, req.Tier)
	if err != nil {
		log.Println("failed to add pre-authentication:", err)
		return c.Status(400).JSON(fiber.Map{
			"error": AddPreAuthError,
		})
	}

	return c.JSON(fiber.Map{"session_id": sessionID})
}

func (h *Handler) GetPreAuth(c *fiber.Ctx) error {
	id := c.Params("id")
	if !utils.ValidateUUID(id) {
		log.Println("invalid UUID format:", id)
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}
	pa, err := h.Service.GetPreAuth(c.Context(), id)
	if err != nil {
		log.Println("failed to retrieve pre-authentication:", err)
		return c.Status(400).JSON(fiber.Map{
			"error": GetPreAuthError,
		})
	}

	return c.JSON(pa)
}
