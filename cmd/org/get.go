package org

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewGetCmd returns a new cobra.Command for getting an organization ruleset.
func NewGetCmd() *cobra.Command {
	return cmdcommon.NewGetCmd(cmdcommon.NewOrgScope())
}
