package approve

import (
	"bytes"
	"context"
	"errors"
	"regexp"
	"strconv"
	"text/template"

	"github.com/d-kuro/approve-bot/cmd/config"
	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

var (
	regURL  = regexp.MustCompile(`https://github.com/(.*)/(.*)/pull/(\d*)`)
	regRepo = regexp.MustCompile("https://github.com/(.*)/(.*)")
)

type Options struct {
	client *github.Client
	*PR
}

type PR struct {
	owner  string
	repo   string
	number int
}

type UnmatchedOwnerErr struct {
	msg string
}

func (e UnmatchedOwnerErr) Error() string {
	return e.msg
}

func Approve(token, prURL string, prNum int, cfg *config.ApproveConfig) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	pr, err := splitPR(prURL, prNum, cfg.Repo)
	if err != nil {
		return err
	}

	o := NewOptions(client, pr)
	owner, err := o.getOwner(ctx)
	if err != nil {
		return err
	}
	ownerFiles, err := getOwnerPatterns(owner, cfg)
	if err != nil {
		return err
	}

	prFiles := make([]string, 0)
	var next = 1
	for {
		if next == 0 {
			break
		}
		f, i, err := o.listFiles(ctx, next)
		if err != nil {
			return err
		}
		prFiles = append(prFiles, f...)
		next = i
	}

	if err := o.matchFiles(ctx, prFiles, ownerFiles); err != nil {
		return err
	}
	if err := o.createPRReview(ctx, ownerFiles); err != nil {
		return err
	}
	return nil
}

func (o *Options) createPRReview(ctx context.Context, ownerFiles []string) error {
	tmpl, err := template.New("template").Parse(msgTemplate)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, ownerFiles); err != nil {
		return err
	}

	event := "APPROVE"
	comment := buf.String()
	review := &github.PullRequestReviewRequest{
		Event: &event,
		Body:  &comment,
	}
	if _, _, err := o.client.PullRequests.CreateReview(ctx, o.owner, o.repo, o.number, review); err != nil {
		return err
	}
	return nil
}

func NewOptions(client *github.Client, pr *PR) *Options {
	return &Options{
		client: client,
		PR:     pr,
	}
}

func splitPR(prURL string, prNum int, repo string) (*PR, error) {
	if prURL != "" {
		bs := []byte(prURL)
		group := regURL.FindSubmatch(bs)
		if len(group) != 4 {
			return nil, errors.New("invalid PR URL: " + prURL)
		}
		num, err := strconv.Atoi(string(group[3]))
		if err != nil {
			return nil, err
		}
		return &PR{
			owner:  string(group[1]),
			repo:   string(group[2]),
			number: num,
		}, nil
	}
	bs := []byte(repo)
	group := regRepo.FindSubmatch(bs)
	if len(group) != 3 {
		return nil, errors.New("")
	}
	return &PR{
		owner:  string(group[1]),
		repo:   string(group[2]),
		number: prNum,
	}, nil
}

const msgTemplate = `
**[APPROVE]** Matched with owner's file:
{{range .}}
* {{.}}{{end}}`
