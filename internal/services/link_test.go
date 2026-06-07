package services

// Table-driven tests: the standard Go pattern.
// One test function, a slice of cases, loop with t.Run.
// Benefits: easy to add cases, each gets its own name in output.

import "testing"

func TestGetUniqueCode(t *testing.T) {
	// Each case is an anonymous struct. Fields: name (shown in output),
	// input, and the expected output.
	tests := []struct {
		name  string
		input int
		want  string
	}{
		// lookup = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		{"zero handled explicitly", 0, "0"},
		{"single digit", 5, "5"},
		{"first lowercase letter", 10, "a"},
		{"last base62 character", 61, "Z"},
		{"rolls over to two chars", 62, "10"}, // 1*62 + 0 = 62
		{"two-char mid-range", 124, "20"},     // 2*62 + 0 = 124
	}

	for _, tt := range tests {
		// t.Run creates a subtest. If one fails the others still run.
		t.Run(tt.name, func(t *testing.T) {
			got := getUniqueCode(tt.input)
			if got != tt.want {
				// t.Errorf marks the test as failed but continues running.
				// Use t.Fatalf when the rest of the test would panic on failure.
				t.Errorf("getUniqueCode(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
