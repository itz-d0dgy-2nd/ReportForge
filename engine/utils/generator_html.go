package Utils

import (
	"os"
	"text/template"
)

func GenerateHTML(_reportTemplatePath, _HTMLTemplatePath string, _frontMatter FrontmatterYML, _reportSummaries []Markdown, _severity SeverityAssessmentYML, _findings []Markdown, _suggestions []Markdown, _appendices []Markdown) {

	currentProject := Report{
		Frontmatter:     _frontMatter,
		ReportSummaries: _reportSummaries,
		Severity:        _severity,
		Findings:        _findings,
		Suggestions:     _suggestions,
		Appendices:      _appendices,
		Path:            _reportTemplatePath,
	}

	templateHTML, ErrTemplateHTML := template.ParseFiles(_HTMLTemplatePath)
	ErrorChecker(ErrTemplateHTML)

	createHTML, ErrCreateHTML := os.Create("Report.html")
	ErrorChecker(ErrCreateHTML)

	defer createHTML.Close()

	ErrGenerateHTML := templateHTML.Execute(createHTML, currentProject)
	ErrorChecker(ErrGenerateHTML)

}
