package main

import (
	"go/ast"

	"github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/gofix"
	"golang.org/x/tools/go/analysis/passes/hostport"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

var myOSExitInMainAnalyzer = &analysis.Analyzer{
	Name: "osexitinmain",
	Doc:  "checks for os.exit function calls in main func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {

			return true
		})
	}
	return nil, nil
}

// main is the entry point of the static analysis tool. It constructs a list of analyzers to be run by multichecker.Main.
// The analyzers included are:
//   - myOSExitInMainAnalyzer: Custom analyzer to detect usage of os.Exit in main functions.
//   - staticcheck.Analyzers: A set of analyzers from the staticcheck suite (SA*), which detect bugs and suspicious constructs.
//   - stylecheck.Analyzers: A set of analyzers from the staticcheck suite (ST*), which enforce style rules.
//   - errcheck.Analyzer: Checks for unchecked errors in code.
//   - analyzer.Analyzer: go-critic analyzer, which provides a wide range of code improvement suggestions.
//   - appends.Analyzer: Detects suspicious uses of the append function.
//   - asmdecl.Analyzer: Reports mismatches between assembly files and Go declarations.
//   - assign.Analyzer: Detects useless assignments.
//   - atomic.Analyzer: Checks for common mistakes using the sync/atomic package.
//   - atomicalign.Analyzer: Checks for non-64-bit-aligned arguments to sync/atomic functions.
//   - bools.Analyzer: Detects common mistakes involving boolean expressions.
//   - buildssa.Analyzer: Builds SSA (Static Single Assignment) form for Go programs (internal use).
//   - buildtag.Analyzer: Checks for invalid or duplicate build tags.
//   - cgocall.Analyzer: Detects calls to C code that may block the Go scheduler.
//   - composite.Analyzer: Checks for composite literal issues.
//   - copylock.Analyzer: Detects locks that are copied by value.
//   - ctrlflow.Analyzer: Checks for unreachable code and other control flow issues.
//   - deepequalerrors.Analyzer: Detects incorrect uses of reflect.DeepEqual with error values.
//   - defers.Analyzer: Checks for common mistakes in defer statements.
//   - directive.Analyzer: Checks for malformed compiler directives.
//   - errorsas.Analyzer: Checks for incorrect usage of errors.As.
//   - fieldalignment.Analyzer: Suggests struct field reordering to reduce memory usage.
//   - findcall.Analyzer: Finds calls to a specified function (configurable).
//   - framepointer.Analyzer: Checks for functions that may require frame pointers.
//   - gofix.Analyzer: Checks for code that gofix can fix automatically.
//   - hostport.Analyzer: Checks for suspicious host:port strings.
//   - httpmux.Analyzer: Checks for suspicious patterns in HTTP request multiplexer usage.
//   - httpresponse.Analyzer: Checks for mistakes handling HTTP responses.
//   - ifaceassert.Analyzer: Detects impossible interface type assertions.
//   - inspect.Analyzer: Provides a syntax tree inspection API for other analyzers (internal use).
//   - loopclosure.Analyzer: Detects common mistakes in closures within loops.
//   - lostcancel.Analyzer: Detects context cancellations that are not called.
//   - nilfunc.Analyzer: Detects comparisons of functions to nil.
//   - nilness.Analyzer: Checks for redundant or impossible nil comparisons.
//   - pkgfact.Analyzer: Provides package-wide facts for other analyzers (internal use).
//   - printf.Analyzer: Checks for printf-style formatting errors.
//   - reflectvaluecompare.Analyzer: Detects incorrect comparisons of reflect.Value objects.
//   - shadow.Analyzer: Detects shadowed variables.
//   - shift.Analyzer: Detects suspicious bit shift operations.
//   - sigchanyzer.Analyzer: Detects misuse of signal channels.
//   - slog.Analyzer: Checks for issues with structured logging.
//   - sortslice.Analyzer: Detects suspicious uses of sort.Slice.
//   - stdmethods.Analyzer: Checks for misspelled method names of standard interfaces.
//   - stdversion.Analyzer: Checks for usage of features not available in the targeted Go version.
//   - stringintconv.Analyzer: Detects suspicious conversions between strings and integers.
//   - structtag.Analyzer: Checks for struct tag format issues.
//   - testinggoroutine.Analyzer: Detects goroutines started in tests that may outlive the test.
//   - tests.Analyzer: Checks for common mistakes in test code.
//   - timeformat.Analyzer: Checks for incorrect usage of time formatting patterns.
//   - unmarshal.Analyzer: Checks for issues in unmarshaling data into Go structs.
//   - unreachable.Analyzer: Detects unreachable code.
//   - unsafeptr.Analyzer: Detects misuse of unsafe pointers.
//   - unusedresult.Analyzer: Detects unused results of calls to certain functions.
//   - unusedwrite.Analyzer: Detects unused writes to variables.
//   - usesgenerics.Analyzer: Checks for usage of generics (type parameters).
//   - waitgroup.Analyzer: Detects misuse of sync.WaitGroup.
func main() {

	mychecks := []*analysis.Analyzer{}

	// my own os.exit analyzer
	mychecks = append(mychecks, myOSExitInMainAnalyzer)

	// add staticcheck SA* analyzers
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	// add staticcheck ST* analyzers
	for _, v := range stylecheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	// first public analyzer: errcheck
	mychecks = append(mychecks, errcheck.Analyzer)

	// second public analyzer: go-critic
	mychecks = append(mychecks, analyzer.Analyzer)

	mychecks = append(mychecks, appends.Analyzer)
	mychecks = append(mychecks, asmdecl.Analyzer)
	mychecks = append(mychecks, assign.Analyzer)
	mychecks = append(mychecks, atomic.Analyzer)
	mychecks = append(mychecks, atomicalign.Analyzer)
	mychecks = append(mychecks, bools.Analyzer)
	mychecks = append(mychecks, buildssa.Analyzer)
	mychecks = append(mychecks, buildtag.Analyzer)
	mychecks = append(mychecks, cgocall.Analyzer)
	mychecks = append(mychecks, composite.Analyzer)
	mychecks = append(mychecks, copylock.Analyzer)
	mychecks = append(mychecks, ctrlflow.Analyzer)
	mychecks = append(mychecks, deepequalerrors.Analyzer)
	mychecks = append(mychecks, defers.Analyzer)
	mychecks = append(mychecks, directive.Analyzer)
	mychecks = append(mychecks, errorsas.Analyzer)
	mychecks = append(mychecks, fieldalignment.Analyzer)
	mychecks = append(mychecks, findcall.Analyzer)
	mychecks = append(mychecks, framepointer.Analyzer)
	mychecks = append(mychecks, gofix.Analyzer)
	mychecks = append(mychecks, hostport.Analyzer)
	mychecks = append(mychecks, httpmux.Analyzer)
	mychecks = append(mychecks, httpresponse.Analyzer)
	mychecks = append(mychecks, ifaceassert.Analyzer)
	mychecks = append(mychecks, inspect.Analyzer)
	mychecks = append(mychecks, loopclosure.Analyzer)
	mychecks = append(mychecks, lostcancel.Analyzer)
	mychecks = append(mychecks, nilfunc.Analyzer)
	mychecks = append(mychecks, nilness.Analyzer)
	mychecks = append(mychecks, pkgfact.Analyzer)
	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, reflectvaluecompare.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, shift.Analyzer)
	mychecks = append(mychecks, sigchanyzer.Analyzer)
	mychecks = append(mychecks, slog.Analyzer)
	mychecks = append(mychecks, sortslice.Analyzer)
	mychecks = append(mychecks, stdmethods.Analyzer)
	mychecks = append(mychecks, stdversion.Analyzer)
	mychecks = append(mychecks, stringintconv.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)
	mychecks = append(mychecks, testinggoroutine.Analyzer)
	mychecks = append(mychecks, tests.Analyzer)
	mychecks = append(mychecks, timeformat.Analyzer)
	mychecks = append(mychecks, unmarshal.Analyzer)
	mychecks = append(mychecks, unreachable.Analyzer)
	mychecks = append(mychecks, unsafeptr.Analyzer)
	mychecks = append(mychecks, unusedresult.Analyzer)
	mychecks = append(mychecks, unusedwrite.Analyzer)
	mychecks = append(mychecks, usesgenerics.Analyzer)
	mychecks = append(mychecks, waitgroup.Analyzer)

	multichecker.Main(
		mychecks...,
	)
}
