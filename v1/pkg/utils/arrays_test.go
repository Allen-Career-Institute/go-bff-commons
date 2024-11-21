package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		e    string
		want bool
	}{
		{
			name: "element exists",
			s:    []string{"apple", "banana", "cherry"},
			e:    "banana",
			want: true,
		},
		{
			name: "element does not exist",
			s:    []string{"apple", "banana", "cherry"},
			e:    "grape",
			want: false,
		},
		{
			name: "empty slice",
			s:    []string{},
			e:    "apple",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.s, tt.e); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want []string
	}{
		{
			name: "empty slice",
			s:    []string{},
			want: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDifferenceList(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "non-empty difference",
			args: args{
				a: []string{"apple", "banana", "cherry"},
				b: []string{"banana", "cherry"},
			},
			want: []string{"apple"},
		},
		{
			name: "empty difference",
			args: args{
				a: []string{"apple", "banana", "cherry"},
				b: []string{"apple", "banana", "cherry"},
			},
			want: []string{},
		},
		{
			name: "difference with empty slice",
			args: args{
				a: []string{},
				b: []string{"apple", "banana", "cherry"},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DifferenceList(tt.args.a, tt.args.b), "DifferenceList(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}
