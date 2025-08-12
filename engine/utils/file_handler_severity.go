package Utils

import (
	"os"
	"path/filepath"
)

/*
FileHandlerSeverity function -> Handles markdown files

	XXX
*/
func FileHandlerSeverity(_severityAssessment SeverityAssessmentYML, _directory string) SeverityAssessmentYML {

	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))
			readFiles, errReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(errReadFiles)

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

	return _severityAssessment

}
