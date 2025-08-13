package handlers

import (
	"ReportForge/engine/utils"
	"ReportForge/engine/utils/processors"
	"os"
	"path/filepath"
)

func SeverityRecursiveScan(_directory string, _severityAssessment *Utils.SeverityAssessmentYML) {

	// Read directory structure and process contents
	//   - Subdirectories: Recursively enter directory
	//   - Markdown: Process .md files

	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	Utils.ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))

		if directoryContents.IsDir() {
			SeverityRecursiveScan(subdirectory, _severityAssessment)

		} else if filepath.Ext(directoryContents.Name()) == ".md" {
			processors.ProcessSeverityMatrix(_directory, directoryContents, _severityAssessment)

		}
	}
}

/*
SeverityFileHandler → Handles markdown files
  - Reads directory structure
  - Filters for .nd files
  - Calls SeverityRecursiveScan() → processors.ProcessSeverityMatrix()
  - Returns processed yaml of type Utils.SeverityAssessmentYML
*/
func SeverityFileHandler(_directory string, _severityAssessment Utils.SeverityAssessmentYML) Utils.SeverityAssessmentYML {
	SeverityRecursiveScan(_directory, &_severityAssessment)
	return _severityAssessment
}
