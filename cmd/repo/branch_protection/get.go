package branch_protection

import (
	"context"
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type GetOptions struct {
	Exporter cmdutil.Exporter
}

// NewGetCmd returns a new cobra.Command for getting branch protection settings
func NewGetCmd() *cobra.Command {
	var opts GetOptions
	var repo string

	cmd := &cobra.Command{
		Use:   "get <branch>",
		Short: "Get branch protection settings",
		Long:  `Get the protection settings for a specific branch. If repo is not specified, the current repository will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := args[0]

			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			protection, err := gh.GetBranchProtection(ctx, client, repository, branch)
			if err != nil {
				return fmt.Errorf("failed to get branch protection for %q: %w", branch, err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			return renderer.RenderBranchProtection(branch, protection)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
