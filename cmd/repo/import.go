package repo

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewImportCmd returns a new cobra.Command for importing a repository ruleset.
func NewImportCmd() *cobra.Command {
	return cmdcommon.NewImportCmd(cmdcommon.NewRepoScope())
}
