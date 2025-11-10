# gh-rule-kit

gh extension of github rules api

## Installation

```sh
gh extension install srz-zumix/gh-rule-kit
```

## Commands

### Repository Rulesets

#### List Repository Rulesets

```sh
gh rule-kit repo list [--repo <repo>] [--includes-parent]
```

List all rulesets for a repository. Displays ruleset information.

**Options:**

- `-R, --repo <repo>`: Repository name (optional, defaults to current repository)
- `-p, --includes-parent`: Include parent rulesets from organization level (default: false)
