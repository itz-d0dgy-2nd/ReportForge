package handlers

import (
	"ReportForge/engine/utils"
	"ReportForge/engine/utils/processors"
	"os"
	"path/filepath"
)

func MarkdownRecursiveScan(_reportPath, _directory string, processedMD *[]Utils.Markdown, _frontmatter Utils.FrontmatterYML, _severityAssessment Utils.SeverityAssessmentYML) {

	// Read directory structure and process contents
	//   - Subdirectories: Recursively enter directory
	//   - Markdown: Process .md files

	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	Utils.ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))

		if directoryContents.IsDir() {
			MarkdownRecursiveScan(_reportPath, subdirectory, processedMD, _frontmatter, _severityAssessment)

		} else if filepath.Ext(directoryContents.Name()) == ".md" {
			processors.ProcessMarkdown(_reportPath, _directory, directoryContents, processedMD, _frontmatter, _severityAssessment)

		}
	}
}

/*
MarkdownFileHandler → Handles markdown files
  - Calls MarkdownRecursiveScan() → processors.ProcessMarkdown()
  - Returns processed markdown of type []Utils.Markdown{}
*/
func MarkdownFileHandler(_reportPath string, _directory string, _frontmatter Utils.FrontmatterYML, _severityAssessment Utils.SeverityAssessmentYML) []Utils.Markdown {
	processedMD := []Utils.Markdown{}
	MarkdownRecursiveScan(_reportPath, _directory, &processedMD, _frontmatter, _severityAssessment)
	return processedMD
}
