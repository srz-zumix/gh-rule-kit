package insight

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type GetOptions struct {
	Exporter cmdutil.Exporter
}

// NewGetCmd returns a new cobra.Command for getting an organization rule suite
func NewGetCmd() *cobra.Command {
	var opts GetOptions
	var owner string

	cmd := &cobra.Command{
		Use:   "get <rule-suite-id>",
		Short: "Get an organization rule suite",
		Long:  `Get detailed information about a specific organization rule suite by its ID. If org is not specified, the current repository's organization will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleSuiteID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid rule suite ID: %w", err)
			}

			repository, err := parser.Repository(parser.RepositoryOwner(owner))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			ghClient, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ruleSuite, err := gh.GetOrgRuleSuite(ctx, ghClient, repository, ruleSuiteID)
			if err != nil {
				return fmt.Errorf("failed to get organization rule suite: %w", err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			renderer.RenderRuleSuiteDetail(ruleSuite)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&owner, "owner", "", "The organization name")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
