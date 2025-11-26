package gcovr

import (
	"strings"
	"testing"
)

func TestCalculateCoverage(t *testing.T) {
	tests := []struct {
		name                       string
		report                     *GcovrReport
		expectedFuncCount          int
		expectedTotalLines         int
		expectedCoveredLines       int
		expectedCoveragePercentage float64
	}{
		{
			name: "Single function fully covered",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 2, FunctionName: "foo", Count: 1},
							{LineNumber: 3, FunctionName: "foo", Count: 1},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			expectedFuncCount:          1,
			expectedTotalLines:         3,
			expectedCoveredLines:       3,
			expectedCoveragePercentage: 100.0,
		},
		{
			name: "Single function partially covered",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 2, FunctionName: "foo", Count: 0},
							{LineNumber: 3, FunctionName: "foo", Count: 1},
							{LineNumber: 4, FunctionName: "foo", Count: 0},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			expectedFuncCount:          1,
			expectedTotalLines:         4,
			expectedCoveredLines:       2,
			expectedCoveragePercentage: 50.0,
		},
		{
			name: "Multiple functions in one file",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 2, FunctionName: "foo", Count: 1},
							{LineNumber: 3, FunctionName: "bar", Count: 0},
							{LineNumber: 4, FunctionName: "bar", Count: 0},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
							{Name: "bar", DemangledName: "bar()"},
						},
					},
				},
			},
			expectedFuncCount:          2,
			expectedTotalLines:         4,
			expectedCoveredLines:       2,
			expectedCoveragePercentage: 50.0,
		},
		{
			name: "Multiple files",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "a.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "func1", Count: 1},
							{LineNumber: 2, FunctionName: "func1", Count: 1},
						},
						Functions: []Function{
							{Name: "func1", DemangledName: "func1()"},
						},
					},
					{
						FilePath: "b.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "func2", Count: 0},
							{LineNumber: 2, FunctionName: "func2", Count: 0},
						},
						Functions: []Function{
							{Name: "func2", DemangledName: "func2()"},
						},
					},
				},
			},
			expectedFuncCount:          2,
			expectedTotalLines:         4,
			expectedCoveredLines:       2,
			expectedCoveragePercentage: 50.0,
		},
		{
			name: "Empty report",
			report: &GcovrReport{
				Files: []File{},
			},
			expectedFuncCount:          0,
			expectedTotalLines:         0,
			expectedCoveredLines:       0,
			expectedCoveragePercentage: 0.0,
		},
		{
			name: "All lines uncovered",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 0},
							{LineNumber: 2, FunctionName: "foo", Count: 0},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			expectedFuncCount:          1,
			expectedTotalLines:         2,
			expectedCoveredLines:       0,
			expectedCoveragePercentage: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateCoverage(tt.report)
			if err != nil {
				t.Fatalf("CalculateCoverage() error = %v", err)
			}

			if len(result.Functions) != tt.expectedFuncCount {
				t.Errorf("Expected %d functions, got %d", tt.expectedFuncCount, len(result.Functions))
			}

			if result.TotalLines != tt.expectedTotalLines {
				t.Errorf("Expected TotalLines=%d, got %d", tt.expectedTotalLines, result.TotalLines)
			}

			if result.TotalCoveredLines != tt.expectedCoveredLines {
				t.Errorf("Expected TotalCoveredLines=%d, got %d", tt.expectedCoveredLines, result.TotalCoveredLines)
			}

			if result.CoveragePercentage != tt.expectedCoveragePercentage {
				t.Errorf("Expected CoveragePercentage=%.1f, got %.1f", tt.expectedCoveragePercentage, result.CoveragePercentage)
			}
		})
	}
}

func TestCalculateCoverage_FunctionDetails(t *testing.T) {
	report := &GcovrReport{
		Files: []File{
			{
				FilePath: "test.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 1},
					{LineNumber: 2, FunctionName: "foo", Count: 0},
					{LineNumber: 3, FunctionName: "foo", Count: 1},
					{LineNumber: 4, FunctionName: "bar", Count: 1},
					{LineNumber: 5, FunctionName: "bar", Count: 1},
				},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
					{Name: "bar", DemangledName: "bar()"},
				},
			},
		},
	}

	result, err := CalculateCoverage(report)
	if err != nil {
		t.Fatalf("CalculateCoverage() error = %v", err)
	}

	if len(result.Functions) != 2 {
		t.Fatalf("Expected 2 functions, got %d", len(result.Functions))
	}

	// Functions should be sorted by name
	barFunc := result.Functions[0]
	fooFunc := result.Functions[1]

	// Check bar function (fully covered)
	if barFunc.FunctionName != "bar" {
		t.Errorf("Expected first function to be 'bar', got '%s'", barFunc.FunctionName)
	}
	if barFunc.TotalLines != 2 {
		t.Errorf("Expected bar to have 2 lines, got %d", barFunc.TotalLines)
	}
	if barFunc.CoveredLines != 2 {
		t.Errorf("Expected bar to have 2 covered lines, got %d", barFunc.CoveredLines)
	}

	// Check foo function (partially covered)
	if fooFunc.FunctionName != "foo" {
		t.Errorf("Expected second function to be 'foo', got '%s'", fooFunc.FunctionName)
	}
	if fooFunc.TotalLines != 3 {
		t.Errorf("Expected foo to have 3 lines, got %d", fooFunc.TotalLines)
	}
	if fooFunc.CoveredLines != 2 {
		t.Errorf("Expected foo to have 2 covered lines, got %d", fooFunc.CoveredLines)
	}
}

