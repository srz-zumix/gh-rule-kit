package org

import (
	"context"
	"fmt"
	"os"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type ImportOptions struct {
	Exporter cmdutil.Exporter
}

// NewImportCmd returns a new cobra.Command for importing an organization ruleset
func NewImportCmd() *cobra.Command {
	var opts ImportOptions
	var owner string
	var input string
	var createIfNotExists bool

	cmd := &cobra.Command{
		Use:   "import <input>",
		Short: "Import an organization ruleset from JSON file",
		Long:  `Import an organization ruleset from a JSON file. If org is not specified, the current repository's organization will be used. Use --create-if-none flag to create a new ruleset if it does not exist.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input = args[0]

			repository, err := parser.Repository(parser.RepositoryOwner(owner))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			var config *gh.RepositoryRulesetConfig
			if input == "-" {
				// Read from stdin
				config, err = gh.LoadRepositoryRulesetConfigFromReader(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
			} else {
				config, err = gh.LoadRepositoryRulesetConfig(input)
				if err != nil {
					return fmt.Errorf("failed to read JSON file: %w", err)
				}
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			found, err := gh.FindOrgRuleset(ctx, client, repository, *config.ID, config.Name)
			if err != nil {
				return fmt.Errorf("failed to find organization ruleset: %w", err)
			}
			if found == nil && !createIfNotExists {
				return fmt.Errorf("ruleset not found with ID %d or name '%s'", *config.ID, config.Name)
			}

			// Convert to RepositoryRuleset
			ruleset := gh.ImportRuleset(config, found)

			resultRuleset := found // nolint
			if found == nil && createIfNotExists {
				// Create new ruleset
				resultRuleset, err = gh.CreateOrgRuleset(ctx, client, repository, ruleset)
				if err != nil {
					return fmt.Errorf("failed to create organization ruleset: %w", err)
				}
				logger.Info("Successfully created ruleset.", "rulesetID", *resultRuleset.ID, "rulesetName", resultRuleset.Name, "organization", repository.Owner)
			} else {
				// Update existing ruleset
				resultRuleset, err = gh.UpdateOrgRuleset(ctx, client, repository, *found.ID, ruleset)
				if err != nil {
					return fmt.Errorf("failed to update organization ruleset: %w", err)
				}
				logger.Info("Successfully updated ruleset.", "rulesetID", *resultRuleset.ID, "rulesetName", resultRuleset.Name, "organization", repository.Owner)
			}

			renderer := render.NewRenderer(opts.Exporter)
			renderer.RenderRepositoryRuleset(resultRuleset, true)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&owner, "owner", "", "Specify the organization name")
	f.BoolVarP(&createIfNotExists, "create-if-none", "c", false, "Create a new ruleset if it does not exist")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
