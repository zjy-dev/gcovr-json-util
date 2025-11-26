# gcovr-json-util

A utility tool for processing and analyzing gcovr JSON coverage reports.

## Features

- **Coverage Diff**: Compare two gcovr JSON reports to identify coverage increases
- **Uncovered Lines Reporting**: Identify which lines, functions, and files lack coverage
- **Coverage Statistics**: Calculate overall coverage percentage for filtered functions
- **Filtering Support**: Filter coverage tracking by specific files and functions using a YAML config
- Reports which functions have improved coverage
- Shows old and new coverage percentages
- Displays newly covered line numbers
- Uses demangled function names for readability
- Can be used both as a CLI tool and as a Go library

## Installation

```bash
go install github.com/zjy-dev/gcovr-json-util@latest
```

Or install a specific version:

```bash
# Install v2.0.0 (with filtering support)
go install github.com/zjy-dev/gcovr-json-util@v2.0.0

# Install v1.0.0 (basic functionality)
go install github.com/zjy-dev/gcovr-json-util@v1.0.0
```

### Updating to v2.0.0 in Your Go Project

If you're using this as a library in your Go project:

**Option 1: Update to latest version**

```bash
go get -u github.com/zjy-dev/gcovr-json-util@latest
go mod tidy
```

**Option 2: Update to specific v2.0.0**

```bash
go get github.com/zjy-dev/gcovr-json-util/v2@v2.0.0
go mod tidy
```

**Option 3: Edit go.mod directly**

```go
require (
    github.com/zjy-dev/gcovr-json-util/v2 v2.0.0
)
```

Then run:

```bash
go mod tidy
```

**Check your current version:**

```bash
go list -m github.com/zjy-dev/gcovr-json-util/v2
```

**Note:** v2.0.0+ uses the `/v2` suffix in the module path, as per Go modules versioning convention.

### Building from Source

Or build from source:

```bash
git clone https://github.com/zjy-dev/gcovr-json-util
cd gcovr-json-util
make build
# Or without make:
go build -o gcovr-util .
```

### Building with Version Information

To build with version information embedded:

```bash
make build
# This automatically embeds version, git commit, and build date

# Or manually:
go build -ldflags "-X main.Version=v1.0.0 -X main.GitCommit=$(git rev-parse --short HEAD) -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o gcovr-util .
```

Check version:

```bash
./gcovr-util --version
```

## Usage

### CLI Tool

#### Coverage Diff Command

Compare two gcovr JSON reports:

```bash
./gcovr-util diff --base base_coverage.json --new new_coverage.json
```

**Options:**

- `--base, -b`: Base gcovr JSON report file (required)
- `--new, -n`: New gcovr JSON report file (required)
- `--filter, -f`: Filter config file (YAML) to specify target files and functions (optional)

#### Uncovered Lines Command

Report which lines are not covered in a gcovr JSON report:

```bash
./gcovr-util uncovered <gcovr-file.json>
```

**Options:**

- `--filter, -f`: Filter config file (YAML) to specify target files and functions (optional)

**Example:**

```bash
# Show all uncovered lines
./gcovr-util uncovered coverage.json

# Show uncovered lines only for filtered functions
./gcovr-util uncovered --filter filter.yaml coverage.json
```

**Example Output:**

```
Uncovered Lines Report
======================

Found 2 function(s) with uncovered lines (4 total uncovered lines):

1. File: demo.cc
   Function: g()
   Coverage: 0/3 lines (0.0%)
   Uncovered Lines (3): [9 10 11]

2. File: demo.cc
   Function: main
   Coverage: 4/5 lines (80.0%)
   Uncovered Lines (1): [17]
```

#### Coverage Command

Calculate overall coverage statistics from a gcovr JSON report:

```bash
./gcovr-util coverage <gcovr-file.json>
```

**Options:**

- `--filter, -f`: Filter config file (YAML) to specify target files and functions (optional)

**Example:**

```bash
# Show coverage for all functions
./gcovr-util coverage coverage.json

# Show coverage only for filtered functions
./gcovr-util coverage --filter filter.yaml coverage.json
```

**Example Output:**

