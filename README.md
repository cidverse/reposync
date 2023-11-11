# RepoSync

> RepoSync allows you to clone all projects you have membership of into a local directory, or mirror specific projects into a local directory.
> You can use complex rules to include or exclude projects based on the namespace or other properties.

## Installation

```bash
curl -L -o ~/.local/bin/reposync https://github.com/cidverse/reposync/releases/download/v0.3.0/linux_amd64
chmod +x ~/.local/bin/reposync
```

## Usage

### Index

> This is work-in-progress and not yet implemented.

Before your first run, you can use `reposync index /old-project-dir`, to add your local projects into the reposync state.
When running `reposync update`, it will then move the projects to the new location, instead of cloning them.

### Update

You can define multiple projects, add a platform to clone all repositories you are a member of, or both. The projects will mirror the structure from the remote server locally (repositories will be moved if necessary).

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/cidverse/reposync/main/configschema/v1.json
servers:
  - url: https://github.com
    type: github
    auth:
      username: YourAccount
      password: <readOnlyPersonalAccessToken>
    mirror:
      dir: /tmp/github
      default-action: exclude
      rules:
        - rule: group == "my-org"
          action: include
```

Supported platforms:
- `github`
- `gitlab`

> The `git clone` will use your local git, so you can use ssh keys or other authentication methods. The personal access tokens are only used to query the repositories you have access to.

### Sync

Sync will check out all defined projects locally in the defined structure, can be grouped to clone specific project groups.

```yaml
sources:
- url: https://github.com/cidverse/reposync.git
  ref: HEAD
  group:
    - test
  target: ~/source/my-project
```

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
