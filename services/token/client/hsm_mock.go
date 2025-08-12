package client

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"ivpn.net/auth/services/token/model"
)

type MockHSM struct{}

func NewMockHSM() (*MockHSM, error) {
	return &MockHSM{}, nil
}

func (h *MockHSM) Token(input string) (*model.HSMToken, error) {
	if input == "" {
		return nil, fmt.Errorf("%s", ErrEmptyInput)
	}

	inputHash := sha512.Sum512([]byte(input))

	return &model.HSMToken{
		Token: base64.StdEncoding.EncodeToString(inputHash[:]),
	}, nil
}
