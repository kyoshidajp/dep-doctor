![GitHub](https://img.shields.io/github/license/kyoshidajp/dep-doctor)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kyoshidajp/dep-doctor)
![GitHub all releases](https://img.shields.io/github/downloads/kyoshidajp/dep-doctor/total)
![GitHub CI Status](https://img.shields.io/github/actions/workflow/status/kyoshidajp/dep-doctor/ci.yaml?branch=main)
![GitHub Release Status](https://img.shields.io/github/actions/workflow/status/kyoshidajp/dep-doctor/release.yaml?branch=main)

# dep-doctor

`dep-doctor` is a tool to diagnose whether your software dependency packages are maintained.

## Support dependencies files

| Language | file | status |
| -------- | --------------- | :----: |
| Ruby | Gemfile.lock | YES |
| Ruby | gemspec | NO (but soon) |
| JavaScript | yarn.lock | YES |
| JavaScript | package.json | NO (but soon) |
| Go | go.sum | NO (but soon) |

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

## Author
Katsuhiko YOSHIDA
