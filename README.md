<img width="100" src="https://raw.githubusercontent.com/loeffel-io/ls-lint/master/ls-lint.png" alt="logo">

# ls-lint

An extremely fast file and directory name linter - Bring some structure to your project directories

[![Build Status](http://ci.loeffel.io/api/badges/loeffel-io/ls-lint/status.svg)](http://ci.loeffel.io/loeffel-io/ls-lint)
[![Go Report Card](https://goreportcard.com/badge/github.com/loeffel-io/ls-lint)](https://goreportcard.com/report/github.com/loeffel-io/ls-lint)
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/v/@ls-lint/ls-lint.svg?sanitize=true" alt="Version"></a>
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/dm/@ls-lint/ls-lint?label=npm%20downloads" alt="NPM Downloads"></a>
<a href="https://www.npmjs.com/package/@ls-lint/ls-lint"><img src="https://img.shields.io/npm/l/@ls-lint/ls-lint.svg?sanitize=true" alt="License"></a>

- Minimal setup with simple rules managed in one single `.ls-lint.yml` file
- Works for directory and file names - all extensions supported - full unicode support
- Incredibly fast - lints thousands of files and directories in milliseconds
- Support for Windows, MacOS and Linux + NPM Package & Docker Image
- Trusted by [Vue.js 3](https://github.com/vuejs/vue-next), [Nuxt.js](https://github.com/nuxt/nuxt.js), [Vant](https://github.com/youzan/vant) and many more
- Almost zero third-party dependencies (only [go-yaml](https://github.com/go-yaml/yaml) and [doublestar](https://github.com/bmatcuk/doublestar))

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

| Package                                              | Mean [s]            | File                                                                                                              | 
| ---------------------------------------------------- | ------------------- | ----------------------------------------------------------------------------------------------------------------- |
| [vuejs/vue](https://github.com/vuejs/vue)            | 283.3 ms ± 19.6 ms  | [examples/vuejs-vue](https://github.com/loeffel-io/ls-lint/tree/master/examples/vuejs-vue/.ls-lint.yml)           |
| [vuejs/vue-next](https://github.com/vuejs/vue-next)  | 267.3 ms ±   9.3 ms | [examples/vuejs-vue-next](https://github.com/loeffel-io/ls-lint/tree/master/examples/vuejs-vue-next/.ls-lint.yml) |

## Logo

Logo created by [Anastasia Marx](https://www.behance.net/AnastasiaMarx)

## License

ls-lint is open-source software licensed under the MIT license.
