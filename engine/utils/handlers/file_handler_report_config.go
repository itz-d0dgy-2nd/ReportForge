package handlers

import (
	"ReportForge/engine/utils"
	"ReportForge/engine/utils/processors"
	"os"
	"path/filepath"
	"strings"
)

func ReportConfigRecursiveScan(_directory string, processedYML *Utils.FrontmatterYML, processedSeverityAssessment *Utils.SeverityAssessmentYML) {

	// Read directory structure and process contents
	//   - Subdirectories: Recursively enter directory
	//   - Markdown: Process .yml files

	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	Utils.ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))

		if directoryContents.IsDir() {
			ReportConfigRecursiveScan(subdirectory, processedYML, processedSeverityAssessment)

		} else if filepath.Ext(directoryContents.Name()) == ".yml" {
			if strings.Contains(directoryContents.Name(), "front_matter") {
				processors.ProcessConfigFrontmatter(_directory, directoryContents, processedYML)

			}
			if strings.Contains(directoryContents.Name(), "severity_assessment") {
				processors.ProcessConfigMatrix(_directory, directoryContents, processedSeverityAssessment)

			}
		}
	}
}

/*
ReportConfigFileHandler → Handles yaml files
  - Reads directory structure
  - Filters for .yml files
  - Calls ReportConfigRecursiveScan() → processors.ProcessConfigMatrix()
  - Returns processed yaml of type Utils.FrontmatterYML{} & Utils.SeverityAssessmentYML{}
*/
func ReportConfigFileHandler(_directory string) (Utils.FrontmatterYML, Utils.SeverityAssessmentYML) {
	processedYML := Utils.FrontmatterYML{}
	processedSeverityAssessment := Utils.SeverityAssessmentYML{}
	ReportConfigRecursiveScan(_directory, &processedYML, &processedSeverityAssessment)
	return processedYML, processedSeverityAssessment
}
