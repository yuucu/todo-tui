# 📦 Release Guide

## 🚀 Overview

このガイドでは `todotui` のリリースプロセスについて説明します。Goreleaser を使用して、クロスプラットフォームバイナリの自動作成とGitHub Releasesの管理を行います。

## 🔧 Prerequisites

1. **Git**: コミット権限
2. **GoReleaser**: ローカルテスト用（オプション）
3. **GitHub Token**: Actions での自動リリース用

## 📝 Release Process

### 1. 開発とテスト

```bash
# 開発とテスト
make all
make test
make lint

# GoReleaser設定のテスト
make release-test

# スナップショットビルドのテスト
make release-snapshot
```

### 2. バージョンの準備

1. **CHANGELOG.md** を更新（該当する場合）
2. バージョン番号を決定（セマンティックバージョニング: `vX.Y.Z`）

### 3. リリースタグの作成

```bash
# タグを作成
git tag v1.0.0

# タグをプッシュ（これによりGitHub Actionsが自動実行される）
git push origin v1.0.0
```

### 4. GitHub Actions の監視

1. [GitHub Actions](https://github.com/yuucu/todotui/actions) でリリースワークフローを確認
2. 成功すると以下が自動生成される：
   - GitHub Release
   - クロスプラットフォームバイナリ
   - チェックサム
   - リリースノート

### 5. リリース後の確認

1. [GitHub Releases](https://github.com/yuucu/todotui/releases) でリリースを確認
2. ダウンロードリンクをテスト
3. インストール方法をREADMEで確認

## 🧪 Local Testing

リリース前のローカルテスト：

```bash
# GoReleaserのインストール
make goreleaser-install

# 設定の検証
make release-test

# スナップショットビルド（実際のリリースなし）
make release-snapshot

# 生成されたファイルを確認
ls dist/
```

## 📋 Supported Platforms

現在サポートしているプラットフォーム：

- **Linux**: x86_64, ARM64
- **macOS**: x86_64 (Intel), ARM64 (Apple Silicon)
- **Windows**: x86_64

## 🔄 Rollback Process

問題があるリリースのロールバック：

1. **GitHub Releases** でリリースを削除
2. Gitタグを削除:
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```
3. 必要に応じて修正版をリリース

## 🍺 Future: Homebrew Distribution

Homebrew配布の設定については `scripts/brew-setup.md` を参照してください。

## 📚 References

- [Semantic Versioning](https://semver.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github) 