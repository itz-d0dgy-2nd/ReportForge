package main

import (
	"ReportForge/engine/generators"
	"ReportForge/engine/handlers"
	"ReportForge/engine/utilities"
	"flag"
	"path/filepath"
)

/*
setupArgumentParser → Configure ReportForge argument parser:
  - Parses command-line arguments using flag.Parse()
    -- Flag (--developmentMode): Enable development mode to account for nested report directory structure due to git submodule
    -- Flag (--customPath): Specify custom report directory path (default: "report")
  - Returns utilities.Arguments struct containing parsed flag values
*/
func setupArgumentParser() utilities.Arguments {
	var argumentsProvided utilities.Arguments

	flag.BoolVar(&argumentsProvided.DevelopmentMode, "developmentMode", false, "Run in development mode")
	flag.StringVar(&argumentsProvided.CustomPath, "customPath", "report", "Custom Path")
	flag.Parse()

	return argumentsProvided
}

/*
setupReportPaths → Configure ReportForge paths, normalising them for Windows/Unix based systems:
  - Determines root path based on provided arguments
  - Constructs all report subdirectory paths using filepath.Clean() and filepath.Join()
  - Returns utilities.ReportPaths struct containing all report directory paths
*/
func setupReportPaths(_argumentsProvided utilities.Arguments) utilities.ReportPaths {
	rootPath := _argumentsProvided.CustomPath

	if _argumentsProvided.DevelopmentMode {
		rootPath = filepath.Clean(filepath.Join(utilities.RootDirectory, "report"))
	}

	return utilities.ReportPaths{
		RootPath:        rootPath,
		ConfigPath:      filepath.Clean(filepath.Join(rootPath, utilities.ConfigDirectory)),
		TemplatePath:    filepath.Clean(filepath.Join(rootPath, utilities.TemplateDirectory, "html", "template.html")),
		SummariesPath:   filepath.Clean(filepath.Join(rootPath, utilities.SummariesDirectory)),
		FindingsPath:    filepath.Clean(filepath.Join(rootPath, utilities.FindingsDirectory)),
		SuggestionsPath: filepath.Clean(filepath.Join(rootPath, utilities.SuggestionsDirectory)),
		RisksPath:       filepath.Clean(filepath.Join(rootPath, utilities.RisksDirectory)),
		AppendicesPath:  filepath.Clean(filepath.Join(rootPath, utilities.AppendicesDirectory)),
		ScreenshotsPath: filepath.Clean(filepath.Join(rootPath, utilities.ScreenshotsDirectory)),
	}
}

/*
setupReportData → Execute ReportForge handlers, process config, modify and process markdown generating report data:
  - Initialises file cache with pre-loaded markdown and config files
  - Processes configuration files (metadata, severity assessment, directory order)
  - Modifies markdown files (severity calculation and identifier assignment)
  - Processes markdown files (summaries, findings, suggestions, risks, appendices)
*/
func setupReportData(_reportPaths utilities.ReportPaths) *utilities.FileCache {
	fileCache := utilities.NewFileCache(_reportPaths.RootPath)
	handlers.HandleConfigs(_reportPaths, fileCache)
	handlers.HandleModifications(_reportPaths, fileCache)
	handlers.HandleProcessing(_reportPaths, fileCache)
	return fileCache
}

func main() {
	argumentsParsed := setupArgumentParser()
	reportPaths := setupReportPaths(argumentsParsed)
	fileCache := setupReportData(reportPaths)

	utilities.SortSeverityMatrix(&fileCache.SeverityMatrix)
	utilities.SortReportData(fileCache.Summaries, utilities.SummariesDirectory, fileCache.ContentConfig)
	utilities.SortReportData(fileCache.Findings, utilities.FindingsDirectory, fileCache.ContentConfig)
	utilities.SortReportData(fileCache.Suggestions, utilities.SuggestionsDirectory, fileCache.ContentConfig)
	utilities.SortReportData(fileCache.Risks, utilities.RisksDirectory, fileCache.ContentConfig)
	utilities.OptimiseImagesForPDF(reportPaths.ScreenshotsPath)

	generators.GenerateHTML(fileCache, reportPaths)
	generators.GeneratePDF(fileCache, reportPaths)
	generators.GenerateXLSX(fileCache)
}
