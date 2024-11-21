package utils

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestGetListOfCustomHeadersAsCommaSeparated(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test GetListOfCustomHeadersAsCommaSeparated",
			want: "X-ACCESS-TOKEN,X-REFRESH-TOKEN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetListOfCustomHeadersAsCommaSeparated(); got != tt.want {
				t.Errorf("GetListOfCustomHeadersAsCommaSeparated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserLocationAsiaKolkata(t *testing.T) {
	tests := []struct {
		name string
		mock gomock.Call
	}{
		{
			name: "Test UserLocationAsiaKolkata",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserLocationAsiaKolkata()
			if got.String() != "Asia/Kolkata" && got.String() != time.Now().Location().String() {
				t.Errorf("UserLocationAsiaKolkata() = %v, want either 'Asia/Kolkata' or current system location", got)
			}
		})
	}
}
