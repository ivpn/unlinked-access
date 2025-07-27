package http

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

type Http struct {
	Cfg config.APIConfig
}

func New(cfg config.APIConfig) *Http {
	return &Http{
		Cfg: cfg,
	}
}

func (h Http) GetManifest() (model.Manifest, error) {
	req := fiber.Get(h.Cfg.ManifestURL)
	req.Set("Accept", "application/json")
	req.Set("Authorization", "Bearer "+h.Cfg.ManifestPSK)

	status, body, errs := req.Bytes()
	if len(errs) > 0 || status != fiber.StatusOK {
		return model.Manifest{}, errs[0]
	}

	var manifest model.Manifest
	err := json.Unmarshal(body, &manifest)
	if err != nil {
		return model.Manifest{}, err
	}

	return manifest, nil
}
