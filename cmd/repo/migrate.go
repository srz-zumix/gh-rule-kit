package repo

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

// NewMigrateCmd returns a new cobra.Command for migrating repository rulesets
func NewMigrateCmd() *cobra.Command {
	var srcRepo string
	var gitHubActionsAppID int64
	var mappings *settings.CompiledMappings
	var dryRun bool
	var opts MigrateOptions

	cmd := &cobra.Command{
		Use:   "migrate <dst-repo> [ruleset-id...]",
		Short: "Migrate repository rulesets to another repository",
		Long:  `Migrate repository rulesets from source repository to destination repository. If ruleset IDs are not specified, all rulesets will be migrated. Source repository is specified with --repo flag, destination repository is specified as the first argument. When --usermap is specified, source user logins in User-type bypass actors are mapped to destination logins using the mapping file (as produced by 'user map' in gh-team-kit). Use --dryrun to preview the rulesets that would be migrated without actually creating them.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse destination repository
			dstRepository, err := parser.Repository(parser.RepositoryInput(args[0]))
			if err != nil {
				return fmt.Errorf("error parsing destination repository: %w", err)
			}

			// Parse source repository
			srcRepository, err := parser.Repository(parser.RepositoryInput(srcRepo))
			if err != nil {
				return fmt.Errorf("error parsing source repository: %w", err)
			}

			// Create clients for source and destination
			srcClient, err := gh.NewGitHubClientWithRepo(srcRepository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client for source repository: %w", err)
			}

			dstClient, err := gh.NewGitHubClientWithRepo(dstRepository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client for destination repository: %w", err)
			}

			ctx := cmd.Context()
			var rulesetIDs []int64
			if len(args) > 1 {
				// Parse specified ruleset IDs
				for _, idStr := range args[1:] {
					id, err := strconv.ParseInt(idStr, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid ruleset ID '%s': %w", idStr, err)
					}
					rulesetIDs = append(rulesetIDs, id)
				}
			} else {
				// Get all rulesets from source repository
				rulesets, err := gh.ListRepositoryRulesets(ctx, srcClient, srcRepository, false)
				if err != nil {
					return fmt.Errorf("failed to list repository rulesets: %w", err)
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

			logger.Info("Starting migration", "source", fmt.Sprintf("%s/%s", srcRepository.Owner, srcRepository.Name), "destination", fmt.Sprintf("%s/%s", dstRepository.Owner, dstRepository.Name), "count", len(rulesetIDs))

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
				transformedRuleset, err := gh.ImportMigrateRuleset(ctx, dstClient, dstRepository, migrateConfig, gitHubActionsAppIDPtr, resolve)
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
	f.StringVarP(&srcRepo, "repo", "R", "", "The source repository in the format 'owner/repo'")
	f.Int64Var(&gitHubActionsAppID, "github-actions-app-id", 0, "The GitHub Actions App ID for integration mapping")
	f.BoolVarP(&dryRun, "dryrun", "n", false, "Print the rulesets that would be migrated without actually creating them")
	cmdflags.AddUsermapFlag(cmd, &mappings, "User mapping file to map source User-type bypass actor logins to destination logins (as produced by 'user map' in gh-team-kit)")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
