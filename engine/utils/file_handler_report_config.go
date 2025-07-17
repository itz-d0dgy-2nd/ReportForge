package Utils

import (
	"os"
	"path/filepath"
	"strings"
)

func FileHandlerReportConfig(directory string) (FrontmatterYML, SeverityAssessmentYML) {

	processedYML := FrontmatterYML{}
	processedSeverityAssessment := SeverityAssessmentYML{}

	readDirectoryContents, ErrReadDirectoryContents := os.ReadDir(directory)
	ErrorChecker(ErrReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := directory + "/" + directoryContents.Name()
			readFiles, ErrReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(ErrReadFiles)

			for _, subdirectoryContents := range readFiles {
				if filepath.Ext(subdirectoryContents.Name()) == ".yml" {
					if strings.Contains(subdirectoryContents.Name(), "front_matter") {
						ProcessConfigFrontmatter(subdirectory, subdirectoryContents, &processedYML)
					}
					if strings.Contains(subdirectoryContents.Name(), "severity_assessment") {
						ProcessConfigMatrix(subdirectory, subdirectoryContents, &processedSeverityAssessment)
					}
				}
			}

		} else if !directoryContents.IsDir() {
			if filepath.Ext(directoryContents.Name()) == ".yml" {
				if strings.Contains(directoryContents.Name(), "front_matter") {
					ProcessConfigFrontmatter(directory, directoryContents, &processedYML)
				}
				if strings.Contains(directoryContents.Name(), "severity_assessment") {
					ProcessConfigMatrix(directory, directoryContents, &processedSeverityAssessment)
				}
			}
		}
	}

	// To do: Add error handling. What if the files dont exist?

	return processedYML, processedSeverityAssessment
}
