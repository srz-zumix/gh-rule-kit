package repo

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewExportCmd returns a new cobra.Command for exporting a repository ruleset.
func NewExportCmd() *cobra.Command {
	return cmdcommon.NewExportCmd(cmdcommon.NewRepoScope())
}
