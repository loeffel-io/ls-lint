<img width="100" src="https://raw.githubusercontent.com/loeffel-io/ls-lint/master/assets/logo/ls-lint.png" alt="logo">

# ls-lint

An extremely fast directory and filename linter - Bring some structure to your project filesystem

[![CI](https://github.com/loeffel-io/ls-lint/actions/workflows/bazel.yml/badge.svg?branch=master)](https://github.com/loeffel-io/ls-lint/actions/workflows/bazel.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/loeffel-io/ls-lint)](https://goreportcard.com/report/github.com/loeffel-io/ls-lint)
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/v/@ls-lint/ls-lint.svg?sanitize=true" alt="Version"></a>
![npm](https://img.shields.io/npm/dy/@ls-lint/ls-lint?label=npm%20downloads%20total)
![npm](https://badgen.net/static/npm%20total%20downloads/4M+/green)
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/l/@ls-lint/ls-lint.svg?sanitize=true" alt="License"></a>

- Minimal setup with simple rules managed in one single or multiple `.ls-lint.yml` files
- Works for directory and file names - all extensions supported - full unicode support
- Incredibly fast - lints thousands of files and directories in milliseconds
- Support for Windows, MacOS and Linux + NPM Package + [GitHub Action](https://github.com/ls-lint/action) + [Homebrew](https://formulae.brew.sh/formula/ls-lint) + & Docker Image
- Almost zero third-party dependencies (only [go-yaml](https://github.com/go-yaml/yaml)
  and [doublestar](https://github.com/bmatcuk/doublestar))

## Documentation

The full documentation can be found at [ls-lint.org](https://ls-lint.org)

- [Installation](https://ls-lint.org/2.2/getting-started/installation.html#curl)
- [The Basics](https://ls-lint.org/2.2/configuration/the-basics.html)
- [The Rules](https://ls-lint.org/2.2/configuration/the-basics.html)
- [Contributions](https://ls-lint.org/2.2/prologue/contributions.html)

## Demo

### Configuration `.ls-lint.yml`

```yaml
ls:
  .js: snake_case
  .ts: snake_case | camelCase
  .d.ts: PascalCase
  .html: regex:[a-z0-9]+

ignore:
  - node_modules
```

### Result

<img src="https://i.imgur.com/pxXkYcl.gif" alt="command" width="600">

## Discord

[Join the ls-lint discord server](https://discord.gg/bsf9q7f2Rh)

## Sponsors

<a href="https://jetbrains.com"><img height="130" src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png?_ga=2.249742848.788370738.1691416665-1384286648.1691416665" alt="jetbrains"></a>

## Logo

Logo created by [Studio Ajot](https://www.studio-ajot.de/)

## License

ls-lint is open-source software licensed under the MIT license.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Floeffel-io%2Fls-lint.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Floeffel-io%2Fls-lint?ref=badge_large)
