<img width="100" src="https://raw.githubusercontent.com/loeffel-io/ls-lint/master/ls-lint.png" alt="logo">

# ls-lint

An extremely fast file and directory name linter

[![Build Status](http://ci.loeffel.io/api/badges/loeffel-io/ls-lint/status.svg)](http://ci.loeffel.io/loeffel-io/ls-lint)
[![Go Report Card](https://goreportcard.com/badge/github.com/loeffel-io/ls-lint)](https://goreportcard.com/report/github.com/loeffel-io/ls-lint)
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/v/@ls-lint/ls-lint.svg?sanitize=true" alt="Version"></a>
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/dm/@ls-lint/ls-lint?label=npm%20downloads" alt="NPM Downloads"></a>
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/l/@ls-lint/ls-lint.svg?sanitize=true" alt="License"></a>

- Works for directory and file names (all extensions supported)
- Incredibly fast
- Full unicode support
- Linux, MacOS & Windows Support
- Docker support & Npm package support
- Part of [Vue.js 3](https://github.com/vuejs/vue-next), [Nuxt.js](https://github.com/nuxt/nuxt.js) and [Vant](https://github.com/youzan/vant)
- Almost zero third-party dependencies (only [go-yaml](https://github.com/go-yaml/yaml) and [doublestar](https://github.com/bmatcuk/doublestar))

## Demo

<img src="https://i.imgur.com/plZml7D.gif" alt="command" width="600">

## Example & How-to ([vuejs/vue-next](https://github.com/vuejs/vue-next))

- `.ls-lint.yml` file must be present in your root directory
- Multiple rules supported by `|` - They are logicly *OR* combined
- `.dir` set rules for the current directory and their subdirectories
- Rules for subdirectories will overwrite the rules for all their subdirectories
- For Windows you can use backslashs `\` or slashs `/` - slashs recommenced

```yaml
# .ls-lint.yml

ls:
  .js: kebab-case
  .ts: camelCase
  .d.ts: kebab-case
  .mock.ts: kebab-case
  .spec.ts: camelCase
  .test-d.ts: kebab-case
  .config.js: kebab-case
  .umd.js: kebab-case
  .spec.ts.snap: camelCase

  scripts:
    .js: camelCase

  packages/**/{components,collections}:
    .ts: PascalCase
    .spec.ts: PascalCase

ignore:
  - node_modules
  - .git
  - .circleci
  - .github
  - .vscode
```

## Install & Run

## Binary

### MacOS

```bash
curl -sL -o ls-lint https://github.com/loeffel-io/ls-lint/releases/download/v1.9.0/ls-lint-darwin && chmod +x ls-lint && ./ls-lint
```

### Linux

```bash
curl -sL -o ls-lint https://github.com/loeffel-io/ls-lint/releases/download/v1.9.0/ls-lint-linux && chmod +x ls-lint && ./ls-lint
```

### Windows

```bash
# (!) First download the .exe from https://github.com/loeffel-io/ls-lint/releases/download/v1.9.0/ls-lint-windows.exe
ls-lint-windows.exe
```

## NPM

### Install

```bash
# global
npm install -g @ls-lint/ls-lint

# local
npm install @ls-lint/ls-lint
```

### Run

```bash
# global
ls-lint

# local
node_modules/.bin/ls-lint # use backslashs for windows

npx @ls-lint/ls-lint
```

## Docker

```bash
docker run -t -v /path/to/files:/data lslintorg/ls-lint:1
```

## Rules

| Rule       | Alias       | Description                                                                        |
| ---------- | ----------- | ---------------------------------------------------------------------------------- | 
| regex      | -           | Checks if string matches regex pattern: ^{pattern}$                                |
| lowercase  | -           | Checks if every letter is lower; Skip non letters                                  |
| camelcase  | camelCase   | Checks if string is camel case; Only letters and digits allowed                    |
| pascalcase | PascalCase  | Checks if string is pascal case; Only letters and digits allowed                   |
| snakecase  | snake_case  | Checks if string is snake case; Only lowercase letters, digits and `_` allowed     | 
| kebabcase  | kebab-case  | Checks if string is kebab case; Only lowercase letters, digits and `-` allowed     |
| pointcase  | point.case  | Checks if string is "point case"; Only lowercase letters, digits and `.` allowed   |

## Roadmap

- [ ] Public and Private Registry to share configurations based on [go-saas](https://github.com/go-saas)
- [x] Npm Windows package (one package for all os)
- [x] Docker support
- [x] Regex Rule
- [x] Windows support
- [x] Npm package
- [x] Add ignore directories and files

## Major changes

**v1.9.0**

- Added path separator replacement: you can now use `/` on any os instead of e.g. using `\` for windows machines in your `ls-lint.yml` file

**v1.8.0**

- Added glob support like `packages/**` or `packages/*/src`

**v1.7.0**

- Rules improved: more tests, more flexibility, digits allowed. Checkout [rules](https://github.com/loeffel-io/ls-lint#rules) for more informations

**v1.6.0**

- Rules are not longer logicly `AND` combined - Now they are logicly `OR` combined by `|`

**v1.5.0**

- Npm packages `ls-lint-darwin` and `ls-lint-linux` are not longer supported. Please use `@ls-lint/ls-lint` instead (linux, windows and macOS support)

## Benchmarks ([hyperfine](https://github.com/sharkdp/hyperfine))

| Package                                              | Mean [s]            | File                                                                                                              | 
| ---------------------------------------------------- | ------------------- | ----------------------------------------------------------------------------------------------------------------- |
| [vuejs/vue](https://github.com/vuejs/vue)            | 283.3 ms ± 19.6 ms  | [examples/vuejs-vue](https://github.com/loeffel-io/ls-lint/tree/master/examples/vuejs-vue/.ls-lint.yml)           |
| [vuejs/vue-next](https://github.com/vuejs/vue-next)  | 267.3 ms ±   9.3 ms | [examples/vuejs-vue-next](https://github.com/loeffel-io/ls-lint/tree/master/examples/vuejs-vue-next/.ls-lint.yml) |

## Logo

Logo created by [Anastasia Marx](https://www.behance.net/AnastasiaMarx)

## License

ls-lint is open-source software licensed under the MIT license.
