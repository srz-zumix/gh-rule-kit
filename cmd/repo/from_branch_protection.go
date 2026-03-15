package repo

import (
	"context"
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

// NewFromBranchProtectionCmd returns a new cobra.Command for migrating branch protection rules to rulesets
func NewFromBranchProtectionCmd() *cobra.Command {
	var opts struct{ Exporter cmdutil.Exporter }
	var repoFlag string
	var dryRun bool
	var deleteAfter bool

	cmd := &cobra.Command{
		Use:   "from-branch-protection <branch>",
		Short: "Convert a branch protection rule to a ruleset",
		Long: `Convert a branch protection rule to a repository ruleset.
The converted ruleset is displayed. Use --dry-run to preview without creating.

Rules that have no direct ruleset equivalent (e.g., push-access restrictions) are reported as warnings.
Use --delete to remove the original branch protection rule after successful conversion.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := args[0]

			repository, err := parser.Repository(parser.RepositoryInput(repoFlag))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			protection, err := gh.GetBranchProtection(ctx, client, repository, branch)
			if err != nil {
				return fmt.Errorf("failed to get branch protection for %q: %w", branch, err)
			}

			ruleset := gh.ConvertBranchProtectionToRuleset(branch, protection)

			renderer := render.NewRenderer(opts.Exporter)

			if !dryRun {
				ruleset, err = gh.CreateRepositoryRuleset(ctx, client, repository, ruleset)
				if err != nil {
					return fmt.Errorf("failed to create ruleset for branch %q: %w", branch, err)
				}

				logger.Info("Successfully created ruleset", "branch", branch, "ruleset_id", *ruleset.ID, "name", ruleset.Name)

				if deleteAfter {
					if err := gh.RemoveBranchProtection(ctx, client, repository, branch); err != nil {
						return fmt.Errorf("failed to delete original branch protection for %q: %w", branch, err)
					} else {
						logger.Info("Deleted branch protection rule", "branch", branch)
					}
				}
			}
			return renderer.RenderRepositoryRuleset(ruleset, true)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repoFlag, "repo", "R", "", "The repository in the format 'owner/repo'")
	f.BoolVarP(&dryRun, "dry-run", "n", false, "Print the ruleset that would be created without actually creating it")
	f.BoolVar(&deleteAfter, "delete", false, "Delete the original branch protection rule after successful conversion")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
