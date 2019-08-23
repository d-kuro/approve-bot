package cmd

import (
	"io"
	"os"
	"strconv"

	"github.com/d-kuro/approve-bot/cmd/config"
	"github.com/d-kuro/approve-bot/pkg/approve"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 1
)

type Option struct {
	PRURL     string
	PRNum     int
	Token     string
	Config    string
	outStream io.Writer
	errStream io.Writer
}

func ExecuteCmd(outStream, errStream io.Writer) int {
	o := NewOption(outStream, errStream)
	cmd := NewRootCommand(o)
	addCommands(cmd, o)

	if err := cmd.Execute(); err != nil {
		red := color.New(color.FgRed)
		switch e := err.(type) {
		case ValidateError:
			red.Fprintf(errStream, "validate error: %s (exit code: 0)\n", e.Error())
			return exitCodeOK
		case approve.UnmatchedOwnerErr:
			red.Fprintf(errStream, "error: %s (exit code: 0)\n", e.Error())
			return exitCodeOK
		case approve.UnmatchedFilesErr:
			red.Fprintf(errStream, "error: %s (exit code: 0)\n", e.Error())
			return exitCodeOK
		default:
			red.Fprintf(errStream, "error: %s\n", err)
			return exitCodeErr
		}
	}
	return exitCodeOK
}

func NewOption(outStream, errStream io.Writer) *Option {
	return &Option{
		outStream: outStream,
		errStream: errStream,
	}
}

func NewRootCommand(o *Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "approve-bot",
		Short:         "Approve Pull Request of the file owner.",
		Example:       example,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := getEnv(o); err != nil {
				return err
			}
			cfg, err := config.LoadConfigFromFile(o.Config)
			if err != nil {
				return err
			}
			if err := Validate(cfg, o); err != nil {
				return err
			}
			return run(cfg, o)
		},
	}

	fset := cmd.Flags()
	fset.StringVar(&o.PRURL, "pr", "", "Pull Request URL. Or environment variable (\"CIRCLE_PULL_REQUEST\")")
	fset.IntVar(&o.PRNum, "prnum", 0, "Pull Request Number. Or environment variable (\"TRAVIS_PULL_REQUEST\")")
	fset.StringVar(&o.Token, "Token", "", "GitHub Token. Or environment variable (\"GITHUB_TOKEN\")")
	fset.StringVar(&o.Config, "Config", ".approve.yaml", "Config YAML file path.")

	return cmd
}

func addCommands(rootCmd *cobra.Command, o *Option) {
	rootCmd.AddCommand(
		NewVersionCmd(o),
	)
}

func getEnv(o *Option) error {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" && o.Token == "" {
		o.Token = token
	}
	// https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables
	if prURL := os.Getenv("CIRCLE_PULL_REQUEST"); prURL != "" && o.PRURL == "" {
		o.PRURL = prURL
	}
	// https://docs.travis-ci.com/user/environment-variables/#default-environment-variables
	if prNum := os.Getenv("TRAVIS_PULL_REQUEST"); prNum != "false" && prNum != "" && o.PRNum == 0 {
		i, err := strconv.Atoi(prNum)
		if err != nil {
			return err
		}
		o.PRNum = i
	}

	return nil
}

func run(config *config.ApproveConfig, o *Option) error {
	option := approve.NewOption(o.PRURL, o.PRNum, o.Token)
	if err := approve.Approve(option, config); err != nil {
		return err
	}

	green := color.New(color.FgGreen)
	if o.PRURL != "" {
		green.Fprintf(o.outStream, "Approved PR: %s\n", o.PRURL)
		return nil
	}
	green.Fprintf(o.outStream, "Approved PR: https://%s/pull/%d\n", config.Repo, o.PRNum)
	return nil
}

const example = `
$ approve-bot --Token <your GitHub Token for repo scope> --pr https://github.com/d-kuro/approve-bot/pull/1

.approve.yaml:
---
owners:
  - name: d-kuro
    patterns:
      - path/to/file
      - regexp
      - cmd/approve.go
      - cmd/[a-z]+.go
---

# Or specify a Pull Request number. "repo" of Config is required, when using Pull Request number.
$ approve-bot --Token <your GitHub Token for repo scope > --prnum 1

.approve.yaml:
---
repo: github.com/d-kuro/approve-bot
owners:
  - name: d-kuro
    patterns:
      - path/to/file
      - regexp
      - cmd/approve.go
      - cmd/[a-z]+.go
---
`
