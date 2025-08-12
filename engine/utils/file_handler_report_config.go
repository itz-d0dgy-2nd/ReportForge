package Utils

import (
	"os"
	"path/filepath"
	"strings"
)

/*
ReportConfigFileHandler function -> Handles yml files

	XXX
*/
func ReportConfigFileHandler(_directory string) (FrontmatterYML, SeverityAssessmentYML) {
	processedYML := FrontmatterYML{}
	processedSeverityAssessment := SeverityAssessmentYML{}
	ReportConfigRecursiveScan(_directory, &processedYML, &processedSeverityAssessment)
	return processedYML, processedSeverityAssessment
}

func ReportConfigRecursiveScan(_directory string, processedYML *FrontmatterYML, processedSeverityAssessment *SeverityAssessmentYML) {
	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))

		if directoryContents.IsDir() {
			ReportConfigRecursiveScan(subdirectory, processedYML, processedSeverityAssessment)

		} else if filepath.Ext(directoryContents.Name()) == ".yml" {
			if strings.Contains(directoryContents.Name(), "front_matter") {
				ProcessConfigFrontmatter(_directory, directoryContents, processedYML)

			}
			if strings.Contains(directoryContents.Name(), "severity_assessment") {
				ProcessConfigMatrix(_directory, directoryContents, processedSeverityAssessment)

			}
		}
	}
}
