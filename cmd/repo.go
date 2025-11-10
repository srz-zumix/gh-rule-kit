package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/cmd/repo"
)

// NewRepoCmd returns a new cobra.Command for repository commands
func NewRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage repository rulesets",
		Long:  `Commands to manage repository rulesets`,
	}

	cmd.AddCommand(repo.NewDeleteCmd())
	cmd.AddCommand(repo.NewExportCmd())
	cmd.AddCommand(repo.NewGetCmd())
	cmd.AddCommand(repo.NewImportCmd())
	cmd.AddCommand(repo.NewListCmd())
	cmd.AddCommand(repo.NewMigrateCmd())

	return cmd
}

func init() {
	rootCmd.AddCommand(NewRepoCmd())
}