```
Coverage Report
===============

Overall Coverage: 7/11 lines (63.6%)

Functions (3):

  File: demo.cc
    1. f(): 3/3 lines (100.0%)
    2. g(): 0/3 lines (0.0%)
    3. main: 4/5 lines (80.0%)
```

#### Using Filter Configuration

You can use a YAML configuration file to filter which files and functions to track:

```bash
./gcovr-util diff --base base.json --new new.json --filter filter.yaml
```

**Filter Config Format** (`filter.yaml`):

```yaml
compiler:
  path: "/usr/bin/gcc"
  gcovr_exec_path: "/path/to/build"

targets:
  - file: "demo.cc"
    functions:
      - "f"
      - "g"
  - file: "another_file.cpp"
    functions:
      - "myFunction"
      - "anotherFunction"
```

This will only report coverage increases for the specified functions in the specified files. All other files and functions will be ignored.

**Note**:

- File paths can be specified as relative paths, absolute paths, or just filenames
- Function names should match the demangled names (e.g., "f" instead of "\_Z1fv")
- The `*.json` files and filter config file paths support both relative and absolute paths

#### Example Output

```
Coverage Increase Report
=========================

Found 2 function(s) with increased coverage:

1. File: demo.cc
   Function: g()
   Old Coverage: 0/3 lines (0.0%)
   New Coverage: 3/3 lines (100.0%)
   Lines Increased: 3
   Newly Covered Line Numbers: [9 10 11]

2. File: demo.cc
   Function: main
   Old Coverage: 4/5 lines (80.0%)
   New Coverage: 5/5 lines (100.0%)
   Lines Increased: 1
   Newly Covered Line Numbers: [17]
```

### Go Library

You can also use this tool as a Go library in your projects:

**Example 1: Coverage Diff**

```go
import "github.com/zjy-dev/gcovr-json-util/v2/pkg/gcovr"

// Parse coverage reports
baseReport, err := gcovr.ParseReport("base.json")
if err != nil {
    log.Fatal(err)
}

newReport, err := gcovr.ParseReport("new.json")
if err != nil {
    log.Fatal(err)
}

// Optional: Apply filtering
filterConfig, err := gcovr.ParseFilterConfig("filter.yaml")
if err != nil {
    log.Fatal(err)
}
baseReport = gcovr.ApplyFilter(baseReport, filterConfig)
newReport = gcovr.ApplyFilter(newReport, filterConfig)

// Compute coverage increase
report, err := gcovr.ComputeCoverageIncrease(baseReport, newReport)
if err != nil {
    log.Fatal(err)
}

// Format and display results
output := gcovr.FormatReport(report)
fmt.Print(output)
```

**Example 2: Find Uncovered Lines**

```go
import "github.com/zjy-dev/gcovr-json-util/v2/pkg/gcovr"

// Parse coverage report
report, err := gcovr.ParseReport("coverage.json")
if err != nil {
    log.Fatal(err)
}

// Optional: Apply filtering
filterConfig, err := gcovr.ParseFilterConfig("filter.yaml")
if err != nil {
    log.Fatal(err)
}
report = gcovr.ApplyFilter(report, filterConfig)

// Find uncovered lines
uncoveredReport, err := gcovr.FindUncoveredLines(report)
if err != nil {
    log.Fatal(err)
}

// Format and display results
output := gcovr.FormatUncoveredReport(uncoveredReport)
fmt.Print(output)
```

**Example 3: Calculate Coverage Statistics**

```go
import "github.com/zjy-dev/gcovr-json-util/v2/pkg/gcovr"

// Parse coverage report
report, err := gcovr.ParseReport("coverage.json")
if err != nil {
    log.Fatal(err)
}

// Optional: Apply filtering to calculate coverage for specific functions
filterConfig, err := gcovr.ParseFilterConfig("filter.yaml")
if err != nil {
    log.Fatal(err)
}
report = gcovr.ApplyFilter(report, filterConfig)

// Calculate coverage statistics
coverageReport, err := gcovr.CalculateCoverage(report)
if err != nil {
    log.Fatal(err)
}

// Access coverage data programmatically
fmt.Printf("Overall Coverage: %.1f%%\n", coverageReport.CoveragePercentage)
fmt.Printf("Total Lines: %d, Covered: %d\n", coverageReport.TotalLines, coverageReport.TotalCoveredLines)

// Or format and display results
output := gcovr.FormatCoverageReport(coverageReport)
fmt.Print(output)
```

