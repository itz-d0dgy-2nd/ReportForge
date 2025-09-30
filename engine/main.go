package main

import (
	Utils "ReportForge/engine/utils"
	"ReportForge/engine/utils/generators"
	"ReportForge/engine/utils/handlers"
	"flag"
	"path/filepath"
	"sort"
)

/*
setupArgumentParser → Configure ReportForge argument parser:
  - Flag ( `--developmentMode` ): Account for the nested report directory due to git submodule
  - Flag ( `--customPath` ): Account for a custom report directory
*/
func setupArgumentParser() Utils.ArgumentsStruct {
	argumentsProvided := Utils.ArgumentsStruct{}

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
  - AppendicesPath: The report markdown files path (Default: report/4_appendices/* || report\4_appendices\*)
*/
func setupReportPaths(_argumentsProvided Utils.ArgumentsStruct) Utils.ReportPathsStruct {
	rootPath := _argumentsProvided.CustomMode

	if _argumentsProvided.DevelopmentMode {
		rootPath = filepath.Clean(filepath.Join(rootPath, "report"))
	}

	return Utils.ReportPathsStruct{
		RootPath:        rootPath,
		ConfigPath:      filepath.Clean(filepath.Join(rootPath, "0_report_config")),
		TemplatePath:    filepath.Clean(filepath.Join(rootPath, "0_report_template", "html", "template.html")),
		SummariesPath:   filepath.Clean(filepath.Join(rootPath, "1_summaries")),
		FindingsPath:    filepath.Clean(filepath.Join(rootPath, "2_findings")),
		SuggestionsPath: filepath.Clean(filepath.Join(rootPath, "3_suggestions")),
		AppendicesPath:  filepath.Clean(filepath.Join(rootPath, "4_appendices")),
	}
}

/*
setupReportData → Execute ReportForge handlers, modifying markdown YAML frontmatter and processing markdown data:
  - Metadata: Recursively iterates over directory structure foreach .yml and process the content
  - Severity: Recursively iterates over directory structure foreach .md and process the severities
  - Summaries: Recursively iterates over directory structure foreach .md and process the summaries
  - Findings: Recursively iterates over directory structure foreach .md and process the findings
  - Suggestions: Recursively iterates over directory structure foreach .md and process the suggestions
  - Appendices: Recursively iterates over directory structure foreach .md and process the appendices
  - Path: The report file path, used for HTML `href=` and `src=`
*/
func setupReportData(_reportPaths Utils.ReportPathsStruct) Utils.ReportDataStruct {
	metadata, severityAssessment := handlers.ConfigFileHandler(_reportPaths.RootPath)
	handlers.ModifierFileHandler(_reportPaths.RootPath, severityAssessment)

	return Utils.ReportDataStruct{
		Metadata:    metadata,
		Severity:    handlers.SeverityFileHandler(_reportPaths.FindingsPath, severityAssessment),
		Summaries:   handlers.MarkdownFileHandler(_reportPaths.RootPath, _reportPaths.SummariesPath, metadata),
		Findings:    handlers.MarkdownFileHandler(_reportPaths.RootPath, _reportPaths.FindingsPath, metadata),
		Suggestions: handlers.MarkdownFileHandler(_reportPaths.RootPath, _reportPaths.SuggestionsPath, metadata),
		Appendices:  handlers.MarkdownFileHandler(_reportPaths.RootPath, _reportPaths.AppendicesPath, metadata),
		Path:        _reportPaths.RootPath,
	}
}

func main() {
	argumentsParsed := setupArgumentParser()
	reportPaths := setupReportPaths(argumentsParsed)
	reportData := setupReportData(reportPaths)

	sort.Slice(reportData.Findings, func(i, j int) bool {
		if reportData.Findings[i].Directory != reportData.Findings[j].Directory {
			return reportData.Findings[i].Directory < reportData.Findings[j].Directory
		}
		return reportData.Findings[i].FileName < reportData.Findings[j].FileName
	})

	sort.Slice(reportData.Suggestions, func(i, j int) bool {
		if reportData.Suggestions[i].Directory != reportData.Suggestions[j].Directory {
			return reportData.Suggestions[i].Directory < reportData.Suggestions[j].Directory
		}
		return reportData.Suggestions[i].FileName < reportData.Suggestions[j].FileName
	})

	generators.GenerateHTML(reportData, reportPaths)
	generators.GeneratePDF(reportPaths)
	generators.GenerateXLSX(reportData)
}
