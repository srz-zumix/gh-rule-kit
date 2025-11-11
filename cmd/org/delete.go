package org

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// NewDeleteCmd returns a new cobra.Command for deleting an organization ruleset
func NewDeleteCmd() *cobra.Command {
	var owner string

	cmd := &cobra.Command{
		Use:   "delete <ruleset-id>",
		Short: "Delete an organization ruleset",
		Long:  `Delete a specific organization ruleset by its ID. If org is not specified, the current repository's organization will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rulesetID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid ruleset ID: %w", err)
			}

			repository, err := parser.Repository(parser.RepositoryOwner(owner))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			err = gh.DeleteOrgRuleset(ctx, client, repository, rulesetID)
			if err != nil {
				return fmt.Errorf("failed to delete organization ruleset: %w", err)
			}

			logger.Info("Deletion completed successfully.", "rulesetID", rulesetID, "organization", repository.Owner)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&owner, "owner", "", "Specify the organization name")

	return cmd
}
