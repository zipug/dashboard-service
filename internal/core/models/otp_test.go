package models

import "testing"

func BenchmarkLen(b *testing.B) {
	code := OTPCode(123456)
	for i := 0; i < b.N; i++ {
		code.len()
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name  string
		input OTPCode
		want  int
	}{
		{"1 should by 1", OTPCode(1), 1},
		{"-4 should be 1", OTPCode(-4), 1},
		{"12 should be 2", OTPCode(12), 2},
		{"123 should be 3", OTPCode(123), 3},
		{"1234 should be 4", OTPCode(1234), 4},
		{"12345 should be 5", OTPCode(12345), 5},
		{"123456 should be 6", OTPCode(123456), 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.len()
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
