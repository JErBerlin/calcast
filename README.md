# Syntactic Calculator

This is a calculator that reads mathematical terms containing floating point numbers, +, -, * and / as well as parenthesis.

The calculator CLI supports various use cases through flags for file input, manual input, evaluation method selection, and profiling. Below are examples on how to use these flags for different scenarios:

## File Input

To process a mathematical expression from a file, use the -f flag followed by the path to the file:
```
./calculator -f ./testdata/10k.txt
```

This will read the expression from the specified file and output the result.

## Manual Input

If you prefer to input the mathematical expression manually via stdin, use the -i flag:
```
./calculator -i
```

After running the command, you can type your expression directly into the console.

## In-place Evaluation

To use the EvalParse function for in-place evaluation, which may improve performance for certain expressions, include the -eval flag:
```
./calculator -f ./testdata/10k.txt -eval
```

## Profiling

Enable heap profiling to analyze memory usage and optimize performance by adding the -profile flag. This is particularly useful for understanding how the calculator handles large expressions:
```
./calculator -f ./testdata/10m.txt -profile
```

The profiler output files can be found in the working directory, named heap_profile_post_parse.prof or heap_profile_post_eval.prof, depending on whether the -eval flag was also used.

## Combining Flags

Flags can be combined for more specific use cases. For example, to manually input an expression and enable in-place evaluation with profiling:
```
./calculator -i -eval -profile
```

This setup allows for maximum flexibility in testing and optimizing the calculator for different scenarios.

# Solving Strategy

## Overview

Creating a calculator that interprets and evaluates mathematical expressions involves several steps: reading input, parsing the input into a structure that can be evaluated, and performing the evaluation. There are two main strategies for parsing and evaluation: constructing an Abstract Syntax Tree (AST) and using Reverse Polish Notation (RPN).

## Data Types

Expressions, numbers, signs, and operations can be represented using interfaces and polymorphism, allowing for a flexible design that can handle various types of expressions.

### Expression Interface

The Expression interface is a fundamental part of the design, encapsulating the concept of an arithmetic expression within the calculator. It defines the essential operations that any type of expression must support to be evaluated and represented as a string.
```
type Expr interface {
    // Eval computes the value of the expression.
    Eval() (float64, error)
    ...
}
```

This interface mandates that any arithmetic expression, whether a simple number, a unary operation (like negation), or a binary operation (such as addition, subtraction, multiplication, or division), must implement the Eval() method, this means, must be able to be evaluated to a number:

`Eval()` computes the numerical result of the expression. For operations, this involves recursively evaluating operands and applying then the corresponding operation.

### Numeric and Operation Types

Implementing the Expr interface, the calculator defines specific types for numbers and operations.

#### Num

`type num float64`

Num represents a numeric constant in an expression. It's a simple wrapper around a floating-point number, allowing it to satisfy the Expr interface. This type directly supports evaluation (`Eval`) by returning its own value.

### Unary

```
type unary struct {
op rune // '+' or '-'
x Expr
}
```

Unary encapsulates unary operations, which are operations with a single operand. It includes the operator (`op`) and the operand (`x`). The `Eval` method computes the result of applying the operator to the operand.

### Binary
```
type binary struct {
op rune // '+', '-', '*', '/'
x, y Expr
}
```

Binary represents binary operations, such as addition, subtraction, multiplication, and division, with two operands (`x` and `y`) and an operator (`op`). Its `Eval` method performs the operation on the operands' values.


## Parsing Algorithm
The parsing algorithm transforms text-based arithmetic expressions into a structured format for easy evaluation, often an Abstract Syntax Tree (AST) or an implicit tree. The key function in this process is parseBinary.

### parseBinary Function
This is a recursive function that plays a central role in the parsing algorithm. It parses binary operations and their operands based on operator precedence. It starts with a unary expression, evaluates the precedence of the following operator, and recursively parses the right operand if the operator has higher precedence. 
```
func parseBinary(lex *lexer, prio0 int) (Expr, error) {
    left, err := parseUnary(lex)
    if err != nil {
        return nil, err
    }
    for prio := priority(lex.token); prio >= prio0; prio-- {
        while priority(lex.token) == prio {
            op := lex.token
            lex.next()
            right, err := parseBinary(lex, prio+1)
            if err != nil {
                return nil, err
            }
            left = binary{op, left, right}
        }
    }
    return left, nil
}
```
This approach allows for accurate handling of complex expressions with varying operator precedence and nested parentheses. The AST or implicit tree structure it creates reflects the input expression's semantics accurately.

