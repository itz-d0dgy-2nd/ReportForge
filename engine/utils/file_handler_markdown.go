package Utils

import (
	"os"
	"path/filepath"
)

func FileHandlerMarkdown(_reportTemplatePath string, _frontMatter FrontmatterYML, _severityAssessment SeverityAssessmentYML, _directory string) []Markdown {

	processedMD := []Markdown{}

	readDirectoryContents, ErrReadDirectoryContents := os.ReadDir(_directory)
	ErrorChecker(ErrReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := filepath.Join(_directory, directoryContents.Name())
			readFiles, ErrReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(ErrReadFiles)

			for _, subdirectoryContents := range readFiles {
				if filepath.Ext(subdirectoryContents.Name()) == ".md" {
					ProcessMarkdown(_reportTemplatePath, _frontMatter, _severityAssessment, subdirectory, subdirectoryContents, &processedMD)
				}
			}

		} else if !directoryContents.IsDir() {
			if filepath.Ext(directoryContents.Name()) == ".md" {
				ProcessMarkdown(_reportTemplatePath, _frontMatter, _severityAssessment, _directory, directoryContents, &processedMD)
			}
		}
	}

	// To do: Add error handling. What if the files dont exist?

	return processedMD

}
