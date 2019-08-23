package cmd

import (
	"github.com/d-kuro/approve-bot/cmd/config"
)

type ValidateError struct {
	msg string
}

func (e ValidateError) Error() string {
	return e.msg
}

func Validate(c *config.ApproveConfig, o *Option) error {
	if err := validatePR(c, o); err != nil {
		return err
	}
	if err := validateToken(o); err != nil {
		return err
	}
	return nil
}

func validatePR(c *config.ApproveConfig, o *Option) error {
	if o.PRURL == "" || c.Repo == "" && o.PRNum == 0 {
		return ValidateError{msg: "PR URL or repo URL and the PR number is required"}
	}
	return nil
}

func validateToken(o *Option) error {
	if o.Token == "" {
		return ValidateError{msg: "Token is required"}
	}
	return nil
}
