package cmdcommon

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

// NewExportCmd returns an export-ruleset command for the given scope.
func NewExportCmd(s Scope) *cobra.Command {
	var exporter cmdutil.Exporter
	var output string

	cmd := &cobra.Command{
		Use:   "export <ruleset-id>",
		Short: fmt.Sprintf("Export %s ruleset to JSON file", s.NounWithArticle()),
		Long: fmt.Sprintf(
			"Export a specific %s ruleset by its ID to a JSON file. %s The exported JSON can be used for backup or to import into another %s.",
			s.Noun(), s.NotSpecifiedHint(), s.Noun(),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rulesetID, err := parseRulesetID(args[0])
			if err != nil {
				return err
			}

			repository, err := s.Parse()
			if err != nil {
				return fmt.Errorf("error parsing %s: %w", s.Noun(), err)
			}

			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ctx := cmd.Context()
			ruleset, err := gh.GetRuleset(ctx, client, repository, rulesetID, s.IncludesParent())
			if err != nil {
				return fmt.Errorf("failed to get %s ruleset: %w", s.Noun(), err)
			}

			config := gh.ExportRuleset(ruleset)
			config.BypassActorsMeta = gh.BuildBypassActorsMeta(ctx, client, repository, ruleset)

			renderer := render.NewRenderer(exporter)
			if exporter != nil {
				return renderer.RenderExportedData(config)
			}

			jsonData, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal ruleset to JSON: %w", err)
			}
			if output == "" || output == "-" {
				fmt.Println(string(jsonData))
				return nil
			}
			if err := os.WriteFile(output, jsonData, 0644); err != nil {
				return fmt.Errorf("failed to write JSON to file: %w", err)
			}
			logger.Info("Export completed successfully.", "output", output)
			return nil
		},
	}
	s.AddTargetFlag(cmd)
	s.AddIncludesParentFlag(cmd)
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")
	cmdutil.AddFormatFlags(cmd, &exporter)
	cmd.MarkFlagsMutuallyExclusive("output", "format")
	return cmd
}
