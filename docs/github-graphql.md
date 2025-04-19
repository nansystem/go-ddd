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

1. 設定ファイルは `internal/infrastructure/github/graphql/gqlgenc.yml` にあります
2. GraphQLクエリは `internal/infrastructure/github/query/*.graphql` に定義されています
3. 以下のコマンドでクライアントコードを生成します:

   ```sh
   ./scripts/generate-github-client.sh
   ```

   これにより以下のファイルが生成されます:
   - `internal/infrastructure/github/graphql/model/models_gen.go`
   - `internal/infrastructure/github/graphql/client/client_gen.go`

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

1. `internal/infrastructure/github/query/` に新しい `.graphql` ファイルを追加
2. クエリを定義（GitHubのGraphQL APIドキュメントを参照）
3. `./scripts/generate-github-client.sh` を実行してクライアントコードを生成
4. 生成されたクライアントを `github.GitHubClient` を通して利用

## GitHub GraphQL APIの詳細情報

- [GitHub GraphQL API v4](https://docs.github.com/ja/graphql)
- [Explorer](https://docs.github.com/ja/graphql/overview/explorer)（クエリのテストに使用できます）