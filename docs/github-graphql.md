# GitHub GraphQL APIの使い方

このプロジェクトでは、[gqlgenc](https://github.com/Yamashou/gqlgenc)を使ってGitHubのGraphQL APIにアクセスする実装例を提供しています。

## セットアップ

1. GitHubのパーソナルアクセストークンを取得する
   - [GitHub Personal Access Tokens](https://github.com/settings/tokens)から取得してください
   - 必要なスコープ: `repo`, `read:user`

2. 環境変数に設定する
   ```sh
   export GITHUB_TOKEN=your_token_here
   ```

## クライアントコードの生成

1. 設定ファイルは `internal/infrastructure/github/generated/.gqlgenc.yml` にあります。
   - この設定ファイルは、ローカルにダウンロードされたGitHub GraphQLスキーマを参照します。
   - スキーマは以下のコマンドで取得・更新できます:
     ```sh
     curl -L -o internal/infrastructure/github/generated/schema/github_schema.graphql https://docs.github.com/public/schema.docs.graphql
     ```
2. GraphQLクエリは `internal/infrastructure/github/generated/query/*.graphql` に定義されています
3. `internal/infrastructure/github/generated` ディレクトリに移動し、以下のコマンドでクライアントコードを生成します:
   ```sh
   gqlgenc generate
   ```

   これにより以下のファイルが生成されます:
   - `internal/infrastructure/github/generated/models_gen.go`
   - `internal/infrastructure/github/generated/client.go`

## 使い方

サンプルコードは `cmd/github/main.go` にあります。実行方法:

```sh
# GITHUB_TOKENを設定していることを確認
export GITHUB_TOKEN=your_token_here

# デフォルトでYamashou/gqlgencリポジトリの情報を取得
go run cmd/github/main.go

# 引数でリポジトリを指定することも可能
go run cmd/github/main.go octocat hello-world
```

## 新しいクエリの追加方法

1. `internal/infrastructure/github/generated/query/` に新しい `.graphql` ファイルを追加
2. クエリを定義（GitHubのGraphQL APIドキュメントを参照）
3. `internal/infrastructure/github/generated` ディレクトリで `gqlgenc generate` を実行してクライアントコードを生成
4. 生成されたクライアントを `github.GitHubClient` を通して利用

## GitHub GraphQL APIの詳細情報

- [GitHub GraphQL API v4](https://docs.github.com/ja/graphql)
- [Explorer](https://docs.github.com/ja/graphql/overview/explorer)（クエリのテストに使用できます）