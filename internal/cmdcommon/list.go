package cmdcommon

import (
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/rulesetdetails"
	"github.com/srz-zumix/go-gh-extension/pkg/cmdflags"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

// NewListCmd returns a list-rulesets command for the given scope.
func NewListCmd(s Scope) *cobra.Command {
	var exporter cmdutil.Exporter
	var rules []string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   fmt.Sprintf("List %s rulesets", s.Noun()),
		Long:    fmt.Sprintf("List all rulesets for %s. %s Use --rule to show selected rule settings for each ruleset.", s.NounWithArticle(), s.NotSpecifiedHint()),
		Aliases: []string{"ls"},
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.NoArgs(cmd, args); err != nil {
				return err
			}
			return rulesetdetails.ValidateRules(rules, cmd.Flags().Changed("rule"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repository, err := s.Parse()
			if err != nil {
				return fmt.Errorf("error parsing %s: %w", s.Noun(), err)
			}

			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ctx := cmd.Context()
			rulesets, err := gh.ListRulesets(ctx, client, repository, s.IncludesParent())
			if err != nil {
				return fmt.Errorf("failed to list %s rulesets: %w", s.Noun(), err)
			}

			renderer := render.NewRenderer(exporter)
			if len(rules) > 0 {
				details, err := rulesetdetails.Fetch(ctx, client, repository, rulesets, s.IncludesParent(), rules)
				if err != nil {
					return fmt.Errorf("failed to retrieve %s ruleset details: %w", s.Noun(), err)
				}
				return rulesetdetails.Render(renderer, details)
			}
			return renderer.RenderRepositoryRulesetsDefault(rulesets)
		},
	}
	s.AddTargetFlag(cmd)
	s.AddIncludesParentFlag(cmd)
	cmdflags.StringSliceEnumVar(cmd, &rules, "rule", nil, rulesetdetails.RuleTypes(), "Show details for rule types (comma-separated or repeated)")
	cmdutil.AddFormatFlags(cmd, &exporter)
	return cmd
}
