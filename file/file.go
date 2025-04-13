package file

import (
	"fmt"
	"io"
	"os"
)

// pathで指定するファイルをfで更新する.
// fはfromの内容を持つファイルを更新してtoに書き込む関数であることを期待する.
// fがエラーを返した場合にはUpdateもエラーを返し、元のファイルには何もしない.
func Update(path string, f func(from io.Reader, to io.Writer) error) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("ファイルを開けません: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ファイルのメタデータを取得できません: %v", err)
	}

	// 整形後の結果を一時ファイルに書き込む
	tmpFile, err := os.CreateTemp("", "textfmt-")
	if err != nil {
		return fmt.Errorf("一時ファイルを作成できません: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	defer tmpFile.Close()

	// ファイルの更新操作を実行する
	if err := f(file, tmpFile); err != nil {
		return fmt.Errorf("ファイルの更新に失敗しました: %v", err)
	}

	// 一時ファイルを閉じる
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("一時ファイルを閉じられません: %v", err)
	}

	// 一時ファイルに出力した整形結果を元のファイルにうつす
	originalFile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("ファイルを開けません: %v", err)
	}
	defer originalFile.Close()

	tmpFile, err = os.Open(tmpPath)
	if err != nil {
		return fmt.Errorf("一時ファイルを開けません: %v", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(originalFile, tmpFile); err != nil {
		return fmt.Errorf("ファイルをコピーできません: %v", err)
	}

	return nil
}
