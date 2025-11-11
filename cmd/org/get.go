package org

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

// NewGetCmd returns a new cobra.Command for getting an organization ruleset
func NewGetCmd() *cobra.Command {
	var opts GetOptions
	var owner string

	cmd := &cobra.Command{
		Use:   "get <ruleset-id>",
		Short: "Get an organization ruleset",
		Long:  `Get detailed information about a specific organization ruleset by its ID. If org is not specified, the current repository's organization will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rulesetID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid ruleset ID: %w", err)
			}

			repository, err := parser.Repository(parser.RepositoryOwner(owner))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ruleset, err := gh.GetOrgRuleset(ctx, client, repository, rulesetID)
			if err != nil {
				return fmt.Errorf("failed to get organization ruleset: %w", err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			renderer.RenderRepositoryRuleset(ruleset, true)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&owner, "owner", "", "Specify the organization name")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
