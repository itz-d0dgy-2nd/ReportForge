package generators

import (
	"ReportForge/engine/utilities"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

/*
groupByDirectory → Groups a slice of MarkdownFiles into DirectoryGroups
*/
func groupByDirectory(_files []utilities.MarkdownFile) []utilities.DirectoryGroup {
	var groups []utilities.DirectoryGroup
	groupIndex := make(map[string]int)

	for _, file := range _files {
		if idx, exists := groupIndex[file.Directory]; exists {
			groups[idx].Items = append(groups[idx].Items, file)
		} else {
			groupIndex[file.Directory] = len(groups)
			groups = append(groups, utilities.DirectoryGroup{
				Directory: file.Directory,
				Items:     []utilities.MarkdownFile{file},
			})
		}
	}

	return groups
}

/*
buildControlRiskMappings → Pre-computes control to risk mappings from processed data
*/
func buildControlRiskMappings(_controlFiles []utilities.MarkdownFile, _riskFiles []utilities.MarkdownFile) []utilities.ControlRiskMapping {
	var mappings []utilities.ControlRiskMapping

	for _, control := range _controlFiles {
		mapping := utilities.ControlRiskMapping{
			ControlID:    control.Headers.ControlID,
			ControlTitle: control.Headers.ControlTitle,
			ControlName:  control.Headers.ControlName,
		}

		for _, risk := range _riskFiles {
			recommendedControls := utilities.RiskPattern.RecommendedControls.FindStringSubmatch(risk.Body)
			if len(recommendedControls) > 1 && strings.Contains(recommendedControls[1], "\"#"+control.Headers.ControlName+"\"") {
				mapping.RiskIDs = append(mapping.RiskIDs, risk.Headers.RiskID)
				mapping.RiskName = append(mapping.RiskName, risk.Headers.RiskName)
			}
		}

		mappings = append(mappings, mapping)
	}

	return mappings
}

/*
buildTemplateData → Constructs TemplateData from the file cache for both template types
*/
func buildTemplateData(_fileCache *utilities.FileCache) utilities.TemplateData {
	return utilities.TemplateData{
		FileCache:           _fileCache,
		FindingGroups:       groupByDirectory(_fileCache.Findings),
		SuggestionGroups:    groupByDirectory(_fileCache.Suggestions),
		AppendixGroups:      groupByDirectory(_fileCache.Appendices),
		RiskGroups:          groupByDirectory(_fileCache.Risks),
		ControlGroups:       groupByDirectory(_fileCache.Controls),
		ControlRiskMappings: buildControlRiskMappings(_fileCache.Controls, _fileCache.Risks),
	}
}

/*
buildTemplateFunctionMap → Constructs shared template functions for both HTML templates
*/
func buildTemplateFunctionMap() template.FuncMap {
	return template.FuncMap{
		"inc":                   func(i int) int { return i + 1 },
		"dec":                   func(i int) int { return i - 1 },
		"add":                   func(a, b int) int { return a + b },
		"sub":                   func(a, b int) int { return a - b },
		"split":                 strings.Split,
		"contains":              strings.Contains,
		"extractRiskSection":    extractRiskSection,
		"extractControlSection": extractControlSection,
	}
}

/*
GenerateHTML → Generate final HTML report from processed report data
*/
func GenerateHTML(_fileCache *utilities.FileCache, _reportPaths utilities.ReportPaths) {
	metadataConfig := _fileCache.MetadataConfig()
	templatePath := filepath.Clean(filepath.Join(_reportPaths.TemplatePath, "html", "template.html.tmpl"))

	templateHTML, errTemplateHTML := template.New("template.html.tmpl").Funcs(buildTemplateFunctionMap()).ParseFiles(templatePath)
	if errTemplateHTML != nil {
		utilities.Check(utilities.NewFileSystemError(
			templatePath,
			fmt.Sprintf("failed to parse HTML template: %s", errTemplateHTML.Error()),
			errTemplateHTML,
		))
	}

	outputFileName := fmt.Sprintf("%s_%s.html", metadataConfig.DocumentName, utilities.DocumentVersion)
	createHTML, errCreateHTML := os.Create(outputFileName)
	if errCreateHTML != nil {
		utilities.Check(utilities.NewFileSystemError(
			outputFileName,
			"failed to create HTML output file",
			errCreateHTML,
		))
	}

	defer createHTML.Close()

	if errExecute := templateHTML.Execute(createHTML, buildTemplateData(_fileCache)); errExecute != nil {
		utilities.Check(utilities.NewProcessingError(
			outputFileName,
			fmt.Sprintf("failed to execute template: %s", errExecute.Error()),
		))
	}
}
