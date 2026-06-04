package main

import (
	"ReportForge/engine/generators"
	"ReportForge/engine/handlers"
	"ReportForge/engine/utilities"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

/*
setupArgumentParser → configures and parses ReportForge command-line arguments.

Flags:
  - --customPath: specifies custom report directory path
  - --debug: enables debug logging
  - --watch: enables file watcher
*/
func setupArgumentParser() utilities.Arguments {
	customPath := flag.String("customPath", "", "Custom path to report directory")
	debug := flag.Bool("debug", false, "Enable debug logging")
	watch := flag.Bool("watch", false, "Enable file watcher")
	flag.Parse()

	return utilities.Arguments{
		CustomPath: *customPath,
		Debug:      *debug,
		Watch:      *watch,
	}
}

/*
setupReportTemplate → detects report template type from directory structure.
*/
func setupReportTemplate(_rootPath string) utilities.TemplateType {
	var templateType utilities.TemplateType

	if _, errStat := os.Stat(filepath.Join(_rootPath, utilities.FindingsDirectory)); errStat == nil {
		templateType.Technical = true
	}

	if _, errStat := os.Stat(filepath.Join(_rootPath, utilities.RisksDirectory)); errStat == nil {
		templateType.Sra = true
	}

	if templateType.Sra && templateType.Technical {
		utilities.Check(utilities.NewConfigError(
			_rootPath,
			fmt.Sprintf("conflicting template types detected - found both '%s' and '%s' directories, but only one template type is allowed per report",
				utilities.FindingsDirectory,
				utilities.RisksDirectory,
			),
		))
	}

	if !templateType.Technical && !templateType.Sra {
		utilities.Check(utilities.NewConfigError(
			_rootPath,
			fmt.Sprintf("unable to detect template type - expected either '%s' (technical) or '%s' (SRA) directory",
				utilities.FindingsDirectory,
				utilities.RisksDirectory,
			),
		))
	}

	return templateType
}

/*
setupReportPaths → constructs and normalises all report directory paths.
*/
func setupReportPaths(_argumentsProvided utilities.Arguments) utilities.ReportPaths {
	rootPath := filepath.Clean(_argumentsProvided.CustomPath)
	template := setupReportTemplate(rootPath)

	conditionalPath := func(condition bool, dir string) string {
		if condition {
			return filepath.Clean(filepath.Join(rootPath, dir))
		}
		return ""
	}

	return utilities.ReportPaths{
		RootPath:        rootPath,
		ConfigPath:      filepath.Clean(filepath.Join(rootPath, utilities.ConfigDirectory)),
		TemplatePath:    filepath.Clean(filepath.Join(rootPath, utilities.TemplateDirectory)),
		SummariesPath:   filepath.Clean(filepath.Join(rootPath, utilities.SummariesDirectory)),
		FindingsPath:    conditionalPath(template.Technical, utilities.FindingsDirectory),
		SuggestionsPath: conditionalPath(template.Technical, utilities.SuggestionsDirectory),
		RisksPath:       conditionalPath(template.Sra, utilities.RisksDirectory),
		ControlsPath:    conditionalPath(template.Sra, utilities.ControlsDirectory),
		ScreenshotsPath: filepath.Clean(filepath.Join(rootPath, utilities.ScreenshotsDirectory)),
		AppendicesPath:  filepath.Clean(filepath.Join(rootPath, utilities.AppendicesDirectory)),
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

	utilities.Logger.Info("report generation complete")
}

func main() {
	argumentsParsed := setupArgumentParser()
	utilities.NewLogger(argumentsParsed.Debug)
	reportPaths := setupReportPaths(argumentsParsed)
	utilities.Logger.Info("initialising file cache", "path", reportPaths.RootPath)
	fileCache := utilities.NewFileCache(reportPaths.RootPath)
	generateReport(reportPaths, fileCache)

	if argumentsParsed.Watch {
		watcher, errWatcher := utilities.NewWatcher()
		if errWatcher != nil {
			utilities.Check(errWatcher)
		}
		defer watcher.Close()

		errWatcher = watcher.WatchForChanges(reportPaths.RootPath, fileCache, func(changedFile string) {
			generateReport(reportPaths, fileCache)
		})
		if errWatcher != nil {
			utilities.Check(errWatcher)
		}
	}
}
