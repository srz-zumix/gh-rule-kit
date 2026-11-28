package rulesetdetails

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v93/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh/client"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

func TestValidateRules(t *testing.T) {
	tests := []struct {
		name      string
		rules     []string
		specified bool
		wantError bool
	}{
		{"omitted", nil, false, false},
		{"multiple", []string{"pull_request", "required_signatures"}, true, false},
		{"empty", nil, true, true},
		{"unknown", []string{"not_a_rule"}, true, true},
		{"empty element", []string{"pull_request", ""}, true, true},
		{"duplicate", []string{"pull_request", "pull_request"}, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRules(tt.rules, tt.specified)
			if (err != nil) != tt.wantError {
				t.Fatalf("ValidateRules() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestSelectRulesAndRender(t *testing.T) {
	var ruleset github.RepositoryRuleset
	data := `{
		"id": 123, "name": "main", "source": "owner/repo", "enforcement": "disabled",
		"rules": [
			{"type": "pull_request", "parameters": {
				"required_approving_review_count": 2,
				"dismiss_stale_reviews_on_push": false,
				"require_code_owner_review": true,
				"require_last_push_approval": false,
				"required_review_thread_resolution": false,
				"allowed_merge_methods": [],
				"required_reviewers": [{"file_patterns": []}],
				"dismissal_restriction": {"enabled": true, "allowed_actors": []}
			}},
			{"type": "update", "parameters": {"update_allows_fetch_and_merge": false}}
		]
	}`
	if err := json.Unmarshal([]byte(data), &ruleset); err != nil {
		t.Fatal(err)
	}
	result, err := selectRules(&ruleset, []string{"required_signatures", "pull_request", "update"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Enforcement != "disabled" || len(result.Rules) != 3 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.Rules[0].Configured || !result.Rules[1].Configured || !result.Rules[2].Configured {
		t.Fatalf("unexpected configured states: %+v", result.Rules)
	}
	if !strings.Contains(string(result.Rules[1].Parameters), `"required_approving_review_count":2`) ||
		!strings.Contains(string(result.Rules[1].Parameters), `"dismiss_stale_reviews_on_push":false`) ||
		!strings.Contains(string(result.Rules[1].Parameters), `"allowed_merge_methods":[]`) ||
		!strings.Contains(string(result.Rules[1].Parameters), `"file_patterns":[]`) ||
		!strings.Contains(string(result.Rules[1].Parameters), `"allowed_actors":[]`) {
		t.Fatalf("pull request parameters missing: %s", result.Rules[1].Parameters)
	}
	if string(result.Rules[2].Parameters) != `{"update_allows_fetch_and_merge":false}` {
		t.Fatalf("update parameters = %s", result.Rules[2].Parameters)
	}

	table := render.NewStringRenderer(nil)
	if err := Render(&table.Renderer, []Ruleset{result}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"main", "disabled", "required_signatures", "NO", "pull_request", "YES", `"required_approving_review_count":2`} {
		if !strings.Contains(table.Stdout.String(), want) {
			t.Errorf("table output missing %q:\n%s", want, table.Stdout.String())
		}
	}
	output := render.NewStringRenderer(cmdutil.NewJSONExporter())
	if err := Render(&output.Renderer, []Ruleset{result}); err != nil {
		t.Fatal(err)
	}
	var decoded []Ruleset
	if err := json.Unmarshal(output.Stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v: %s", err, output.Stdout.String())
	}
	if len(decoded) != 1 || len(decoded[0].Rules) != 3 || decoded[0].Rules[0].Configured ||
		!decoded[0].Rules[1].Configured || string(decoded[0].Rules[2].Parameters) != `{"update_allows_fetch_and_merge":false}` {
		t.Fatalf("unexpected JSON output: %s", output.Stdout.String())
	}
}

func TestSelectRulesMissingDetails(t *testing.T) {
	for _, ruleset := range []*github.RepositoryRuleset{nil, {}, {ID: github.Ptr(int64(1))}} {
		if _, err := selectRules(ruleset, []string{"creation"}); err == nil {
			t.Fatalf("expected error for missing details: %+v", ruleset)
		}
	}
}

func TestAllSupportedRuleTypesAreSelectable(t *testing.T) {
	apiRules := make([]map[string]string, 0, len(ruleTypes))
	for _, name := range ruleTypes {
		apiRules = append(apiRules, map[string]string{"type": name})
	}
	data, err := json.Marshal(map[string]any{"id": 1, "rules": apiRules})
	if err != nil {
		t.Fatal(err)
	}
	var ruleset github.RepositoryRuleset
	if err := json.Unmarshal(data, &ruleset); err != nil {
		t.Fatal(err)
	}
	selected, err := selectRules(&ruleset, ruleTypes)
	if err != nil {
		t.Fatal(err)
	}
	for _, rule := range selected.Rules {
		if !rule.Configured {
			t.Errorf("configured rule %q was not selected", rule.Type)
		}
	}
}

func TestFetch(t *testing.T) {
	tests := []struct {
		name           string
		repo           repository.Repository
		includesParent bool
		path           string
		query          string
		status         int
		wantError      bool
	}{
		{"repository inherited", repository.Repository{Owner: "acme", Name: "app"}, true, "/repos/acme/app/rulesets/5", "true", http.StatusOK, false},
		{"repository only", repository.Repository{Owner: "acme", Name: "app"}, false, "/repos/acme/app/rulesets/5", "false", http.StatusOK, false},
		{"organization", repository.Repository{Owner: "acme"}, false, "/orgs/acme/rulesets/5", "", http.StatusOK, false},
		{"api failure", repository.Repository{Owner: "acme"}, false, "/orgs/acme/rulesets/5", "", http.StatusForbidden, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count++
				if r.URL.Path != tt.path || r.URL.Query().Get("includes_parents") != tt.query {
					t.Errorf("unexpected request: %s", r.URL)
				}
				w.WriteHeader(tt.status)
				fmt.Fprint(w, `{"id":5,"name":"test","enforcement":"active","rules":[{"type":"creation"}]}`)
			}))
			defer server.Close()
			baseURL := server.URL + "/"
			githubClient, err := github.NewClient(github.WithHTTPClient(server.Client()), github.WithURLs(&baseURL, nil))
			if err != nil {
				t.Fatal(err)
			}
			ghClient, err := client.NewClient(githubClient)
			if err != nil {
				t.Fatal(err)
			}
			got, err := Fetch(context.Background(), ghClient, tt.repo,
				[]*github.RepositoryRuleset{{ID: github.Ptr(int64(5))}}, tt.includesParent, []string{"creation"})
			if (err != nil) != tt.wantError {
				t.Fatalf("Fetch() error = %v, wantError %v", err, tt.wantError)
			}
			if count != 1 {
				t.Fatalf("got %d requests, want 1", count)
			}
			if !tt.wantError && (len(got) != 1 || !got[0].Rules[0].Configured) {
				t.Fatalf("unexpected results: %+v", got)
			}
		})
	}
}

func TestRenderEmpty(t *testing.T) {
	table := render.NewStringRenderer(nil)
	if err := Render(&table.Renderer, []Ruleset{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(table.Stdout.String(), "No repository rulesets.") {
		t.Fatal(table.Stdout.String())
	}
	output := render.NewStringRenderer(cmdutil.NewJSONExporter())
	if err := Render(&output.Renderer, []Ruleset{}); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(output.Stdout.String()) != "[]" {
		t.Fatal(output.Stdout.String())
	}
}
