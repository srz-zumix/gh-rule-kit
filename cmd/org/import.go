package org

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewImportCmd returns a new cobra.Command for importing an organization ruleset.
func NewImportCmd() *cobra.Command {
	return cmdcommon.NewImportCmd(cmdcommon.NewOrgScope())
}
