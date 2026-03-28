package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	delimiter *string = flag.String("delimiter", ",", "specify separator (e.g. \"\\t\")")
	lazyQuote *bool   = flag.Bool("lazyQuote", true, "allow lazyQuote")
	output    *string = flag.String("o", "", "output file path (default: stdout)")
	jsonl     *bool   = flag.Bool("jsonl", false, "output as JSON Lines")
)

type JSON map[string]interface{}

func parseDelimiter(s string) (rune, error) {
	r, _, _, err := strconv.UnquoteChar(s, '"')
	if err != nil {
		return 0, fmt.Errorf("invalid delimiter '%s': %w", s, err)
	}
	return r, nil
}

func convert(r io.Reader, comma rune, lazyQuotes bool) ([]JSON, error) {
	reader := csv.NewReader(r)
	reader.Comma = comma
	reader.LazyQuotes = lazyQuotes

	header, err := reader.Read()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	reader.FieldsPerRecord = len(header)

	var results []JSON
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		jsonData := make(JSON)
		for i := range row {
			jsonData[header[i]] = row[i]
		}
		results = append(results, jsonData)
	}

	return results, nil
}

func formatJSON(results []JSON) ([]byte, error) {
	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return output, nil
}

func formatJSONL(results []JSON) ([]byte, error) {
	var buf []byte
	for _, r := range results {
		line, err := json.Marshal(r)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		buf = append(buf, line...)
		buf = append(buf, '\n')
	}
	return buf, nil
}

func run(r io.Reader, w io.Writer, comma rune, lazyQuotes bool, asJSONL bool) error {
	results, err := convert(r, comma, lazyQuotes)
	if err != nil {
		return err
	}
	if results == nil {
		return nil
	}

	var out []byte
	if asJSONL {
		out, err = formatJSONL(results)
	} else {
		out, err = formatJSON(results)
	}
	if err != nil {
		return err
	}

	if !asJSONL {
		out = append(out, '\n')
	}
	_, err = w.Write(out)
	return err
}

func main() {
	flag.Parse()

	comma, err := parseDelimiter(*delimiter)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Determine input source
	var r io.Reader
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}

	// Determine output destination
	var w io.Writer
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer f.Close()
		w = f
	} else {
		w = os.Stdout
	}

	if err := run(r, w, comma, *lazyQuote, *jsonl); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
