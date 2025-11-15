package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/cmd/org"
)

// NewOrgCmd returns a new cobra.Command for organization commands
func NewOrgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Manage organization rulesets",
		Long:  `Commands to manage organization rulesets`,
	}

	cmd.AddCommand(org.NewDeleteCmd())
	cmd.AddCommand(org.NewExportCmd())
	cmd.AddCommand(org.NewGetCmd())
	cmd.AddCommand(org.NewImportCmd())
	cmd.AddCommand(org.NewInsightCmd())
	cmd.AddCommand(org.NewListCmd())
	cmd.AddCommand(org.NewMigrateCmd())

	return cmd
}

func init() {
	rootCmd.AddCommand(NewOrgCmd())
}
