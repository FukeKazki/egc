# egc(emoji-generator-cli)

Slack などで使用する絵文字を生成する CLI ツール

# Install

## mise

```
mise use -g go:github.com/FukeKazki/egc@latest
```

## go install

```
go install github.com/FukeKazki/egc@latest
```

# Usage

## 絵文字を作成(デフォルトでは文字の色が pink)

```
egc 完全に理解した
```

生成物

![完全に理解した](https://github.com/apple-yagi/egc/assets/57742720/b5f676a1-2b49-470f-9e64-612357942034)

## 文字の色を指定して絵文字を作成

```
egc 完全に理解した -c yellow
```

生成物

![完全に理解した](https://github.com/apple-yagi/egc/assets/57742720/23256e8d-51ec-43f7-995d-aee0d793286a)

指定できる色

- pink
- yellow
- black
- red
- green
- blue

## フォントを指定して絵文字を作成

```
egc 完全に理解した -f mplus
```

指定できるフォント

- mono — Noto Sans Mono CJK JP Bold (default)
- noto — Noto Sans JP Bold
- mplus — M+ 1 Black

## 改行を含む絵文字を作成

`\n` を入力に含めるか、シェルの改行入りクォート (`$'...\n...'`) を渡すと複数行になります。

```
egc '絵文\n字。'
```

## 自動 2 行レイアウト

改行を含まない 3 文字以上の入力は、Slack の絵文字メーカー風に
自動的に中央で 2 行に分割されます (例: `完全に理解した` → `完全に` / `理解した`)。
分割位置を自分で指定したい場合は `\n` を明示的に入れてください。
