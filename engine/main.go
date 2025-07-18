package main

import (
	Utils "ReportForge/engine/utils"
)

func main() {

	// Yes i know i can improve this but for now its working. Ill come back and refactor later
	frontMatter, severityAssessment := Utils.FileHandlerReportConfig("report/0_report_config")
	reportSummaries := Utils.FileHandlerMarkdown(frontMatter, severityAssessment, "report/1_report_summaries")
	severity := Utils.FileHandlerSeverity(severityAssessment, "report/2_findings")
	findings := Utils.FileHandlerMarkdown(frontMatter, severityAssessment, "report/2_findings")
	suggestions := Utils.FileHandlerMarkdown(frontMatter, severityAssessment, "report/3_suggestions")
	appendices := Utils.FileHandlerMarkdown(frontMatter, severityAssessment, "report/4_appendices")

	Utils.GenerateHTML(frontMatter, reportSummaries, severity, findings, suggestions, appendices)
	Utils.GeneratePDF()
	Utils.GenerateXSLX(findings, suggestions)

}
