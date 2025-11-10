package repo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// NewDeleteCmd returns a new cobra.Command for deleting a repository ruleset
func NewDeleteCmd() *cobra.Command {
	var repo string

	cmd := &cobra.Command{
		Use:   "delete <ruleset-id>",
		Short: "Delete a repository ruleset",
		Long:  `Delete a specific repository ruleset by its ID. If repo is not specified, the current repository will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rulesetID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid ruleset ID: %w", err)
			}

			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			err = gh.DeleteRepositoryRuleset(ctx, client, repository, rulesetID)
			if err != nil {
				return fmt.Errorf("failed to delete repository ruleset: %w", err)
			}

			logger.Info("Deletion completed successfully.", "rulesetID", rulesetID, "repository", parser.GetRepositoryFullName(repository))
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")

	return cmd
}