**Example 4: Working with Grouped Data**

The `FindUncoveredLines()` function returns data already grouped by file:

```go
uncoveredReport, err := gcovr.FindUncoveredLines(report)
if err != nil {
    log.Fatal(err)
}

// Iterate through files
for _, file := range uncoveredReport.Files {
    fmt.Printf("File: %s\n", file.FilePath)

    // Iterate through uncovered functions in this file
    for _, fn := range file.UncoveredFunctions {
        fmt.Printf("  Function: %s\n", fn.DemangledName)
        fmt.Printf("  Coverage: %d/%d lines\n", fn.CoveredLines, fn.TotalLines)
        fmt.Printf("  Uncovered Lines: %v\n", fn.UncoveredLineNumbers)
    }
}
```

### API Data Structures

**UncoveredReport** - Grouped by file structure:

```go
type UncoveredReport struct {
    Files []FileUncovered  // Files containing uncovered lines
}

type FileUncovered struct {
    FilePath           string               // Path to the source file
    UncoveredFunctions []FunctionUncovered  // Uncovered functions in this file
}

type FunctionUncovered struct {
    FunctionName         string  // Mangled function name
    DemangledName        string  // Human-readable function name
    UncoveredLineNumbers []int   // Line numbers without coverage
    TotalLines           int     // Total lines in function
    CoveredLines         int     // Number of covered lines
}
```

**CoverageReport** - Coverage statistics structure:

```go
type CoverageReport struct {
    Functions          []FunctionCoverage  // Per-function coverage data
    TotalLines         int                 // Total lines across all functions
    TotalCoveredLines  int                 // Total covered lines
    CoveragePercentage float64             // Overall coverage percentage
}

type FunctionCoverage struct {
    FilePath      string  // Path to the source file
    FunctionName  string  // Mangled function name
    DemangledName string  // Human-readable function name
    TotalLines    int     // Total lines in function
    CoveredLines  int     // Number of covered lines
}
```

## Version History & Migration

### v2.2.0 - November 26, 2025

**New Features:**

- ✨ Coverage statistics feature
- 📊 Calculate overall coverage percentage for filtered functions
- 📋 New `coverage` command (alias: `cov`)
- 🔧 New library functions: `CalculateCoverage()`, `FormatCoverageReport()`
- 📦 New data structures: `CoverageReport`, `FunctionCoverage`

**Migration from v2.1.0:**

```bash
go get github.com/zjy-dev/gcovr-json-util/v2@v2.2.0
go mod tidy
```

### v2.1.0 - November 19, 2025

**New Features:**

- ✨ Uncovered lines reporting feature
- 🎯 Identify files, functions, and specific lines lacking coverage
- 📋 New `uncovered` command (alias: `un`)
- 🔧 New library functions: `FindUncoveredLines()`, `FormatUncoveredReport()`
- 📊 Data returned grouped by file for easier programmatic access

**API Changes:**

- New types: `FileUncovered`, updated `UncoveredReport` structure
- `FindUncoveredLines()` now returns data grouped by file (breaking change for API users)

**Migration from v2.0.0:**

```bash
go get github.com/zjy-dev/gcovr-json-util/v2@v2.1.0
go mod tidy
```

If you were using the uncovered lines API in v2.0.0, update your code:

```go
// Old (v2.0.0 - if used)
for _, fn := range uncoveredReport.UncoveredFunctions {
    fmt.Printf("File: %s, Function: %s\n", fn.File, fn.FunctionName)
}

// New (v2.1.0)
for _, file := range uncoveredReport.Files {
    fmt.Printf("File: %s\n", file.FilePath)
    for _, fn := range file.UncoveredFunctions {
        fmt.Printf("  Function: %s\n", fn.FunctionName)
    }
}
```

### v2.0.0 - November 12, 2025

**New Features:**

