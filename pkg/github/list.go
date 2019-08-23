package github

import (
	"context"

	"github.com/google/go-github/v26/github"
)

func (c *Client) ListPRFiles(ctx context.Context, pr *PR) ([]string, error) {
	prFiles := make([]string, 0)
	var next = 1
	for {
		if next == 0 {
			break
		}
		f, i, err := listFiles(ctx, c.client, pr, next)
		if err != nil {
			return nil, err
		}
		prFiles = append(prFiles, f...)
		next = i
	}
	return prFiles, nil
}

func listFiles(ctx context.Context, client *github.Client, pr *PR, nextPage int) ([]string, int, error) {
	listOps := &github.ListOptions{
		PerPage: 100,
		Page:    nextPage,
	}
	files, res, err := client.PullRequests.ListFiles(ctx, pr.Owner, pr.Repo, pr.Number, listOps)
	if err != nil {
		return nil, 0, err
	}

	f := make([]string, 0, len(files))
	for _, v := range files {
		f = append(f, *v.Filename)
	}
	return f, res.NextPage, nil
}
