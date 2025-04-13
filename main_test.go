package main

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestFormatFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
		mode     os.FileMode
	}{
		{
			name:     "通常のテキスト",
			content:  "これはテストです。",
			expected: "これはテストです。\n",
			mode:     0644,
		},
		{
			name:     "末尾に空行がある場合",
			content:  "これはテストです。\n\n\n",
			expected: "これはテストです。\n",
			mode:     0644,
		},
		{
			name:     "実行可能ファイル",
			content:  "#!/bin/sh\necho test",
			expected: "#!/bin/sh\necho test\n",
			mode:     0755,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用の一時ディレクトリを作成
			tmpDir := t.TempDir()

			// テスト用のファイルを作成
			path := filepath.Join(tmpDir, "test.txt")
			if err := os.WriteFile(path, []byte(tt.content), tt.mode); err != nil {
				t.Fatalf("テストファイルを作成できません: %v", err)
			}

			// ファイルをフォーマット
			if err := formatFile(path); err != nil {
				t.Errorf("formatFile() error = %v", err)
				return
			}

			// フォーマット後の内容を確認
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("フォーマット後のファイルを読み込めません: %v", err)
				return
			}

			if got := string(content); got != tt.expected {
				t.Errorf("formatFile() = %q, want %q", got, tt.expected)
			}

			// ファイルのモードを確認
			fileInfo, err := os.Stat(path)
			if err != nil {
				t.Errorf("ファイルのメタデータを取得できません: %v", err)
				return
			}

			if got := fileInfo.Mode(); got != tt.mode {
				t.Errorf("formatFile() mode = %v, want %v", got, tt.mode)
			}
		})
	}
}
