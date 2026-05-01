package repo

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// NewMigrateCmd returns a new cobra.Command for migrating repository rulesets.
func NewMigrateCmd() *cobra.Command {
	var srcRepo string
	mf := &cmdcommon.MigrateFlags{}

	cmd := &cobra.Command{
		Use:   "migrate <dst-repo> [ruleset-id...]",
		Short: "Migrate repository rulesets to another repository",
		Long:  `Migrate repository rulesets from source repository to destination repository. If ruleset IDs are not specified, all rulesets will be migrated. Source repository is specified with --repo flag, destination repository is specified as the first argument. When --usermap is specified, source user logins in User-type bypass actors are mapped to destination logins using the mapping file (as produced by 'user map' in gh-team-kit). Use --dryrun to preview the rulesets that would be migrated without actually creating them.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dstRepository, err := parser.Repository(parser.RepositoryInput(args[0]))
			if err != nil {
				return fmt.Errorf("error parsing destination repository: %w", err)
			}
			srcRepository, err := parser.Repository(parser.RepositoryInput(srcRepo))
			if err != nil {
				return fmt.Errorf("error parsing source repository: %w", err)
			}

			srcClient, err := gh.NewGitHubClientWithRepo(srcRepository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client for source repository: %w", err)
			}
			dstClient, err := gh.NewGitHubClientWithRepo(dstRepository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client for destination repository: %w", err)
			}

			rulesetIDs, err := cmdcommon.ParseRulesetIDs(args[1:])
			if err != nil {
				return err
			}

			return cmdcommon.RunMigration(cmd.Context(), cmdcommon.NewRepoScope(),
				srcClient, srcRepository, dstClient, dstRepository,
				rulesetIDs, mf.ToRunOptions())
		},
	}

	cmd.Flags().StringVarP(&srcRepo, "repo", "R", "", "The source repository in the format 'owner/repo'")
	cmdcommon.AddMigrateFlags(cmd, mf)
	return cmd
}
