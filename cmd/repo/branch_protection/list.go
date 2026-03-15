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

type ListOptions struct {
	Exporter cmdutil.Exporter
}

// NewListCmd returns a new cobra.Command for listing protected branches
func NewListCmd() *cobra.Command {
	var opts ListOptions
	var repo string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List protected branches",
		Long:    `List all protected branches for a repository. If repo is not specified, the current repository will be used.`,
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			branches, err := gh.ListProtectedBranches(ctx, client, repository)
			if err != nil {
				return fmt.Errorf("failed to list protected branches: %w", err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			return renderer.RenderBranchProtections(branches)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