### Challenges with Recursion

Recursion, while elegant for parsing expressions, faces speed and memory efficiency issues. Deep recursion can slow down execution and lead to stack overflow errors. To mitigate these problems, converting recursive algorithms to iterative ones is a common solution. This approach reduces memory consumption and improves performance by eliminating the overhead associated with recursive function calls.

## Worked out Examples
[EXAMPLES.md](EXAMPLES.md)
 provides a detailed look at how the parsing algorithm works across different expressions, with different operators and parenthesis levels.

# Problem Setting

## Calculator Kata

Design a calculator capable of interpreting and evaluating expressions that include basic arithmetic operations (+, -, *, /) and parentheses, with support for both integers and floating point numbers. It must accurately handle errors such as division by zero, output results in a readable format, efficiently process long single expressions in terms of both memory and speed, and provide error reporting.

## Input Specification

The calculator accepts UTF-8 encoded input defined by an EBNF grammar, supporting expressions with nested parentheses, various operators with defined precedence, and numbers that can be either integers or floating points.

## Implementation Guidelines

- Limit dependencies to the standard library, excluding test and benchmark utilities.
- Adhere to best practices in code structure and idiomatic Go.
- Include comprehensive unit tests.

## Test Cases

Functional and benchmark tests are available in the `testdata` directory, with file names indicating approximate sizes. The expected results for the provided test files range from simple numbers to complex expressions, demonstrating the calculator's capability to handle inputs of varying complexity and size.


# Performance Testing Documentation

We use three distinct types of performance tests: Normal Tests, Benchmark Tests, and Profiling. We try to evaluate the efficiency of `Parse`, `Eval`, and `EvalParse` functions using a suite of test files of varying sizes named according the number of symbols to be parsed or the number of bytes to be parsed: 1k, 10k, 100k, 1M, and 10M.

Using different testing approaches allow us to understand performance testing from different perspectives. Each one is suitable for a different use case.

1. **Normal Tests:** Conducted within `main_test.go`, these tests manually time the `Parse` and `Eval` functions, alongside the `EvalParse` function, for immediate evaluation during parsing. We can see the results calling go test with the verbose flag 
`$ go test -v`.

    ```go
    func TestParseAndEvalPerformance(t *testing.T) {
        path := "./testdata/"
        for _, file := range benchmarkFiles {
            t.Run(fmt.Sprintf("Test data %s", strings.Split(file, ".")[0]), func(t *testing.T) {
                file = path + file
                f, err := os.Open(file)
                // ...

                startParse := time.Now()
                expr, err := Parse(f)
                durationParse := time.Since(startParse)

                // ...

                startEval := time.Now()
                _, err = expr.Eval()
                durationEval := time.Since(startEval)

                // ...

                t.Logf("Parse time for %s: %v", file, durationParse)
                t.Logf("Eval time for %s: %v", file, durationEval)
            })
        }
    }
    ```

2. **Benchmark Tests:** Implemented in `evalparse_bench.go`, these tests use Go's testing.B to automatically run performance tests across our dataset, allowing for standardized benchmarking across different sizes of arithmetic expressions.

    ```go
    func BenchmarkParseAndEval_1k(b *testing.B) {
        benchmarkParseAndEval("./testdata/1k.txt", b)
    }
    func benchmarkParseAndEval(fileName string, b *testing.B) {
        fileContent, err := os.ReadFile(fileName)
        // ...

        for i := 0; i < b.N; i++ {
            reader := bytes.NewReader(fileContent)
            expr, _ := Parse(reader) // Ignore errors while benchmarking
            expr.Eval()
        }
    }
    ```

3. **Profiling:** Profiling in `main.go` captures CPU and heap memory usage before and after parsing and evaluation. This is crucial for identifying performance hotspots and optimizing memory allocation. 

We can write the proff files calling the application with the `-profile` flag and then using the go tool for profiling to read the insights:

`$ go tool pprof heap_profile_post_eval.prof` 
followed for instance by the command `> top`.

    ```go
    if *profile {
        if *evalFlag {
            fileName = "cpu_profile.prof"
        }
        // ...
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()

        // Parse and evaluate expression here

        if *profile {
            fileName = "heap_profile_post_eval.prof"
            pprof.WriteHeapProfile(f)
        }
    }
    ```

