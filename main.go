package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/naoyafurudono/textfmt/file"
)

// rの内容を整形してwに書き込む
func formatText(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		// 行末尾の空白を削除（半角スペース、タブ、全角スペース）
		line := strings.TrimRight(scanner.Text(), " \t　")
		lines = append(lines, line)
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

// pathのファイルを整形する
func formatFile(path string) error {
	return file.Update(formatText, path)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		// 標準入力を処理
		if err := formatText(os.Stdin, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else {
		// ファイルを処理
		for _, path := range args {
			if err := formatFile(path); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
				os.Exit(1)
			}
		}
	}
}
