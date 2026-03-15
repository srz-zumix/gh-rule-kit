package branch_protection

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// NewDeleteCmd returns a new cobra.Command for removing branch protection
func NewDeleteCmd() *cobra.Command {
	var repo string

	cmd := &cobra.Command{
		Use:   "delete <branch>",
		Short: "Delete branch protection settings",
		Long:  `Remove the protection settings from a branch. If repo is not specified, the current repository will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := args[0]

			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			if err := gh.RemoveBranchProtection(ctx, client, repository, branch); err != nil {
				return fmt.Errorf("failed to delete branch protection for %q: %w", branch, err)
			}

			logger.Info("Branch protection deleted", "branch", branch)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")

	return cmd
}