Through a combination of detailed testing approaches, we aim to comprehensively understand and optimize the performance of our system across a wide range of scenarios and data sizes.

# Performance Test Results Analysis

## Outputs
```

$ go test -v
=== RUN   TestParseAndEvalPerformance
    main_test.go:45: Parse time for ./testdata/1k.txt: 46.003µs
    main_test.go:46: Eval time for ./testdata/1k.txt: 13.098µs

    main_test.go:45: Parse time for ./testdata/10k.txt: 216.206µs
    main_test.go:46: Eval time for ./testdata/10k.txt: 120.768µs

    main_test.go:45: Parse time for ./testdata/100k.txt: 1.992199ms
    main_test.go:46: Eval time for ./testdata/100k.txt: 1.362884ms

    main_test.go:45: Parse time for ./testdata/1m.txt: 23.125631ms
    main_test.go:46: Eval time for ./testdata/1m.txt: 13.578191ms

    main_test.go:45: Parse time for ./testdata/10m.txt: 207.793154ms
    main_test.go:46: Eval time for ./testdata/10m.txt: 208.643661ms

--- PASS: TestParseAndEvalPerformance (0.46s)
    --- PASS: TestParseAndEvalPerformance/Test_data_10m (0.42s)

=== RUN   TestEvalParseAndEvalPerformance
    main_test.go:83: EvalParse time for ./testdata/1k.txt: 42.529µs
    main_test.go:84: Eval time for ./testdata/1k.txt: 36ns

    main_test.go:83: EvalParse time for ./testdata/10k.txt: 251.635µs
    main_test.go:84: Eval time for ./testdata/10k.txt: 32ns

    main_test.go:83: EvalParse time for ./testdata/100k.txt: 2.147516ms
    main_test.go:84: Eval time for ./testdata/100k.txt: 29ns

    main_test.go:83: EvalParse time for ./testdata/1m.txt: 21.218545ms
    main_test.go:84: Eval time for ./testdata/1m.txt: 40ns

    main_test.go:83: EvalParse time for ./testdata/10m.txt: 261.075206ms
    main_test.go:84: Eval time for ./testdata/10m.txt: 41ns
--- PASS: TestEvalParseAndEvalPerformance (0.29s)
    --- PASS: TestEvalParseAndEvalPerformance/Test_data_10m (0.26s)


$ go test -benchmem -run=^$ -bench ^(BenchmarkParseAndEval_1k|BenchmarkParseAndEval_10k|BenchmarkParseAndEval_100k|BenchmarkParseAndEval_1m|BenchmarkParseAndEval_10m|BenchmarkEvalParseAndEval_1k|BenchmarkEvalParseAndEval_10k|BenchmarkEvalParseAndEval_100k|BenchmarkEvalParseAndEval_1m|BenchmarkEvalParseAndEval_10m)$ github.com/jerberlin/calcast

goos: linux
goarch: amd64
pkg: github.com/jerberlin/calcast
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkParseAndEval_1k-16          	   60585	     20543 ns/op	   10368 B/op	     425 allocs/op
BenchmarkParseAndEval_10k-16         	    5712	    204693 ns/op	   90818 B/op	    4196 allocs/op
BenchmarkParseAndEval_100k-16        	     520	   2212963 ns/op	  885472 B/op	   41435 allocs/op
BenchmarkParseAndEval_1m-16          	      42	  29312575 ns/op	 9099348 B/op	  425216 allocs/op
BenchmarkParseAndEval_10m-16         	       4	307 308870 ns/op   93 503056 B/op	4 251669 allocs/op
BenchmarkEvalParseAndEval_1k-16      	   50121	     23510 ns/op	   13224 B/op	     782 allocs/op
BenchmarkEvalParseAndEval_10k-16     	    4821	    231631 ns/op	  117275 B/op	    7503 allocs/op
BenchmarkEvalParseAndEval_100k-16    	     535	   2272704 ns/op	 1151122 B/op	   74642 allocs/op
BenchmarkEvalParseAndEval_1m-16      	      51	  24285284 ns/op	11813699 B/op	  765096 allocs/op
BenchmarkEvalParseAndEval_10m-16     	       5	223 377430 ns/op  120 132124 B/op	7 650251 allocs/op
PASS
ok  	github.com/jerberlin/calcast	17.066s

$ go tool pprof heap_profile_post_eval.prof
File: calcast
Type: inuse_space
Time: Mar 11, 2024 at 9:46pm (CET)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 39.22MB, 100% of 39.22MB total
Showing top 10 nodes out of 11
      flat  flat%   sum%        cum   cum%
   29.50MB 75.21% 75.21%    37.50MB 95.61%  main.parseBinary
    6.50MB 16.57% 91.79%     6.50MB 16.57%  text/scanner.(*Scanner).TokenText
    1.72MB  4.39% 96.18%     1.72MB  4.39%  runtime/pprof.StartCPUProfile
    1.50MB  3.82%   100%    10.50MB 26.77%  main.parsePrimary
         0     0%   100%     6.50MB 16.57%  main.(*lexer).text (inline)
         0     0%   100%    37.50MB 95.61%  main.Parse
         0     0%   100%    39.22MB   100%  main.main
         0     0%   100%    37.50MB 95.61%  main.parseExpr (inline)
         0     0%   100%    37.50MB 95.61%  main.parseInput
         0     0%   100%    10.50MB 26.77%  main.parseUnary

$ go tool pprof heap_profile_post_eval_evalFlag.prof
File: calcast
Type: inuse_space
Time: Mar 11, 2024 at 9:46pm (CET)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1.72MB, 100% of 1.72MB total
      flat  flat%   sum%        cum   cum%
    1.72MB   100%   100%     1.72MB   100%  runtime/pprof.StartCPUProfile
         0     0%   100%     1.72MB   100%  main.main
         0     0%   100%     1.72MB   100%  runtime.main

$ go tool pprof cpu_profile.prof
File: calcast
Type: cpu
Time: Mar 11, 2024 at 9:46pm (CET)
Duration: 603.41ms, Total samples = 480ms (79.55%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 300ms, 62.50% of 480ms total
Showing top 10 nodes out of 62
      flat  flat%   sum%        cum   cum%
      60ms 12.50% 12.50%       60ms 12.50%  runtime.memmove
      50ms 10.42% 22.92%      200ms 41.67%  main.binary.Eval
      30ms  6.25% 29.17%       60ms 12.50%  runtime.adjustframe
      30ms  6.25% 35.42%       70ms 14.58%  runtime.scanobject
      30ms  6.25% 41.67%      100ms 20.83%  text/scanner.(*Scanner).Scan
      20ms  4.17% 45.83%       20ms  4.17%  main.binary.Len
      20ms  4.17% 50.00%       20ms  4.17%  runtime.(*mspan).base (inline)
      20ms  4.17% 54.17%       90ms 18.75%  runtime.gcDrain
      20ms  4.17% 58.33%       20ms  4.17%  runtime.memclrNoHeapPointers
      20ms  4.17% 62.50%       20ms  4.17%  runtime/internal/syscall.Syscall6

$ go tool pprof cpu_profile_evalFlag.prof
File: calcast
Type: cpu
Time: Mar 11, 2024 at 9:46pm (CET)
Duration: 403.94ms, Total samples = 180ms (44.56%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 160ms, 88.89% of 180ms total
Showing top 10 nodes out of 47
      flat  flat%   sum%        cum   cum%
      50ms 27.78% 27.78%       80ms 44.44%  runtime.mallocgc
      20ms 11.11% 38.89%       20ms 11.11%  main.binary.Eval
      20ms 11.11% 50.00%       50ms 27.78%  text/scanner.(*Scanner).Scan
      10ms  5.56% 55.56%      180ms   100%  main.evalparseBinary
      10ms  5.56% 61.11%       10ms  5.56%  runtime.(*mspan).writeHeapBitsSmall
      10ms  5.56% 66.67%       10ms  5.56%  runtime.exitsyscallfast
      10ms  5.56% 72.22%       10ms  5.56%  runtime.getMCache (inline)
      10ms  5.56% 77.78%       10ms  5.56%  runtime.memmove
      10ms  5.56% 83.33%       10ms  5.56%  runtime.scanobject
      10ms  5.56% 88.89%       10ms  5.56%  runtime/internal/syscall.Syscall6
```
## Normal Tests

