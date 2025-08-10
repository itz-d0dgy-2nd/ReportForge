package main

import (
	Utils "ReportForge/engine/utils"
	"flag"
	"path/filepath"
)

func main() {
	devMode := flag.Bool("dev", false, "Run in development mode (e.g., `go run engine/main.go --dev`)")
	flag.Parse()

	// Yes i know i can improve this but for now its working. Ill come back and refactor later
	reportTemplatePath := filepath.Join("report")
	if *devMode {
		reportTemplatePath = filepath.Join(reportTemplatePath, "report")
	}
	reportConfigPath := filepath.Join(reportTemplatePath, "0_report_config")
	reportSummariesPath := filepath.Join(reportTemplatePath, "1_summaries")
	reportFindingsPath := filepath.Join(reportTemplatePath, "2_findings")
	reportSuggestionsPath := filepath.Join(reportTemplatePath, "3_suggestions")
	reportAppendicesPath := filepath.Join(reportTemplatePath, "4_appendices")
	HTMLTemplatePath := filepath.Join(reportTemplatePath, "0_report_template", "html", "template.html")

	// Yes i know i can improve this but for now its working. Ill come back and refactor later
	frontMatter, severityAssessment := Utils.FileHandlerReportConfig(reportConfigPath)
	reportSummaries := Utils.FileHandlerMarkdown(reportTemplatePath, frontMatter, severityAssessment, reportSummariesPath)
	severity := Utils.FileHandlerSeverity(severityAssessment, reportFindingsPath)
	findings := Utils.FileHandlerMarkdown(reportTemplatePath, frontMatter, severityAssessment, reportFindingsPath)
	suggestions := Utils.FileHandlerMarkdown(reportTemplatePath, frontMatter, severityAssessment, reportSuggestionsPath)
	appendices := Utils.FileHandlerMarkdown(reportTemplatePath, frontMatter, severityAssessment, reportAppendicesPath)

	// Yes i know i can improve this but for now its working. Ill come back and refactor later
	Utils.GenerateHTML(reportTemplatePath, HTMLTemplatePath, frontMatter, reportSummaries, severity, findings, suggestions, appendices)
	Utils.GeneratePDF()
	Utils.GenerateXSLX(findings, suggestions)
}
