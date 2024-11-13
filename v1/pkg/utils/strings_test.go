package utils

import "testing"

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want string
	}{
		{
			name: "Test with multiple strings",
			strs: []string{"Hello", " ", "World"},
			want: "Hello World",
		},
		{
			name: "Test with single string",
			strs: []string{"HelloWorld"},
			want: "HelloWorld",
		},
		{
			name: "Test with no strings",
			strs: []string{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.strs...); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinRunes(t *testing.T) {
	tests := []struct {
		name  string
		runes []rune
		want  string
	}{
		{
			name:  "Test with multiple runes",
			runes: []rune{'H', 'e', 'l', 'l', 'o'},
			want:  "Hello",
		},
		{
			name:  "Test with single rune",
			runes: []rune{'A'},
			want:  "A",
		},
		{
			name:  "Test with no runes",
			runes: []rune{},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinRunes(tt.runes...); got != tt.want {
				t.Errorf("JoinRunes() = %v, want %v", got, tt.want)
			}
		})
	}
}
