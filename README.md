# example-finder

## Build:
Dependencies:
```
gcc (required for sqlite3 support)

Fedora:
dnf install -y gcc

Windows:
http://mingw-w64.org/doku.php/download#mingw-builds
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

Example usage:
./example-finder search go-pg -r 100 -t ABC
```

## Limitations
* Github API rate limiting currently prevents us from getting all possible results. We are limited to 2K request per hour, making this the upper bound to the amount of results we can actually get. We are also limited by 20 searches per minute, which can hinder the user experience as well.
This can potentially be solved by using the GraphQL API, however, code search is not implemented yet.
