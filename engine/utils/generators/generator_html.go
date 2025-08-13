package generators

import (
	"ReportForge/engine/utils"
	"os"
	"text/template"
)

func GenerateHTML(_reportData Utils.ReportDataStruct, _reportPaths Utils.ReportPathsStruct) {
	templateHTML, errTemplateHTML := template.ParseFiles(_reportPaths.TemplatePath)
	Utils.ErrorChecker(errTemplateHTML)

	createHTML, errCreateHTML := os.Create("Report.html")
	Utils.ErrorChecker(errCreateHTML)

	defer createHTML.Close()

	errGenerateHTML := templateHTML.Execute(createHTML, _reportData)
	Utils.ErrorChecker(errGenerateHTML)
}
