package org

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/cmdcommon"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// NewMigrateCmd returns a new cobra.Command for migrating organization rulesets.
func NewMigrateCmd() *cobra.Command {
	mf := &cmdcommon.MigrateFlags{}

	cmd := &cobra.Command{
		Use:   "migrate <[HOST/]src-org> <[HOST/]dst-org> [ruleset-id...]",
		Short: "Migrate organization rulesets to another organization",
		Long:  `Migrate organization rulesets from source organization to destination organization. If ruleset IDs are not specified, all rulesets will be migrated. Source organization is specified as the first argument, destination organization is specified as the second argument. When --usermap is specified, source user logins in User-type bypass actors are mapped to destination logins using the mapping file (as produced by 'user map' in gh-team-kit). Use --dryrun to preview the rulesets that would be migrated without actually creating them.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcRepository, err := parser.Repository(parser.RepositoryOwnerWithHost(args[0]))
			if err != nil {
				return fmt.Errorf("error parsing source organization: %w", err)
			}
			dstRepository, err := parser.Repository(parser.RepositoryOwnerWithHost(args[1]))
			if err != nil {
				return fmt.Errorf("error parsing destination organization: %w", err)
			}

			srcClient, err := gh.NewGitHubClientWithRepo(srcRepository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client for source organization: %w", err)
			}
			dstClient, err := gh.NewGitHubClientWithRepo(dstRepository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client for destination organization: %w", err)
			}

			rulesetIDs, err := cmdcommon.ParseRulesetIDs(args[2:])
			if err != nil {
				return err
			}

			return cmdcommon.RunMigration(cmd.Context(), cmdcommon.NewOrgScope(),
				srcClient, srcRepository, dstClient, dstRepository,
				rulesetIDs, mf.ToRunOptions())
		},
	}

	cmdcommon.AddMigrateFlags(cmd, mf)
	return cmd
}
