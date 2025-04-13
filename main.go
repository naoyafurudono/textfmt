package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func formatText(r io.Reader, w io.Writer) error {
	// 標準入力からテキストを読み込む
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("エラー: %v", err)
	}

	// 空の入力の場合は何も出力しない
	if len(lines) == 0 {
		return nil
	}

	// 末尾の空行を削除
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// テキストを出力
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	return nil
}

func main() {
	if err := formatText(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
