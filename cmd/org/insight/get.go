package insight

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewGetCmd returns a new cobra.Command for getting an organization rule suite.
func NewGetCmd() *cobra.Command {
	return cmdcommon.NewInsightGetCmd(cmdcommon.NewOrgScope())
}
