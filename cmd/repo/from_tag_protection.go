package repo

import (
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

// NewFromTagProtectionCmd returns a new cobra.Command for converting tag protection rules to rulesets.
func NewFromTagProtectionCmd() *cobra.Command {
	var opts struct{ Exporter cmdutil.Exporter }
	var repoFlag string
	var dryRun bool
	var deleteAfter bool

	cmd := &cobra.Command{
		Use:   "from-tag-protection <pattern>",
		Short: "Convert a tag protection rule to a ruleset",
		Long: `Convert a tag protection rule to a repository ruleset.
The converted ruleset is displayed. Use --dry-run to preview without creating.

Use --delete to remove the original tag protection rule after successful conversion.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]

			repository, err := parser.Repository(parser.RepositoryInput(repoFlag))
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

			ruleset := gh.ConvertTagProtectionToRuleset(pattern)

			renderer := render.NewRenderer(opts.Exporter)

			if !dryRun {
				ruleset, err = gh.CreateRepositoryRuleset(ctx, client, repository, ruleset)
				if err != nil {
					return fmt.Errorf("failed to create ruleset for tag pattern %q: %w", pattern, err)
				}

				logger.Info("Successfully created ruleset", "tag_pattern", pattern, "ruleset_id", *ruleset.ID, "name", ruleset.Name)

				if deleteAfter {
					if tagProtection.ID == nil {
						return fmt.Errorf("failed to delete original tag protection for pattern %q: missing tag protection ID", pattern)
					}
					if err := gh.RemoveTagProtection(ctx, client, repository, *tagProtection.ID); err != nil {
						return fmt.Errorf("failed to delete original tag protection for pattern %q: %w", pattern, err)
					}
					logger.Info("Deleted tag protection rule", "tag_pattern", pattern)
				}
			}

			return renderer.RenderRepositoryRuleset(ruleset, true)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repoFlag, "repo", "R", "", "The repository in the format 'owner/repo'")
	f.BoolVarP(&dryRun, "dry-run", "n", false, "Print the ruleset that would be created without actually creating it")
	f.BoolVar(&deleteAfter, "delete", false, "Delete the original tag protection rule after successful conversion")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
