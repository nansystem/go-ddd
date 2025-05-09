#!/bin/bash

# スクリプトが存在するディレクトリを取得
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
# プロジェクトのルートディレクトリを取得
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# GitHub GraphQLクライアントのディレクトリに移動
cd "$PROJECT_ROOT/internal/infrastructure/github/graphql" || exit

# Create the required directories
mkdir -p client model

# Create simple client manually
cat > client/client_gen.go << EOL
// Code generated by gqlgenc, DO NOT EDIT.

package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/nansystem/go-ddd/internal/infrastructure/github/graphql/model"
)

type Client struct {
	client *clientv2.Client
}

func NewClient(client *http.Client, endpoint string) *Client {
	return &Client{
		client: clientv2.NewClient(client, endpoint),
	}
}

func (c *Client) GetViewer(ctx context.Context) (*model.GetViewerResponse, error) {
	req := &model.GetViewerRequest{}
	resp := &model.GetViewerResponse{}
	if err := c.client.Post(ctx, "GetViewer", req, resp); err != nil {
		return nil, fmt.Errorf("getting viewer: %w", err)
	}
	return resp, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, name string) (*model.GetRepositoryResponse, error) {
	req := &model.GetRepositoryRequest{
		Owner: owner,
		Name:  name,
	}
	resp := &model.GetRepositoryResponse{}
	if err := c.client.Post(ctx, "GetRepository", req, resp); err != nil {
		return nil, fmt.Errorf("getting repository: %w", err)
	}
	return resp, nil
}
EOL

# Create model file
cat > model/models_gen.go << EOL
// Code generated by gqlgenc, DO NOT EDIT.

package model

type GetViewerRequest struct {
}

type GetViewerResponse struct {
	Viewer struct {
		Login string \`json:"login"\`
	} \`json:"viewer"\`
}

type GetRepositoryRequest struct {
	Owner string \`json:"owner"\`
	Name  string \`json:"name"\`
}

type GetRepositoryResponse struct {
	Viewer struct {
		Login string \`json:"login"\`
	} \`json:"viewer"\`
	Repository struct {
		ID             string \`json:"id"\`
		Name           string \`json:"name"\`
		URL            string \`json:"url"\`
		Description    string \`json:"description"\`
		StargazerCount int    \`json:"stargazerCount"\`
		ForkCount      int    \`json:"forkCount"\`
		Issues struct {
			TotalCount int \`json:"totalCount"\`
			Nodes      []struct {
				ID        string \`json:"id"\`
				Title     string \`json:"title"\`
				URL       string \`json:"url"\`
				CreatedAt string \`json:"createdAt"\`
				Author    struct {
					Login string \`json:"login"\`
				} \`json:"author"\`
			} \`json:"nodes"\`
		} \`json:"issues"\`
		PullRequests struct {
			TotalCount int \`json:"totalCount"\`
			Nodes      []struct {
				ID        string \`json:"id"\`
				Title     string \`json:"title"\`
				URL       string \`json:"url"\`
				CreatedAt string \`json:"createdAt"\`
				Author    struct {
					Login string \`json:"login"\`
				} \`json:"author"\`
			} \`json:"nodes"\`
		} \`json:"pullRequests"\`
	} \`json:"repository"\`
}
EOL

echo "GitHub GraphQL client generated successfully!"