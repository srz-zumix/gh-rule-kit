package repo

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/cmd/repo/insight"
)

// NewInsightCmd returns a new cobra.Command for repository rule suite insights commands
func NewInsightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insight",
		Short: "Manage repository rule suite insights",
		Long:  `Commands to view repository rule suite insights and evaluations`,
	}

	cmd.AddCommand(insight.NewGetCmd())
	cmd.AddCommand(insight.NewListCmd())

	return cmd
}
