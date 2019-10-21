package github_test

import (
	"reflect"
	"testing"

	"github.com/d-kuro/approve-bot/pkg/github"
)

func TestSplitPR(t *testing.T) {
	cases := []struct {
		name string
		url  string
		num  int
		repo string
		pr   *github.PR
	}{
		{
			name: "Specify PR URL",
			url:  "https://github.com/d-kuro/approve-bot/pull/1",
			pr:   &github.PR{Owner: "d-kuro", Repo: "approve-bot", Number: 1},
		},
		{
			name: "Specify PR number and repo",
			repo: "github.com/d-kuro/approve-bot",
			num:  1,
			pr:   &github.PR{Owner: "d-kuro", Repo: "approve-bot", Number: 1},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			pr, err := github.SplitPR(tt.url, tt.num, tt.repo)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			if !reflect.DeepEqual(pr, tt.pr) {
				t.Errorf("got: %#v, want: %#v", pr, tt.pr)
			}
		})
	}
}
