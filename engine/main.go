package main

import (
	Utils "ReportForge/engine/utils"
	"flag"
	"path/filepath"
)

/*
Main function -> Sets up ReportForge:
  - Configures arguments
  - Configures operating specific file paths
  - Executes file handlers
  - Generates reports
*/
func main() {

	// Create an argument parser:
	//   - Flag ( `--dev` ):
	//     Account for the nested report directory due to git submodule

	devMode := flag.Bool("dev", false, "Run in development mode")
	flag.Parse()

	// Create operating system specific file paths:
	//   - Linux/MacOS:
	//     report/
	//
	//   - Windows:
	//     report\

	reportPath := filepath.Join("report")
	if *devMode {
		reportPath = filepath.Join(reportPath, "report")
	}
	reportConfigPath := filepath.Join(reportPath, "0_report_config")
	reportSummariesPath := filepath.Join(reportPath, "1_summaries")
	reportFindingsPath := filepath.Join(reportPath, "2_findings")
	reportSuggestionsPath := filepath.Join(reportPath, "3_suggestions")
	reportAppendicesPath := filepath.Join(reportPath, "4_appendices")
	HTMLTemplatePath := filepath.Join(reportPath, "0_report_template", "html", "template.html")

	// Execute ReportForge functionality
	//   - FileHandlerReportConfig( string ):
	// 	   Iterates over directory structure foreach .yml file call ProcessConfigFrontmatter()
	//     Returns FrontmatterYML, SeverityAssessmentYML
	//
	//   - FileHandlerSeverity( SeverityAssessmentYML, string ):
	//     Iterates over directory structure foreach .md file call ProcessSeverityMatrix()
	//	   Returns SeverityAssessmentYML
	//
	//   - FileHandlerMarkdown( string, FrontmatterYML, SeverityAssessmentYML. string ):
	//     Iterates over directory structure foreach .md file call ProcessMarkdown()
	//     Returns processedMD

	frontmatter, severityAssessment := Utils.FileHandlerReportConfig(reportConfigPath)
	severity := Utils.FileHandlerSeverity(severityAssessment, reportFindingsPath)
	summaries := Utils.FileHandlerMarkdown(reportPath, frontmatter, severityAssessment, reportSummariesPath)
	findings := Utils.FileHandlerMarkdown(reportPath, frontmatter, severityAssessment, reportFindingsPath)
	suggestions := Utils.FileHandlerMarkdown(reportPath, frontmatter, severityAssessment, reportSuggestionsPath)
	appendices := Utils.FileHandlerMarkdown(reportPath, frontmatter, severityAssessment, reportAppendicesPath)

	// Execute ReportForge functionality
	//   - GenerateHTML( FrontmatterYML,  SeverityAssessmentYML, []Markdown, []Markdown, []Markdown, []Markdown,, string, string ):
	//     Create HTML report
	//
	//   - GenerateXSLX( []Markdown, []Markdown ):
	//     Create XSLX report
	//
	//   - GeneratePDF():
	//     Create PDF report

	Utils.GenerateHTML(frontmatter, severity, summaries, findings, suggestions, appendices, reportPath, HTMLTemplatePath)
	Utils.GenerateXSLX(findings, suggestions)
	Utils.GeneratePDF()
}
