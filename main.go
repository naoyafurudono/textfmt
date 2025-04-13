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

	// ファイルのメタデータを取得
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ファイルのメタデータを取得できません: %v", err)
	}

	// 一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "textfmt-")
	if err != nil {
		return fmt.Errorf("一時ファイルを作成できません: %v", err)
	}
	tmpPath := tmpFile.Name()
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

	// 元のファイルを開く（書き込みモード）
	originalFile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("ファイルを開けません: %v", err)
	}
	defer originalFile.Close()

	// 一時ファイルを開く
	tmpFile, err = os.Open(tmpPath)
	if err != nil {
		return fmt.Errorf("一時ファイルを開けません: %v", err)
	}
	defer tmpFile.Close()

	// 一時ファイルの内容を元のファイルにコピー
	if _, err := io.Copy(originalFile, tmpFile); err != nil {
		return fmt.Errorf("ファイルをコピーできません: %v", err)
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
