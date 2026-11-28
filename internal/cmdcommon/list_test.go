package cmdcommon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-rule-kit/internal/rulesetdetails"
)

type testListScope struct {
	Scope
	repo repository.Repository
}

func (s testListScope) Parse() (repository.Repository, error) { return s.repo, nil }

func TestListRuleFlags(t *testing.T) {
	for _, scope := range []Scope{NewRepoScope(), NewOrgScope()} {
		t.Run(scope.Noun(), func(t *testing.T) {
			cmd := NewListCmd(scope)
			if err := cmd.ParseFlags([]string{"--rule", "pull_request,update", "--rule", "required_signatures"}); err != nil {
				t.Fatal(err)
			}
			if err := cmd.Args(cmd, nil); err != nil {
				t.Fatalf("valid rules rejected: %v", err)
			}
			cmd = NewListCmd(scope)
			if err := cmd.ParseFlags([]string{"--rule", "invalid"}); err == nil || !strings.Contains(err.Error(), `invalid value "invalid"`) {
				t.Fatalf("unexpected error: %v", err)
			}
			cmd = NewListCmd(scope)
			if err := cmd.ParseFlags([]string{"--rule", "update,update"}); err != nil {
				t.Fatal(err)
			}
			if err := cmd.Args(cmd, nil); err == nil || !strings.Contains(err.Error(), "duplicate rule type") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestListRuleCompletion(t *testing.T) {
	var afterComma []string
	for _, rule := range rulesetdetails.RuleTypes() {
		if rule != "pull_request" {
			afterComma = append(afterComma, "pull_request,"+rule)
		}
	}
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"all", []string{"--rule", ""}, rulesetdetails.RuleTypes()},
		{"prefix", []string{"--rule", "required_s"}, []string{"required_signatures", "required_status_checks"}},
		{"comma separated", []string{"--rule", "pull_request,required_s"}, []string{"pull_request,required_signatures", "pull_request,required_status_checks"}},
		{"trailing comma", []string{"--rule", "pull_request,"}, afterComma},
		{"equals form", []string{"--rule=pull_request,required_s"}, []string{"pull_request,required_signatures", "pull_request,required_status_checks"}},
		{"repeated flag", []string{"--rule", "required_signatures", "--rule", "required_s"}, []string{"required_status_checks"}},
		{"multiple selections", []string{"--rule", "creation", "--rule", "pull_request,required_signatures,required_"}, []string{"pull_request,required_signatures,required_deployments", "pull_request,required_signatures,required_linear_history", "pull_request,required_signatures,required_status_checks"}},
		{"duplicate in current value", []string{"--rule", "pull_request,pull_"}, nil},
		{"unknown prefix", []string{"--rule", "not_a_rule"}, nil},
	}
	scopes := []struct {
		name string
		new  func() Scope
	}{
		{"repo", func() Scope { return NewRepoScope() }},
		{"org", func() Scope { return NewOrgScope() }},
	}
	for _, scope := range scopes {
		t.Run(scope.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					root := &cobra.Command{Use: "gh-rule-kit"}
					group := &cobra.Command{Use: scope.name}
					group.AddCommand(NewListCmd(scope.new()))
					root.AddCommand(group)
					var out, errOut bytes.Buffer
					root.SetOut(&out)
					root.SetErr(&errOut)
					root.SetArgs(append([]string{"__complete", scope.name, "list"}, tt.args...))
					if err := root.Execute(); err != nil {
						t.Fatalf("completion failed: %v: %s", err, errOut.String())
					}
					lines := strings.Split(strings.TrimSpace(out.String()), "\n")
					if lines[len(lines)-1] != ":6" {
						t.Fatalf("unexpected completion directive: %s", out.String())
					}
					got := lines[:len(lines)-1]
					if !slices.Equal(got, tt.want) {
						t.Errorf("completions = %q, want %q", got, tt.want)
					}
				})
			}
		})
	}
}

func TestListCommandRequestsAndJSON(t *testing.T) {
	tests := []struct {
		name        string
		scope       Scope
		repo        repository.Repository
		args        []string
		listPath    string
		detailsPath string
	}{
		{"default repository", NewRepoScope(), repository.Repository{Owner: "acme", Name: "app"}, nil, "/repos/acme/app/rulesets", ""},
		{"repository details with parent", NewRepoScope(), repository.Repository{Owner: "acme", Name: "app"}, []string{"-p", "--rule", "pull_request,required_signatures"}, "/repos/acme/app/rulesets", "/repos/acme/app/rulesets/5"},
		{"organization details", NewOrgScope(), repository.Repository{Owner: "acme"}, []string{"--rule", "pull_request,required_signatures"}, "/orgs/acme/rulesets", "/orgs/acme/rulesets/5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests := []string{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests = append(requests, r.URL.Path)
				switch r.URL.Path {
				case tt.listPath:
					fmt.Fprint(w, `[{"id":5,"name":"main","enforcement":"active","source":"acme"}]`)
				case tt.detailsPath:
					if tt.name == "repository details with parent" && r.URL.Query().Get("includes_parents") != "true" {
						t.Errorf("parent flag missing from request: %s", r.URL)
					}
					fmt.Fprint(w, `{"id":5,"name":"main","enforcement":"active","source":"acme","rules":[{"type":"pull_request","parameters":{"required_approving_review_count":2}}]}`)
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			t.Setenv("GH_HOST", "github.com")
			t.Setenv("GH_TOKEN", "test-token")
			t.Setenv("GITHUB_API_URL", server.URL)

			cmd := NewListCmd(testListScope{Scope: tt.scope, repo: tt.repo})
			cmd.SetArgs(append(tt.args, "--format", "json"))
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			oldStdout := os.Stdout
			os.Stdout = writer
			err = cmd.Execute()
			os.Stdout = oldStdout
			writer.Close()
			output, readErr := io.ReadAll(reader)
			reader.Close()
			if err != nil || readErr != nil {
				t.Fatalf("command error = %v, read error = %v", err, readErr)
			}
			var data []map[string]any
			if err := json.Unmarshal(output, &data); err != nil || len(data) != 1 {
				t.Fatalf("unexpected JSON output %q: %v", output, err)
			}
			if tt.detailsPath == "" {
				if len(requests) != 1 || requests[0] != tt.listPath {
					t.Fatalf("unexpected default requests: %v", requests)
				}
				return
			}
			if len(requests) != 2 || requests[0] != tt.listPath || requests[1] != tt.detailsPath {
				t.Fatalf("unexpected detail requests: %v", requests)
			}
			rules, ok := data[0]["rules"].([]any)
			if !ok || len(rules) != 2 {
				t.Fatalf("unexpected selected rules: %s", output)
			}
			first, ok := rules[0].(map[string]any)
			if !ok || first["type"] != "pull_request" || first["configured"] != true {
				t.Fatalf("unexpected first rule: %s", output)
			}
			params, ok := first["parameters"].(map[string]any)
			if !ok || params["required_approving_review_count"] != float64(2) {
				t.Fatalf("unexpected parameters: %s", output)
			}
			second, ok := rules[1].(map[string]any)
			if !ok || second["type"] != "required_signatures" || second["configured"] != false {
				t.Fatalf("unexpected second rule: %s", output)
			}
		})
	}
}
