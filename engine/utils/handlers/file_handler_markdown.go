package handlers

import (
	"ReportForge/engine/utils"
	"os"
	"path/filepath"
)

/*
MarkdownFileHandler function -> Handles markdown files

	XXX
*/
func MarkdownFileHandler(_reportPath string, _directory string, _frontmatter Utils.FrontmatterYML, _severityAssessment Utils.SeverityAssessmentYML) []Utils.Markdown {
	processedMD := []Utils.Markdown{}
	MarkdownRecursiveScan(_reportPath, _directory, &processedMD, _frontmatter, _severityAssessment)
	return processedMD
}

func MarkdownRecursiveScan(_reportPath, _directory string, processedMD *[]Utils.Markdown, _frontmatter Utils.FrontmatterYML, _severityAssessment Utils.SeverityAssessmentYML) {
	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	Utils.ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))

		if directoryContents.IsDir() {
			MarkdownRecursiveScan(_reportPath, subdirectory, processedMD, _frontmatter, _severityAssessment)

		} else if filepath.Ext(directoryContents.Name()) == ".md" {
			Utils.ProcessMarkdown(_reportPath, _directory, directoryContents, processedMD, _frontmatter, _severityAssessment)

		}

	}
}
