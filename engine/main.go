package main

import (
	Utils "ReportForge/engine/utils"
	"ReportForge/engine/utils/generators"
	"ReportForge/engine/utils/handlers"
	"flag"
	"path/filepath"
)

func SetupArgumentParser() Utils.ArgumentsStruct {

	// Configure an argument parser:
	//   - Flag ( `--developmentMode` ): Account for the nested report directory due to git submodule
	//   - Flag ( `--customPath` ): Account for a custom report directory

	argumentsProvided := Utils.ArgumentsStruct{}
	flag.BoolVar(&argumentsProvided.DevelopmentMode, "developmentMode", false, "Run in development mode")
	flag.StringVar(&argumentsProvided.CustomMode, "customPath", "report", "Custom Path")
	flag.Parse()
	return argumentsProvided
}

func SetupReportPaths(_argumentsProvided Utils.ArgumentsStruct) Utils.ReportPathsStruct {

	// Configure report paths:
	//   - RootPath: The report file path (Default: report/* || report\*)
	//   - ConfigPath: The report yaml files path (Default: report/0_report_config/* || report\0_report_config\*)
	//   - TemplatePath: The report template file path (Default: report/0_report_template/* || report\0_report_template\*)
	//   - SummariesPath: The report markdown files path (Default: report/1_summaries/* || report\1_summaries\*)
	//   - FindingsPath: The report markdown files path (Default: report/2_findings/* || report\2_findings\*)
	//   - SuggestionsPath: The report markdown files path (Default: report/3_suggestions/* || report\3_suggestions\*)
	//   - AppendicesPath: The report markdown files path (Default: report/4_appendices/* || report\4_appendices\*)

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

func SetupReportData(_reportPaths Utils.ReportPathsStruct) Utils.ReportDataStruct {

	// Execute ReportForge functionality:
	//   - Frontmatter: Recursively iterates over directory structure foreach .yml and process the content
	//   - Severity: Recursively iterates over directory structure foreach .md and process the severities
	//   - Summaries: Recursively iterates over directory structure foreach .md and process the executive and technical summaries
	//   - Findings: Recursively iterates over directory structure foreach .md and process the findings
	//   - Suggestions: Recursively iterates over directory structure foreach .md and process the suggestions
	//   - Appendices: Recursively iterates over directory structure foreach .md and process the appendices
	//   - Path: The report file path, used for HTML `href=` and `src=`

	frontmatter, severityAssessment := handlers.ReportConfigFileHandler(_reportPaths.RootPath)

	return Utils.ReportDataStruct{
		Frontmatter: frontmatter,
		Severity:    handlers.SeverityFileHandler(_reportPaths.FindingsPath, severityAssessment),
		Summaries:   handlers.MarkdownFileHandler(_reportPaths.RootPath, _reportPaths.SummariesPath, frontmatter, severityAssessment),
		Findings:    handlers.MarkdownFileHandler(_reportPaths.RootPath, _reportPaths.FindingsPath, frontmatter, severityAssessment),
		Suggestions: handlers.MarkdownFileHandler(_reportPaths.RootPath, _reportPaths.SuggestionsPath, frontmatter, severityAssessment),
		Appendices:  handlers.MarkdownFileHandler(_reportPaths.RootPath, _reportPaths.AppendicesPath, frontmatter, severityAssessment),
		Path:        _reportPaths.RootPath,
	}
}

func reportGeneration(_reportData Utils.ReportDataStruct, _reportPaths Utils.ReportPathsStruct) {

	// Execute ReportForge functionality:
	//   - GenerateHTML(): Create HTML report
	//   - GeneratePDF(): Create PDF report
	//   - GenerateXSLX(): Create XSLX report

	generators.GenerateHTML(_reportData, _reportPaths)
	generators.GeneratePDF(_reportPaths)
	generators.GenerateXSLX(_reportData.Findings, _reportData.Suggestions)

}

func main() {
	argumentsParsed := SetupArgumentParser()
	reportPaths := SetupReportPaths(argumentsParsed)
	reportData := SetupReportData(reportPaths)
	reportGeneration(reportData, reportPaths)
}
