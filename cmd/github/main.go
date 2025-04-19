package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Yamashou/gqlgenc/clientv2"

	"github.com/nansystem/go-ddd/internal/config"
	"github.com/nansystem/go-ddd/internal/infrastructure/github/gen"
)

func main() {
	// 環境変数GITHUB_TOKENが設定されていることを確認してください
	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	// 設定を読み込む
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// GitHubクライアントを初期化
	token := os.Getenv("GITHUB_TOKEN") // トークンは既に上でチェック済み

	// Authorizationヘッダーを追加するためのインターセプターを作成
	authInterceptor := clientv2.RequestInterceptor(func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		return next(ctx, req, gqlInfo, res)
	})

	// clientv2.Options 構造体には AddHeader フィールドが存在しません。
	// デフォルトのHTTPクライアントと標準のGitHub GraphQLエンドポイントを使用
	// gen.NewClientはエラーを返さないため、エラーチェックは不要
	githubClient := gen.NewClient(http.DefaultClient, "https://api.github.com/graphql", nil, authInterceptor)

	// まずは自分のユーザー情報を取得
	ctx := context.Background()
	viewer, err := githubClient.GetViewer(ctx)
	if err != nil {
		log.Fatalf("ユーザー情報の取得に失敗しました: %v", err)
	}
	fmt.Printf("Logged in as: %s\n\n", viewer.Viewer.Login)

	// コマンドライン引数からリポジトリのオーナー名とリポジトリ名を取得
	owner := "Yamashou"
	repo := "gqlgenc"

	// コマンドライン引数があれば上書き
	if len(os.Args) >= 3 {
		owner = os.Args[1]
		repo = os.Args[2]
	}

	// リポジトリ情報を取得
	response, err := githubClient.GetRepository(ctx, owner, repo)
	if err != nil {
		log.Fatalf("リポジトリ情報の取得に失敗しました: %v", err)
	}

	// 結果を表示
	fmt.Printf("Repository: %s/%s\n", owner, repo)
	description := response.Repository.Description
	if description == nil {
		defaultDescription := "(説明なし)" // 説明が空の場合のデフォルトテキスト
		description = &defaultDescription
	}
	fmt.Printf("Description: %s\n", *description)
	fmt.Printf("Stars: %d\n", response.Repository.StargazerCount)
	fmt.Printf("Forks: %d\n", response.Repository.ForkCount)
	fmt.Printf("URL: %s\n", response.Repository.URL)

	fmt.Printf("\nOpen Issues (%d total):\n", response.Repository.Issues.TotalCount)
	for _, issue := range response.Repository.Issues.Nodes {
		fmt.Printf("- %s (by @%s): %s\n", issue.Title, issue.Author.Login, issue.URL)
	}

	fmt.Printf("\nOpen Pull Requests (%d total):\n", response.Repository.PullRequests.TotalCount)
	for _, pr := range response.Repository.PullRequests.Nodes {
		fmt.Printf("- %s (by @%s): %s\n", pr.Title, pr.Author.Login, pr.URL)
	}

	// JSON形式で出力
	fmt.Println("\nFull JSON response:")
	jsonData, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(jsonData))
}
