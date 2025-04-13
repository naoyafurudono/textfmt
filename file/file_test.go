package file_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/naoyafurudono/textfmt/file"
)

func TestUpdate(t *testing.T) {
	// テスト用ディレクトリを作成
	testDir, err := os.MkdirTemp("", "file-test-")
	if err != nil {
		t.Fatalf("テスト用ディレクトリを作成できません: %v", err)
	}
	defer os.RemoveAll(testDir)

	t.Run("正常に更新できる場合", func(t *testing.T) {
		// テストファイルを作成
		testPath := filepath.Join(testDir, "test1.txt")
		if err := os.WriteFile(testPath, []byte("hello world"), 0644); err != nil {
			t.Fatalf("テストファイルを作成できません: %v", err)
		}

		// 大文字に変換する更新関数
		uppercase := func(from io.Reader, to io.Writer) error {
			content, err := io.ReadAll(from)
			if err != nil {
				return err
			}
			_, err = to.Write([]byte(strings.ToUpper(string(content))))
			return err
		}

		// ファイルを更新
		if err := file.Update(uppercase, testPath); err != nil {
			t.Fatalf("Update関数が失敗しました: %v", err)
		}

		// 結果を確認
		content, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatalf("更新後のファイルを読み込めません: %v", err)
		}
		if got, want := string(content), "HELLO WORLD"; got != want {
			t.Errorf("更新後の内容が異なります。\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("更新関数がエラーを返す場合", func(t *testing.T) {
		// テストファイルを作成
		testPath := filepath.Join(testDir, "test2.txt")
		originalContent := "test content"
		if err := os.WriteFile(testPath, []byte(originalContent), 0644); err != nil {
			t.Fatalf("テストファイルを作成できません: %v", err)
		}

		// エラーを返す更新関数
		errFunc := func(from io.Reader, to io.Writer) error {
			return errors.New("更新エラー")
		}

		// ファイルを更新（エラーが発生するはず）
		err := file.Update(errFunc, testPath)
		if err == nil {
			t.Fatalf("エラーが期待されましたが、発生しませんでした")
		}

		// ファイルが変更されていないことを確認
		content, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatalf("ファイルを読み込めません: %v", err)
		}
		if got, want := string(content), originalContent; got != want {
			t.Errorf("ファイルが変更されています。\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("存在しないファイルの場合", func(t *testing.T) {
		nonExistentPath := filepath.Join(testDir, "non-existent.txt")

		// 空の更新関数
		nopFunc := func(from io.Reader, to io.Writer) error {
			return nil
		}

		// 存在しないファイルを更新しようとする
		err := file.Update(nopFunc, nonExistentPath)
		if err == nil {
			t.Fatalf("存在しないファイルに対してエラーが期待されましたが、発生しませんでした")
		}
	})

	t.Run("ファイルパーミッションが保持される", func(t *testing.T) {
		// 特定のパーミッションでテストファイルを作成
		testPath := filepath.Join(testDir, "test3.txt")
		permissions := os.FileMode(0640)
		f, err := os.OpenFile(testPath, os.O_CREATE|os.O_WRONLY, permissions)
		if err != nil {
			t.Fatalf("テストファイルを作成できません: %v", err)
		}
		f.WriteString("permission test")
		f.Close()

		// 単純な更新関数
		updateFunc := func(from io.Reader, to io.Writer) error {
			content, err := io.ReadAll(from)
			if err != nil {
				return err
			}
			_, err = to.Write(append(content, []byte(" - updated")...))
			return err
		}

		// ファイルを更新
		if err := file.Update(updateFunc, testPath); err != nil {
			t.Fatalf("Update関数が失敗しました: %v", err)
		}

		// パーミッションが保持されていることを確認
		info, err := os.Stat(testPath)
		if err != nil {
			t.Fatalf("ファイル情報を取得できません: %v", err)
		}
		if info.Mode().Perm() != permissions {
			t.Errorf("ファイルパーミッションが保持されていません。\n got: %v\nwant: %v",
				info.Mode().Perm(), permissions)
		}
	})
}
