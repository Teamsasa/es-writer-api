# ES Writer API

## 概要

ES Writer API は Go で実装されたES自動生成のための REST API サービスです。ユーザーの経験談を保存・取得するための機能を提供します。

主な特徴:

- Go 言語（v1.23）を使用
- Echo フレームワークによる REST API 実装
- PostgreSQL データベースを使用したデータ永続化
- Clerk 認証システムとの統合
- Docker/Docker Compose による開発・デプロイ環境
- クリーンアーキテクチャに基づいた設計

## ディレクトリ構造

```
.
├── app/                    # アプリケーションコード
│   ├── cmd/                # エントリーポイント
│   │   ├── migrate/        # DBマイグレーションツール
│   │   └── server/         # APIサーバー
│   ├── infrastructure/     # インフラストラクチャ層
│   ├── internal/           # 内部パッケージ
│   │   ├── entity/         # ドメインエンティティ
│   │   ├── handler/        # HTTPハンドラー
│   │   ├── repository/     # データリポジトリ
│   │   ├── router/         # ルーティング
│   │   └── usecase/        # ユースケース
│   ├── middleware/         # ミドルウェア
│   └── test/               # テストコード
├── docs/                   # ドキュメント
├── schema/                 # OpenAPI仕様
│   └── openapi.yml         # API定義
├── .env                    # 環境変数
├── docker-compose.yml      # Docker構成
├── Dockerfile              # Docker定義
├── go.mod                  # Goモジュール定義
└── Makefile                # タスク自動化
```

## セットアップ

### 必要条件

- Go 1.23 以上
- Docker と Docker Compose

### 開発環境の起動

```bash
# リポジトリをクローン
git clone [repository-url]
cd es-writer-api

# .envファイルを編集して必要な変数を設定

# Dockerコンテナの起動
# APIサーバーの起動
make up
```

## API ドキュメント

API 仕様は `/schema/openapi.yml` に定義されています。

## テスト

```bash
make test
```
