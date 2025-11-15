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

// NewGetCmd returns a new cobra.Command for getting a repository rule suite
func NewGetCmd() *cobra.Command {
	var opts GetOptions
	var repo string

	cmd := &cobra.Command{
		Use:   "get <rule-suite-id>",
		Short: "Get a repository rule suite",
		Long:  `Get detailed information about a specific repository rule suite by its ID. If repo is not specified, the current repository will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleSuiteID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid rule suite ID: %w", err)
			}

			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			ghClient, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ruleSuite, err := gh.GetRepositoryRuleSuite(ctx, ghClient, repository, ruleSuiteID)
			if err != nil {
				return fmt.Errorf("failed to get repository rule suite: %w", err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			renderer.RenderRuleSuiteDetail(ruleSuite)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
