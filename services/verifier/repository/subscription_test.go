package repository

import (
	"testing"
)

func TestJoinInt64s(t *testing.T) {
	tests := []struct {
		name     string
		ids      []int64
		expected string
	}{
		{
			name:     "empty slice",
			ids:      []int64{},
			expected: "",
		},
		{
			name:     "single element",
			ids:      []int64{123},
			expected: "123",
		},
		{
			name:     "multiple elements",
			ids:      []int64{1, 2, 3},
			expected: "1,2,3",
		},
		{
			name:     "large numbers",
			ids:      []int64{9223372036854775807, 1000000000000},
			expected: "9223372036854775807,1000000000000",
		},
		{
			name:     "negative numbers",
			ids:      []int64{-1, -100, 42},
			expected: "-1,-100,42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinInt64s(tt.ids)
			if result != tt.expected {
				t.Errorf("joinInt64s(%v) = %q, want %q", tt.ids, result, tt.expected)
			}
		})
	}
}
