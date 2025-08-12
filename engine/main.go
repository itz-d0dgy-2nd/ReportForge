package main

import (
	"ReportForge/engine/utils/generators"
	"ReportForge/engine/utils/handlers"
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

	reportPath := filepath.Clean(filepath.Join("report"))
	if *devMode {
		reportPath = filepath.Clean(filepath.Join(reportPath, "report"))
	}
	reportConfigPath := filepath.Clean(filepath.Join(reportPath, "0_report_config"))
	reportSummariesPath := filepath.Clean(filepath.Join(reportPath, "1_summaries"))
	reportFindingsPath := filepath.Clean(filepath.Join(reportPath, "2_findings"))
	reportSuggestionsPath := filepath.Clean(filepath.Join(reportPath, "3_suggestions"))
	reportAppendicesPath := filepath.Clean(filepath.Join(reportPath, "4_appendices"))
	HTMLTemplatePath := filepath.Clean(filepath.Join(reportPath, "0_report_template", "html", "template.html"))

	// Execute ReportForge functionality
	//   - ReportConfigFileHandler( string ):
	// 	   Iterates over directory structure foreach .yml file call ProcessConfigFrontmatter()
	//     Returns FrontmatterYML, SeverityAssessmentYML
	//
	//   - SeverityFileHandler( string, SeverityAssessmentYML ):
	//     Iterates over directory structure foreach .md file call ProcessSeverityMatrix()
	//	   Returns SeverityAssessmentYML
	//
	//   - MarkdownFileHandler( string, string, FrontmatterYML, SeverityAssessmentYML ):
	//     Iterates over directory structure foreach .md file call ProcessMarkdown()
	//     Returns processedMD

	frontmatter, severityAssessment := handlers.ReportConfigFileHandler(reportConfigPath)
	severity := handlers.SeverityFileHandler(reportFindingsPath, severityAssessment)
	summaries := handlers.MarkdownFileHandler(reportPath, reportSummariesPath, frontmatter, severityAssessment)
	findings := handlers.MarkdownFileHandler(reportPath, reportFindingsPath, frontmatter, severityAssessment)
	suggestions := handlers.MarkdownFileHandler(reportPath, reportSuggestionsPath, frontmatter, severityAssessment)
	appendices := handlers.MarkdownFileHandler(reportPath, reportAppendicesPath, frontmatter, severityAssessment)

	// Execute ReportForge functionality
	//   - GenerateHTML( FrontmatterYML,  SeverityAssessmentYML, []Markdown, []Markdown, []Markdown, []Markdown,, string, string ):
	//     Create HTML report
	//
	//   - GenerateXSLX( []Markdown, []Markdown ):
	//     Create XSLX report
	//
	//   - GeneratePDF():
	//     Create PDF report

	generators.GenerateHTML(frontmatter, severity, summaries, findings, suggestions, appendices, reportPath, HTMLTemplatePath)
	generators.GenerateXSLX(findings, suggestions)
	generators.GeneratePDF()
}
