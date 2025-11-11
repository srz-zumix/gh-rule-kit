package repo

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

// NewImportCmd returns a new cobra.Command for importing a repository ruleset
func NewImportCmd() *cobra.Command {
	var opts ImportOptions
	var repo string
	var input string
	var createIfNotExists bool

	cmd := &cobra.Command{
		Use:   "import <input>",
		Short: "Import a repository ruleset from JSON file",
		Long:  `Import a repository ruleset from a JSON file. If repo is not specified, the current repository will be used. Use --update flag with --ruleset-id to update an existing ruleset.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input = args[0]

			repository, err := parser.Repository(parser.RepositoryInput(repo))
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
			found, err := gh.FindRepositoryRuleset(ctx, client, repository, *config.ID, config.Name, false)
			if err != nil {
				return fmt.Errorf("failed to find repository ruleset: %w", err)
			}
			if found == nil && !createIfNotExists {
				return fmt.Errorf("ruleset not found with ID %d or name '%s'", *config.ID, config.Name)
			}

			// Convert to RepositoryRuleset
			ruleset := gh.ImportRuleset(config, found)

			resultRuleset := found // nolint
			if found == nil && createIfNotExists {
				// Create new ruleset
				resultRuleset, err = gh.CreateRepositoryRuleset(ctx, client, repository, ruleset)
				if err != nil {
					return fmt.Errorf("failed to create repository ruleset: %w", err)
				}
				logger.Info("Successfully created ruleset.", "rulesetID", *resultRuleset.ID, "rulesetName", resultRuleset.Name, "repository", parser.GetRepositoryFullName(repository))
			} else {
				// Update existing ruleset
				resultRuleset, err = gh.UpdateRepositoryRuleset(ctx, client, repository, *found.ID, ruleset)
				if err != nil {
					return fmt.Errorf("failed to update repository ruleset: %w", err)
				}
				logger.Info("Successfully updated ruleset.", "rulesetID", *resultRuleset.ID, "rulesetName", resultRuleset.Name, "repository", parser.GetRepositoryFullName(repository))
			}

			renderer := render.NewRenderer(opts.Exporter)
			renderer.RenderRepositoryRuleset(resultRuleset, true)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")
	f.BoolVarP(&createIfNotExists, "create-if-none", "c", false, "Create a new ruleset if it does not exist")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
