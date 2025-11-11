package org

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

// NewListCmd returns a new cobra.Command for listing organization rulesets
func NewListCmd() *cobra.Command {
	var opts ListOptions
	var owner string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List organization rulesets",
		Long:    `List all rulesets for an organization. If org is not specified, the current repository's organization will be used.`,
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repository, err := parser.Repository(parser.RepositoryOwner(owner))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			rulesets, err := gh.ListOrgRulesets(ctx, client, repository)
			if err != nil {
				return fmt.Errorf("failed to list organization rulesets: %w", err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			renderer.RenderRepositoryRulesetsDefault(rulesets)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&owner, "owner", "", "Specify the organization name")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
