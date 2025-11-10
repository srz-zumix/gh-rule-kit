package repo

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

// NewListCmd returns a new cobra.Command for listing repository rulesets
func NewListCmd() *cobra.Command {
	var opts ListOptions
	var repo string
	var listIncludesParent bool

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List repository rulesets",
		Long:    `List all rulesets for a repository. If repo is not specified, the current repository will be used.`,
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

			rulesets, err := gh.ListRepositoryRulesets(ctx, client, repository, listIncludesParent)
			if err != nil {
				return fmt.Errorf("failed to list repository rulesets: %w", err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			renderer.RenderRepositoryRulesetsDefault(rulesets)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")
	f.BoolVarP(&listIncludesParent, "includes-parent", "p", false, "Include parent rulesets")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
