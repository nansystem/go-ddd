package github

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/nansystem/go-ddd/internal/config"
	"github.com/nansystem/go-ddd/internal/infrastructure/github/graphql/client"
	"github.com/nansystem/go-ddd/internal/infrastructure/github/graphql/model"
)

// Client GitHubのGraphQL APIクライアント
type Client struct {
	client *client.Client
}

// NewClient GitHubクライアントを作成する
func NewClient(cfg *config.Config) (*Client, error) {
	token := cfg.GitHub.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("GitHub token is not set")
	}

	httpClient := &http.Client{
		Transport: &tokenTransport{
			token: token,
		},
	}

	// 生成されたクライアントを初期化
	c := client.NewClient(httpClient, "https://api.github.com/graphql")

	return &Client{
		client: c,
	}, nil
}

// GetRepository リポジトリ情報を取得する
func (g *Client) GetRepository(ctx context.Context, owner, name string) (*model.GetRepositoryResponse, error) {
	// GraphQLクエリを実行
	resp, err := g.client.GetRepository(ctx, owner, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return resp, nil
}

// GetViewer ログインユーザー情報を取得する
func (g *Client) GetViewer(ctx context.Context) (*model.GetViewerResponse, error) {
	resp, err := g.client.GetViewer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get viewer: %w", err)
	}

	return resp, nil
}

// tokenTransport GitHub APIにトークンを付与するトランスポート
type tokenTransport struct {
	token string
}

// RoundTrip トークンをヘッダーに追加
func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
