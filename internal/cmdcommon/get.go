package cmdcommon

import (
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

// NewGetCmd returns a get-ruleset command for the given scope.
func NewGetCmd(s Scope) *cobra.Command {
	var exporter cmdutil.Exporter

	cmd := &cobra.Command{
		Use:   "get <ruleset-id>",
		Short: fmt.Sprintf("Get %s ruleset", s.NounWithArticle()),
		Long:  fmt.Sprintf("Get detailed information about a specific %s ruleset by its ID. %s", s.Noun(), s.NotSpecifiedHint()),
		Args:  cobra.ExactArgs(1),
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

			renderer := render.NewRenderer(exporter)
			return renderer.RenderRepositoryRuleset(ruleset, true)
		},
	}
	s.AddTargetFlag(cmd)
	s.AddIncludesParentFlag(cmd)
	cmdutil.AddFormatFlags(cmd, &exporter)
	return cmd
}
