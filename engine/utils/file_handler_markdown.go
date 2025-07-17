package Utils

import (
	"os"
	"path/filepath"
)

func FileHandlerMarkdown(frontmatter FrontmatterYML, matrix SeverityMatrix, directory string) []Markdown {

	processedMD := []Markdown{}

	readDirectoryContents, ErrReadDirectoryContents := os.ReadDir(directory)
	ErrorChecker(ErrReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := directory + "/" + directoryContents.Name()
			readFiles, ErrReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(ErrReadFiles)

			for _, subdirectoryContents := range readFiles {
				if filepath.Ext(subdirectoryContents.Name()) == ".md" {
					ProcessMarkdown(frontmatter, matrix, subdirectory, subdirectoryContents, &processedMD)
				}
			}

		} else if !directoryContents.IsDir() {
			if filepath.Ext(directoryContents.Name()) == ".md" {
				ProcessMarkdown(frontmatter, matrix, directory, directoryContents, &processedMD)
			}
		}
	}

	return processedMD

}
