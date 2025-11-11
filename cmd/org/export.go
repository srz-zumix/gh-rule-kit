package org

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// NewExportCmd returns a new cobra.Command for exporting an organization ruleset
func NewExportCmd() *cobra.Command {
	var owner string
	var output string

	cmd := &cobra.Command{
		Use:   "export <ruleset-id>",
		Short: "Export an organization ruleset to JSON file",
		Long:  `Export a specific organization ruleset by its ID to a JSON file. If org is not specified, the current repository's organization will be used. The exported JSON can be used for backup or to import into another organization.`,
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

			config := gh.ExportRuleset(ruleset)

			var jsonData []byte
			jsonData, err = json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal ruleset to JSON: %w", err)
			}

			if output == "" || output == "-" {
				// Output to stdout
				fmt.Println(string(jsonData))
			} else {
				// Output to file
				err = os.WriteFile(output, jsonData, 0644)
				if err != nil {
					return fmt.Errorf("failed to write JSON to file: %w", err)
				}
				logger.Info("Export completed successfully.", "output", output)
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&owner, "owner", "", "Specify the organization name")
	f.StringVarP(&output, "output", "o", "", "Output file path (default: stdout)")

	return cmd
}
