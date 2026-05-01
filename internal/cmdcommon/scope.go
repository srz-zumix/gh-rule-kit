// Package cmdcommon provides shared building blocks for ruleset commands
// that are duplicated between the organization-scoped and repository-scoped
// command trees (cmd/org and cmd/repo).
package cmdcommon

import (
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// Scope abstracts whether a command operates on an organization or a
// repository so that shared command implementations can stay agnostic of
// the target kind.
type Scope interface {
	// Noun returns the bare noun for the scope (e.g. "organization", "repository").
	Noun() string
	// NounWithArticle returns the noun with an indefinite article
	// (e.g. "an organization", "a repository").
	NounWithArticle() string
	// NotSpecifiedHint returns the help text fragment describing the
	// fallback behaviour when the target flag is not specified.
	NotSpecifiedHint() string
	// LabelKey returns the structured-log key used to identify the target
	// (e.g. "organization", "repository").
	LabelKey() string
	// Label returns a printable identifier for the resolved repository.
	Label(r repository.Repository) string

	// AddTargetFlag registers the flag that identifies the target
	// (e.g. --owner for org, -R/--repo for repo).
	AddTargetFlag(cmd *cobra.Command)
	// AddIncludesParentFlag registers --includes-parent for repository
	// scopes; it is a no-op for organization scopes.
	AddIncludesParentFlag(cmd *cobra.Command)
	// IncludesParent reports whether parent rulesets should be included.
	IncludesParent() bool
	// Parse resolves the target repository.Repository from current flag state.
	Parse() (repository.Repository, error)
}

// OrgScope implements Scope for organization-level rulesets.
type OrgScope struct {
	owner string
}

// NewOrgScope returns a fresh OrgScope.
func NewOrgScope() *OrgScope { return &OrgScope{} }

func (s *OrgScope) Noun() string            { return "organization" }
func (s *OrgScope) NounWithArticle() string { return "an organization" }
func (s *OrgScope) NotSpecifiedHint() string {
	return "If --owner is not specified, the current repository's organization will be used."
}
func (s *OrgScope) LabelKey() string                       { return "organization" }
func (s *OrgScope) Label(r repository.Repository) string   { return r.Owner }
func (s *OrgScope) IncludesParent() bool                   { return false }
func (s *OrgScope) AddIncludesParentFlag(_ *cobra.Command) {}
func (s *OrgScope) AddTargetFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&s.owner, "owner", "", "Specify the organization name")
}
func (s *OrgScope) Parse() (repository.Repository, error) {
	return parser.Repository(parser.RepositoryOwner(s.owner))
}

// RepoScope implements Scope for repository-level rulesets.
type RepoScope struct {
	repo           string
	includesParent bool
}

// NewRepoScope returns a fresh RepoScope.
func NewRepoScope() *RepoScope { return &RepoScope{} }

func (s *RepoScope) Noun() string            { return "repository" }
func (s *RepoScope) NounWithArticle() string { return "a repository" }
func (s *RepoScope) NotSpecifiedHint() string {
	return "If repo is not specified, the current repository will be used."
}
func (s *RepoScope) LabelKey() string                     { return "repository" }
func (s *RepoScope) Label(r repository.Repository) string { return parser.GetRepositoryFullName(r) }
func (s *RepoScope) IncludesParent() bool                 { return s.includesParent }
func (s *RepoScope) AddTargetFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&s.repo, "repo", "R", "", "The repository in the format 'owner/repo'")
}
func (s *RepoScope) AddIncludesParentFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&s.includesParent, "includes-parent", "p", false, "Include parent rulesets")
}
func (s *RepoScope) Parse() (repository.Repository, error) {
	return parser.Repository(parser.RepositoryInput(s.repo))
}
