package Utils

import (
	"os"
	"text/template"
)

func GenerateHTML(frontMatter FrontmatterYML, severity SeverityAssessmentYML, reportsummaries []Markdown, findings []Markdown, suggestions []Markdown, appendices []Markdown) {

	currentProject := Report{
		Frontmatter:     frontMatter,
		Severity:        severity,
		ReportSummaries: reportsummaries,
		Findings:        findings,
		Suggestions:     suggestions,
		Appendices:      appendices,
	}

	templateHTML, ErrTemplateHTML := template.ParseFiles("engine/template/html/template.html")
	ErrorChecker(ErrTemplateHTML)

	createHTML, ErrCreateHTML := os.Create("Report.html")
	ErrorChecker(ErrCreateHTML)

	defer createHTML.Close()

	ErrGenerateHTML := templateHTML.Execute(createHTML, currentProject)
	ErrorChecker(ErrGenerateHTML)

}
