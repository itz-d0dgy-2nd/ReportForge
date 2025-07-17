package main

import (
	Utils "ReportForge/engine/utils"
)

func main() {

	// Yes i know i can improve this but for now its working. Ill come back and refactor later
	frontmatter, matrix := Utils.FileHandlerReportConfig("report/0_report_config")
	severity := Utils.FileHandlerSeverity(matrix, "report/2_findings")
	reportBody := Utils.FileHandlerMarkdown(frontmatter, matrix, "report/1_report_body")
	findings := Utils.FileHandlerMarkdown(frontmatter, matrix, "report/2_findings")
	suggestions := Utils.FileHandlerMarkdown(frontmatter, matrix, "report/3_suggestions")
	appendices := Utils.FileHandlerMarkdown(frontmatter, matrix, "report/4_appendices")

	Utils.GenerateHTML(frontmatter, severity, reportBody, findings, suggestions, appendices)
	Utils.GeneratePDF()
	Utils.GenerateXSLX(findings, suggestions)

}
