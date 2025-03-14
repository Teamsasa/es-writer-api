## 開発環境と操作方法

このプロジェクトでは、開発とデプロイを効率化するために Makefile を使用しています。以下のコマンドが利用可能です。

### Make コマンド一覧

| コマンド | 説明 |
| --- | --- |
| `make all` | アプリケーションを起動し、マイグレーションを実行します |
| `make up` | Docker Compose でアプリケーション環境を起動します |
| `make down` | Docker Compose 環境を停止し、ボリュームを削除します |
| `make migrate` | データベースマイグレーションを実行します |
| `make prune` | 未使用の Docker イメージを削除します |
| `make fmt` | すべての Go コードをフォーマットします |
| `make test-setup` | テスト環境をセットアップします |
| `make test-repository` | リポジトリ層のテストを実行します |
| `make test-usecase` | ユースケース層のテストを実行します |
| `make test-handler` | ハンドラー層のテストを実行します |
| `make test` | すべてのテストを実行します |
| `make help` | 利用可能なコマンドの一覧とその説明を表示します |

### 使用例

### 開発環境のセットアップと起動

初めて環境をセットアップする場合や、アプリケーションを完全に起動する場合は以下のコマンドを使用します：

```bash
make all
```

このコマンドは、Docker 環境を起動し、データベースマイグレーションを実行します。

### 開発中のよく使うコマンド

開発中は以下のコマンドをよく使用します：

- アプリケーションの起動：

```bash
make up
```

- コードの変更後、フォーマットの適用：

```bash
make fmt
```

- テストの実行：

```bash
make test
```

- 開発終了時にアプリケーションを停止：

```bash
make down
```

### テスト関連の操作

特定のレイヤーだけテストしたい場合は、以下のコマンドを使用します：

```bash
# リポジトリ層のテスト
make test-repository

# ユースケース層のテスト
make test-usecase

# ハンドラー層のテスト
make test-handler
```

テスト環境のセットアップだけを行う場合：

```bash
make test-setup
```

### 注意事項

- `make down` コマンドはボリュームも削除するため、データベース内のデータも失われます
- テストを実行する前に、`make test-setup` が自動的に実行されますが、明示的に実行することもできます
- 困ったときは `make help` コマンドで利用可能なコマンドの一覧を確認できます
- 個別ファイルでtestを実行したい場合は手動でtestコマンドを叩いてください