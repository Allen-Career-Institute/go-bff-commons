package commons

import (
	"reflect"
	"testing"
)

func TestRemoveEmptyStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "No empty strings",
			input:    []string{"apple", "banana", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "Some empty strings",
			input:    []string{"apple", "", "banana", "", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "Slice with spaces",
			input:    []string{"apple", " ", "banana", ""},
			expected: []string{"apple", " ", "banana"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveEmptyStrings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveEmptyStrings(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
