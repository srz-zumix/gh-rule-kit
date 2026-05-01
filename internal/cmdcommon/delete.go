package cmdcommon

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

// NewDeleteCmd returns a delete-ruleset command for the given scope.
func NewDeleteCmd(s Scope) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <ruleset-id>",
		Short: fmt.Sprintf("Delete %s ruleset", s.NounWithArticle()),
		Long:  fmt.Sprintf("Delete a specific %s ruleset by its ID. %s", s.Noun(), s.NotSpecifiedHint()),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rulesetID, err := parseRulesetID(args[0])
			if err != nil {
				return err
			}

			repository, err := s.Parse()
			if err != nil {
				return fmt.Errorf("error parsing %s: %w", s.Noun(), err)
			}

			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ctx := cmd.Context()
			if err := gh.DeleteRuleset(ctx, client, repository, rulesetID); err != nil {
				return fmt.Errorf("failed to delete %s ruleset: %w", s.Noun(), err)
			}

			logger.Info("Deletion completed successfully.", "rulesetID", rulesetID, s.LabelKey(), s.Label(repository))
			return nil
		},
	}
	s.AddTargetFlag(cmd)
	return cmd
}
