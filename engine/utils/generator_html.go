package Utils

import (
	"os"
	"text/template"
)

func GenerateHTML(_frontmatter FrontmatterYML, _severity SeverityAssessmentYML, _summaries []Markdown, _findings []Markdown, _suggestions []Markdown, _appendices []Markdown, _reportTemplatePath string, _HTMLTemplatePath string) {

	currentProject := Report{
		Frontmatter: _frontmatter,
		Severity:    _severity,
		Summaries:   _summaries,
		Findings:    _findings,
		Suggestions: _suggestions,
		Appendices:  _appendices,
		Path:        _reportTemplatePath,
	}

	templateHTML, errTemplateHTML := template.ParseFiles(_HTMLTemplatePath)
	ErrorChecker(errTemplateHTML)

	createHTML, errCreateHTML := os.Create("Report.html")
	ErrorChecker(errCreateHTML)

	defer createHTML.Close()

	errGenerateHTML := templateHTML.Execute(createHTML, currentProject)
	ErrorChecker(errGenerateHTML)

}
