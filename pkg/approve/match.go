package approve

import (
	"context"
	"regexp"
	"sync"

	"github.com/d-kuro/approve-bot/cmd/config"
	"github.com/google/go-github/v26/github"
)

type UnmatchedFilesErr struct {
	msg string
}

func (e UnmatchedFilesErr) Error() string {
	return e.msg
}

func (o *Options) getOwner(ctx context.Context) (string, error) {
	pr, _, err := o.client.PullRequests.Get(ctx, o.owner, o.repo, o.number)
	if err != nil {
		return "", err
	}
	return *pr.User.Login, nil
}

func (o *Options) listFiles(ctx context.Context, nextPage int) ([]string, int, error) {
	listOps := &github.ListOptions{
		PerPage: 100,
		Page:    nextPage,
	}
	files, res, err := o.client.PullRequests.ListFiles(ctx, o.owner, o.repo, o.number, listOps)
	if err != nil {
		return nil, 0, err
	}

	f := make([]string, 0, len(files))
	for _, v := range files {
		f = append(f, *v.Filename)
	}
	return f, res.NextPage, nil
}

func (o *Options) matchFiles(ctx context.Context, prFiles []string, ownerPatterns []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error)
	wg := &sync.WaitGroup{}
	once := sync.Once{}

	regexpMap, err := compileOwnerPatterns(ownerPatterns)
	if err != nil {
		return err
	}

	for _, prf := range prFiles {
		wg.Add(1)
		go func(prf string) {
			defer wg.Done()
			if !isDone(ctx) {
				for _, p := range ownerPatterns {
					r := regexpMap[p]
					if r.MatchString(prf) {
						return
					}
				}
				once.Do(func() {
					errCh <- UnmatchedFilesErr{msg: "unmatched file: " + prf}
					cancel()
				})
			}
		}(prf)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func getOwnerPatterns(owner string, cfg *config.ApproveConfig) ([]string, error) {
	for _, o := range cfg.Owners {
		if o.Name == owner {
			return o.Patterns, nil
		}
	}
	return nil, UnmatchedOwnerErr{msg: "unmatched owner"}
}

func compileOwnerPatterns(ownerFiles []string) (map[string]*regexp.Regexp, error) {
	compiles := make(map[string]*regexp.Regexp, len(ownerFiles))
	for _, f := range ownerFiles {
		r, err := regexp.Compile(f)
		if err != nil {
			return nil, err
		}
		compiles[f] = r
	}
	return compiles, nil
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
