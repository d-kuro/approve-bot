package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/d-kuro/approve-bot/cmd/config"
	"github.com/d-kuro/approve-bot/pkg/approve"
	"github.com/spf13/cobra"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 1
)

type Option struct {
	prURL     string
	prNum     int
	token     string
	config    string
	outStream io.Writer
	errStream io.Writer
}

func Execute(outStream, errStream io.Writer) int {
	o := NewOption(outStream, errStream)
	cmd := NewRootCommand(o)
	addCommands(cmd, o)

	if err := cmd.Execute(); err != nil {
		switch e := err.(type) {
		case ValidateError:
			fmt.Fprintf(errStream, "validate error: %s (exit code: 0)\n", e.Error())
			return exitCodeOK
		case approve.UnmatchedOwnerErr:
			fmt.Fprintf(errStream, "error: %s (exit code: 0)\n", e.Error())
			return exitCodeOK
		case approve.UnmatchedFilesErr:
			fmt.Fprintf(errStream, "error: %s (exit code: 0)\n", e.Error())
			return exitCodeOK
		default:
			fmt.Fprintf(errStream, "error: %s\n", err)
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
			cfg, err := config.GetConfig(o.config)
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
	fset.StringVar(&o.prURL, "pr", "", "Pull Request URL. Or environment variable (\"CIRCLE_PULL_REQUEST\")")
	fset.IntVar(&o.prNum, "prnum", 0, "Pull Request Number. Or environment variable (\"TRAVIS_PULL_REQUEST\")")
	fset.StringVar(&o.token, "token", "", "GitHub token. Or environment variable (\"GITHUB_TOKEN\")")
	fset.StringVar(&o.config, "config", ".approve.yaml", "Config YAML file path.")

	return cmd
}

func addCommands(rootCmd *cobra.Command, o *Option) {
	rootCmd.AddCommand(
		NewVersionCmd(o),
	)
}

func getEnv(o *Option) error {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		o.token = token
	}
	// https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables
	if prURL := os.Getenv("CIRCLE_PULL_REQUEST"); prURL != "" {
		o.prURL = prURL
	}
	// https://docs.travis-ci.com/user/environment-variables/#default-environment-variables
	if prNum := os.Getenv("TRAVIS_PULL_REQUEST"); prNum != "false" && prNum != "" {
		i, err := strconv.Atoi(prNum)
		if err != nil {
			return err
		}
		o.prNum = i
	}

	return nil
}

func run(cfg *config.ApproveConfig, o *Option) error {
	if err := approve.Approve(o.token, o.prURL, o.prNum, cfg); err != nil {
		return err
	}

	if o.prURL != "" {
		fmt.Fprintf(o.outStream, "approved PR: %s", o.prURL)
		return nil
	}
	fmt.Fprintf(o.outStream, "approved PR: https://%s/pull/%d", cfg.Repo, o.prNum)
	return nil
}

const example = `
$ approve-bot --config approve.yaml --token <your GitHub token for repo scope> \
  --pr https://github.com/d-kuro/approve-bot/pull/1

approve.yaml:
---
owners:
  - name: d-kuro
    files:
      - /path/to/file1
      - /path/to/file2
---

# Or specify a Pull Request number. "repo" of config is required, when using Pull Request number.
$ approve-bot --config approve.yaml --token <your GitHub token for repo scope > --prnum 1

approve.yaml:
---
repo: github.com/d-kuro/approve-bot
owners:
  - name: d-kuro
    files:
      - /path/to/file1
      - /path/to/file2
---
`
