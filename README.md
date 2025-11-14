# gh-rule-kit

gh extension of github rules api

## Installation

```sh
gh extension install srz-zumix/gh-rule-kit
```

## Commands

### Repository Rulesets

#### List repository rulesets

```sh
gh rule-kit repo list [-R <repo>] [-p]
```

List all rulesets for a repository. If repo is not specified, the current repository will be used.

**Options:**

- `-p, --includes-parent`: Include parent rulesets (default: false)
- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

#### Get a repository ruleset

```sh
gh rule-kit repo get <ruleset-id> [-R <repo>] [-p]
```

Get detailed information about a specific repository ruleset by its ID. If repo is not specified, the current repository will be used.

**Options:**

- `-p, --includes-parent`: Include parent rulesets (default: false)
- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

#### Export a repository ruleset to JSON file

```sh
gh rule-kit repo export <ruleset-id> [-R <repo>] [-o <output>] [-p]
```

Export a specific repository ruleset by its ID to a JSON file. If repo is not specified, the current repository will be used. The exported JSON can be used for backup or to import into another repository.

**Options:**

- `-p, --includes-parent`: Include parent rulesets (default: false)
- `-o, --output <output>`: Output file path (optional, defaults to stdout)
- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

#### Import a repository ruleset from JSON file

```sh
gh rule-kit repo import <input> [-R <repo>] [-c]
```

Import a repository ruleset from a JSON file. If repo is not specified, the current repository will be used. Use --create-if-none flag to create a new ruleset if it does not exist.

**Options:**

- `-c, --create-if-none`: Create a new ruleset if it does not exist (default: false)
- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

#### Migrate repository rulesets to another repository

```sh
gh rule-kit repo migrate <dst-repo> [ruleset-id...] [-R <repo>] [--github-actions-app-id <id>]
```

Migrate repository rulesets from source repository to destination repository. If ruleset IDs are not specified, all rulesets will be migrated. Source repository is specified with --repo flag, destination repository is specified as the first argument.

**Options:**

- `--github-actions-app-id <id>`: The GitHub Actions App ID for integration mapping (optional, default: 0)
- `-R, --repo <repo>`: The source repository in the format 'owner/repo' (optional, defaults to current repository)

#### Delete a repository ruleset

```sh
gh rule-kit repo delete <ruleset-id> [-R <repo>]
```

Delete a specific repository ruleset by its ID. If repo is not specified, the current repository will be used.

**Options:**

- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

### Repository Rule Suite Insights

#### List repository rule suites

```sh
gh rule-kit repo insight list [-R <repo>] [--ref <ref>] [--time-period <period>] [--actor-name <name>] [--result <result>]
```

List all rule suites for a repository. If repo is not specified, the current repository will be used. Rule suites represent evaluations of repository rules.

**Options:**

- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)
- `--ref <ref>`: Filter by ref name (e.g., 'main', 'refs/heads/main') (optional)
- `--time-period <period>`: Filter by time period (e.g., 'hour', 'day', 'week', 'month') (optional)
- `--actor-name <name>`: Filter by actor name (optional)
- `--result <result>`: Filter by rule suite result (e.g., 'pass', 'fail', 'bypass') (optional)

**Note:** This feature requires the Rule Suites API which is not yet fully implemented in go-github v73. The command structure is prepared for future implementation.

#### Get a repository rule suite

```sh
gh rule-kit repo insight get <rule-suite-id> [-R <repo>]
```

Get detailed information about a specific repository rule suite by its ID. If repo is not specified, the current repository will be used.

**Options:**

- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

**Note:** This feature requires the Rule Suites API which is not yet fully implemented in go-github v73. The command structure is prepared for future implementation.

### Organization Rulesets

#### List organization rulesets

```sh
gh rule-kit org list [--owner <owner>]
```

List all rulesets for an organization. If org is not specified, the current repository's organization will be used.

**Options:**

- `--owner <owner>`: The organization name (optional, defaults to current repository's organization)

#### Get an organization ruleset

```sh
gh rule-kit org get <ruleset-id> [--owner <owner>]
```

Get detailed information about a specific organization ruleset by its ID. If org is not specified, the current repository's organization will be used.

**Options:**

- `--owner <owner>`: The organization name (optional, defaults to current repository's organization)

#### Export an organization ruleset to JSON file

```sh
gh rule-kit org export <ruleset-id> [--owner <owner>] [-o <output>]
```

Export a specific organization ruleset by its ID to a JSON file. If org is not specified, the current repository's organization will be used. The exported JSON can be used for backup or to import into another organization.

**Options:**

- `-o, --output <output>`: Output file path (optional, defaults to stdout)
- `--owner <owner>`: The organization name (optional, defaults to current repository's organization)

#### Import an organization ruleset from JSON file

```sh
gh rule-kit org import <input> [--owner <owner>] [-c]
```

Import an organization ruleset from a JSON file. If org is not specified, the current repository's organization will be used. Use --create-if-none flag to create a new ruleset if it does not exist.

**Options:**

- `-c, --create-if-none`: Create a new ruleset if it does not exist (default: false)
- `--owner <owner>`: The organization name (optional, defaults to current repository's organization)

#### Migrate organization rulesets to another organization

```sh
gh rule-kit org migrate <[HOST/]src-org> <[HOST/]dst-org> [ruleset-id...]
```

Migrate organization rulesets from source organization to destination organization. If ruleset IDs are not specified, all rulesets will be migrated.

**Options:**

- `--github-actions-app-id <id>`: The GitHub Actions App ID for integration mapping (optional, default: 0)

#### Delete an organization ruleset

```sh
gh rule-kit org delete <ruleset-id> [--owner <owner>]
```

Delete a specific organization ruleset by its ID. If org is not specified, the current repository's organization will be used.

**Options:**

- `--owner <owner>`: The organization name (optional, defaults to current repository's organization)

### Organization Rule Suite Insights

#### List organization rule suites

```sh
gh rule-kit org insight list [--owner <owner>] [--ref <ref>] [--time-period <period>] [--actor-name <name>] [--result <result>]
```

List all rule suites for an organization. If org is not specified, the current repository's organization will be used. Rule suites represent evaluations of organization rules.

**Options:**

- `--owner <owner>`: The organization name (optional, defaults to current repository's organization)
- `--ref <ref>`: Filter by ref name (e.g., 'main', 'refs/heads/main') (optional)
- `--time-period <period>`: Filter by time period (e.g., 'hour', 'day', 'week', 'month') (optional)
- `--actor-name <name>`: Filter by actor name (optional)
- `--result <result>`: Filter by rule suite result (e.g., 'pass', 'fail', 'bypass') (optional)

**Note:** This feature requires the Rule Suites API which is not yet fully implemented in go-github v73. The command structure is prepared for future implementation.

#### Get an organization rule suite

```sh
gh rule-kit org insight get <rule-suite-id> [--owner <owner>]
```

Get detailed information about a specific organization rule suite by its ID. If org is not specified, the current repository's organization will be used.

**Options:**

- `--owner <owner>`: The organization name (optional, defaults to current repository's organization)

**Note:** This feature requires the Rule Suites API which is not yet fully implemented in go-github v73. The command structure is prepared for future implementation.

