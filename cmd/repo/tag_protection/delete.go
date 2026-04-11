package tag_protection

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// NewDeleteCmd returns a new cobra.Command for removing tag protection.
func NewDeleteCmd() *cobra.Command {
	var repo string

	cmd := &cobra.Command{
		Use:   "delete <pattern>",
		Short: "Delete tag protection settings",
		Long:  `Remove the protection settings from a tag pattern. If repo is not specified, the current repository will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]

			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ctx := cmd.Context()
			tagProtection, err := gh.GetTagProtection(ctx, client, repository, pattern)
			if err != nil {
				return fmt.Errorf("failed to get tag protection for pattern %q: %w", pattern, err)
			}
			if tagProtection.ID == nil {
				return fmt.Errorf("failed to delete tag protection for pattern %q: missing tag protection ID", pattern)
			}

			if err := gh.RemoveTagProtection(ctx, client, repository, *tagProtection.ID); err != nil {
				return fmt.Errorf("failed to delete tag protection for pattern %q: %w", pattern, err)
			}

			logger.Info("Tag protection deleted", "pattern", pattern)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")

	return cmd
}
