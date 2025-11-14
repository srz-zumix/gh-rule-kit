package org

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/cmd/org/insight"
)

// NewInsightCmd returns a new cobra.Command for organization rule suite insights commands
func NewInsightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insight",
		Short: "Manage organization rule suite insights",
		Long:  `Commands to view organization rule suite insights and evaluations`,
	}

	cmd.AddCommand(insight.NewGetCmd())
	cmd.AddCommand(insight.NewListCmd())

	return cmd
}
