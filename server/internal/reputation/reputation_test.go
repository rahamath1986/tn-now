package reputation

import (
	"testing"
)

func TestCalculateLevel(t *testing.T) {
	tests := []struct {
		name          string
		approvedCount int32
		trustScore    int32
		wantLevel     string
	}{
		{
			name:          "New Contributor default",
			approvedCount: 2,
			trustScore:    10,
			wantLevel:     "New Contributor",
		},
		{
			name:          "Trusted Contributor threshold",
			approvedCount: 5,
			trustScore:    20,
			wantLevel:     "Trusted Contributor",
		},
		{
			name:          "Verified Contributor threshold",
			approvedCount: 20,
			trustScore:    50,
			wantLevel:     "Verified Contributor",
		},
		{
			name:          "Core Contributor threshold",
			approvedCount: 50,
			trustScore:    85,
			wantLevel:     "Core Contributor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateLevel(tt.approvedCount, tt.trustScore)
			if got != tt.wantLevel {
				t.Errorf("CalculateLevel() = %v, want %v", got, tt.wantLevel)
			}
		})
	}
}
