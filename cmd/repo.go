package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/cmd/repo"
)

// NewRepoCmd returns a new cobra.Command for repository commands
func NewRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage repository rulesets, branch protection, and related conversions",
		Long:  `Commands to manage repository rulesets, branch protection, and related conversions`,
	}

	cmd.AddCommand(repo.NewBranchProtectionCmd())
	cmd.AddCommand(repo.NewDeleteCmd())
	cmd.AddCommand(repo.NewExportCmd())
	cmd.AddCommand(repo.NewFromBranchProtectionCmd())
	cmd.AddCommand(repo.NewFromTagProtectionCmd())
	cmd.AddCommand(repo.NewGetCmd())
	cmd.AddCommand(repo.NewImportCmd())
	cmd.AddCommand(repo.NewInsightCmd())
	cmd.AddCommand(repo.NewListCmd())
	cmd.AddCommand(repo.NewMigrateCmd())
	cmd.AddCommand(repo.NewTagProtectionCmd())

	return cmd
}

func init() {
	rootCmd.AddCommand(NewRepoCmd())
}
