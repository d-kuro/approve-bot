package approve

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"github.com/d-kuro/approve-bot/cmd/config"
)

type UnmatchedFilesErr struct {
	msg string
}

func (e UnmatchedFilesErr) Error() string {
	return e.msg
}

type UnmatchedOwnerErr struct {
	msg string
}

func (e UnmatchedOwnerErr) Error() string {
	return e.msg
}

func matchFiles(ctx context.Context, prFiles []string, ownerPatterns []string) error {
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

func getOwnerPatterns(owner string, config *config.ApproveConfig) ([]string, error) {
	for _, o := range config.Owners {
		if o.Name == owner {
			return o.Patterns, nil
		}
	}
	return nil, UnmatchedOwnerErr{msg: fmt.Sprintf("unmatched owner: %s", owner)}
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
