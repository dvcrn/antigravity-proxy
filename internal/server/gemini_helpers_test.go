package server

import "testing"

func TestNormalizeModelName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "gemini-3-pro-preview",
			input:    "gemini-3-pro-preview",
			expected: "gemini-3-pro-preview",
		},
		{
			name:     "gemini-3-pro",
			input:    "gemini-3-pro",
			expected: "gemini-3-pro",
		},
		{
			name:     "gemini-3-pro-low",
			input:    "gemini-3-pro-low",
			expected: "gemini-3-pro-low",
		},
		{
			name:     "gemini-3-flash",
			input:    "gemini-3-flash",
			expected: "gemini-3-flash",
		},
		{
			name:     "gemini-2.5-pro",
			input:    "gemini-2.5-pro",
			expected: "gemini-2.5-pro",
		},
		{
			name:     "gemini-2.5-flash",
			input:    "gemini-2.5-flash",
			expected: "gemini-2.5-flash",
		},
		{
			name:     "unknown model passes through",
			input:    "unknown-model",
			expected: "unknown-model",
		},
		{
			name:     "custom model variant",
			input:    "custom-variant-123",
			expected: "custom-variant-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalized := normalizeModelName(tc.input)
			if normalized != tc.expected {
				t.Errorf("Expected model name %s, but got %s", tc.expected, normalized)
			}
		})
	}
}
