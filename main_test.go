package main

import (
	"bytes"
	"os"
	"path/filepath"
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
		want    int
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

	err := run(strings.NewReader(input), &buf, ',', true, false)
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

	err := run(strings.NewReader(""), &buf, ',', true, false)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty input, got: %s", buf.String())
	}
}

func TestRunJSONL(t *testing.T) {
	input := "id,name\n1,taro\n2,hanako\n"
	var buf bytes.Buffer

	err := run(strings.NewReader(input), &buf, ',', true, true)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %s", len(lines), buf.String())
	}
	if !strings.Contains(lines[0], `"id":"1"`) {
		t.Errorf("line 0 missing id: %s", lines[0])
	}
	if !strings.Contains(lines[1], `"id":"2"`) {
		t.Errorf("line 1 missing id: %s", lines[1])
	}
}

func TestRunJSONLEmptyInput(t *testing.T) {
	var buf bytes.Buffer

	err := run(strings.NewReader(""), &buf, ',', true, true)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty input, got: %s", buf.String())
	}
}

func TestFileInput(t *testing.T) {
	// Create temp CSV file
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "test.csv")
	err := os.WriteFile(csvPath, []byte("id,name\n1,taro\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	f, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("failed to open temp file: %v", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	err = run(f, &buf, ',', true, false)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if !strings.Contains(buf.String(), `"id": "1"`) {
		t.Errorf("output missing id field: %s", buf.String())
	}
}

func TestOutputFile(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "output.json")

	f, err := os.Create(outPath)
	if err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}

	input := "id,name\n1,taro\n"
	err = run(strings.NewReader(input), f, ',', true, false)
	f.Close()
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !strings.Contains(string(data), `"id": "1"`) {
		t.Errorf("output file missing id field: %s", string(data))
	}
}
