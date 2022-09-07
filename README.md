# reposync

> A cli tool to mirror/sync many projects onto the local file system (and/or merge content of specific folders to aggregate ie. doc files)

## Installation

```bash
curl -L -o /usr/local/bin/reposync https://github.com/cidverse/reposync/releases/download/v0.1.0/linux_amd64
chmod +x /usr/local/bin/reposync
```

## Usage

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

### Bundle

Bundle is a utility command to merge specific folders of multiple projects, this could for example be used to merge the doc folders from many projects to generate a multi-project documentation.

```yaml
bundle:
  docs:
    target: ~/source/my-docs
    sources:
      - url: https://github.com/cidverse/reposync.git
        ref: HEAD
        bundle:
          source-prefix: docs
          target-prefix: reposync
          extensions: [".go"]
```

## License

Released under the [MIT license](./LICENSE).
