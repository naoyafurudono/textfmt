# textfmt

テキストの末尾の改行を整えるシンプルなコマンドラインツールです。

## インストール

```bash
go install github.com/naoyafurudono/textfmt@latest
```

## 使い方

標準入力からテキストを受け取り、末尾の改行を整えて出力します。

```bash
# 基本的な使い方
cat input.txt | textfmt
```

## 機能

- 末尾の改行の整理
  - ファイルの末尾は必ず改行で終わります
  - 末尾の空行は削除されます

## 例

入力:
```
これはテストです。



```

出力:
```
これはテストです。
