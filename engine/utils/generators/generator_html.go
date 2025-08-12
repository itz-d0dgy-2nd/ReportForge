package generators

import (
	"ReportForge/engine/utils"
	"os"
	"text/template"
)

func GenerateHTML(_frontmatter Utils.FrontmatterYML, _severity Utils.SeverityAssessmentYML, _summaries []Utils.Markdown, _findings []Utils.Markdown, _suggestions []Utils.Markdown, _appendices []Utils.Markdown, _reportTemplatePath string, _HTMLTemplatePath string) {

	currentProject := Utils.Report{
		Frontmatter: _frontmatter,
		Severity:    _severity,
		Summaries:   _summaries,
		Findings:    _findings,
		Suggestions: _suggestions,
		Appendices:  _appendices,
		Path:        _reportTemplatePath,
	}

	templateHTML, errTemplateHTML := template.ParseFiles(_HTMLTemplatePath)
	Utils.ErrorChecker(errTemplateHTML)

	createHTML, errCreateHTML := os.Create("Report.html")
	Utils.ErrorChecker(errCreateHTML)

	defer createHTML.Close()

	errGenerateHTML := templateHTML.Execute(createHTML, currentProject)
	Utils.ErrorChecker(errGenerateHTML)

}
