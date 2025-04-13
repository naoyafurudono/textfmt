package main

import (
	"bufio"
	"flag"
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

func formatFile(path string) error {
	// ファイルを開く
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("ファイルを開けません: %v", err)
	}
	defer file.Close()

	// 一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "textfmt-")
	if err != nil {
		return fmt.Errorf("一時ファイルを作成できません: %v", err)
	}
	tmpPath := tmpFile.Name()

	// 関数が正常に実行できた場合はエラーを返すが問題ない
	defer os.Remove(tmpPath)
	defer tmpFile.Close()

	// ファイルをフォーマット
	if err := formatText(file, tmpFile); err != nil {
		return err
	}

	// 一時ファイルを閉じる
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("一時ファイルを閉じられません: %v", err)
	}

	// 元のファイルを一時ファイルで置き換え
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("ファイルを置き換えられません: %v", err)
	}

	return nil
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