The Normal Tests, focusing on the largest data file (10m), show significant execution times with the `Parse` function taking `216.100776ms` and the `Eval` function taking `212.288621ms`. This indicates a nearly equal distribution of time between parsing and evaluation stages.

```plaintext
Parse time for ./testdata/10m.txt: 216.100776ms
Eval time for ./testdata/10m.txt: 212.288621ms
```

When testing the EvalParse function, a notable performance improvement is observed, with the entire operation taking 236.652888ms, demonstrating that integrating parsing and evaluation can slightly increase overall execution time due to the complexity of simultaneous operations, but drastically reduces the evaluation time to negligible levels (45 ns).

EvalParse time for ./testdata/10m.txt: 236.652888ms
Eval time for ./testdata/10m.txt: 45ns

## Benchmark Tests
Benchmark Tests provide a different perspective, emphasizing the performance under repeated execution. For the 10m data file, ParseAndEval required 307308870 ns/op (307.31s), with a considerable amount of allocations (4,251,669 allocs). Conversely, EvalParseAndEval showed a slightly lower time of 223377430 ns/op (223.38s) with a higher number of allocations (7,650,251 allocs). This reveals that while EvalParse can be faster, it demands more memory resources, as indicated by the increased allocations.

```plaintext
BenchmarkParseAndEval_10m-16          307308870 ns/op    93503056 B/op    4251669 allocs/op
BenchmarkEvalParseAndEval_10m-16      223377430 ns/op    120132124 B/op   7650251 allocs/op
```

