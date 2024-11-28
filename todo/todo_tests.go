package todo

import "testing"

func TestCycleState(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"empty line": {
			input:    "",
			expected: "",
		},
		"line without state": {
			input:    "test line",
			expected: "TODO test line",
		},
		"line with TODO": {
			input:    "TODO test line",
			expected: "DOING test line",
		},
		"line with DOING": {
			input:    "DOING test line",
			expected: "DONE test line",
		},
		"line with DONE": {
			input:    "DONE test line",
			expected: "TODO test line",
		},
		"indented line without state": {
			input:    "    test line",
			expected: "    TODO test line",
		},
		"indented line with TODO": {
			input:    "    TODO test line",
			expected: "    DOING test line",
		},
		"line with priority": {
			input:    "TODO [#A] test line",
			expected: "DOING [#A] test line",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := CycleState(tt.input)
			if result != tt.expected {
				t.Errorf("CycleState(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCyclePriority(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"empty line": {
			input:    "",
			expected: "",
		},
		"line without TODO state": {
			input:    "test line",
			expected: "test line", // No change if no TODO state
		},
		"TODO line without priority": {
			input:    "TODO test line",
			expected: "TODO [#A] test line",
		},
		"TODO line with [#A]": {
			input:    "TODO [#A] test line",
			expected: "TODO [#B] test line",
		},
		"TODO line with [#B]": {
			input:    "TODO [#B] test line",
			expected: "TODO [#C] test line",
		},
		"TODO line with [#C]": {
			input:    "TODO [#C] test line",
			expected: "TODO test line",
		},
		"DOING line without priority": {
			input:    "DOING test line",
			expected: "DOING [#A] test line",
		},
		"DONE line without priority": {
			input:    "DONE test line",
			expected: "DONE [#A] test line",
		},
		"indented TODO without priority": {
			input:    "    TODO test line",
			expected: "    TODO [#A] test line",
		},
		"indented line without TODO": {
			input:    "    test line",
			expected: "    test line", // No change if no TODO state
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := CyclePriority(tt.input)
			if result != tt.expected {
				t.Errorf("CyclePriority(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
