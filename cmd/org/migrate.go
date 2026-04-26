package org

import (
	"fmt"
	"strconv"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/cmdflags"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
	"github.com/srz-zumix/go-gh-extension/pkg/settings"
)

type MigrateOptions struct {
	Exporter cmdutil.Exporter
}

// NewMigrateCmd returns a new cobra.Command for migrating organization rulesets
func NewMigrateCmd() *cobra.Command {
	var gitHubActionsAppID int64
	var mappings *settings.CompiledMappings
	var dryRun bool
	var opts MigrateOptions

	cmd := &cobra.Command{
		Use:   "migrate <[HOST/]src-org> <[HOST/]dst-org> [ruleset-id...]",
		Short: "Migrate organization rulesets to another organization",
		Long:  `Migrate organization rulesets from source organization to destination organization. If ruleset IDs are not specified, all rulesets will be migrated. Source organization is specified as the first argument, destination organization is specified as the second argument. When --usermap is specified, source user logins in User-type bypass actors are mapped to destination logins using the mapping file (as produced by 'user map' in gh-team-kit). Use --dryrun to preview the rulesets that would be migrated without actually creating them.`,
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
			var successCount int

			// Resolve user mapping if specified
			var resolve func(string) (string, bool)
			if mappings != nil {
				resolve = mappings.ResolveSrc
			}

			renderer := render.NewRenderer(opts.Exporter)

			for _, rulesetID := range rulesetIDs {
				logger.Info("Migrating ruleset", "id", rulesetID)

				// Export ruleset from source (includes team information for actor mapping)
				migrateConfig, err := gh.ExportMigrateRuleset(ctx, srcClient, srcRepository, rulesetID)
				if err != nil {
					logger.Error("Failed to export ruleset", "id", rulesetID, "error", err)
					continue
				}

				// Transform ruleset for the destination (bypass actors, conditions, rules remapping)
				transformedRuleset, err := gh.TransformMigrateRuleset(ctx, dstClient, dstRepository, migrateConfig, gitHubActionsAppIDPtr, resolve)
				if err != nil {
					logger.Error("Failed to transform ruleset", "name", migrateConfig.Ruleset.Name, "error", err)
					continue
				}
				if transformedRuleset == nil {
					// Skipped by ImportMigrateRuleset (e.g. push target on unsupported platform)
					continue
				}

				if dryRun {
					if err := renderer.RenderRepositoryRuleset(transformedRuleset, true); err != nil {
						return fmt.Errorf("failed to render ruleset %d: %w", rulesetID, err)
					}
					successCount++
					continue
				}

				// Create or update ruleset in destination
				createdRuleset, err := gh.CreateOrUpdateRuleset(ctx, dstClient, dstRepository, transformedRuleset)
				if err != nil {
					logger.Error("Failed to create or update ruleset", "name", migrateConfig.Ruleset.Name, "error", err)
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
	f.BoolVarP(&dryRun, "dryrun", "n", false, "Print the rulesets that would be migrated without actually creating them")
	cmdflags.AddUsermapFlag(cmd, &mappings, "User mapping file to map source User-type bypass actor logins to destination logins (as produced by 'user map' in gh-team-kit)")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
