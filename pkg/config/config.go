package config

import (
	"os/user"

	"github.com/cidverse/go-rules/pkg/expr"
	"github.com/rs/zerolog/log"
)

type DotfilesConfig struct {
	Directories []Dir `yaml:"directories"`
}

type Dir struct {
	Path       string      `yaml:"path"`
	Target     string      `yaml:"target"`
	Conditions []Condition `yaml:"conditions"` // At least one condition must match for the rule to apply
}

type Condition struct {
	Rule string `yaml:"rule"`
}

func EvaluateConditions(conditions []Condition) bool {
	if len(conditions) == 0 {
		return true
	}

	// context
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get current user")
	}
	ctx := map[string]interface{}{
		"user": currentUser.Username,
	}

	// evaluate
	for _, c := range conditions {
		match, cErr := expr.EvaluateRule(c.Rule, ctx)
		if cErr != nil {
			log.Fatal().Err(cErr).Str("rule", c.Rule).Msg("failed to evaluate condition, check your configuration file syntax")
		}
		if match {
			return true
		}
	}

	return false
}
