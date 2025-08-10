package Utils

import (
	"os"
	"path/filepath"
)

func FileHandlerSeverity(_severityAssessment SeverityAssessmentYML, _directory string) SeverityAssessmentYML {

	// [^DIA]: https://www.digital.govt.nz/standards-and-guidance/privacy-security-and-risk/risk-management/risk-assessments/analyse/initial-risk-ratings#table-1

	readDirectoryContents, ErrReadDirectoryContents := os.ReadDir(_directory)
	ErrorChecker(ErrReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := _directory + "/" + directoryContents.Name()
			readFiles, ErrReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(ErrReadFiles)

			for _, subdirectoryContents := range readFiles {
				if filepath.Ext(subdirectoryContents.Name()) == ".md" {
					ProcessSeverityMatrix(subdirectory, subdirectoryContents, &_severityAssessment)
				}
			}

		} else if !directoryContents.IsDir() {
			if filepath.Ext(directoryContents.Name()) == ".md" {
				ProcessSeverityMatrix(_directory, directoryContents, &_severityAssessment)
			}
		}
	}

	// To do: Add error handling. What if the files dont exist?

	return _severityAssessment
}
