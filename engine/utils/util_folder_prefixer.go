package Utils

import (
	"os"
	"path/filepath"
)

func SetupPrefixMap(_directory string) map[string]string {
	identifierPrefixMap := make(map[string]string)
	currentIndex := 0
	for _, subdirectoryName := range []string{"2_findings", "3_suggestions"} {
		subdirectoryPath := filepath.Join(_directory, subdirectoryName)
		if subdirectoryContents, errReadDirectory := os.ReadDir(subdirectoryPath); errReadDirectory == nil {
			for _, directoryEntry := range subdirectoryContents {
				if directoryEntry.IsDir() {
					identifierPrefixMap[filepath.Join(subdirectoryPath, directoryEntry.Name())] = string('A' + currentIndex)
					currentIndex++
				}
			}
		}
	}
	return identifierPrefixMap
}