- ✨ Filter configuration support via YAML files
- 🎯 Selective coverage tracking by file and function
- 📋 New `--filter` flag
- 🔧 New library functions: `ParseFilterConfig()`, `ApplyFilter()`

**Dependencies:**

- Added: `gopkg.in/yaml.v3`

**Migration from v1.0.0:**

```bash
# Update dependency (note the /v2 suffix for v2.0.0+)
go get github.com/zjy-dev/gcovr-json-util/v2@v2.0.0
go mod tidy
```

**Important:** For v2.0.0+, the import path includes `/v2` suffix:

```go
// Old (v1.0.0)
import "github.com/zjy-dev/gcovr-json-util/pkg/gcovr"

// New (v2.0.0+)
import "github.com/zjy-dev/gcovr-json-util/v2/pkg/gcovr"
```

**Code changes (optional, for new filtering feature):**

```go
// v1.0.0 code (still works in v2.0.0)
baseReport, _ := gcovr.ParseReport("base.json")
newReport, _ := gcovr.ParseReport("new.json")
report, _ := gcovr.ComputeCoverageIncrease(baseReport, newReport)

// v2.0.0 new feature (optional)
filterConfig, _ := gcovr.ParseFilterConfig("filter.yaml")
baseReport = gcovr.ApplyFilter(baseReport, filterConfig)
newReport = gcovr.ApplyFilter(newReport, filterConfig)
report, _ := gcovr.ComputeCoverageIncrease(baseReport, newReport)
```

**Backward Compatibility:** ✅ Fully backward compatible - all v1.0.0 code works in v2.0.0 without changes.

### v1.0.0 - November 11, 2025

**Features:**

- Basic coverage diff functionality
- Version information (`--version`)
- CLI and library usage

## Project Structure

```
.
├── main.go              # CLI entry point
├── version.go           # Version information
├── cmd/                 # CLI commands
│   ├── root.go         # Root command
│   ├── diff.go         # Diff command implementation
│   ├── uncovered.go    # Uncovered lines command
│   └── coverage.go     # Coverage statistics command
├── pkg/
│   └── gcovr/          # Public library package
│       ├── types.go    # Data structures
│       ├── parser.go   # JSON parsing
│       ├── diff.go     # Coverage diff logic
│       ├── filter.go   # Filter configuration
│       ├── uncovered.go # Uncovered lines logic
│       └── coverage.go # Coverage statistics logic
├── test_data/          # Sample test files
│   ├── f.json
│   ├── g.json
│   ├── m.json
│   ├── filter.yaml              # Example filter config
│   └── filter-f-only.yaml       # Another filter example
├── Makefile            # Build automation
├── CHANGELOG.md        # Version history
└── README.md           # This file
```

## How It Works

### Coverage Diff

1. **Parse**: Reads and parses two gcovr JSON reports
2. **Compare**: Compares line-by-line coverage for each function
3. **Identify**: Identifies lines that were uncovered in base but are covered in new report
4. **Report**: Generates a detailed report with:
   - Function names (demangled for C++)
   - Old coverage percentage and line count
   - New coverage percentage and line count
   - Number of newly covered lines
   - Specific line numbers that gained coverage

### Uncovered Lines

1. **Parse**: Reads and parses a gcovr JSON report
2. **Analyze**: Identifies all lines with zero coverage count
3. **Group**: Groups uncovered lines by file and function
4. **Report**: Generates a detailed report with:
   - Function names (demangled for C++)
   - Coverage percentage and line count
   - Specific line numbers that lack coverage

### Coverage Statistics

1. **Parse**: Reads and parses a gcovr JSON report
2. **Calculate**: Computes coverage statistics for each function
3. **Aggregate**: Calculates overall coverage percentage
4. **Report**: Generates a detailed report with:
   - Overall coverage percentage
   - Per-function coverage statistics
   - Total and covered line counts

## Use Cases

- **CI/CD Pipelines**: Automatically track coverage improvements in pull requests
- **Code Review**: Verify that new tests actually improve coverage
- **Quality Gates**: Ensure new code increases overall test coverage
- **Test Analysis**: Identify which functions benefited from new test cases

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
