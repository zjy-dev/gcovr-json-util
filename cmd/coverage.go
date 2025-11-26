package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zjy-dev/gcovr-json-util/v2/pkg/gcovr"
)

var (
	coverageFilterFile string
)

// coverageCmd represents the coverage command
var coverageCmd = &cobra.Command{
	Use:     "coverage [gcovr-file]",
	Aliases: []string{"cov"},
	Short:   "Calculate overall coverage from a gcovr JSON report",
	Long: `Analyze a gcovr JSON report and display the overall coverage statistics
for all functions. When used with a filter, calculates coverage for only
the specified files and functions.

The tool will show:
- Overall coverage percentage
- Coverage statistics for each function
- Total lines and covered lines`,
	Args: cobra.ExactArgs(1),
	RunE: runCoverage,
}

func init() {
	rootCmd.AddCommand(coverageCmd)

	coverageCmd.Flags().StringVarP(&coverageFilterFile, "filter", "f", "",
		"Filter config file (YAML) to specify target files and functions")
}

func runCoverage(cmd *cobra.Command, args []string) error {
	reportFile := args[0]

	// Parse the gcovr JSON report
	fmt.Printf("Reading report: %s\n", reportFile)
	report, err := gcovr.ParseReport(reportFile)
	if err != nil {
		return fmt.Errorf("failed to parse report: %w", err)
	}

	// Apply filter if specified
	if coverageFilterFile != "" {
		fmt.Printf("Reading filter config: %s\n", coverageFilterFile)
		filterConfig, err := gcovr.ParseFilterConfig(coverageFilterFile)
		if err != nil {
			return fmt.Errorf("failed to parse filter config: %w", err)
		}

		fmt.Printf("Filtering enabled: tracking %d file(s)\n", len(filterConfig.Targets))
		report = gcovr.ApplyFilter(report, filterConfig)
		fmt.Println("Applying filters...")
	}

	// Calculate coverage
	fmt.Println("Calculating coverage...")
	coverageReport, err := gcovr.CalculateCoverage(report)
	if err != nil {
		return fmt.Errorf("failed to calculate coverage: %w", err)
	}

	// Display results
	fmt.Println()
	output := gcovr.FormatCoverageReport(coverageReport)
	fmt.Print(output)

	return nil
}
