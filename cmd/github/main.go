package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nansystem/go-ddd/internal/config"
	"github.com/nansystem/go-ddd/internal/infrastructure/github"
)

func main() {
	// 環境変数GITHUB_TOKENが設定されていることを確認してください
	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	// 設定を読み込む
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// GitHubクライアントを初期化
	githubClient, err := github.NewClient(cfg)
	if err != nil {
		log.Fatalf("GitHubクライアントの初期化に失敗しました: %v", err)
	}

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
	fmt.Printf("Description: %s\n", response.Repository.Description)
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
