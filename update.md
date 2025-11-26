# Plan: Uncovered Lines Reporter

This document outlines the plan to add a new feature that identifies and reports uncovered lines from a `gcovr` JSON report.

## 1. Feature Overview

The goal is to create a tool that can parse a `gcovr` JSON file and report which lines, within which functions and files, have a coverage count of zero.

This feature will be exposed in two ways:

1.  **Public API**: A new set of functions in the `pkg/gcovr` library.
2.  **CLI Command**: A new subcommand, `gcovr-util uncovered`, to display the report.

The new feature must be compatible with the existing `--filter` functionality.

## 2. API Changes (`pkg/gcovr/`)

### a. `pkg/gcovr/types.go`

Add new data structures to represent the uncovered lines report.

```go
// ... existing code ...

// FunctionUncovered represents the uncovered lines within a single function.
type FunctionUncovered struct {
	File                 string
	FunctionName         string // Mangled name
	DemangledName        string
	UncoveredLineNumbers []int
	TotalLines           int
	CoveredLines         int
}

// UncoveredReport represents a complete report of all uncovered functions and lines.
type UncoveredReport struct {
	UncoveredFunctions []FunctionUncovered
}
```

### b. `pkg/gcovr/uncovered.go` (New File)

Create a new file to house the logic for finding and formatting the uncovered lines report.

#### `FindUncoveredLines` function:

This function will be the core of the new API.

- **Signature**: `func FindUncoveredLines(report *GcovrReport) (*UncoveredReport, error)`
- **Logic**:
  1.  Initialize an empty `UncoveredReport`.
  2.  Create a map `map[string]map[string][]int` to store uncovered lines grouped by file and then by function name.
  3.  Iterate through each `file` in `report.Files`.
  4.  Iterate through each `line` in `file.Lines`.
  5.  If `line.Count == 0`, add the `line.LineNumber` to our map for that file and function.
  6.  After collecting all uncovered lines, iterate through the map to build the `FunctionUncovered` structs.
  7.  For each function, calculate `TotalLines` and `CoveredLines` by re-iterating through the original report's lines for that function.
  8.  Populate and return the `UncoveredReport`.

#### `FormatUncoveredReport` function:

This function will generate a human-readable string from the report.

- **Signature**: `func FormatUncoveredReport(report *UncoveredReport) string`
- **Logic**:
  1.  Check if `report.UncoveredFunctions` is empty. If so, return a "No uncovered lines found" message.
  2.  Iterate through `report.UncoveredFunctions`, grouping by file.
  3.  For each file, print the file name.
  4.  For each function within that file, print:
      - The demangled function name.
      - Coverage summary (e.g., "Coverage: 5/10 lines (50.0%)").
      - The list of uncovered line numbers.

## 3. CLI Changes (`cmd/`)

### a. `cmd/uncovered.go` (New File)

Create a new Cobra command for the `uncovered` feature.

- **Command**: `gcovr-util uncovered [flags] <gcovr-file>`
- **Alias**: `un`
- **Flags**:
  - It will reuse the existing global `--filter, -f` flag.
- **Argument**: It will accept exactly one argument: the path to the gcovr JSON file.
- **Logic (`runUncovered`)**:
  1.  Validate that one argument is provided.
  2.  Read the gcovr JSON file using `gcovr.ParseReport`.
  3.  **Filtering**: If the `--filter` flag is present, parse the filter config and apply it to the report using `gcovr.ApplyFilter`. This ensures we only check for uncovered lines in the targeted files/functions.
  4.  Call `gcovr.FindUncoveredLines` with the (potentially filtered) report.
  5.  Call `gcovr.FormatUncoveredReport` to get the formatted output string.
  6.  Print the result to standard output.

### b. `cmd/root.go`

- Add the new `uncoveredCmd` to the `rootCmd` in the `init()` function.

## 4. Documentation

### a. `README.md`

- Add the `uncovered` command to the "Usage" section with an example.
- Update the "Go Library" section to include an example of using `FindUncoveredLines`.
- Add "Uncovered Lines Reporting" to the "Features" list.

### b. `CHANGELOG.md`

- Create a new `[v2.1.0]` or `[v3.0.0]` section.
- Document the new `uncovered` command and API functions.

## 5. Testing Plan

1.  **Unit Tests**:
    - Create `pkg/gcovr/uncovered_test.go`.
    - Test `FindUncoveredLines` with a sample report to ensure it correctly identifies uncovered lines.
    - Test `FormatUncoveredReport` to ensure the output is formatted as expected.
2.  **Manual CLI Tests**:
    - Run `gcovr-util uncovered test_data/f.json` and verify the output.
    - Run `gcovr-util uncovered test_data/m.json` (which has full coverage for some functions) and verify the output.
    - Run `gcovr-util uncovered --filter gcc-x64-canary.yaml test_data/f.json` to ensure filtering works correctly with the new command.
    - Test with a file that has 100% coverage to see the "No uncovered lines" message.
