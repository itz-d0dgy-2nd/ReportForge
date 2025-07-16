package main

import (
	Utils "ReportForge/engine/utils"
)

func main() {

	// Yes i know i can improve this but for now its working. Ill come back and refactor later
	frontmatter := Utils.ProcessFrontmatter("report/frontmatter.json")
	severity := Utils.SeverityFileHandler("report/2_findings")
	reportBody := Utils.MarkdownFileHandler(frontmatter, "report/1_report_body")
	findings := Utils.MarkdownFileHandler(frontmatter, "report/2_findings")
	suggestions := Utils.MarkdownFileHandler(frontmatter, "report/3_suggestions")
	appendices := Utils.MarkdownFileHandler(frontmatter, "report/4_appendices")

	Utils.GenerateHTML(frontmatter, severity, reportBody, findings, suggestions, appendices)
	Utils.GeneratePDF()
	Utils.GenerateXSLX(findings, suggestions)

}
