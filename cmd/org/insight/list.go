package insight

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
)

// NewListCmd returns a new cobra.Command for listing organization rule suites.
func NewListCmd() *cobra.Command {
	return cmdcommon.NewInsightListCmd(cmdcommon.NewOrgScope())
}
