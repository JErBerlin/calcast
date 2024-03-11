package main

import (
	"bytes"
	"os"
	"testing"
)

func benchmarkParseAndEval(fileName string, b *testing.B) {
	// Read the entire file content into memory
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		b.Fatalf("could not read file %s: %v", fileName, err)
	}

	for i := 0; i < b.N; i++ {
		// Create a new bytes.Reader for each iteration
		// this is a faster way than opening the file or resetting the reader index every time
		reader := bytes.NewReader(fileContent)

		expr, _ := Parse(reader) // Ignore errors while benchmarking
		expr.Eval()
	}
}

// Series of benchmark functions for each file size.
func BenchmarkParseAndEval_1k(b *testing.B)   { benchmarkParseAndEval("./testdata/1k.txt", b) }
func BenchmarkParseAndEval_10k(b *testing.B)  { benchmarkParseAndEval("./testdata/10k.txt", b) }
func BenchmarkParseAndEval_100k(b *testing.B) { benchmarkParseAndEval("./testdata/100k.txt", b) }
func BenchmarkParseAndEval_1m(b *testing.B)   { benchmarkParseAndEval("./testdata/1m.txt", b) }
func BenchmarkParseAndEval_10m(b *testing.B)  { benchmarkParseAndEval("./testdata/10m.txt", b) }

func benchmarkEvalParseAndEval(fileName string, b *testing.B) {
	// Read the entire file content into memory
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		b.Fatalf("could not read file %s: %v", fileName, err)
	}

	for i := 0; i < b.N; i++ {
		// Create a new bytes.Reader for each iteration
		// this is a faster way than opening the file or resetting the reader index every time
		reader := bytes.NewReader(fileContent)

		// Ignore errors while benchmarking
		EvalParse(reader)
	}
}

// Series of benchmark functions for each file size.
func BenchmarkEvalParseAndEval_1k(b *testing.B)  { benchmarkEvalParseAndEval("./testdata/1k.txt", b) }
func BenchmarkEvalParseAndEval_10k(b *testing.B) { benchmarkEvalParseAndEval("./testdata/10k.txt", b) }
func BenchmarkEvalParseAndEval_100k(b *testing.B) {
	benchmarkEvalParseAndEval("./testdata/100k.txt", b)
}
func BenchmarkEvalParseAndEval_1m(b *testing.B)  { benchmarkEvalParseAndEval("./testdata/1m.txt", b) }
func BenchmarkEvalParseAndEval_10m(b *testing.B) { benchmarkEvalParseAndEval("./testdata/10m.txt", b) }
