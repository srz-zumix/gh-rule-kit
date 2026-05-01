package org

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewExportCmd returns a new cobra.Command for exporting an organization ruleset.
func NewExportCmd() *cobra.Command {
	return cmdcommon.NewExportCmd(cmdcommon.NewOrgScope())
}
