package repo

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewGetCmd returns a new cobra.Command for getting a repository ruleset.
func NewGetCmd() *cobra.Command {
	return cmdcommon.NewGetCmd(cmdcommon.NewRepoScope())
}
