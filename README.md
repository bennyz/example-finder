# example-finder

## Build:
Dependencies:
```
gcc (required for sqlite3 support)
```
build:

```
go build
```

test:
```
go test ./...
```

*Work in progress*
## Usage:
```
./example-finder search TEXT [-l|--lang LANG] [-t|--token TOKEN] [-m|--mode MODE] [-r|--results num]
  lang    - language to return results in (not implemented yet)
  token   - Github API token, see https://github.com/settings/tokens - should be inside .token file
  mode    - currently only REST mode is supported, a future plan is to support GraphQL as well
  results - number of results per page to return (default 30)
```
