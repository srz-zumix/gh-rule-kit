package cmdcommon

import (
	"context"
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/cmdflags"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
	"github.com/srz-zumix/go-gh-extension/pkg/settings"
)

// MigrateOptions holds runtime options shared by migrate commands.
type MigrateOptions struct {
	GitHubActionsAppID *int64
	DryRun             bool
	Resolve            func(string) (string, bool)
	Renderer           *render.Renderer
}

// MigrateFlags bundles the cobra flags shared by the migrate commands.
type MigrateFlags struct {
	GitHubActionsAppID int64
	DryRun             bool
	Mappings           *settings.CompiledMappings
	Exporter           cmdutil.Exporter
}

// AddMigrateFlags registers the flags shared between org and repo migrate commands.
func AddMigrateFlags(cmd *cobra.Command, mf *MigrateFlags) {
	f := cmd.Flags()
	f.Int64Var(&mf.GitHubActionsAppID, "github-actions-app-id", 0, "The GitHub Actions App ID for integration mapping")
	f.BoolVarP(&mf.DryRun, "dryrun", "n", false, "Print the rulesets that would be migrated without actually creating them")
	cmdflags.AddUsermapFlag(cmd, &mf.Mappings, "User mapping file to map source User-type bypass actor logins to destination logins (as produced by 'user map' in gh-team-kit)")
	cmdutil.AddFormatFlags(cmd, &mf.Exporter)
}

// ToRunOptions converts MigrateFlags to MigrateOptions for RunMigration.
func (mf *MigrateFlags) ToRunOptions() MigrateOptions {
	opts := MigrateOptions{
		DryRun:   mf.DryRun,
		Renderer: render.NewRenderer(mf.Exporter),
	}
	if mf.GitHubActionsAppID != 0 {
		appID := mf.GitHubActionsAppID
		opts.GitHubActionsAppID = &appID
	}
	if mf.Mappings != nil {
		opts.Resolve = mf.Mappings.ResolveSrc
	}
	return opts
}

// RunMigration migrates the given ruleset IDs from src to dst. When
// rulesetIDs is empty, all rulesets from the source are listed and migrated.
func RunMigration(
	ctx context.Context,
	scope Scope,
	srcClient *gh.GitHubClient, srcRepo repository.Repository,
	dstClient *gh.GitHubClient, dstRepo repository.Repository,
	rulesetIDs []int64,
	opts MigrateOptions,
) error {
	if len(rulesetIDs) == 0 {
		rulesets, err := gh.ListRulesets(ctx, srcClient, srcRepo, false)
		if err != nil {
			return fmt.Errorf("failed to list %s rulesets: %w", scope.Noun(), err)
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

	logger.Info("Starting migration", "source", scope.Label(srcRepo), "destination", scope.Label(dstRepo), "count", len(rulesetIDs))

	var successCount int
	for _, rulesetID := range rulesetIDs {
		logger.Info("Migrating ruleset", "id", rulesetID)

		// Export ruleset from source (includes team information for actor mapping)
		migrateConfig, err := gh.ExportMigrateRuleset(ctx, srcClient, srcRepo, rulesetID)
		if err != nil {
			logger.Error("Failed to export ruleset", "id", rulesetID, "error", err)
			continue
		}

		// Transform ruleset for the destination (bypass actors, conditions, rules remapping)
		transformedRuleset, err := gh.TransformMigrateRuleset(ctx, dstClient, dstRepo, migrateConfig, opts.GitHubActionsAppID, opts.Resolve)
		if err != nil {
			logger.Error("Failed to transform ruleset", "name", migrateConfig.Ruleset.Name, "error", err)
			continue
		}
		if transformedRuleset == nil {
			// Skipped by TransformMigrateRuleset (e.g. push target on unsupported platform)
			continue
		}

		if opts.DryRun {
			if err := opts.Renderer.RenderRepositoryRuleset(transformedRuleset, true); err != nil {
				return fmt.Errorf("failed to render ruleset %d: %w", rulesetID, err)
			}
			successCount++
			continue
		}

		// Create or update ruleset in destination
		createdRuleset, err := gh.CreateOrUpdateRuleset(ctx, dstClient, dstRepo, transformedRuleset)
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
}
