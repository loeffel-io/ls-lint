<img width="100" src="https://raw.githubusercontent.com/loeffel-io/ls-lint/master/ls-lint.png" alt="logo">

# ls-lint

An extremely fast directory and filename linter - Bring some structure to your project directories

[![Build Status](https://cloud.drone.io/api/badges/loeffel-io/ls-lint/status.svg)](https://cloud.drone.io/loeffel-io/ls-lint)
[![Go Report Card](https://goreportcard.com/badge/github.com/loeffel-io/ls-lint)](https://goreportcard.com/report/github.com/loeffel-io/ls-lint)
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/v/@ls-lint/ls-lint.svg?sanitize=true" alt="Version"></a>
[![codecov](https://codecov.io/gh/loeffel-io/ls-lint/branch/master/graph/badge.svg?token=76N7CNQT2G)](https://codecov.io/gh/loeffel-io/ls-lint)
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/dm/@ls-lint/ls-lint?label=npm%20downloads" alt="NPM Downloads"></a>
![npm](https://img.shields.io/npm/dt/@ls-lint/ls-lint?label=npm%20downloads%20total)
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/l/@ls-lint/ls-lint.svg?sanitize=true" alt="License"></a>

- Minimal setup with simple rules managed in one single `.ls-lint.yml` file
- Works for directory and file names - all extensions supported - full unicode support
- Incredibly fast - lints thousands of files and directories in milliseconds
- Support for Windows, MacOS and Linux + NPM Package & Docker Image
- ARM support
- Almost zero third-party dependencies (only [go-yaml](https://github.com/go-yaml/yaml)
  and [doublestar](https://github.com/bmatcuk/doublestar))

## Documentation

The full documentation can be found at [ls-lint.org](https://ls-lint.org)

- [Installation](https://ls-lint.org/1.x/getting-started/installation.html#curl)
- [The Basics](https://ls-lint.org/1.x/configuration/the-basics.html)
- [The Rules](https://ls-lint.org/1.x/configuration/the-basics.html)
- [Contributions](https://ls-lint.org/1.x/prologue/contributions.html)

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

## Benchmarks ([hyperfine](https://github.com/sharkdp/hyperfine))

| Package                                          | v1.10.x: Mean [s]  | v1.11.x: Mean [s]  | File                                                                                                          | 
|--------------------------------------------------|--------------------|--------------------|---------------------------------------------------------------------------------------------------------------|
| [nuxt/nuxt.js](https://github.com/nuxt/nuxt.js)  | 627.2 ms ± 19.0 ms | 51.9 ms ± 3.2 ms   | [examples/nuxt-nuxt-js](https://github.com/loeffel-io/ls-lint/tree/master/examples/nuxt-nuxt-js/.ls-lint.yml) |

## Discord

[Join the ls-lint discord server](https://discord.gg/bsf9q7f2Rh)

## Sponsors

- [Makeless - Saas Framework](https://github.com/makeless)

## Logo

Logo created by [Anastasia Marx](https://www.behance.net/AnastasiaMarx)

## License

ls-lint is open-source software licensed under the MIT license.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Floeffel-io%2Fls-lint.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Floeffel-io%2Fls-lint?ref=badge_large)