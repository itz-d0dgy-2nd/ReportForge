package Utils

import (
	"os"
	"path/filepath"
)

/*
FileHandlerMarkdown function -> Handles markdown files
  - Instantiate a variable of type `[]Markdown`
  - Read provided directory contents
  - Iterate over directory structure
  - foreach `.md` file call `ProcessMarkdown()`
  - Return variable of type `[]Markdown`
*/
func FileHandlerMarkdown(_reportPath string, _frontmatter FrontmatterYML, _severityAssessment SeverityAssessmentYML, _directory string) []Markdown {

	processedMD := []Markdown{}

	readDirectoryContents, errReadDirectoryContents := os.ReadDir(_directory)
	ErrorChecker(errReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := filepath.Join(_directory, directoryContents.Name())
			readFiles, errReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(errReadFiles)

			for _, subdirectoryContents := range readFiles {
				if filepath.Ext(subdirectoryContents.Name()) == ".md" {
					ProcessMarkdown(_reportPath, _frontmatter, _severityAssessment, subdirectory, subdirectoryContents, &processedMD)
				}
			}

		} else if !directoryContents.IsDir() {
			if filepath.Ext(directoryContents.Name()) == ".md" {
				ProcessMarkdown(_reportPath, _frontmatter, _severityAssessment, _directory, directoryContents, &processedMD)
			}
		}
	}

	// To do: Add error handling. What if the files dont exist?

	return processedMD

}
