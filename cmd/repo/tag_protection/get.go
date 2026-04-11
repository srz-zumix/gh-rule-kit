package tag_protection

import (
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

// NewGetCmd returns a new cobra.Command for getting tag protection settings.
func NewGetCmd() *cobra.Command {
	var opts GetOptions
	var repo string

	cmd := &cobra.Command{
		Use:   "get <pattern>",
		Short: "Get tag protection settings",
		Long:  `Get the protection settings for a specific tag pattern. If repo is not specified, the current repository will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]

			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ctx := cmd.Context()
			tagProtection, err := gh.GetTagProtection(ctx, client, repository, pattern)
			if err != nil {
				return fmt.Errorf("failed to get tag protection for pattern %q: %w", pattern, err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			return renderer.RenderTagProtection(tagProtection)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
