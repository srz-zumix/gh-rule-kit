---
name: gh-rule-kit
description: gh-rule-kit GitHub CLI extension for managing GitHub repository and organization rulesets, branch protection rules, tag protection rules, and rule suite insights. Use when listing/getting/exporting/importing/migrating rulesets, converting branch or tag protection rules to rulesets, or inspecting rule suite evaluations.
license: MIT
compatibility:
  - Requires gh CLI (https://cli.github.com) with gh-rule-kit extension installed (`gh extension install srz-zumix/gh-rule-kit`)
---

# gh-rule-kit

`gh-rule-kit` is a GitHub CLI extension for managing GitHub Rules API resources: repository rulesets, organization rulesets, branch protection rules, tag protection rules, and rule suite insights.

## Prerequisites

```bash
# Install gh CLI
brew install gh          # macOS
# or: https://cli.github.com/

# Install gh-rule-kit extension
gh extension install srz-zumix/gh-rule-kit

# Authenticate
gh auth login

# Verify
gh rule-kit --version
```

## Persistent Global Flags

| Flag | Description |
| --- | --- |
| `--read-only` | Prevent any write operations (create, update, delete, import, migrate are blocked) |
| `-L`, `--log-level` | Log level (`debug`, `info`, `warn`, `error`; default `info`) |

Common flags such as `-R`/`--repo`, `--owner`, `-o`/`--output`, `-p`/`--includes-parent`, `-n`/`--dry-run`, `--usermap`, `--github-actions-app-id` are only available on specific subcommands. Check each subcommand's help for supported options.

If `-R`/`--repo` is omitted, the current repository is used. If `--owner` is omitted, the current repository's organization is used.

---

## CLI Structure

```text
gh rule-kit                              # Root command
├── repo                                 # Repository rulesets
│   ├── list                             # List repository rulesets
│   ├── get                              # Get a repository ruleset
│   ├── export                           # Export a ruleset to JSON
│   ├── import                           # Import a ruleset from JSON
│   ├── migrate                          # Migrate rulesets to another repo
│   ├── delete                           # Delete a ruleset
│   ├── from-branch-protection           # Convert branch protection to ruleset
│   ├── from-tag-protection              # Convert tag protection to ruleset
│   ├── branch-protection                # Branch protection rules
│   │   ├── list
│   │   ├── get
│   │   └── delete
│   ├── tag-protection                   # Tag protection rules
│   │   ├── list
│   │   ├── get
│   │   └── delete
│   └── insight                          # Repository rule suite insights
│       ├── list
│       └── get
└── org                                  # Organization rulesets
    ├── list                             # List organization rulesets
    ├── get                              # Get an organization ruleset
    ├── export                           # Export a ruleset to JSON
    ├── import                           # Import a ruleset from JSON
    ├── migrate                          # Migrate rulesets to another org
    ├── delete                           # Delete a ruleset
    └── insight                          # Organization rule suite insights
        ├── list
        └── get
```

---

## Repository Ruleset Commands (`gh rule-kit repo`)

### `list`

```bash
# List rulesets for the current repository
gh rule-kit repo list

# List rulesets for a specific repository
gh rule-kit repo list -R owner/repo

# Include parent (organization) rulesets
gh rule-kit repo list -p
```

### `get`

```bash
# Get a repository ruleset by ID
gh rule-kit repo get <ruleset-id>

# Specify repository
gh rule-kit repo get <ruleset-id> -R owner/repo

# Include parent rulesets
gh rule-kit repo get <ruleset-id> -p
```

### `export`

```bash
# Export ruleset to stdout
gh rule-kit repo export <ruleset-id>

# Export ruleset to file
gh rule-kit repo export <ruleset-id> -o ruleset.json

# Export from a specific repository, including parent rulesets
gh rule-kit repo export <ruleset-id> -R owner/repo -p -o ruleset.json
```

### `import`

```bash
# Import (update) a ruleset from a JSON file
gh rule-kit repo import ruleset.json

# Create the ruleset if no matching one exists
gh rule-kit repo import ruleset.json -c

# Import to a specific repository
gh rule-kit repo import ruleset.json -R owner/repo

# Map User-type bypass actor logins using a usermap (from gh-team-kit)
gh rule-kit repo import ruleset.json --usermap users.csv
```

### `migrate`

```bash
# Migrate all rulesets from current repo to destination repo
gh rule-kit repo migrate dst-owner/dst-repo

# Migrate specific rulesets by ID
gh rule-kit repo migrate dst-owner/dst-repo 12345 67890

# Specify source repository
gh rule-kit repo migrate dst-owner/dst-repo -R src-owner/src-repo

# Map GitHub Actions App ID and User-type bypass actors
gh rule-kit repo migrate dst-owner/dst-repo \
  --github-actions-app-id 1234 \
  --usermap users.csv
```

### `delete`

```bash
# Delete a repository ruleset by ID
gh rule-kit repo delete <ruleset-id>

# Delete from a specific repository
gh rule-kit repo delete <ruleset-id> -R owner/repo
```

### `from-branch-protection`

Converts a branch protection rule to a repository ruleset.

```bash
# Preview the converted ruleset (no write)
gh rule-kit repo from-branch-protection main -n

# Convert and create the ruleset
gh rule-kit repo from-branch-protection main

# Convert, then delete the original branch protection rule
gh rule-kit repo from-branch-protection main --delete

# Convert in a specific repository
gh rule-kit repo from-branch-protection main -R owner/repo
```

Conversion mapping:

| Branch Protection Setting | Ruleset Equivalent |
|---|---|
| Require linear history | `required_linear_history` |
| Allow force pushes (disabled) | `non_fast_forward` |
| Allow deletions (disabled) | `deletion` |
| Block creations | `creation` |
| Require signed commits | `required_signatures` |
| Required status checks | `required_status_checks` |
| Required pull request reviews | `pull_request` |
| Require conversation resolution | `required_review_thread_resolution` (in `pull_request`) |
| Do not include administrators | Repository Admin bypass actor |

Push-access restrictions (`Restrictions`) cannot be directly converted and are reported as warnings.

### `from-tag-protection`

Converts a tag protection rule to a repository ruleset targeting `refs/tags/<pattern>` with `creation`, `update`, and `deletion` rules.

```bash
# Preview the converted ruleset
gh rule-kit repo from-tag-protection 'v*' -n

# Convert and create the ruleset
gh rule-kit repo from-tag-protection 'v*'

# Convert, then delete the original tag protection rule
gh rule-kit repo from-tag-protection 'v*' --delete
```

---

## Branch Protection Commands (`gh rule-kit repo branch-protection`)

### `list`

```bash
# List protected branches in the current repository
gh rule-kit repo branch-protection list

# Specify repository
gh rule-kit repo branch-protection list -R owner/repo
```

### `get`

```bash
# Get protection settings for a branch
gh rule-kit repo branch-protection get main

# Specify repository
gh rule-kit repo branch-protection get main -R owner/repo
```

### `delete`

```bash
# Remove branch protection from a branch
gh rule-kit repo branch-protection delete main

# Specify repository
gh rule-kit repo branch-protection delete main -R owner/repo
```

---

## Tag Protection Commands (`gh rule-kit repo tag-protection`)

### `list`

```bash
# List tag protection settings
gh rule-kit repo tag-protection list

# Specify repository
gh rule-kit repo tag-protection list -R owner/repo
```

### `get`

```bash
# Get tag protection for a pattern
gh rule-kit repo tag-protection get 'v*'
```

### `delete`

```bash
# Remove tag protection by pattern
gh rule-kit repo tag-protection delete 'v*'
```

---

## Repository Rule Suite Insights (`gh rule-kit repo insight`)

Rule suites represent evaluations of repository rules.

### `list`

```bash
# List rule suites for the current repository
gh rule-kit repo insight list

# Filter by ref, time period, actor, result
gh rule-kit repo insight list \
  --ref refs/heads/main \
  --time-period day \
  --actor-name octocat \
  --result fail
```

### `get`

```bash
# Get a rule suite by ID
gh rule-kit repo insight get <rule-suite-id>

# Specify repository
gh rule-kit repo insight get <rule-suite-id> -R owner/repo
```

> Note: rule suite commands depend on the Rule Suites API and may be limited until fully supported by the underlying SDK.

---

## Organization Ruleset Commands (`gh rule-kit org`)

### `list`

```bash
# List rulesets for the current org
gh rule-kit org list

# Specify owner
gh rule-kit org list --owner myorg
```

### `get`

```bash
# Get an organization ruleset by ID
gh rule-kit org get <ruleset-id>

# Specify owner
gh rule-kit org get <ruleset-id> --owner myorg
```

### `export`

```bash
# Export to stdout
gh rule-kit org export <ruleset-id>

# Export to file
gh rule-kit org export <ruleset-id> -o ruleset.json --owner myorg
```

### `import`

```bash
# Import (update) an organization ruleset from JSON
gh rule-kit org import ruleset.json --owner myorg

# Create the ruleset if none exists
gh rule-kit org import ruleset.json -c --owner myorg

# Use usermap for User-type bypass actors
gh rule-kit org import ruleset.json --usermap users.csv --owner myorg
```

### `migrate`

```bash
# Migrate all org rulesets from src-org to dst-org
gh rule-kit org migrate src-org dst-org

# Cross-host migration (HOST/owner)
gh rule-kit org migrate github.com/src-org enterprise.internal/dst-org

# Migrate specific rulesets
gh rule-kit org migrate src-org dst-org 12345 67890

# Map GitHub Actions App ID and User-type bypass actors
gh rule-kit org migrate src-org dst-org \
  --github-actions-app-id 1234 \
  --usermap users.csv
```

### `delete`

```bash
# Delete an organization ruleset by ID
gh rule-kit org delete <ruleset-id> --owner myorg
```

---

## Organization Rule Suite Insights (`gh rule-kit org insight`)

### `list`

```bash
# List rule suites for the current org
gh rule-kit org insight list

# Filter results
gh rule-kit org insight list \
  --owner myorg \
  --ref refs/heads/main \
  --time-period week \
  --actor-name octocat \
  --result bypass
```

### `get`

```bash
# Get an org rule suite by ID
gh rule-kit org insight get <rule-suite-id> --owner myorg
```

> Note: rule suite commands depend on the Rule Suites API and may be limited until fully supported by the underlying SDK.

---

## Common Workflows

### Back up and restore a repository ruleset

```bash
# Export
gh rule-kit repo export 12345 -R owner/repo -o ruleset.json

# Restore (update if exists, create with -c)
gh rule-kit repo import ruleset.json -R owner/repo -c
```

### Migrate all rulesets between repositories

```bash
gh rule-kit repo migrate dst-owner/dst-repo \
  -R src-owner/src-repo \
  --usermap users.csv
```

### Convert legacy branch protection to rulesets

```bash
# Preview
gh rule-kit repo from-branch-protection main -n

# Apply and clean up the legacy rule
gh rule-kit repo from-branch-protection main --delete
```

### Safe inspection in CI

```bash
# Block all writes regardless of subcommand
gh rule-kit --read-only repo import ruleset.json -c
```

---

## Best Practices

1. Use `--read-only` in CI or exploratory sessions to guarantee no writes occur.
2. Preview conversions with `-n`/`--dry-run` (`from-branch-protection`, `from-tag-protection`) before applying.
3. When migrating across orgs/hosts, prepare a usermap (e.g. produced by `gh team-kit user map`) and pass it via `--usermap` so User-type bypass actors are remapped correctly.
4. Provide `--github-actions-app-id` when the destination uses a different GitHub Actions App installation ID.
5. Use `-p`/`--includes-parent` to inspect the effective ruleset stack (repo + org).

## References

- Repository: <https://github.com/srz-zumix/gh-rule-kit>
- GitHub Rules API: <https://docs.github.com/en/rest/repos/rules>
- Repository rulesets: <https://docs.github.com/en/rest/repos/rules#get-all-repository-rulesets>
- Organization rulesets: <https://docs.github.com/en/rest/orgs/rules>
