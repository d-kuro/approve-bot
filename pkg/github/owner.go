package github

import "context"

func (c *Client) GetPROwner(ctx context.Context, pr *PR) (string, string, error) {
	prInfo, _, err := c.client.PullRequests.Get(ctx, pr.Owner, pr.Repo, pr.Number)
	if err != nil {
		return "", "", err
	}

	return *prInfo.User.Login, *prInfo.User.HTMLURL, nil
}
