package gcovr

import (
	"fmt"
	"sort"
)

// CalculateCoverage analyzes a gcovr report and returns the overall coverage statistics
// for all functions in the report
func CalculateCoverage(report *GcovrReport) (*CoverageReport, error) {
	result := &CoverageReport{
		Functions: make([]FunctionCoverage, 0),
	}

	// Process each file
	for _, file := range report.Files {
		// Build a map of function metadata
		funcMetadata := make(map[string]string) // funcName -> demangledName
		for _, fn := range file.Functions {
			funcMetadata[fn.Name] = fn.DemangledName
		}

		// Group lines by function
		funcLines := make(map[string][]Line)
		for _, line := range file.Lines {
			funcLines[line.FunctionName] = append(funcLines[line.FunctionName], line)
		}

		// Calculate coverage for each function
		for funcName, lines := range funcLines {
			totalLines := len(lines)
			coveredLines := 0
			for _, line := range lines {
				if line.Count > 0 {
					coveredLines++
				}
			}

			demangledName := funcMetadata[funcName]
			if demangledName == "" {
				demangledName = funcName
			}

			result.Functions = append(result.Functions, FunctionCoverage{
				FilePath:      file.FilePath,
				FunctionName:  funcName,
				DemangledName: demangledName,
				TotalLines:    totalLines,
				CoveredLines:  coveredLines,
			})

			result.TotalLines += totalLines
			result.TotalCoveredLines += coveredLines
		}
	}

	// Sort functions by file path and then by function name for consistent output
	sort.Slice(result.Functions, func(i, j int) bool {
		if result.Functions[i].FilePath != result.Functions[j].FilePath {
			return result.Functions[i].FilePath < result.Functions[j].FilePath
		}
		return result.Functions[i].FunctionName < result.Functions[j].FunctionName
	})

	// Calculate overall coverage percentage
	if result.TotalLines > 0 {
		result.CoveragePercentage = float64(result.TotalCoveredLines) * 100.0 / float64(result.TotalLines)
	}

	return result, nil
}

// FormatCoverageReport formats the coverage report as a human-readable string
func FormatCoverageReport(report *CoverageReport) string {
	result := fmt.Sprintf("Coverage Report\n")
	result += fmt.Sprintf("===============\n\n")

	result += fmt.Sprintf("Overall Coverage: %d/%d lines (%.1f%%)\n\n",
		report.TotalCoveredLines, report.TotalLines, report.CoveragePercentage)

	if len(report.Functions) == 0 {
		result += "No functions found in the report.\n"
		return result
	}

	result += fmt.Sprintf("Functions (%d):\n", len(report.Functions))

	currentFile := ""
	funcIdx := 1
	for _, fn := range report.Functions {
		if fn.FilePath != currentFile {
			currentFile = fn.FilePath
			result += fmt.Sprintf("\n  File: %s\n", currentFile)
		}

		coveragePercent := 0.0
		if fn.TotalLines > 0 {
			coveragePercent = float64(fn.CoveredLines) * 100.0 / float64(fn.TotalLines)
		}

		result += fmt.Sprintf("    %d. %s: %d/%d lines (%.1f%%)\n",
			funcIdx, fn.DemangledName, fn.CoveredLines, fn.TotalLines, coveragePercent)
		funcIdx++
	}

	return result
}
