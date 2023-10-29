# dep-doctor

`dep-doctor` is a tool to diagnose whether your software dependency packages are maintained.

Today, most software relies heavily on external packages. Vulnerabilities in those packages can be detected by vulnerability scanners ([dependabot](https://docs.github.com/en/code-security/dependabot), [trivy](https://aquasecurity.github.io/trivy), [Grype](https://github.com/anchore/grype), etc) if they are publicly available.

However, some packages have archived their source code repositories or have had their development stopped, although not explicitly. `dep-doctor` will notify you of those packages in the dependencies files.

![overview](doc/images/dep-doctor_overview.png "dep-doctor overview")

## Support dependencies files

| language | package manager | file (e.g.) | status |
| -------- | ------------- | -- | :----: |
| Ruby | bundler | Gemfile.lock | :heavy_check_mark: |
| JavaScript | yarn | yarn.lock | :heavy_check_mark: |
| JavaScript | npm | package-lock.json | :heavy_check_mark: |
| Python | pip | requirements.txt | :heavy_check_mark: |
| Python | poetry | poetry.lock | (later) |
| Python | pipenv | Pipfile.lock | (later) |
| PHP | composer | composer.lock | :heavy_check_mark: |
| Go | | go.sum | (later) |
| Rust | cargo | Cargo.lock | (later) |

## Install

### Homebrew (macOS and Linux)

```console
$ brew tap kyoshidajp/dep-doctor
$ brew install kyoshidajp/dep-doctor/dep-doctor
```

### Binary packages

[Releases](https://github.com/kyoshidajp/dep-doctor/releases)

## How to use

Set GitHub access token as `GITHUB_TOKEN` to your environment variable.

For example:

```console
$ dep-doctor diagnose -p bundler -file /path/to/Gemfile.lock
```

## How works

## Author
Katsuhiko YOSHIDA
