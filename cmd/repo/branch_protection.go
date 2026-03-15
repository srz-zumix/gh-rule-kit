package repo

import (
	"github.com/spf13/cobra"
	branchprotection "github.com/srz-zumix/gh-rule-kit/cmd/repo/branch_protection"
)

// NewBranchProtectionCmd returns a new cobra.Command for branch protection commands
func NewBranchProtectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch-protection",
		Short: "Manage branch protection rules",
		Long:  `Commands to manage branch protection rules for a repository`,
	}

	cmd.AddCommand(branchprotection.NewDeleteCmd())
	cmd.AddCommand(branchprotection.NewGetCmd())
	cmd.AddCommand(branchprotection.NewListCmd())

	return cmd
}