func TestCalculateCoverage_WithFilter(t *testing.T) {
	// Test that CalculateCoverage works correctly after filtering
	report := &GcovrReport{
		Files: []File{
			{
				FilePath: "test.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 1},
					{LineNumber: 2, FunctionName: "foo", Count: 1},
					{LineNumber: 3, FunctionName: "bar", Count: 0},
					{LineNumber: 4, FunctionName: "bar", Count: 0},
					{LineNumber: 5, FunctionName: "baz", Count: 1},
				},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
					{Name: "bar", DemangledName: "bar()"},
					{Name: "baz", DemangledName: "baz()"},
				},
			},
		},
	}

	// Create a filter to only include foo and bar
	filterConfig := &FilterConfig{
		Targets: []TargetFile{
			{
				File:      "test.cpp",
				Functions: []string{"foo", "bar"},
			},
		},
	}

	// Apply filter
	filteredReport := ApplyFilter(report, filterConfig)

	// Calculate coverage on filtered report
	result, err := CalculateCoverage(filteredReport)
	if err != nil {
		t.Fatalf("CalculateCoverage() error = %v", err)
	}

	// Should only have foo and bar functions
	if len(result.Functions) != 2 {
		t.Errorf("Expected 2 functions after filter, got %d", len(result.Functions))
	}

	// Total lines should be 4 (2 from foo + 2 from bar)
	if result.TotalLines != 4 {
		t.Errorf("Expected 4 total lines, got %d", result.TotalLines)
	}

	// Covered lines should be 2 (both from foo)
	if result.TotalCoveredLines != 2 {
		t.Errorf("Expected 2 covered lines, got %d", result.TotalCoveredLines)
	}

	// Coverage should be 50%
	if result.CoveragePercentage != 50.0 {
		t.Errorf("Expected 50%% coverage, got %.1f%%", result.CoveragePercentage)
	}
}

func TestFormatCoverageReport(t *testing.T) {
	tests := []struct {
		name     string
		report   *CoverageReport
		contains []string
	}{
		{
			name: "Report with functions",
			report: &CoverageReport{
				Functions: []FunctionCoverage{
					{
						FilePath:      "test.cpp",
						FunctionName:  "foo",
						DemangledName: "foo()",
						TotalLines:    5,
						CoveredLines:  3,
					},
				},
				TotalLines:         5,
				TotalCoveredLines:  3,
				CoveragePercentage: 60.0,
			},
			contains: []string{
				"Coverage Report",
				"Overall Coverage: 3/5 lines (60.0%)",
				"Functions (1):",
				"File: test.cpp",
				"foo()",
				"3/5 lines",
			},
		},
		{
			name: "Empty report",
			report: &CoverageReport{
				Functions:          []FunctionCoverage{},
				TotalLines:         0,
				TotalCoveredLines:  0,
				CoveragePercentage: 0.0,
			},
			contains: []string{
				"Coverage Report",
				"Overall Coverage: 0/0 lines (0.0%)",
				"No functions found",
			},
		},
		{
			name: "Multiple files",
			report: &CoverageReport{
				Functions: []FunctionCoverage{
					{
						FilePath:      "a.cpp",
						FunctionName:  "func1",
						DemangledName: "func1()",
						TotalLines:    3,
						CoveredLines:  3,
					},
					{
						FilePath:      "b.cpp",
						FunctionName:  "func2",
						DemangledName: "func2()",
						TotalLines:    2,
						CoveredLines:  0,
					},
				},
				TotalLines:         5,
				TotalCoveredLines:  3,
				CoveragePercentage: 60.0,
			},
			contains: []string{
				"File: a.cpp",
				"func1()",
				"File: b.cpp",
				"func2()",
				"0/2 lines",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCoverageReport(tt.report)

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected output to contain %q, but it doesn't.\nOutput: %s", substr, result)
				}
			}
		})
	}
}
