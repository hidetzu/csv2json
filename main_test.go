package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseDelimiter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    rune
		wantErr bool
	}{
		{"comma", ",", ',', false},
		{"tab", `\t`, '\t', false},
		{"pipe", "|", '|', false},
		{"semicolon", ";", ';', false},
		{"invalid", `\z`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDelimiter(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDelimiter(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDelimiter(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		comma   rune
		want    int // expected number of JSON objects
		wantNil bool
		wantErr bool
	}{
		{
			name:  "basic csv",
			input: "id,name,email\n1,taro,taro@sample.com\n2,hanako,hanako@sample.com\n",
			comma: ',',
			want:  2,
		},
		{
			name:  "tsv",
			input: "id\tname\n1\ttaro\n",
			comma: '\t',
			want:  1,
		},
		{
			name:    "empty input",
			input:   "",
			comma:   ',',
			wantNil: true,
		},
		{
			name:  "header only",
			input: "id,name,email\n",
			comma: ',',
			want:  0,
		},
		{
			name:    "column count mismatch",
			input:   "id,name\n1,taro,extra\n",
			comma:   ',',
			wantErr: true,
		},
		{
			name:  "quoted fields",
			input: "id,name\n1,\"taro, the great\"\n",
			comma: ',',
			want:  1,
		},
		{
			name:  "single column",
			input: "name\ntaro\nhanako\n",
			comma: ',',
			want:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert(strings.NewReader(tt.input), tt.comma, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.wantNil {
				if got != nil {
					t.Errorf("convert() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("convert() returned %d results, want %d", len(got), tt.want)
			}
		})
	}
}

func TestConvertFieldValues(t *testing.T) {
	input := "id,name\n1,taro\n"
	got, err := convert(strings.NewReader(input), ',', true)
	if err != nil {
		t.Fatalf("convert() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("convert() returned %d results, want 1", len(got))
	}
	if got[0]["id"] != "1" {
		t.Errorf("got[0][\"id\"] = %v, want \"1\"", got[0]["id"])
	}
	if got[0]["name"] != "taro" {
		t.Errorf("got[0][\"name\"] = %v, want \"taro\"", got[0]["name"])
	}
}

func TestRun(t *testing.T) {
	input := "id,name\n1,taro\n"
	var buf bytes.Buffer

	err := run(strings.NewReader(input), &buf, ',', true)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"id": "1"`) {
		t.Errorf("output missing id field: %s", output)
	}
	if !strings.Contains(output, `"name": "taro"`) {
		t.Errorf("output missing name field: %s", output)
	}
}

func TestRunEmptyInput(t *testing.T) {
	var buf bytes.Buffer

	err := run(strings.NewReader(""), &buf, ',', true)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty input, got: %s", buf.String())
	}
}
