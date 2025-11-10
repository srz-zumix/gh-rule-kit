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
gh rule-kit repo list [--repo <repo>] [--includes-parent]
```

List all rulesets for a repository. If repo is not specified, the current repository will be used.

**Options:**

- `--includes-parent`: Include parent rulesets (default: false)
- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

#### Get a repository ruleset

```sh
gh rule-kit repo get <ruleset-id> [--repo <repo>] [--includes-parent]
```

Get detailed information about a specific repository ruleset by its ID. If repo is not specified, the current repository will be used.

**Options:**

- `--includes-parent`: Include parent rulesets (default: false)
- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

#### Export a repository ruleset to JSON file

```sh
gh rule-kit repo export <ruleset-id> [--repo <repo>] [--output <output>] [--includes-parent]
```

Export a specific repository ruleset by its ID to a JSON file. If repo is not specified, the current repository will be used. The exported JSON can be used for backup or to import into another repository.

**Options:**

- `--includes-parent`: Include parent rulesets (default: false)
- `-o, --output <output>`: Output file path (optional, defaults to stdout)
- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

#### Import a repository ruleset from JSON file

```sh
gh rule-kit repo import <input> [--repo <repo>] [--create-if-none]
```

Import a repository ruleset from a JSON file. If repo is not specified, the current repository will be used. Use --update flag with --ruleset-id to update an existing ruleset.

**Options:**

- `-c, --create-if-none`: Create a new ruleset if it does not exist (default: false)
- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)

#### Migrate repository rulesets to another repository

```sh
gh rule-kit repo migrate <dst-repo> [ruleset-id...] [--repo <repo>] [--github-actions-app-id <id>]
```

Migrate repository rulesets from source repository to destination repository. If ruleset IDs are not specified, all rulesets will be migrated. Source repository is specified with --repo flag, destination repository is specified as the first argument.

**Options:**

- `--github-actions-app-id <id>`: The GitHub Actions App ID for integration mapping (optional, default: 0)
- `-R, --repo <repo>`: The source repository in the format 'owner/repo' (optional, defaults to current repository)

#### Delete a repository ruleset

```sh
gh rule-kit repo delete <ruleset-id> [--repo <repo>]
```

Delete a specific repository ruleset by its ID. If repo is not specified, the current repository will be used.

**Options:**

- `-R, --repo <repo>`: The repository in the format 'owner/repo' (optional, defaults to current repository)
