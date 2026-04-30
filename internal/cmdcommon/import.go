package cmdcommon

import (
	"fmt"
	"os"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/cmdflags"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
	"github.com/srz-zumix/go-gh-extension/pkg/settings"
)

// NewImportCmd returns an import-ruleset command for the given scope.
func NewImportCmd(s Scope) *cobra.Command {
	var exporter cmdutil.Exporter
	var createIfNotExists bool
	var dryRun bool
	var mappings *settings.CompiledMappings

	cmd := &cobra.Command{
		Use:   "import <input>",
		Short: fmt.Sprintf("Import %s ruleset from JSON file", s.NounWithArticle()),
		Long: fmt.Sprintf(
			"Import %s ruleset from a JSON file. %s Use --create-if-none flag to create a new ruleset if it does not exist. When --usermap is specified, source user logins in User-type bypass actors are automatically converted to destination logins using the mapping file (as produced by 'user map' in gh-team-kit). Use --dryrun to preview the ruleset that would be written without actually creating or updating it.",
			s.NounWithArticle(), s.NotSpecifiedHint(),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]

			repository, err := s.Parse()
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			var config *gh.RepositoryRulesetConfig
			if input == "-" {
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

			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ctx := cmd.Context()

			// Apply user mapping if specified
			if mappings != nil {
				if err := gh.ApplyUserMappingToRulesetConfig(ctx, client, repository, config, mappings.ResolveSrc); err != nil {
					return fmt.Errorf("error applying user mapping: %w", err)
				}
			}

			found, err := gh.FindRuleset(ctx, client, repository, *config.ID, config.Name, s.IncludesParent())
			if err != nil {
				return fmt.Errorf("failed to find %s ruleset: %w", s.Noun(), err)
			}
			if found == nil && !createIfNotExists {
				return fmt.Errorf("ruleset not found with ID %d or name '%s'", *config.ID, config.Name)
			}

			ruleset := config.ToRepositoryRuleset(found)
			if !dryRun {
				if found == nil && createIfNotExists {
					ruleset, err = gh.CreateRuleset(ctx, client, repository, ruleset)
					if err != nil {
						return fmt.Errorf("failed to create %s ruleset: %w", s.Noun(), err)
					}
					logger.Info("Successfully created ruleset.", "rulesetID", *ruleset.ID, "rulesetName", ruleset.Name, s.LabelKey(), s.Label(repository))
				} else {
					ruleset, err = gh.UpdateRuleset(ctx, client, repository, *found.ID, ruleset)
					if err != nil {
						return fmt.Errorf("failed to update %s ruleset: %w", s.Noun(), err)
					}
					logger.Info("Successfully updated ruleset.", "rulesetID", *ruleset.ID, "rulesetName", ruleset.Name, s.LabelKey(), s.Label(repository))
				}
			}

			renderer := render.NewRenderer(exporter)
			return renderer.RenderRepositoryRuleset(ruleset, true)
		},
	}
	s.AddTargetFlag(cmd)
	cmd.Flags().BoolVarP(&createIfNotExists, "create-if-none", "c", false, "Create a new ruleset if it does not exist")
	cmd.Flags().BoolVarP(&dryRun, "dryrun", "n", false, "Print the ruleset that would be written without actually creating or updating it")
	cmdflags.AddUsermapFlag(cmd, &mappings, fmt.Sprintf("User mapping file for User-type bypass actors in the target %s ruleset (as produced by 'user map' in gh-team-kit)", s.Noun()))
	cmdutil.AddFormatFlags(cmd, &exporter)
	return cmd
}
