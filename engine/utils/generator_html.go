package Utils

import (
	"os"
	"text/template"
)

func GenerateHTML(frontMatter FrontmatterYML, reportSummaries []Markdown, severity SeverityAssessmentYML, findings []Markdown, suggestions []Markdown, appendices []Markdown) {

	currentProject := Report{
		Frontmatter:     frontMatter,
		ReportSummaries: reportSummaries,
		Severity:        severity,
		Findings:        findings,
		Suggestions:     suggestions,
		Appendices:      appendices,
	}

	templateHTML, ErrTemplateHTML := template.ParseFiles("report/0_report_template/html/template.html")
	ErrorChecker(ErrTemplateHTML)

	createHTML, ErrCreateHTML := os.Create("Report.html")
	ErrorChecker(ErrCreateHTML)

	defer createHTML.Close()

	ErrGenerateHTML := templateHTML.Execute(createHTML, currentProject)
	ErrorChecker(ErrGenerateHTML)

}
