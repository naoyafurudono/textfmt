package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空の入力",
			input:    "",
			expected: "",
		},
		{
			name:     "通常のテキスト",
			input:    "これはテストです。",
			expected: "これはテストです。\n",
		},
		{
			name:     "末尾に空行がある場合",
			input:    "これはテストです。\n\n\n",
			expected: "これはテストです。\n",
		},
		{
			name:     "複数行のテキスト",
			input:    "1行目\n2行目\n3行目\n\n",
			expected: "1行目\n2行目\n3行目\n",
		},
		{
			name:     "空行のみ",
			input:    "\n\n\n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatText(strings.NewReader(tt.input), &buf)
			if err != nil {
				t.Errorf("formatText() error = %v", err)
				return
			}
			if got := buf.String(); got != tt.expected {
				t.Errorf("formatText() = %q, want %q", got, tt.expected)
			}
		})
	}
}
