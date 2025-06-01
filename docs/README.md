# Todo TUI - ドキュメント

Todo TUIは、ターミナル上で動作するtodo.txt形式対応のタスク管理ツールです。

## 📚 ドキュメント構成

- **[アーキテクチャ概要](ARCHITECTURE.md)** - プロジェクト構造と設計方針
- **[開発ガイド](DEVELOPMENT.md)** - 開発環境と貢献方法

## 🚀 クイックスタート

### インストール
```bash
go install github.com/yuucu/todotui/cmd/todotui@latest
```

### 基本使用方法
```bash
todotui ~/todo.txt
```

詳細は[メインREADME](../README.md)を参照してください。

## 🔧 開発

```bash
git clone https://github.com/yuucu/todotui.git
cd todotui
make build
```

詳細は[開発ガイド](DEVELOPMENT.md)を参照してください。
