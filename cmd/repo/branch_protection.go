package repo

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/cmd/repo/branch_protection"
)

// NewBranchProtectionCmd returns a new cobra.Command for branch protection commands
func NewBranchProtectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch-protection",
		Short: "Manage branch protection rules",
		Long:  `Commands to manage branch protection rules for a repository`,
	}

	cmd.AddCommand(branch_protection.NewDeleteCmd())
	cmd.AddCommand(branch_protection.NewGetCmd())
	cmd.AddCommand(branch_protection.NewListCmd())

	return cmd
}
