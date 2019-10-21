package github

import (
	"context"

	"github.com/google/go-github/v26/github"
)

func (c *Client) CreatePRReview(ctx context.Context, comment string, pr *PR) error {
	event := "APPROVE"
	review := &github.PullRequestReviewRequest{
		Event: &event,
		Body:  &comment,
	}

	if _, _, err := c.client.PullRequests.CreateReview(ctx, pr.Owner, pr.Repo, pr.Number, review); err != nil {
		return err
	}

	return nil
}
