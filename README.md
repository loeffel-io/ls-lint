# ls-lint

File and Directory name linter

[![Build Status](http://ci.loeffel.io/api/badges/loeffel-io/ls-lint/status.svg)](http://ci.loeffel.io/loeffel-io/ls-lint)

- Works for directory and file names (all extensions supported)
- Linux & Mac Support (Windows coming soon)
- Incredibly fast
- Full unicode support
- Almost zero third-party dependencies (only [go-yaml](https://github.com/go-yaml/yaml))

## Demo

<img src="https://i.imgur.com/plZml7D.gif" alt="command" width="600">

## Example

- `.ls-lint.yml` file must be present in your root directory
- Multiple rules supported by `,`
- `.dir` set rules for the current directory and their subdirectories
- Rules for subdirectories will overwrite the rules for all their subdirectories

```yaml
# .ls-lint.yml

ls:
  .dir: lowercase, kebab-case 
  .js: snake_case
  .vue: PascalCase

  src: # set new rules for the src subdirectory and all their subdirectories
      .dir: lowercase, kebab-case
      .js: kebab-case
      .vue: PascalCase

ignore:
  - .idea
  - .git
```

## Rules

| Rule       | Alias       | Description                                                    |
| ---------- | ----------- | -------------------------------------------------------------- |
| lowercase  | -           | Checks if every letter is lower; Skip non letters              |
| camelcase  | camelCase   | Checks if string is camel case; Only letters allowed           |
| pascalcase | PascalCase  | Checks if string is pascal case; Only letters allowed          |
| snakecase  | snake_case  | Checks if string is snake case; Only letters and `_` allowed   |
| kebabcase  | kebab-case  | Checks if string is kebab case; Only letters and `-` allowed   |
| pointcase  | point.case  | Checks if string is "point case"; Only letters and `.` allowed |

## Roadmap

- Windows support
- Docker support
- Add more rules
- Add ignore directories and files

## License

ls-lint is open-source software licensed under the MIT license.