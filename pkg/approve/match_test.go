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
	ownerMatchPatterns = []string{
		"cmd/approve.go",
		"cmd/version.go",
		"cmd/validate.go",
	}
	ownerUnmatchedPatterns = []string{
		"cmd/foo.go",
	}
	ownerUnmatchedRegexpPatterns = []string{
		"cmd/[A-Z]+",
	}
)

func Test_matchFiles(t *testing.T) {
	cases := []struct {
		name          string
		prFiles       []string
		ownerPatterns []string
		success       bool
	}{
		{
			name:          "match files",
			prFiles:       prFiles,
			ownerPatterns: ownerMatchPatterns,
			success:       true,
		},
		{
			name:          "unmatched files",
			prFiles:       prFiles,
			ownerPatterns: ownerUnmatchedPatterns,
			success:       false,
		},
		{
			name:          "unmatched regexp files",
			prFiles:       prFiles,
			ownerPatterns: ownerUnmatchedRegexpPatterns,
			success:       false,
		},
	}

	o := Options{}
	ctx := context.Background()

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := o.matchFiles(ctx, tt.prFiles, tt.ownerPatterns)
			if err != nil && tt.success {
				t.Errorf("error: %s", err)
			} else if err == nil && !tt.success {
				t.Error("expect errors")
			}
		})
	}
}
