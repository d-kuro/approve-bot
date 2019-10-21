package github

import (
	"context"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

func NewClient(ctx context.Context, token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{client: github.NewClient(tc)}
}
