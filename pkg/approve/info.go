package approve

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

type Info struct {
	Owner     string
	OwnerURL  string
	PRFiles   []string
	Patterns  []string
	outStream io.Writer
}

func NewInfo(owner string, ownerURL string, prFiles []string, patterns []string, outStream io.Writer) Info {
	return Info{
		Owner:     owner,
		OwnerURL:  ownerURL,
		PRFiles:   prFiles,
		Patterns:  patterns,
		outStream: outStream,
	}
}

func (i Info) printInfo() error {
	color.New(color.FgCyan).Fprintf(i.outStream, "Owner: %s\n\n", i.Owner)

	pr, err := template.New("PRFiles").Parse(prFileTemplate)
	if err != nil {
		return err
	}
	prbuf := &bytes.Buffer{}
	if err := pr.Execute(prbuf, i.PRFiles); err != nil {
		return err
	}
	color.New(color.FgGreen).Fprint(i.outStream, prbuf.String())

	ptn, err := template.New("Patterns").Parse(matchPatternsTemplate)
	if err != nil {
		return err
	}
	ptnbuf := &bytes.Buffer{}
	if err := ptn.Execute(ptnbuf, i.Patterns); err != nil {
		return err
	}
	color.New(color.FgYellow).Fprint(i.outStream, ptnbuf.String())

	return nil
}

func (i Info) generateReviewComment() (string, error) {
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

owner: [{{ .Owner }}]({{ .OwnerURL }})

PR Files:
{{range .PRFiles }}
* #{{ . }}#{{end}}

Patterns:
{{range .Patterns }}
* #{{ . }}#{{end}}
`

	prFileTemplate = `PR Files:
{{range .}}* {{.}}
{{end}}
`

	matchPatternsTemplate = `Patterns:
{{range .}}* {{.}}
{{end}}
`
)
