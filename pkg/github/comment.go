package github

import (
	"bytes"
	"strings"
	"text/template"
)

type Info struct {
	Owner    string
	OwnerURL string
	PRFiles  []string
	Patterns []string
}

func NewInfo(owner string, ownerURL string, prFiles []string, patterns []string) Info {
	return Info{
		Owner:    owner,
		OwnerURL: ownerURL,
		PRFiles:  prFiles,
		Patterns: patterns,
	}
}

func (i Info) GenerateReviewComment() (string, error) {
	tmpl, err := template.New("template").Parse(commentTemplate)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, i); err != nil {
		return "", err
	}
	return strings.ReplaceAll(buf.String(), "#", "`"), nil
}

const (
	commentTemplate = `
**[APPROVE]** Matched with Owner's Patterns

Owner: [{{ .Owner }}]({{ .OwnerURL }})

PR Files:
{{range .PRFiles }}
* #{{ . }}#{{end}}

Patterns:
{{range .Patterns }}
* #{{ . }}#{{end}}
`
)
