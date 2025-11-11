package gcovr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FilterConfig represents the filter configuration file structure
type FilterConfig struct {
	Compiler struct {
		Path          string `yaml:"path"`
		GcovrExecPath string `yaml:"gcovr_exec_path"`
	} `yaml:"compiler"`
	Targets []TargetFile `yaml:"targets"`
}

// TargetFile represents a file and its target functions to track
type TargetFile struct {
	File      string   `yaml:"file"`
	Functions []string `yaml:"functions"`
}

// ParseFilterConfig reads and parses a filter configuration file
func ParseFilterConfig(filePath string) (*FilterConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read filter config %s: %w", filePath, err)
	}

	var config FilterConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML from %s: %w", filePath, err)
	}

	return &config, nil
}

// ApplyFilter filters a GcovrReport based on the filter configuration
// It only keeps files and functions specified in the targets
func ApplyFilter(report *GcovrReport, config *FilterConfig) *GcovrReport {
	if config == nil || len(config.Targets) == 0 {
		return report
	}

	// Build a map of file -> allowed functions
	filterMap := make(map[string]map[string]bool)
	for _, target := range config.Targets {
		// Normalize file paths for comparison
		normalizedFile := normalizeFilePath(target.File)

		funcMap := make(map[string]bool)
		for _, fn := range target.Functions {
			funcMap[fn] = true
		}
		filterMap[normalizedFile] = funcMap
	}

	// Filter the report
	filteredReport := &GcovrReport{
		FormatVersion: report.FormatVersion,
		Files:         make([]File, 0),
	}

	for _, file := range report.Files {
		normalizedFilePath := normalizeFilePath(file.FilePath)

		// Check if this file is in the filter
		allowedFunctions, fileInFilter := filterMap[normalizedFilePath]
		if !fileInFilter {
			// Try matching just the filename
			fileName := filepath.Base(file.FilePath)
			allowedFunctions, fileInFilter = filterMap[fileName]
			if !fileInFilter {
				continue // Skip this file
			}
		}

		// Filter lines and functions
		filteredFile := File{
			FilePath:  file.FilePath,
			Lines:     make([]Line, 0),
			Functions: make([]Function, 0),
		}

		// Filter functions
		for _, fn := range file.Functions {
			if shouldIncludeFunction(fn.DemangledName, fn.Name, allowedFunctions) {
				filteredFile.Functions = append(filteredFile.Functions, fn)
			}
		}

		// Filter lines to only include those from allowed functions
		allowedFuncNames := make(map[string]bool)
		for _, fn := range filteredFile.Functions {
			allowedFuncNames[fn.Name] = true
		}

		for _, line := range file.Lines {
			if allowedFuncNames[line.FunctionName] {
				filteredFile.Lines = append(filteredFile.Lines, line)
			}
		}

		// Only add the file if it has content after filtering
		if len(filteredFile.Functions) > 0 {
			filteredReport.Files = append(filteredReport.Files, filteredFile)
		}
	}

	return filteredReport
}

// normalizeFilePath normalizes file paths for comparison
func normalizeFilePath(path string) string {
	// Clean the path and convert to forward slashes
	cleaned := filepath.Clean(path)
	cleaned = filepath.ToSlash(cleaned)
	return cleaned
}

// shouldIncludeFunction checks if a function should be included based on filter
func shouldIncludeFunction(demangledName, mangledName string, allowedFunctions map[string]bool) bool {
	// Check demangled name (strip parentheses for matching)
	simpleName := demangledName
	if idx := strings.Index(simpleName, "("); idx != -1 {
		simpleName = simpleName[:idx]
	}

	if allowedFunctions[simpleName] {
		return true
	}

	// Also check full demangled name
	if allowedFunctions[demangledName] {
		return true
	}

	// Check mangled name
	if allowedFunctions[mangledName] {
		return true
	}

	return false
}
