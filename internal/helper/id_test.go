package helper

import "testing"

func TestParseIDParam(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      int64
		wantError bool
	}{
		{
			name:      "valid positive integer",
			input:     "123",
			want:      123,
			wantError: false,
		},
		{
			name:      "zero value",
			input:     "0",
			want:      0,
			wantError: false,
		},
		{
			name:      "invalid string",
			input:     "abc",
			want:      0,
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			want:      0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIDParam(tt.input)

			if tt.wantError && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Fatalf("did not expect error, got %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}