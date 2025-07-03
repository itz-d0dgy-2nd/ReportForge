package Utils

import (
	"os"
	"path/filepath"
)

func SeverityFileHandler(directory string) [5][5]string {

	// [^DIA]: https://www.digital.govt.nz/standards-and-guidance/privacy-security-and-risk/risk-management/risk-assessments/analyse/initial-risk-ratings#table-1
	processedMatrix := SeverityMatrix{
		map[string]int{
			"Severe":      0,
			"Significant": 1,
			"Moderate":    2,
			"Minor":       3,
			"Minimal":     4,
		},

		map[string]int{
			"Almost Never":          0,
			"Possible but Unlikely": 1,
			"Possible":              2,
			"Highly Probable":       3,
			"Almost Certain":        4,
		},

		[5][5]string{
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
			{"", "", "", "", ""},
		},
	}

	readDirectoryContents, ErrReadDirectoryContents := os.ReadDir(directory)
	ErrorChecker(ErrReadDirectoryContents)

	for _, directoryContents := range readDirectoryContents {
		if directoryContents.IsDir() {
			subdirectory := directory + "/" + directoryContents.Name()
			readFiles, ErrReadFiles := os.ReadDir(subdirectory)
			ErrorChecker(ErrReadFiles)

			for _, subdirectoryContents := range readFiles {
				if filepath.Ext(subdirectoryContents.Name()) == ".md" {
					ProcessSeverityMatrix(subdirectory, subdirectoryContents, &processedMatrix)
				}
			}

		} else if !directoryContents.IsDir() {
			if filepath.Ext(directoryContents.Name()) == ".md" {
				ProcessSeverityMatrix(directory, directoryContents, &processedMatrix)
			}
		}
	}

	return processedMatrix.Matrix

}
