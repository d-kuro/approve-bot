package github

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	regURL  = regexp.MustCompile(`https://github.com/(.*)/(.*)/pull/(\d*)`)
	regRepo = regexp.MustCompile("https://github.com/(.*)/(.*)")
)

type PR struct {
	owner  string
	repo   string
	number int
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
