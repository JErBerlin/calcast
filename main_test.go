package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// Benchmark files represent different sizes of arithmetic expressions.
var benchmarkFiles = []string{"1k.txt", "10k.txt", "100k.txt", "1m.txt", "10m.txt"}

func TestParseAndEvalPerformance(t *testing.T) {
	path := "./testdata/"
	for _, file := range benchmarkFiles {
		t.Run(fmt.Sprintf("Test data %s", strings.Split(file, ".")[0]), func(t *testing.T) {
			file = path + file
			f, err := os.Open(file)
			if err != nil {
				t.Fatalf("could not open file %s: %v", file, err)
			}
			defer f.Close()

			// Measure the time taken by Parse function
			startParse := time.Now()
			expr, err := Parse(f)
			// expr, err := EvalParse(f)
			durationParse := time.Since(startParse)

			if err != nil {
				t.Fatalf("could not parse expression from %s: %v", file, err)
			}

			// Measure the time taken by Eval function
			startEval := time.Now()
			_, err = expr.Eval()
			durationEval := time.Since(startEval)

			if err != nil {
				t.Fatalf("could not evaluate expression from %s: %v", file, err)
			}

			// Print out the times
			t.Logf("Parse time for %s: %v", file, durationParse)
			t.Logf("Eval time for %s: %v", file, durationEval)
		})
	}
}

func TestEvalParseAndEvalPerformance(t *testing.T) {
	path := "./testdata/"
	for _, file := range benchmarkFiles {
		t.Run(fmt.Sprintf("Test data %s", strings.Split(file, ".")[0]), func(t *testing.T) {
			file = path + file

			f, err := os.Open(file)
			if err != nil {
				t.Fatalf("could not open file %s: %v", file, err)
			}
			defer f.Close()

			// Measure the time taken by Parse function
			startEvalParse := time.Now()
			// expr, err := Parse(f)
			expr, err := EvalParse(f)
			durationEvalParse := time.Since(startEvalParse)

			if err != nil {
				t.Fatalf("could not parse expression from %s: %v", file, err)
			}

			// Measure the time taken by Eval function
			startEval := time.Now()
			_, err = expr.Eval()
			durationEval := time.Since(startEval)

			if err != nil {
				t.Fatalf("could not evaluate expression from %s: %v", file, err)
			}

			// Print out the times
			t.Logf("EvalParse time for %s: %v", file, durationEvalParse)
			t.Logf("Eval time for %s: %v", file, durationEval)
		})
	}
}
