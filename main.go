package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	defaultPath := "./testdata/1k.txt"

	filePath := flag.String("f", defaultPath, "Path to the file containing the math expression.")
	evalFlag := flag.Bool("eval", false, "Use EvalParse function for in-place evaluation.")
	profile := flag.Bool("profile", false, "Enable heap profiling.") // for mem analysis and optimisation purposes
	manualInput := flag.Bool("i", false, "Read input manually from stdin instead of from a file.")

	flag.Parse()

	var reader io.Reader
	var err error

	// Input is optionally from stdin or from a file
	if *manualInput {
		fmt.Println("Enter your math expression (CTRL+D to submit):")
		reader = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(*filePath)
		if err != nil {
			log.Fatalf("Could not open file %s: %v", *filePath, err)
		}
		defer file.Close()
		reader = file
	}

	// ** CPU Profiling **
	// Start cpu profiling for before parsing
	if *profile {
		var fileName string
		if *evalFlag {
			fileName = "cpu_profile_evalFlag.prof"
		} else {
			fileName = "cpu_profile.prof"
		}
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal("could not create heap profile:", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	exp, err := parseInput(reader, *evalFlag)
	if err != nil {
		log.Fatalf("Could not parse expression: %v", err)
	}

	// ** Mem Profiling **
	// Write heap profile after parsing
	if *profile {
		var fileName string
		if *evalFlag {
			fileName = "heap_profile_post_parse_evalFlag.prof"
		} else {
			fileName = "heap_profile_post_parse.prof"
		}
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal("could not create heap profile:", err)
		}
		defer f.Close()
		pprof.WriteHeapProfile(f)
	}

	res, err := exp.Eval()
	if err != nil {
		log.Fatalf("Failed evaluation: %v", err)
	}

	// ** Mem Profiling **
	// Write heap profile after evaluation
	if *profile {
		var fileName string
		if *evalFlag {
			fileName = "heap_profile_post_eval_evalFlag.prof"
		} else {
			fileName = "heap_profile_post_eval.prof"
		}
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal("could not create heap profile:", err)
		}
		defer f.Close()
		pprof.WriteHeapProfile(f)
	}

	printResult(exp, res)
}

func parseInput(reader io.Reader, useEval bool) (Expr, error) {
	if useEval {
		return EvalParse(reader)
	} else {
		return Parse(reader)
	}
}

func printResult(exp Expr, res float64) {
	// we use a new (English) printer for outputting thousands comma
	p := message.NewPrinter(language.English)

	if exp.Len() <= 1000 {
		p.Printf("Eval(%v) = %.2f\n", exp, res)
	} else {
		p.Printf("Eval() = %.2f\n", res)
	}
}
