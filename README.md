# RepoSync

> You can configure one or more platforms as source to mirror the structure from the remote server onto your local system (repositories will be moved if necessary).
> You can use complex rules to include or exclude projects based on the namespace or other properties.

## Configuration

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/cidverse/reposync/main/configschema/v1.json
servers:
  - url: https://github.com
    type: github
    auth:
      username: YourAccount
      password: <readOnlyPersonalAccessToken>
      password-file: /path/to/file
      password-command: pass show personal/github.com/read-only-pat
    mirror:
      dir: /tmp/github
      default-action: exclude
      naming-style: slug
      rules:
        - rule: group == "my-org"
          action: include
hooks:
  project-added:
    - zoxide add {{projectDir}}
```

The configuration is read from `~/.config/reposync/config.yaml` by default, but you can also specify a custom path by setting the `REPOSYNC_CONFIG` environment variable.

Supported platforms:

- `github`
- `gitlab`

> The `git` commands will use your local git installation, so you can use ssh keys or other authentication methods.
> The personal access tokens are only used to query the repositories you have access to and not to clone them.

### Hooks

Hooks allow you to run commands when certain events occur. Each hook is defined as a list of commands, which are executed in order.

| Event           | Description                    | Variables                                                              |
|-----------------|--------------------------------|------------------------------------------------------------------------|
| `project-added` | a new project was added/cloned | `projectDir`, `projectName`, `namespace`, `projectId`, `projectUrl`                  |

Example, to add newly cloned projects to your `zoxide` database:

```yaml
hooks:
  project-added:
    - zoxide add {{projectDir}}
```

The following variables are available inside hook commands:

| Variable      | Description                            |
|---------------|----------------------------------------|
| `projectDir`  | The local directory of the project     |
| `projectName` | The name of the project                |
| `namespace`   | The namespace / group of the project   |
| `projectId`   | The unique id of the project           |
| `projectUrl`  | The remote url of the project          |

Hooks are executed with your shell (`sh -c` on unix, `cmd /C` on windows) and can contain environment variables and shell syntax.

## Installation

```bash
curl -L -o ~/.local/bin/reposync https://github.com/cidverse/reposync/releases/download/v0.5.0/linux_amd64
chmod +x ~/.local/bin/reposync
```

## Usage

### Clone

`reposync clone` will clone all repositories you have access to into the target directory, keeping the remote structure as close as possible.

**Notes:**

- tracked repositories might be moved to track the remote structure changes (`<namespace>/<projectName>`).

### Pull

`reposync pull` will pull the latest changes for all tracked repositories.

- `git fetch`
- `git pull --ff-only`

### List

`reposync list` will list all tracked repositories.

Additional options:

- `--format` (`-f`) -> specify the output format (table, json, csv, ...)
- `--columns` (`-c`) -> filter output by columns

### HouseKeeping

`reposync housekeeping` (`reposync hk`) will run the following tasks for all repositories:

- `repack` -> `git repack -a -d --write-bitmap-index`
- `prune` -> `git prune --expire now`
- `gc` -> `git gc --auto`
- `fsck` -> `git fsck --full --unreachable --strict`
- `commit-graph` -> `git commit-graph write --reachable`

You can choose to only run a subset of the tasks by specifying them as arguments. For example, `reposync hk --jobs repack prune`.

### Hooks

`reposync hooks run <event>` runs the hooks for the given event for all tracked projects. This is useful if you added a new hook that should also apply to already known projects (by default hooks only run for newly added ones).

```bash
reposync hooks run project-added
```

- `--dry-run` (`-d`) -> print the commands that would be executed without running them

### Version

`reposync version` prints build information.

- `--short` (`-s`) -> print only the version number

### Rules

Rules support the following variables to match against:

| Variable   | Example           | Description                             |
|------------|-------------------|-----------------------------------------|
| `uniqueId` | github-com/123456 | The unique id of the project            |
| `id`       | 123456            | The id of the project                   |
| `group`    | my-org            | The group / namespace the project is in |
| `name`     | my-project        | The name of the project                 |
| `path`     | my-org/my-project | The full path of the project            |
| `is_fork`  | false             | Whether the project is a fork           |

Rules follow the [Common Expression Language](https://github.com/google/cel-spec) syntax.

## License

Released under the [MIT license](./LICENSE).
