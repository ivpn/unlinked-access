package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/auth/services/preauth/config"
	"ivpn.net/auth/services/preauth/model"
	"ivpn.net/auth/services/preauth/utils"
)

var (
	ErrInvalidRequest = "The request is invalid."
	AddPreAuthSuccess = "Pre-authentication added successfully."
	AddPreAuthError   = "Failed to add pre-authentication."
	GetPreAuthError   = "Failed to retrieve pre-authentication."
)

type Service interface {
	AddPreAuth(string) error
	GetPreAuth(string) (model.PreAuth, error)
}

type Handler struct {
	Cfg       config.APIConfig
	Server    *fiber.App
	Service   Service
	Validator utils.Validator
}

func Start(cfg config.APIConfig, service Service) error {
	log.Printf("preauth server starting on :%s", cfg.Port)

	app := fiber.New()

	h := &Handler{
		Cfg:       cfg,
		Server:    app,
		Service:   service,
		Validator: utils.NewValidator(),
	}

	h.SetupRoutes(h.Cfg)

	return h.Server.Listen(":" + h.Cfg.Port)
}

func (h *Handler) AddPreAuth(c *fiber.Ctx) error {
	req := PreauthReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	err = h.Validator.Struct(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	err = h.Service.AddPreAuth(req.AccountID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": AddPreAuthError,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": AddPreAuthSuccess,
	})
}

func (h *Handler) GetPreAuth(c *fiber.Ctx) error {
	id := c.Params("id")
	if !utils.ValidateUUID(id) {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}
	preAuth, err := h.Service.GetPreAuth(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": GetPreAuthError,
		})
	}

	return c.JSON(preAuth)
}