## Profiling
### CPU Profiling
CPU profiling for the Parse operation reveals a significant portion of time spent on memory-related operations (memmove and memclrNoHeapPointers) and garbage collection (gcDrain). The binary.Eval function also consumes a notable amount of CPU time.

```
60ms 12.50%  runtime.memmove
50ms 10.42%  main.binary.Eval
```

In contrast, CPU profiling with the evalFlag highlights mallocgc (memory allocation) as the major consumer of CPU resources, along with Scan operations. This suggests that the integrated EvalParse operation intensifies memory allocation pressures

```
50ms 27.78%  runtime.mallocgc
20ms 11.11%  text/scanner.(*Scanner).Scan
```

## Heap Memory Profiling
Heap memory profiling after the Parse operation (without evalFlag) showed the largest memory allocations within parseBinary and text/scanner.(*Scanner).TokenText, accounting for over 90% of the total memory used.
```
29.50MB 75.21%  main.parseBinary
6.50MB  16.57%  text/scanner.(*Scanner).TokenText
```

# Additional Insights
This section delves into noteworthy aspects of the code.

## Stream Readers and Lazy Loading
It is often necessary to use lazy loading to efficiently process data from large files. This means that we don't load all the file in memory, but we just read a small part of the file at a time, when we need it. 

Our approach involves using a reader to feed a scanner, significantly optimizing memory usage and performance. The text/scanner package in Go provides a powerful tool for this purpose, allowing us to read data incrementally as needed, while using the in-built capabilities of the package to read tokens (i.e. symbols like '+' or strings of digits) and ignoring blank characters.

The lexer is a critical component of our parsing system. It uses its scanner to tokenize the input stream. The scanner captures each token using the next() method and retrieves its textual representation with text().

```
type lexer struct {
    scan  scanner.Scanner
    token rune // current token, used as lookahead
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func Parse(r io.Reader) (Expr, error) {
    lex := new(lexer)
    lex.scan.Init(r)
}
```

### Benchmark Testing and Memory Management

For benchmark tests, on the contrary, we adopt a strategy to minimize IO operation impact on performance measures. Data is preloaded into memory, eliminating IO latency from affecting our benchmarks.

```
fileContent, err := os.ReadFile(fileName)

for i := 0; i < b.N; i++ {
    reader := bytes.NewReader(fileContent)

    expr, _ := Parse(reader)
    expr.Eval()
}
```

## The message printer for pretty printing results

To enhance the readability of numerical results, particularly for large numbers, our calculator uses the `golang.org/x/text/message` package for pretty printing. This package provides internationalization features, allowing numbers to be formatted according to different local norms. We use it just to insert commas as thousands separators in English format.

```
func printResult(exp Expr, res float64) {
    // Create a new printer for English locale to use thousands comma separator
    p := message.NewPrinter(language.English)

    // Conditionally format the output based on the expression length
    if exp.Len() <= 1000 {
        p.Printf("Eval(%v) = %.2f\n", exp, res)
    } else {
        p.Printf("Eval() = %.2f\n", res)
    }
}
```

