# textfmt

テキストの末尾の改行と空白を整えるコマンドラインツールです。

## インストール

```bash
go install github.com/naoyafurudono/textfmt@latest
```

## 使い方

```bash
# 標準入力から
cat input.txt | textfmt

# ファイルを直接処理
textfmt file.txt
textfmt file1.txt file2.txt
```

## 機能

- 行末尾に空白文字がないようにする
- 入力末尾を空でない改行で終わる行にする

## 例

```bash
# 標準入力から
echo -e "行1  \n行2　\n行3\t\n" | textfmt
# 出力:
# 行1
# 行2
# 行3

# ファイルを直接処理
textfmt test.txt
```
