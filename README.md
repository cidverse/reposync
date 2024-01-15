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
      # you can provide your personal access token directly or in a file
      username: YourAccount
      password: <readOnlyPersonalAccessToken>
      password-file: /path/to/file
    mirror:
      dir: /tmp/github
      default-action: exclude
      naming-style: slug
      rules:
        - rule: group == "my-org"
          action: include
```

The configuration is read from `~/.config/reposync/config.yaml` by default, but you can also specify a custom path by setting the `REPOSYNC_CONFIG` environment variable.

Supported platforms:

- `github`
- `gitlab`

> The `git` commands will use your local git installation, so you can use ssh keys or other authentication methods.
> The personal access tokens are only used to query the repositories you have access to and not to clone them.

## Installation

```bash
curl -L -o ~/.local/bin/reposync https://github.com/cidverse/reposync/releases/download/v0.3.0/linux_amd64
chmod +x ~/.local/bin/reposync
```

## Usage

### Index

> This is work-in-progress and not yet implemented.

Before your first run, you can use `reposync index /old-project-dir`, to add your local projects into the reposync state.
When running `reposync clone`, it will then move the projects to the new location, instead of cloning them.

### Clone

`reposync clone` will clone all projects you have access to into a target directory.
It will also move existing projects to mirror the remote structure (`<namespace>/<projectName>`).

### Update

`reposync update` will pull the latest changes from the remote server for tracked projects.

### HouseKeeping

`reposync housekeeping` (`reposync hk`) will run the following tasks for all repositories:

- `git prune --expire now`
- `git gc --auto`
- `git fsck --full --unreachable --strict`

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
