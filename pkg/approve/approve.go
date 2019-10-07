package approve

import (
	"context"

	"github.com/d-kuro/approve-bot/pkg/github"

	"github.com/d-kuro/approve-bot/cmd/config"
)

type Option struct {
	PRURL string
	PRNum int
	Token string
}

func NewOption(url string, num int, token string) Option {
	return Option{
		PRURL: url,
		PRNum: num,
		Token: token,
	}
}

func Approve(o Option, c *config.ApproveConfig) error {
	ctx := context.Background()
	client := github.NewClient(ctx, o.Token)
	pr, err := github.SplitPR(o.PRURL, o.PRNum, c.Repo)
	if err != nil {
		return err
	}

	owner, ownerURL, err := client.GetPROwner(ctx, pr)
	if err != nil {
		return err
	}
	ownerPatterns, err := getOwnerPatterns(owner, c)
	if err != nil {
		return err
	}
	prFiles, err := client.ListPRFiles(ctx, pr)
	if err != nil {
		return err
	}
	if err := matchFiles(ctx, prFiles, ownerPatterns); err != nil {
		return err
	}

	info := github.NewInfo(owner, ownerURL, prFiles, ownerPatterns)
	comment, err := info.GenerateReviewComment()
	if err != nil {
		return err
	}
	if err := client.CreatePRReview(ctx, comment, pr); err != nil {
		return err
	}
	return nil
}
