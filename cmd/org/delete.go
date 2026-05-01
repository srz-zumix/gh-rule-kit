package org

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewDeleteCmd returns a new cobra.Command for deleting an organization ruleset.
func NewDeleteCmd() *cobra.Command {
	return cmdcommon.NewDeleteCmd(cmdcommon.NewOrgScope())
}
