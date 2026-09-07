package main

import (
	"ReportForge/engine/generators"
	"ReportForge/engine/handlers"
	"ReportForge/engine/utilities"
	"flag"
	"path/filepath"
)

/*
setupArgumentParser → configures and parses ReportForge command-line arguments.

Flags:
  - --customPath: specifies custom report directory path
  - --debug: enables debug logging
  - --verbose: enables verbose logging
  - --watch: enables file watcher
*/
func setupArgumentParser() utilities.Arguments {
	customPath := flag.String("customPath", utilities.RootDirectory, "Custom path to report directory")
	debug := flag.Bool("debug", false, "Enable debug logging")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	watch := flag.Bool("watch", false, "Enable file watcher")
	flag.Parse()

	return utilities.Arguments{
		CustomPath: *customPath,
		Debug:      *debug,
		Verbose:    *verbose,
		Watch:      *watch,
	}
}

/*
setupReportPaths → constructs and normalises all report directory paths.
*/
func setupReportPaths(_customPath string) utilities.ReportPaths {
	rootPath := filepath.Clean(_customPath)

	return utilities.ReportPaths{
		RootPath:        rootPath,
		ConfigPath:      filepath.Join(rootPath, utilities.ConfigDirectory),
		TemplatePath:    filepath.Join(rootPath, utilities.TemplateDirectory),
		SummariesPath:   filepath.Join(rootPath, utilities.SummariesDirectory),
		FindingsPath:    filepath.Join(rootPath, utilities.FindingsDirectory),
		SuggestionsPath: filepath.Join(rootPath, utilities.SuggestionsDirectory),
		RisksPath:       filepath.Join(rootPath, utilities.RisksDirectory),
		ControlsPath:    filepath.Join(rootPath, utilities.ControlsDirectory),
		ScreenshotsPath: filepath.Join(rootPath, utilities.ScreenshotsDirectory),
		AppendicesPath:  filepath.Join(rootPath, utilities.AppendicesDirectory),
	}
}

/*
generateReport → generates the complete report from the current file cache state.
*/
func generateReport(_reportPaths utilities.ReportPaths, _fileCache *utilities.FileCache) {
	_fileCache.ClearProcessedData()
	contentConfig := _fileCache.ContentConfig()

	utilities.Logger.Info("Modifying markdown")
	handlers.HandleModifications(_reportPaths, _fileCache)

	utilities.Logger.Info("Processing markdown")
	handlers.HandleProcessing(_reportPaths, _fileCache)

	utilities.Logger.Debug("Sorting matrices")
	utilities.SortSeverityMatrix(&_fileCache.SeverityMatrix)
	utilities.SortRiskMatrices(&_fileCache.RiskMatrices)

	utilities.Logger.Debug("Sorting sections")
	utilities.SortReportData(_fileCache.Summaries, utilities.SummariesDirectory, contentConfig)
	utilities.SortReportData(_fileCache.Findings, utilities.FindingsDirectory, contentConfig)
	utilities.SortReportData(_fileCache.Suggestions, utilities.SuggestionsDirectory, contentConfig)
	utilities.SortReportData(_fileCache.Risks, utilities.RisksDirectory, contentConfig)
	utilities.SortReportData(_fileCache.Controls, utilities.ControlsDirectory, contentConfig)

	utilities.Logger.Info("Optimising images")
	utilities.OptimiseImagesForPDF(_reportPaths.ScreenshotsPath)

	utilities.Logger.Info("Generating HTML")
	generators.GenerateHTML(_fileCache, _reportPaths)

	utilities.Logger.Info("Generating PDF")
	generators.GeneratePDF(_fileCache, _reportPaths)

	utilities.Logger.Info("Generating XLSX")
	generators.GenerateXLSX(_fileCache)

	utilities.Logger.Info("Report generation complete")
}

func main() {
	argumentsParsed := setupArgumentParser()
	utilities.NewLogger(argumentsParsed.Debug, argumentsParsed.Verbose)
	reportPaths := setupReportPaths(argumentsParsed.CustomPath)
	fileCache := utilities.NewFileCache(reportPaths.RootPath)

	generateReport(reportPaths, fileCache)

	if argumentsParsed.Watch {
		watcher, errWatcher := utilities.NewWatcher()
		utilities.Check(errWatcher)
		defer watcher.Close()

		errWatcher = watcher.WatchForChanges(reportPaths.RootPath, fileCache, func() {
			generateReport(reportPaths, fileCache)
		})
		utilities.Check(errWatcher)
	}
}
