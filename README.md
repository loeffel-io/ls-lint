# ls-lint

File and Directory name linter

[![Build Status](http://ci.loeffel.io/api/badges/loeffel-io/ls-lint/status.svg)](http://ci.loeffel.io/loeffel-io/ls-lint)

- Works for directory and file names
- Incredibly fast
- Linux & Mac Support (Windows coming soon)
 
## Example

```yaml
# .ls-lint.yml

ls:
  src:
    .dir: lowercase
    .js: snake_case
    .json: snake_case
    .vue: PascalCase
```

## Rules 

| Rule       | Alias       | Description                                                  |
| ---------- | ----------- | ------------------------------------------------------------ |
| lowercase  | -           | Checks if every letter is lower; Skip non letters            |
| camelcase  | camelCase   | Checks if string is camel case; Only letters allowed         |
| pascalcase | PascalCase  | Checks if string is pascal case; Only letters allowed        |
| snakecase  | snake_case  | Checks if string is snake case; Only letters and `_` allowed |


