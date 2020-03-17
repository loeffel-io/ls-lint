<img width="100" src="https://raw.githubusercontent.com/loeffel-io/ls-lint/master/ls-lint.png" alt="logo">

# ls-lint

An extremely fast file and directory name linter

[![Build Status](http://ci.loeffel.io/api/badges/loeffel-io/ls-lint/status.svg)](http://ci.loeffel.io/loeffel-io/ls-lint)
[![Go Report Card](https://goreportcard.com/badge/github.com/loeffel-io/ls-lint)](https://goreportcard.com/report/github.com/loeffel-io/ls-lint)
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/v/@ls-lint/ls-lint.svg?sanitize=true" alt="Version"></a>
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/l/@ls-lint/ls-lint.svg?sanitize=true" alt="License"></a>

- Works for directory and file names (all extensions supported)
- Linux, MacOS & Windows Support
- Docker support
- Incredibly fast
- Full unicode support
- Almost zero third-party dependencies (only [go-yaml](https://github.com/go-yaml/yaml))

## Demo

<img src="https://i.imgur.com/plZml7D.gif" alt="command" width="600">

## Example & How-to ([vue.js](https://github.com/vuejs/vue))

- `.ls-lint.yml` file must be present in your root directory
- Multiple rules supported by `|` - They are logicly *OR* combined
- `.dir` set rules for the current directory and their subdirectories
- Rules for subdirectories will overwrite the rules for all their subdirectories
- For Windows you must use backslashs `\` instead of slashs `/` 

```yaml
# .ls-lint.yml

ls:
  .dir: regex:[a-z0-9\-]+
  .js: kebab-case
  .css: kebab-case
  .html: kebab-case
  .json: kebab-case
  .ts: kebab-case
  .sh: kebab-case
  .dev.js: kebab-case
  .prod.js: kebab-case
  .d.ts: kebab-case
  .vdom.js: kebab-case
  .spec.js: kebab-case

  dist:
    .js: point.case

  benchmarks/ssr:
    .js: camelCase

ignore:
  - test
  - benchmarks/dbmon/ENV.js
  - .babelrc.js
  - .eslintrc.js
  - .github
  - .circleci
  - .git
```

## Install & Run

## Binary

### MacOS

```bash
curl -sL -o ls-lint https://github.com/loeffel-io/ls-lint/releases/download/v1.7.0/ls-lint-darwin && chmod +x ls-lint && ./ls-lint
```

### Linux

```bash
curl -sL -o ls-lint https://github.com/loeffel-io/ls-lint/releases/download/v1.7.0/ls-lint-linux && chmod +x ls-lint && ./ls-lint
```

### Windows

```bash
# (!) First download the .exe from https://github.com/loeffel-io/ls-lint/releases/download/v1.7.0/ls-lint-windows.exe
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

- [ ] Public and Private Registry to share configurations
- [x] Npm Windows package (one package for all os)
- [x] Docker support
- [x] Regex Rule
- [x] Windows support
- [x] Npm package
- [x] Add ignore directories and files

## Major changes

**v1.7.0**

- Rules improved: more tests, more flexibility, digits allowed. Checkout [rules](https://github.com/loeffel-io/ls-lint#rules) for more informations

**v1.6.0**

- Rules are not longer logicly `AND` combined - Now they are logicly `OR` combined by `|`

**v1.5.0**

- Npm packages `ls-lint-darwin` and `ls-lint-linux` are not longer supported. Please use `@ls-lint/ls-lint` instead (linux, windows and macOS support)

## Benchmarks ([hyperfine](https://github.com/sharkdp/hyperfine))

| Package                                    | Mean [s]           | File                                                                                                    | 
| ------------------------------------------ | ------------------ | ------------------------------------------------------------------------------------------------------- |
| [vuejs/vue](https://github.com/vuejs/vue)  | 14.6 ms ± 1.1 ms   | [examples/vuejs-vue](https://github.com/loeffel-io/ls-lint/tree/master/examples/vuejs-vue/.ls-lint.yml) |

## Logo

Logo created by [Anastasia Marx](https://www.behance.net/AnastasiaMarx)

## License

ls-lint is open-source software licensed under the MIT license.
