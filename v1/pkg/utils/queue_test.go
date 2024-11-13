package utils

import (
	"reflect"
	"testing"
)

func TestQueue_Enqueue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Test Enqueue single element",
			input:    "element1",
			expected: []string{"element1"},
		},
		{
			name:     "Test Enqueue multiple elements",
			input:    "element2",
			expected: []string{"element1", "element2"},
		},
		{
			name:     "Test Enqueue with empty string",
			input:    "",
			expected: []string{"element1", "element2", ""},
		},
	}

	q := &Queue{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q.Enqueue(tt.input)

			if !reflect.DeepEqual(q.Elements, tt.expected) {
				t.Errorf("Queue.Enqueue() = %v, want %v", q.Elements, tt.expected)
			}
		})
	}
}

func TestQueue_Dequeue(t *testing.T) {
	tests := []struct {
		name     string
		queue    *Queue
		expected string
	}{
		{
			name:     "Test Dequeue from empty queue",
			queue:    &Queue{},
			expected: "",
		},
		{
			name:     "Test Dequeue from queue with one element",
			queue:    &Queue{Elements: []string{"element1"}},
			expected: "element1",
		},
		{
			name:     "Test Dequeue from queue with multiple elements",
			queue:    &Queue{Elements: []string{"element1", "element2"}},
			expected: "element1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.queue.Dequeue()
			if result != tt.expected {
				t.Errorf("Queue.Dequeue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQueue_GetLength(t *testing.T) {
	tests := []struct {
		name     string
		queue    *Queue
		expected int
	}{
		{
			name:     "Test GetLength with empty queue",
			queue:    &Queue{},
			expected: 0,
		},
		{
			name:     "Test GetLength with one element",
			queue:    &Queue{Elements: []string{"element1"}},
			expected: 1,
		},
		{
			name:     "Test GetLength with multiple elements",
			queue:    &Queue{Elements: []string{"element1", "element2"}},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.queue.GetLength()
			if result != tt.expected {
				t.Errorf("Queue.GetLength() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		queue    *Queue
		expected bool
	}{
		{
			name:     "Test IsEmpty with empty queue",
			queue:    &Queue{},
			expected: true,
		},
		{
			name:     "Test IsEmpty with one element",
			queue:    &Queue{Elements: []string{"element1"}},
			expected: false,
		},
		{
			name:     "Test IsEmpty with multiple elements",
			queue:    &Queue{Elements: []string{"element1", "element2"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.queue.IsEmpty()
			if result != tt.expected {
				t.Errorf("Queue.IsEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}
