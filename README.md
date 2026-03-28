# csv2json

[![CI](https://github.com/hidetzu/csv2json/actions/workflows/ci.yml/badge.svg)](https://github.com/hidetzu/csv2json/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hidetzu/csv2json)](https://goreportcard.com/report/github.com/hidetzu/csv2json)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A simple CSV to JSON converter for the command line.

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
csv2json input.csv
cat input.csv | csv2json
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-delimiter` | `,` | Field separator (supports escape sequences like `\t`) |
| `-lazyQuote` | `true` | Allow lazy quoting in CSV fields |
| `-o` | (stdout) | Output file path |
| `-jsonl` | `false` | Output as JSON Lines |

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

### File input

```sh
csv2json data/sample.csv
```

### Output to file

```sh
csv2json -o output.json data/sample.csv
```

### JSON Lines output

```sh
csv2json -jsonl data/sample.csv
```

### Pipe with other tools

```sh
cat data.csv | csv2json | jq '.[0]'
```

## Input / Output Specification

- **Input**: CSV format from stdin or a file argument. The first row is treated as the header (JSON keys).
- **Output**: Pretty-printed JSON array (default) or JSON Lines (`-jsonl`) to stdout or a file (`-o`).

## License

MIT
