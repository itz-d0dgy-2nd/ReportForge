package Utils

import (
	"os"
	"path/filepath"
	"strings"
)

func FileHandlerReportConfig(directory string) (FrontmatterYML, SeverityMatrix) {

	processedYML := FrontmatterYML{}
	processedMatrix := SeverityMatrix{}

	readDirectoryContents, ErrReadDirectoryContents := os.ReadDir(directory)
	ErrorChecker(ErrReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := directory + "/" + directoryContents.Name()
			readFiles, ErrReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(ErrReadFiles)

			for _, subdirectoryContents := range readFiles {
				if filepath.Ext(subdirectoryContents.Name()) == ".yml" {
					if strings.Contains(subdirectoryContents.Name(), "frontmatter") {
						ProcessConfigFrontmatter(subdirectory, subdirectoryContents, &processedYML)
					}
					if strings.Contains(subdirectoryContents.Name(), "matrix") {
						ProcessConfigMatrix(subdirectory, subdirectoryContents, &processedMatrix)
					}

				}
			}

		} else if !directoryContents.IsDir() {
			if filepath.Ext(directoryContents.Name()) == ".yml" {
				if strings.Contains(directoryContents.Name(), "frontmatter") {
					ProcessConfigFrontmatter(directory, directoryContents, &processedYML)
				}
				if strings.Contains(directoryContents.Name(), "matrix") {
					ProcessConfigMatrix(directory, directoryContents, &processedMatrix)
				}
			}
		}
	}

	return processedYML, processedMatrix
}
