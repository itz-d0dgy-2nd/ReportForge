package Utils

import (
	"os"
	"path/filepath"
)

/*
FileHandlerMarkdown function -> Handles markdown files
  - Read provided directory contents
  - Iterate over directory structure
  - foreach `.md` file call `ProcessSeverityMatrix()`
*/
func FileHandlerSeverity(_severityAssessment SeverityAssessmentYML, _directory string) SeverityAssessmentYML {

	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := filepath.Join(_directory, directoryContents.Name())
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

	// To do: Add error handling. What if the files dont exist?

	return _severityAssessment
}
