package approve

import (
	"context"
	"testing"
)

var (
	prFiles = []string{
		"cmd/approve.go",
		"cmd/version.go",
		"cmd/validate.go",
	}
	ownerMatchFiles = []string{
		"cmd/approve.go",
		"cmd/version.go",
		"cmd/validate.go",
	}
	ownerUnmatchedFiles = []string{
		"cmd/foo.go",
	}
	ownerUnmatchedRegexpFiles = []string{
		"cmd/[A-Z]+",
	}
)

func Test_matchFiles(t *testing.T) {
	cases := []struct {
		name       string
		prFiles    []string
		ownerFiles []string
		success    bool
	}{
		{
			name:       "match files",
			prFiles:    prFiles,
			ownerFiles: ownerMatchFiles,
			success:    true,
		},
		{
			name:       "unmatched files",
			prFiles:    prFiles,
			ownerFiles: ownerUnmatchedFiles,
			success:    false,
		},
		{
			name:       "unmatched regexp files",
			prFiles:    prFiles,
			ownerFiles: ownerUnmatchedRegexpFiles,
			success:    false,
		},
	}

	o := Options{}
	ctx := context.Background()

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := o.matchFiles(ctx, tt.prFiles, tt.ownerFiles)
			if err != nil && tt.success {
				t.Errorf("error: %s", err)
			} else if err == nil && !tt.success {
				t.Error("expect errors")
			}
		})
	}
}
