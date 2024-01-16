package login

import (
	"testing"
)

func Test_checkInternetAccess(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Test_checkInternetAccess",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkInternetAccess(); got != tt.want {
				t.Errorf("checkInternetAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
