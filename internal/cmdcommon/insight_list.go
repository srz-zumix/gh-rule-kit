package cmdcommon

import (
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

// NewInsightListCmd returns a list rule-suites command for the given scope.
func NewInsightListCmd(s Scope) *cobra.Command {
	var exporter cmdutil.Exporter
	var ref string
	var timePeriod string
	var actorName string
	var result string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   fmt.Sprintf("List %s rule suites", s.Noun()),
		Long:    fmt.Sprintf("List all rule suites for %s. %s Rule suites represent evaluations of %s rules.", s.NounWithArticle(), s.NotSpecifiedHint(), s.Noun()),
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

			listOpts := &gh.ListRuleSuitesOptions{
				Ref:             ref,
				TimePeriod:      timePeriod,
				ActorName:       actorName,
				RuleSuiteResult: result,
			}

			ctx := cmd.Context()
			ruleSuites, err := gh.ListRuleSuites(ctx, client, repository, listOpts)
			if err != nil {
				return fmt.Errorf("failed to list %s rule suites: %w", s.Noun(), err)
			}

			renderer := render.NewRenderer(exporter)
			return renderer.RenderRuleSuites(ruleSuites, nil)
		},
	}
	s.AddTargetFlag(cmd)
	cmd.Flags().StringVar(&ref, "ref", "", "Filter by ref name (e.g., 'main', 'refs/heads/main')")
	cmdutil.StringEnumFlag(cmd, &timePeriod, "time-period", "", "", gh.RuleSuiteTimePeriodList, "Filter by time period")
	cmd.Flags().StringVar(&actorName, "actor-name", "", "Filter by actor name")
	cmdutil.StringEnumFlag(cmd, &result, "result", "", "", gh.RuleSuiteResultList, "Filter by rule suite result")
	cmdutil.AddFormatFlags(cmd, &exporter)
	return cmd
}
