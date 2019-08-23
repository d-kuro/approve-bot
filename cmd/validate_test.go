package cmd_test

import (
	"testing"

	"github.com/d-kuro/approve-bot/cmd"
	"github.com/d-kuro/approve-bot/cmd/config"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name   string
		config *config.ApproveConfig
		option *cmd.Option
		msg    string
	}{
		{
			name:   "PR URL is required",
			config: &config.ApproveConfig{},
			option: &cmd.Option{Token: "token"},
			msg:    "PR URL or repo URL and the PR number is required",
		},
		{
			name:   "Repo URL and the PR number is required",
			config: &config.ApproveConfig{Repo: "github.com/d-kuro/approve-bot"},
			option: &cmd.Option{PRNum: 1, Token: "token"},
			msg:    "PR URL or repo URL and the PR number is required",
		},
		{
			name:   "Token is required",
			config: &config.ApproveConfig{Repo: "github.com/d-kuro/approve-bot"},
			option: &cmd.Option{PRURL: "https://github.com/d-kuro/approve-bot/pull/1"},
			msg:    "Token is required",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Validate(tt.config, tt.option)
			if err == nil {
				t.Fatalf("expect an error")
			}
			switch e := err.(type) {
			case cmd.ValidateError:
				if e.Error() != tt.msg {
					t.Errorf("got: %v, want: %v", tt.msg, e.Error())
				}
			}
		})
	}
}
