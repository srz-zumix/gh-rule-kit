package cmdcommon

import (
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

// NewListCmd returns a list-rulesets command for the given scope.
func NewListCmd(s Scope) *cobra.Command {
	var exporter cmdutil.Exporter

	cmd := &cobra.Command{
		Use:     "list",
		Short:   fmt.Sprintf("List %s rulesets", s.Noun()),
		Long:    fmt.Sprintf("List all rulesets for %s. %s", s.NounWithArticle(), s.NotSpecifiedHint()),
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
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
			return renderer.RenderRepositoryRulesetsDefault(rulesets)
		},
	}
	s.AddTargetFlag(cmd)
	s.AddIncludesParentFlag(cmd)
	cmdutil.AddFormatFlags(cmd, &exporter)
	return cmd
}
