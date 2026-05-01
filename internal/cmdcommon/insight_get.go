package cmdcommon

import (
	"fmt"
	"strconv"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

// NewInsightGetCmd returns a get rule-suite command for the given scope.
func NewInsightGetCmd(s Scope) *cobra.Command {
	var exporter cmdutil.Exporter

	cmd := &cobra.Command{
		Use:   "get <rule-suite-id>",
		Short: fmt.Sprintf("Get %s rule suite", s.NounWithArticle()),
		Long:  fmt.Sprintf("Get detailed information about a specific %s rule suite by its ID. %s", s.Noun(), s.NotSpecifiedHint()),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleSuiteID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid rule suite ID: %w", err)
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
			ruleSuite, err := gh.GetRuleSuite(ctx, client, repository, ruleSuiteID)
			if err != nil {
				return fmt.Errorf("failed to get %s rule suite: %w", s.Noun(), err)
			}

			renderer := render.NewRenderer(exporter)
			return renderer.RenderRuleSuiteDetail(ruleSuite)
		},
	}
	s.AddTargetFlag(cmd)
	cmdutil.AddFormatFlags(cmd, &exporter)
	return cmd
}
