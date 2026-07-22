package service

import (
	"fmt"
	"io"
	"log"
	"os"

	jsonv2 "github.com/go-json-experiment/json"
	"ivpn.net/auth/services/distributor/config"
	"ivpn.net/auth/services/distributor/model"
)

const CURRENT_MANIFEST = "current.json"
const BASE_PATH = "/app/data"

type Service struct {
	Cfg config.Config
}

func New(cfg config.Config) *Service {
	return &Service{
		Cfg: cfg,
	}
}

func (s *Service) GetManifest() (model.Manifest, error) {
	path := BASE_PATH + "/" + CURRENT_MANIFEST

	log.Println("fetching manifest from", path)

	// Open the JSON file
	file, err := os.Open(path)
	if err != nil {
		return model.Manifest{}, fmt.Errorf("failed to open manifest: %w", err)
	}
	defer file.Close()

	// Read file contents
	bytes, err := io.ReadAll(file)
	if err != nil {
		return model.Manifest{}, fmt.Errorf("failed to read manifest: %w", err)
	}

	// Unmarshal JSON into Manifest struct
	var manifest model.Manifest
	if err = jsonv2.Unmarshal(bytes, &manifest, jsonv2.RejectUnknownMembers(true)); err != nil {
		return model.Manifest{}, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	return manifest, nil
}
