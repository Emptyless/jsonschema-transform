# JSON-Schema Transform

_This CLI tool is under active development and must be considered alpha. It's API may be changed in a breaking way until a 1.0 version is released. Submit issues to the Github issue tracker if found._

`jsonschema-transform` is a CLI utility that can transform [JSON Schema](https://json-schema.org) into [D2](https://d2lang.com). This can be useful to visualize a large set of JSON Schema's that are interconnected (and possibly have remote references).

### Quickstart

With default arguments, running `jsonschema-transform d2` to create the `d2` representation of the `testdata` directory: 

```
$ jsonschema-transform d2 --globs ./testdata/*.json
```

A new output file named `diagram.d2` is created in the current working directory. To create an SVG, change the output to a SVG filetype:

```
$ jsonschema-transform d2 --globs ./testdata/*.json --output diagram.svg
```

<img src="./diagram.svg" width=500 height=500>

### Installation

```
go install github.com/Emptyless/jsonschema-transform
```

Also make sure that `d2` is installed:

```
brew install d2
```

### Usage

- `--globs`: to match containing JSON Schema documents, e.g. `*/*.json` or `./testdata/pet.json`
- `--base-uri`: to use for fetching relative $refs, including `file://` based $refs
- `--overwrite`: allow overwrite of output file if the file exists already
- `--output`: name of the output file (extension must be either 'svg' or 'd2')

### TODO's

- [x] get basic structure of CLI working
- [ ] test against more complex JSON Schema's
- [ ] generate markdown (MD) files from the JSON schemas with clickable links

### Mentions

- [kaptinlin/jsonschema](github.com/kaptinlin/jsonschema): a Go JSON Schema library

