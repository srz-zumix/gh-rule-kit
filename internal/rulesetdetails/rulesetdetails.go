package rulesetdetails

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v93/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

var ruleTypes = []string{
	"branch_name_pattern", "code_coverage", "code_scanning", "commit_author_email_pattern",
	"commit_message_pattern", "committer_email_pattern", "copilot_code_review", "creation",
	"deletion", "file_extension_restriction", "file_path_restriction", "max_file_path_length",
	"max_file_size", "merge_queue", "non_fast_forward", "pull_request", "repository_create",
	"repository_delete", "repository_name", "repository_transfer", "repository_visibility",
	"required_deployments", "required_linear_history", "required_signatures",
	"required_status_checks", "tag_name_pattern", "update", "workflows",
}

// RuleTypes returns the supported rule types.
func RuleTypes() []string {
	return slices.Clone(ruleTypes)
}

type Rule struct {
	Type       string          `json:"type"`
	Configured bool            `json:"configured"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

type Ruleset struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Enforcement string `json:"enforcement"`
	Source      string `json:"source"`
	Rules       []Rule `json:"rules"`
}

func ValidateRules(rules []string, specified bool) error {
	if specified && len(rules) == 0 {
		return fmt.Errorf("--rule requires at least one rule type")
	}
	seen := make(map[string]bool, len(rules))
	for _, rule := range rules {
		found := false
		for _, candidate := range ruleTypes {
			if rule == candidate {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown rule type %q; available types: %s", rule, strings.Join(ruleTypes, ", "))
		}
		if seen[rule] {
			return fmt.Errorf("duplicate rule type %q", rule)
		}
		seen[rule] = true
	}
	return nil
}

func Fetch(ctx context.Context, client *gh.GitHubClient, repo repository.Repository, summaries []*github.RepositoryRuleset, includesParent bool, selected []string) ([]Ruleset, error) {
	results := make([]Ruleset, 0, len(summaries))
	for _, summary := range summaries {
		if summary == nil || summary.ID == nil {
			return nil, fmt.Errorf("listed ruleset has no ID")
		}
		detail, err := gh.GetRuleset(ctx, client, repo, *summary.ID, includesParent)
		if err != nil {
			return nil, fmt.Errorf("get ruleset %d: %w", *summary.ID, err)
		}
		result, err := selectRules(detail, selected)
		if err != nil {
			return nil, fmt.Errorf("select rules for ruleset %d: %w", *summary.ID, err)
		}
		results = append(results, result)
	}
	return results, nil
}

func selectRules(ruleset *github.RepositoryRuleset, selected []string) (Ruleset, error) {
	if ruleset == nil || ruleset.ID == nil || ruleset.Rules == nil {
		return Ruleset{}, fmt.Errorf("ruleset details are missing ID or rules")
	}
	data, err := json.Marshal(ruleset.Rules)
	if err != nil {
		return Ruleset{}, fmt.Errorf("encode rules: %w", err)
	}
	var configured []struct {
		Type       string          `json:"type"`
		Parameters json.RawMessage `json:"parameters"`
	}
	if err := json.Unmarshal(data, &configured); err != nil {
		return Ruleset{}, fmt.Errorf("decode rules: %w", err)
	}
	byType := make(map[string]json.RawMessage, len(configured))
	for _, rule := range configured {
		byType[rule.Type] = rule.Parameters
	}
	// go-github omits update's false parameter when serializing rule settings.
	if ruleset.Rules.Update != nil {
		params := struct {
			UpdateAllowsFetchAndMerge bool `json:"update_allows_fetch_and_merge"`
		}{ruleset.Rules.Update.UpdateAllowsFetchAndMerge}
		byType["update"], err = json.Marshal(params)
		if err != nil {
			return Ruleset{}, fmt.Errorf("encode update rule: %w", err)
		}
	}
	if pr := ruleset.Rules.PullRequest; pr != nil {
		byType["pull_request"], err = marshalPullRequest(pr)
		if err != nil {
			return Ruleset{}, fmt.Errorf("encode pull request rule: %w", err)
		}
	}
	result := Ruleset{
		ID:          *ruleset.ID,
		Name:        ruleset.Name,
		Enforcement: string(ruleset.Enforcement),
		Source:      ruleset.Source,
		Rules:       make([]Rule, 0, len(selected)),
	}
	for _, name := range selected {
		params, ok := byType[name]
		result.Rules = append(result.Rules, Rule{Type: name, Configured: ok, Parameters: params})
	}
	return result, nil
}

func marshalPullRequest(pr *github.PullRequestRuleParameters) ([]byte, error) {
	type reviewer struct {
		*github.RulesetRequiredReviewer
		FilePatterns *[]string `json:"file_patterns,omitempty"`
	}
	type dismissal struct {
		*github.DismissalRestriction
		AllowedActors *[]*github.DismissalRestrictionActor `json:"allowed_actors,omitempty"`
	}
	params := struct {
		*github.PullRequestRuleParameters
		AllowedMergeMethods  *[]github.PullRequestMergeMethod `json:"allowed_merge_methods,omitempty"`
		RequiredReviewers    *[]*reviewer                     `json:"required_reviewers,omitempty"`
		DismissalRestriction *dismissal                       `json:"dismissal_restriction,omitempty"`
	}{PullRequestRuleParameters: pr}
	if pr.AllowedMergeMethods != nil {
		params.AllowedMergeMethods = &pr.AllowedMergeMethods
	}
	if pr.RequiredReviewers != nil {
		reviewers := make([]*reviewer, 0, len(pr.RequiredReviewers))
		for _, r := range pr.RequiredReviewers {
			if r == nil {
				reviewers = append(reviewers, nil)
				continue
			}
			entry := &reviewer{RulesetRequiredReviewer: r}
			if r.FilePatterns != nil {
				entry.FilePatterns = &r.FilePatterns
			}
			reviewers = append(reviewers, entry)
		}
		params.RequiredReviewers = &reviewers
	}
	if pr.DismissalRestriction != nil {
		params.DismissalRestriction = &dismissal{DismissalRestriction: pr.DismissalRestriction}
		if pr.DismissalRestriction.AllowedActors != nil {
			params.DismissalRestriction.AllowedActors = &pr.DismissalRestriction.AllowedActors
		}
	}
	return json.Marshal(params)
}

func Render(renderer *render.Renderer, rulesets []Ruleset) error {
	if renderer.HasExporter() {
		return renderer.RenderExportedData(rulesets)
	}
	if len(rulesets) == 0 {
		_, err := fmt.Fprintln(renderer.IO.Out, "No repository rulesets.")
		return err
	}
	for _, ruleset := range rulesets {
		if _, err := fmt.Fprintf(renderer.IO.Out, "Ruleset %d: %s (source: %s, enforcement: %s)\n", ruleset.ID, ruleset.Name, ruleset.Source, ruleset.Enforcement); err != nil {
			return err
		}
		table := renderer.NewTableWriter([]string{"RULE", "CONFIGURED", "PARAMETERS"})
		for _, rule := range ruleset.Rules {
			table.Append([]string{rule.Type, render.ToString(rule.Configured), string(rule.Parameters)})
		}
		if err := table.Render(); err != nil {
			return err
		}
	}
	return nil
}
