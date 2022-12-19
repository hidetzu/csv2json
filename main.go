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
)

type JSON map[string]interface{}

func main() {
	flag.Parse()

	// Unquote Escaped Character (e.g. \t)
	rune, _, _, err := strconv.UnquoteChar(*delimiter, '"')
	if err != nil {
		log.Fatalf("Error: UnquoteChar fail: Input: '%s', Message: %v\n", *delimiter, err)
	}

	// create CSV reader from stdin
	r := csv.NewReader(os.Stdin)
	r.Comma = rune
	r.LazyQuotes = *lazyQuote

	results := []JSON{}

	// read header
	header, err := r.Read()
	if err == io.EOF {
		return
	}

	for {
		// read csv body
		rows, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error: csv Read fail: %v\n", err)
		}

		jsonData := make(JSON)
		for i := range rows {
			jsonData[header[i]] = string(rows[i])
		}
		results = append(results, jsonData)
	}

	// output json file
	json, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error: json.Marshal fail: Input: %v, Message: %v", results, err)
	}

	fmt.Printf("%s\n", json)
}
