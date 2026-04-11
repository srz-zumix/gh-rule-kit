package repo

import (
	"github.com/spf13/cobra"
	tagprotection "github.com/srz-zumix/gh-rule-kit/cmd/repo/tag_protection"
)

// NewTagProtectionCmd returns a new cobra.Command for tag protection commands.
func NewTagProtectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag-protection",
		Short: "Manage tag protection rules",
		Long:  `Commands to manage tag protection rules for a repository`,
	}

	cmd.AddCommand(tagprotection.NewDeleteCmd())
	cmd.AddCommand(tagprotection.NewGetCmd())
	cmd.AddCommand(tagprotection.NewListCmd())

	return cmd
}
