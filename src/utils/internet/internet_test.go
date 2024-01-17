package internet

import (
	"testing"
)

func TestCheckInternetAccess(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "check internet access",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckInternetAccess(); got != tt.want {
				t.Errorf("CheckInternetAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
