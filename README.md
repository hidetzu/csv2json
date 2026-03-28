# csv2json

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

CSV to JSON converter for the command line. Reads CSV from stdin and outputs JSON to stdout.

## Install

```sh
go install github.com/hidetzu/csv2json@latest
```

Or build from source:

```sh
git clone https://github.com/hidetzu/csv2json.git
cd csv2json
go build -o csv2json .
```

## Usage

```sh
cat input.csv | csv2json
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-delimiter` | `,` | Field separator (supports escape sequences like `\t`) |
| `-lazyQuote` | `true` | Allow lazy quoting in CSV fields |

## Examples

### Basic conversion

```sh
$ cat data/sample.csv
id,name,email
1,taro,taro@sample.com
2,hanako,hanako@sample.com

$ cat data/sample.csv | csv2json
[
  {
    "email": "taro@sample.com",
    "id": "1",
    "name": "taro"
  },
  {
    "email": "hanako@sample.com",
    "id": "2",
    "name": "hanako"
  }
]
```

### TSV input

```sh
cat data.tsv | csv2json -delimiter '\t'
```

### Pipe with other tools

```sh
cat data.csv | csv2json | jq '.[0]'
```

## Input / Output Specification

- **Input**: CSV format from stdin. The first row is treated as the header (JSON keys).
- **Output**: Pretty-printed JSON array to stdout. All values are output as strings.

## License

MIT
