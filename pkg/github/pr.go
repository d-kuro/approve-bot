package github

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	regURL  = regexp.MustCompile(`https://github.com/(.*)/(.*)/pull/(\d*)`)
	regRepo = regexp.MustCompile("github.com/(.*)/(.*)")
)

type PR struct {
	Owner  string
	Repo   string
	Number int
}

func SplitPR(prURL string, prNum int, repo string) (*PR, error) {
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
			Owner:  string(group[1]),
			Repo:   string(group[2]),
			Number: num,
		}, nil
	}
	bs := []byte(repo)
	group := regRepo.FindSubmatch(bs)
	if len(group) != 3 {
		return nil, errors.New("invalid PR format")
	}
	return &PR{
		Owner:  string(group[1]),
		Repo:   string(group[2]),
		Number: prNum,
	}, nil
}
