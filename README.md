# yamlfmt
Yet another yaml formatter

## Use

usage examples
```
# print help
yamlfmt -h

# quietly format .yaml and .yml files in the current dir
yamlfmt ./ -q

# recurisvely print differences in .yaml and .yml, dry run
yamlfmt ./... -v -d

```

### formatting 
changes applied to the files are as follows:

* the indentation is set to 2 spaces.
* line breaks are reduced to only one.
* quotes in keys and values are removed whenever it is possible.


## Development

#### Requirements

* go
* make
* goreleaser
* golangci-lint
* git

#### Release

make sure you have your gh token stored locally in ~/.goreleaser/gh_token

to release a new version:
```bash 
make release  version="v0.1.2"
```
