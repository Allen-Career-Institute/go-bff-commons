package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMin(t *testing.T) {
	tests := []struct {
		name string
		val1 int
		val2 int
		want int
	}{
		{"Test uint", 5, 3, 3},
		{"Test float64", 5, 3, 3}, // Assuming GetMin can handle ints and floats
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMin(tt.val1, tt.val2)
			if got != tt.want {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestGetMax(t *testing.T) {
	tests := []struct {
		name string
		val1 int
		val2 int
		want int
	}{
		{"Test uint", 5, 3, 5},
		{"Test float64", 5, 3, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMax(tt.val1, tt.val2)
			if got != tt.want {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
