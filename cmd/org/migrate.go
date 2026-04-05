package org

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/settings"
)

// NewMigrateCmd returns a new cobra.Command for migrating organization rulesets
func NewMigrateCmd() *cobra.Command {
	var gitHubActionsAppID int64
	var mapFile string

	cmd := &cobra.Command{
		Use:   "migrate <[HOST/]src-org> <[HOST/]dst-org> [ruleset-id...]",
		Short: "Migrate organization rulesets to another organization",
		Long:  `Migrate organization rulesets from source organization to destination organization. If ruleset IDs are not specified, all rulesets will be migrated. Source organization is specified as the first argument, destination organization is specified as the second argument. When --usermap is specified, source user logins in User-type bypass actors are mapped to destination logins using the mapping file (as produced by 'user map' in gh-team-kit).`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcOrg := args[0]
			dstOrg := args[1]

			// Parse source organization
			srcRepository, err := parser.Repository(parser.RepositoryOwnerWithHost(srcOrg))
			if err != nil {
				return fmt.Errorf("error parsing source organization: %w", err)
			}

			// Parse destination organization
			dstRepository, err := parser.Repository(parser.RepositoryOwnerWithHost(dstOrg))
			if err != nil {
				return fmt.Errorf("error parsing destination organization: %w", err)
			}

			// Create clients for source and destination
			srcClient, err := gh.NewGitHubClientWithRepo(srcRepository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client for source organization: %w", err)
			}

			dstClient, err := gh.NewGitHubClientWithRepo(dstRepository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client for destination organization: %w", err)
			}

			ctx := cmd.Context()
			var rulesetIDs []int64
			if len(args) > 2 {
				// Parse specified ruleset IDs
				for _, idStr := range args[2:] {
					id, err := strconv.ParseInt(idStr, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid ruleset ID '%s': %w", idStr, err)
					}
					rulesetIDs = append(rulesetIDs, id)
				}
			} else {
				// Get all rulesets from source organization
				rulesets, err := gh.ListOrgRulesets(ctx, srcClient, srcRepository)
				if err != nil {
					return fmt.Errorf("failed to list organization rulesets: %w", err)
				}
				for _, ruleset := range rulesets {
					if ruleset.ID != nil {
						rulesetIDs = append(rulesetIDs, *ruleset.ID)
					}
				}
			}

			if len(rulesetIDs) == 0 {
				logger.Info("No rulesets to migrate")
				return nil
			}

			logger.Info("Starting migration", "source", srcRepository.Owner, "destination", dstRepository.Owner, "count", len(rulesetIDs))

			var gitHubActionsAppIDPtr *int64
			if gitHubActionsAppID != 0 {
				gitHubActionsAppIDPtr = &gitHubActionsAppID
			}

			// Load mapping file if specified
			var resolve func(string) (string, bool)
			if mapFile != "" {
				compiledMappings, err := settings.NewCompiledMappingsFromFile(mapFile)
				if err != nil {
					return fmt.Errorf("error loading mapping file '%s': %w", mapFile, err)
				}
				resolve = compiledMappings.ResolveSrc
			}

			// Migrate each ruleset
			successCount := 0
			for _, rulesetID := range rulesetIDs {
				logger.Info("Migrating ruleset", "id", rulesetID)

				// Export ruleset from source (includes team information for actor mapping)
				migrateConfig, err := gh.ExportMigrateRuleset(ctx, srcClient, srcRepository, rulesetID)
				if err != nil {
					logger.Error("Failed to export ruleset", "id", rulesetID, "error", err)
					continue
				}

				// Import ruleset to destination (handles team actor ID mapping)
				createdRuleset, err := gh.ImportMigrateRuleset(ctx, dstClient, dstRepository, migrateConfig, gitHubActionsAppIDPtr, resolve)
				if err != nil {
					logger.Error("Failed to import ruleset", "name", migrateConfig.Ruleset.Name, "error", err)
					continue
				}

				logger.Info("Successfully migrated ruleset", "src_id", rulesetID, "dst_id", *createdRuleset.ID, "name", createdRuleset.Name)
				successCount++
			}

			logger.Info("Migration completed", "total", len(rulesetIDs), "success", successCount, "failed", len(rulesetIDs)-successCount)

			if successCount == 0 {
				return fmt.Errorf("failed to migrate any rulesets")
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.Int64Var(&gitHubActionsAppID, "github-actions-app-id", 0, "The GitHub Actions App ID for integration mapping")
	f.StringVar(&mapFile, "usermap", "", "User mapping file for User-type bypass actor login conversion (as produced by 'user map' in gh-team-kit)")

	return cmd
}
