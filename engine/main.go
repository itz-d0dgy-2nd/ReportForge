package main

import (
	"ReportForge/engine/generators"
	"ReportForge/engine/handlers"
	"ReportForge/engine/utilities"
	"flag"
	"path/filepath"
	"sort"
)

/*
setupArgumentParser → Configure ReportForge argument parser:
  - Flag ( `--developmentMode` ): Account for the nested report directory due to git submodule
  - Flag ( `--customPath` ): Account for a custom report directory
*/
func setupArgumentParser() utilities.ArgumentsStruct {
	var argumentsProvided utilities.ArgumentsStruct

	flag.BoolVar(&argumentsProvided.DevelopmentMode, "developmentMode", false, "Run in development mode")
	flag.StringVar(&argumentsProvided.CustomMode, "customPath", "report", "Custom Path")
	flag.Parse()

	return argumentsProvided
}

/*
setupReportPaths → Configure ReportForge paths, normalising them for Windows/Unix based systems:
  - RootPath: The report file path (Default: report/* || report\*)
  - ConfigPath: The report YAML files path (Default: report/0_report_config/* || report\0_report_config\*)
  - TemplatePath: The report template file path (Default: report/0_report_template/* || report\0_report_template\*)
  - SummariesPath: The report markdown files path (Default: report/1_summaries/* || report\1_summaries\*)
  - FindingsPath: The report markdown files path (Default: report/2_findings/* || report\2_findings\*)
  - SuggestionsPath: The report markdown files path (Default: report/3_suggestions/* || report\3_suggestions\*)
  - RisksPath: The report markdown files path (Default: report/4_risks/* || report\4_risks\*)
  - AppendicesPath: The report markdown files path (Default: report/5_appendices/* || report\5_appendices\*)
*/
func setupReportPaths(_argumentsProvided utilities.ArgumentsStruct) utilities.ReportPathsStruct {
	rootPath := _argumentsProvided.CustomMode

	if _argumentsProvided.DevelopmentMode {
		rootPath = filepath.Clean(filepath.Join(rootPath, "report"))
	}

	return utilities.ReportPathsStruct{
		RootPath:        rootPath,
		ConfigPath:      filepath.Clean(filepath.Join(rootPath, "0_report_config")),
		TemplatePath:    filepath.Clean(filepath.Join(rootPath, "0_report_template", "html", "template.html")),
		SummariesPath:   filepath.Clean(filepath.Join(rootPath, "1_summaries")),
		FindingsPath:    filepath.Clean(filepath.Join(rootPath, "2_findings")),
		SuggestionsPath: filepath.Clean(filepath.Join(rootPath, "3_suggestions")),
		RisksPath:       filepath.Clean(filepath.Join(rootPath, "4_risks")),
		AppendicesPath:  filepath.Clean(filepath.Join(rootPath, "5_appendices")),
	}
}

/*
setupReportData → Execute ReportForge handlers, modifying markdown YAML frontmatter and processing markdown data:
  - Metadata: Recursively iterates over directory structure foreach .yml and process the content
  - Severity: Recursively iterates over directory structure foreach .md and process the severities
  - Summaries: Recursively iterates over directory structure foreach .md and process the summaries
  - Findings: Recursively iterates over directory structure foreach .md and process the findings
  - Suggestions: Recursively iterates over directory structure foreach .md and process the suggestions
  - Risks: Recursively iterates over directory structure foreach .md and process the risks
  - Appendices: Recursively iterates over directory structure foreach .md and process the appendices
  - Path: The report file path, used for HTML `href=` and `src=`
*/
func setupReportData(_reportPaths utilities.ReportPathsStruct) utilities.ReportDataStruct {
	metadata, severityAssessment := handlers.HandleConfigProcessor(_reportPaths.ConfigPath)

	handlers.HandleSeverityModifier(_reportPaths.FindingsPath, severityAssessment)
	handlers.HandleIdentifierModifier(_reportPaths.RootPath, metadata)

	return utilities.ReportDataStruct{
		Metadata:    metadata,
		Severity:    handlers.HandleSeverityProcessor(_reportPaths.FindingsPath, severityAssessment),
		Summaries:   handlers.HandleMarkdownProcessor(_reportPaths.RootPath, _reportPaths.SummariesPath, metadata),
		Findings:    handlers.HandleMarkdownProcessor(_reportPaths.RootPath, _reportPaths.FindingsPath, metadata),
		Suggestions: handlers.HandleMarkdownProcessor(_reportPaths.RootPath, _reportPaths.SuggestionsPath, metadata),
		Risks:       handlers.HandleMarkdownProcessor(_reportPaths.RootPath, _reportPaths.RisksPath, metadata),
		Appendices:  handlers.HandleMarkdownProcessor(_reportPaths.RootPath, _reportPaths.AppendicesPath, metadata),
		Path:        _reportPaths.RootPath,
	}
}

/*
setupMarkdownSorting → Sorts ReportForge data alphabetically by directory and then by filename:
  - Sorts Directory field in ascending order
  - Sorts FileName field in ascending order
  - Modifies the slice in-place
*/
func setupMarkdownSorting(_markdown []utilities.Markdown) {
	sort.Slice(_markdown, func(i, j int) bool {
		if _markdown[i].Directory != _markdown[j].Directory {
			return _markdown[i].Directory < _markdown[j].Directory
		}
		return _markdown[i].FileName < _markdown[j].FileName
	})
}

func main() {
	argumentsParsed := setupArgumentParser()
	reportPaths := setupReportPaths(argumentsParsed)
	reportData := setupReportData(reportPaths)

	setupMarkdownSorting(reportData.Findings)
	setupMarkdownSorting(reportData.Suggestions)
	setupMarkdownSorting(reportData.Risks)

	generators.GenerateHTML(reportData, reportPaths)
	generators.GeneratePDF(reportPaths)
	generators.GenerateXLSX(reportData)
}
