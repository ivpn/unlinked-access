package utils

import (
	"testing"
)

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name     string
		uuid     string
		expected bool
	}{
		{
			name:     "valid UUID v4",
			uuid:     "550e8400-e29b-41d4-a716-446655440000",
			expected: true,
		},
		{
			name:     "valid UUID v1",
			uuid:     "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			expected: true,
		},
		{
			name:     "valid UUID v3",
			uuid:     "6ba7b811-9dad-11d1-80b4-00c04fd430c8",
			expected: true,
		},
		{
			name:     "empty string",
			uuid:     "",
			expected: false,
		},
		{
			name:     "invalid UUID - missing hyphens",
			uuid:     "550e8400e29b41d4a716446655440000",
			expected: false,
		},
		{
			name:     "invalid UUID - wrong length",
			uuid:     "550e8400-e29b-41d4-a716-44665544000",
			expected: false,
		},
		{
			name:     "invalid UUID - invalid characters",
			uuid:     "550e8400-e29b-41d4-a716-44665544000g",
			expected: false,
		},
		{
			name:     "invalid UUID - wrong format",
			uuid:     "not-a-uuid",
			expected: false,
		},
		{
			name:     "invalid UUID - extra characters",
			uuid:     "550e8400-e29b-41d4-a716-446655440000-extra",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUUID(tt.uuid)
			if result != tt.expected {
				t.Errorf("ValidateUUID(%q) = %v, want %v", tt.uuid, result, tt.expected)
			}
		})
	}
}
