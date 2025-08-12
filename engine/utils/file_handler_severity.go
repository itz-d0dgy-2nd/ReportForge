package Utils

import (
	"os"
	"path/filepath"
)

/*
SeverityFileHandler function -> Handles markdown files

	XXX
*/
func SeverityFileHandler(_directory string, _severityAssessment SeverityAssessmentYML) SeverityAssessmentYML {
	SeverityRecursiveScan(_directory, &_severityAssessment)
	return _severityAssessment
}

func SeverityRecursiveScan(_directory string, _severityAssessment *SeverityAssessmentYML) {
	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))

		if directoryContents.IsDir() {
			SeverityRecursiveScan(subdirectory, _severityAssessment)

		} else if filepath.Ext(directoryContents.Name()) == ".md" {
			ProcessSeverityMatrix(_directory, directoryContents, _severityAssessment)

		}
	}
}
