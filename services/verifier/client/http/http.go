package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/auth/services/verifier/config"
	"ivpn.net/auth/services/verifier/model"
)

const (
	// maxCompressedBodySize caps the raw HTTP response body before decompression.
	// Sized for 100k+ subscriptions (~3.5 MB compressed at 5:1 ratio) with ~2× headroom.
	maxCompressedBodySize = 8 * 1024 * 1024 // 8 MB

	// maxDecompressedBodySize caps the gzip-decompressed payload.
	// Sized for 100k+ subscriptions (~18 MB at ~180 bytes/record) with ~1.75× headroom.
	maxDecompressedBodySize = 32 * 1024 * 1024 // 32 MB
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

	// Bind the compressed response body before any allocation.
	// fiber.Get calls Parse() internally, so HostClient is already initialised here.
	req.HostClient.MaxResponseBodySize = maxCompressedBodySize

	status, body, errs := req.Bytes()
	if len(errs) > 0 {
		return model.Manifest{}, errs[0]
	}
	if status != fiber.StatusOK {
		return model.Manifest{}, fmt.Errorf("unexpected status code: %d", status)
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

	// Read one byte past the limit so that an over-size payload is detected
	// rather than silently truncated.
	decompressed, err := io.ReadAll(io.LimitReader(reader, maxDecompressedBodySize+1))
	if err != nil {
		return model.Manifest{}, err
	}
	if int64(len(decompressed)) > maxDecompressedBodySize {
		return model.Manifest{}, fmt.Errorf("decompressed manifest exceeds maximum allowed size (%d bytes)", maxDecompressedBodySize)
	}

	var manifest model.Manifest
	err = json.Unmarshal(decompressed, &manifest)
	if err != nil {
		return model.Manifest{}, err
	}

	return manifest, nil
}
