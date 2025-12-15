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
func setupReportData(_reportPaths utilities.ReportPaths) utilities.ReportData {
	fileCache := utilities.NewFileCache(_reportPaths.RootPath)
	metadataConfig, severityConfig, directoryConfig := handlers.HandleConfigProcessor(_reportPaths, fileCache)
	utilities.DocumentStatus = metadataConfig.DocumentInformation[len(metadataConfig.DocumentInformation)-1].DocumentVersioning["DocumentStatus"]
	handlers.HandleModifications(_reportPaths, fileCache)
	severityMatrix, severityBarGraph, summaries, findings, suggestions, risks, appendices := handlers.HandleProcessing(_reportPaths, fileCache)

	return utilities.ReportData{
		MetadataConfig:   metadataConfig,
		SeverityConfig:   severityConfig,
		DirectoryConfig:  directoryConfig,
		SeverityMatrix:   severityMatrix,
		SeverityBarGraph: severityBarGraph,
		Summaries:        summaries,
		Findings:         findings,
		Suggestions:      suggestions,
		Risks:            risks,
		Appendices:       appendices,
		Path:             _reportPaths.RootPath,
	}
}

func main() {
	argumentsParsed := setupArgumentParser()
	reportPaths := setupReportPaths(argumentsParsed)
	reportData := setupReportData(reportPaths)

	utilities.SortSeverityMatrix(&reportData.SeverityMatrix)
	utilities.SortReportData(reportData.Summaries, utilities.SummariesDirectory, reportData.DirectoryConfig)
	utilities.SortReportData(reportData.Findings, utilities.FindingsDirectory, reportData.DirectoryConfig)
	utilities.SortReportData(reportData.Suggestions, utilities.SuggestionsDirectory, reportData.DirectoryConfig)
	utilities.SortReportData(reportData.Risks, utilities.RisksDirectory, reportData.DirectoryConfig)
	utilities.OptimiseImagesForPDF(reportPaths.ScreenshotsPath)

	generators.GenerateHTML(reportData, reportPaths)
	generators.GeneratePDF(reportData, reportPaths)
	generators.GenerateXLSX(reportData)
}
