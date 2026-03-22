package insight

import (
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type ListOptions struct {
	Exporter cmdutil.Exporter
}

// NewListCmd returns a new cobra.Command for listing organization rule suites
func NewListCmd() *cobra.Command {
	var opts ListOptions
	var owner string
	var ref string
	var timePeriod string
	var actorName string
	var result string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List organization rule suites",
		Long:    `List all rule suites for an organization. If org is not specified, the current repository's organization will be used. Rule suites represent evaluations of organization rules.`,
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repository, err := parser.Repository(parser.RepositoryOwner(owner))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ghClient, err := gh.NewGitHubClientWithRepo(repository)
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
			ruleSuites, err := gh.ListOrgRuleSuites(ctx, ghClient, repository, listOpts)
			if err != nil {
				return fmt.Errorf("failed to list organization rule suites: %w", err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			return renderer.RenderRuleSuites(ruleSuites, nil)
		},
	}

	f := cmd.Flags()
	f.StringVar(&owner, "owner", "", "Specify the organization name")
	f.StringVar(&ref, "ref", "", "Filter by ref name (e.g., 'main', 'refs/heads/main')")
	cmdutil.StringEnumFlag(cmd, &timePeriod, "time-period", "", "", gh.RuleSuiteTimePeriodList, "Filter by time period")
	f.StringVar(&actorName, "actor-name", "", "Filter by actor name")
	cmdutil.StringEnumFlag(cmd, &result, "result", "", "", gh.RuleSuiteResultList, "Filter by rule suite result")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
