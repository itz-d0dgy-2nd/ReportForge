package Utils

import (
	"os"
	"path/filepath"
)

/*
FileHandlerMarkdown function -> Handles markdown files

	XXX
*/
func FileHandlerMarkdown(_reportPath string, _directory string, _frontmatter FrontmatterYML, _severityAssessment SeverityAssessmentYML) []Markdown {

	processedMD := []Markdown{}

	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := filepath.Clean(filepath.Join(_directory, directoryContents.Name()))
			readFiles, errReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(errReadFiles)

			for _, subdirectoryContents := range readFiles {
				if filepath.Ext(subdirectoryContents.Name()) == ".md" {
					ProcessMarkdown(_reportPath, subdirectory, subdirectoryContents, &processedMD, _frontmatter, _severityAssessment)
				}
			}

		} else if !directoryContents.IsDir() {
			if filepath.Ext(directoryContents.Name()) == ".md" {
				ProcessMarkdown(_reportPath, _directory, directoryContents, &processedMD, _frontmatter, _severityAssessment)
			}
		}
	}

	return processedMD

}
