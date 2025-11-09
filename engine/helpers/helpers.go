package helpers

import (
	"ReportForge/engine/utilities"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

/*
- TODO: Remove this function? There must be a better way of prefixing rather then a double pass.
-- So inefficient...
*/
func PrefixMapHelper(_directory string) map[string]string {
	identifierPrefixMap := make(map[string]string)
	currentIndex := 0

	for _, subdirectoryName := range []string{"2_findings", "3_suggestions"} {
		subdirectoryPath := filepath.Join(_directory, subdirectoryName)
		subdirectoryContents, errReadDirectory := os.ReadDir(subdirectoryPath)

		if errReadDirectory != nil {
			continue
		}

		for _, directoryEntry := range subdirectoryContents {
			if directoryEntry.IsDir() {
				if currentIndex >= 26 {
					utilities.ErrorChecker(fmt.Errorf("too many subdirectories: maximum 26 allowed (A-Z)"))
				}
				identifierPrefixMap[filepath.Join(subdirectoryPath, directoryEntry.Name())] = string(rune('A' + currentIndex))
				currentIndex++
			}
		}
	}

	return identifierPrefixMap
}

/*
- TODO: Remove this function! There must be a better way of tracking locked rather then a double pass.
-- So inefficient...
*/
func TrackedLockedHelper(_directory string, _identifierPrefixMap map[string]string, _identifierCounterMap map[string]*int) {
	filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" {
			return errAnonymousFunction
		}

		if !strings.Contains(filePath, "2_findings") && !strings.Contains(filePath, "3_suggestions") {
			return nil
		}

		var unprocessedYaml utilities.MarkdownYML
		readMarkdownContent, _ := os.ReadFile(filePath)
		rawMarkdownContent := strings.ReplaceAll(string(readMarkdownContent), "\r\n", "\n")
		regexMatches := utilities.RegexYamlMatch.FindStringSubmatch(rawMarkdownContent)

		if len(regexMatches) < 2 || yaml.Unmarshal([]byte(regexMatches[1]), &unprocessedYaml) != nil {
			return nil
		}

		identifierPrefix := _identifierPrefixMap[filepath.Dir(filePath)]
		identifier := unprocessedYaml.SuggestionID
		identifierLocked := unprocessedYaml.SuggestionIDLocked

		if strings.Contains(filePath, "2_findings") {
			identifier = unprocessedYaml.FindingID
			identifierLocked = unprocessedYaml.FindingIDLocked
		}

		if identifierLocked && strings.HasPrefix(strings.TrimSpace(identifier), identifierPrefix) {
			var identifierNumber int
			fmt.Sscanf(strings.TrimPrefix(strings.TrimSpace(identifier), identifierPrefix), "%d", &identifierNumber)
			if identifierNumber > *_identifierCounterMap[identifierPrefix] {
				*_identifierCounterMap[identifierPrefix] = identifierNumber
			}
		}

		return nil
	})
}
