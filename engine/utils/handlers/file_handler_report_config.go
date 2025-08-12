package handlers

import (
	"ReportForge/engine/utils"
	"os"
	"path/filepath"
	"strings"
)

/*
ReportConfigFileHandler function -> Handles yml files

	XXX
*/
func ReportConfigFileHandler(_directory string) (Utils.FrontmatterYML, Utils.SeverityAssessmentYML) {
	processedYML := Utils.FrontmatterYML{}
	processedSeverityAssessment := Utils.SeverityAssessmentYML{}
	ReportConfigRecursiveScan(_directory, &processedYML, &processedSeverityAssessment)
	return processedYML, processedSeverityAssessment
}

func ReportConfigRecursiveScan(_directory string, processedYML *Utils.FrontmatterYML, processedSeverityAssessment *Utils.SeverityAssessmentYML) {
	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	Utils.ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))

		if directoryContents.IsDir() {
			ReportConfigRecursiveScan(subdirectory, processedYML, processedSeverityAssessment)

		} else if filepath.Ext(directoryContents.Name()) == ".yml" {
			if strings.Contains(directoryContents.Name(), "front_matter") {
				Utils.ProcessConfigFrontmatter(_directory, directoryContents, processedYML)

			}
			if strings.Contains(directoryContents.Name(), "severity_assessment") {
				Utils.ProcessConfigMatrix(_directory, directoryContents, processedSeverityAssessment)

			}
		}
	}
}
