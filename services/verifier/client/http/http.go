package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"

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
	req.Set("Accept-Encoding", "gzip")
	req.Set("Authorization", "Bearer "+h.Cfg.ManifestPSK)
	req.Set("Accept", "application/json")

	status, body, errs := req.Bytes()
	if len(errs) > 0 || status != fiber.StatusOK {
		return model.Manifest{}, errs[0]
	}

	// Handle gzip decompression if needed
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		// If not gzip, use original body
		var manifest model.Manifest
		err := json.Unmarshal(body, &manifest)
		if err != nil {
			return model.Manifest{}, err
		}
		return manifest, nil
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return model.Manifest{}, err
	}

	var manifest model.Manifest
	err = json.Unmarshal(decompressed, &manifest)
	if err != nil {
		return model.Manifest{}, err
	}

	return manifest, nil
}
