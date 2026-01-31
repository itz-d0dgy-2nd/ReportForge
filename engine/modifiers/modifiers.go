package modifiers

import (
	"ReportForge/engine/utilities"
	"ReportForge/engine/validators"
	"fmt"
	_ "image/png"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

/*
parseFile → Decorator that wraps file reading with YAML validation
  - Calls utilities.ReadAndCleanMarkdownFile() to read, strip BOM, and normalise line endings
  - Validates YAML frontmatter and populates MarkdownYML struct
  - Returns raw markdown content, regex matches, and parsed YAML
*/
func parseFile(_filePath string, _fileCache *utilities.FileCache) (string, []string, utilities.MarkdownYML) {
	var unprocessedYaml utilities.MarkdownYML

	rawMarkdownContent, regexMatches, err := utilities.ReadAndCleanMarkdownFile(_filePath, _fileCache)
	utilities.ErrorChecker(err)

	validators.ValidateYamlFrontmatter(regexMatches, _filePath, &unprocessedYaml)

	return rawMarkdownContent, regexMatches, unprocessedYaml
}

/*
ModifySeverity → Updates the FindingSeverity YAML field in markdown files based on calculated severity and renames files
  - Reads markdown file and normalises line endings to Unix format (Thanks to @bstlaurentnz for using windows)
  - Calls ValidateYamlFrontmatter()
  - Parses YAML into utilities.MarkdownYML
  - For findings specifically:
    -- If ConductSeverityAssessment
    --- Calls validators.ValidateImpactIndex()
    --- Calls validators.ValidateLikelihoodIndex()
    --- Calculates severity from impact/likelihood matrix and Updates FindingSeverity
    -- Validates severity index using validators.ValidateSeverityIndex()
    -- Constructs new filename with severity prefix (e.g., "0_finding_name.md")
  - For suggestions specifically:
    -- Constructs new filename with standard prefix (e.g., "5_suggestion_name.md")
  - Writes modified content to file if changes were made
  - Renames file if filename changed, updating file cache accordingly
  - Handles errors via utilities.ErrorChecker()
*/
func ModifySeverity(_filePath string, _fileCache *utilities.FileCache) string {
	var calculatedSeverity string
	var newFileName string
	var fileModified bool

	rawMarkdownContent, _, unprocessedYaml := parseFile(_filePath, _fileCache)
	directoryType := utilities.GetDirectoryType(_filePath)

	switch directoryType {
	case utilities.FindingsDirectory:
		if _fileCache.SeverityConfig.ConductSeverityAssessment {
			impactIndex := slices.Index(_fileCache.SeverityConfig.Impacts, unprocessedYaml.FindingImpact)
			validators.ValidateImpactIndex(impactIndex, _filePath)

			likelihoodIndex := slices.Index(_fileCache.SeverityConfig.Likelihoods, unprocessedYaml.FindingLikelihood)
			validators.ValidateLikelihoodIndex(likelihoodIndex, _filePath)

			if _fileCache.SeverityConfig.SwapImpactLikelihoodAxis {
				calculatedSeverity = _fileCache.SeverityConfig.CalculatedMatrix[likelihoodIndex][impactIndex]
			} else {
				calculatedSeverity = _fileCache.SeverityConfig.CalculatedMatrix[impactIndex][likelihoodIndex]
			}

			if unprocessedYaml.FindingSeverity != calculatedSeverity {
				fileModified = true
				rawMarkdownContent = utilities.RegexFindingSeverity.ReplaceAllString(rawMarkdownContent, "FindingSeverity: "+calculatedSeverity)
			}
		} else {
			calculatedSeverity = unprocessedYaml.FindingSeverity
		}

		severityScaleIndex := slices.Index(_fileCache.SeverityConfig.Severities, calculatedSeverity)
		validators.ValidateSeverityIndex(severityScaleIndex, _filePath)
		newFileName = strconv.Itoa(severityScaleIndex) + "_" + unprocessedYaml.FindingName + ".md"

	case utilities.SuggestionsDirectory:
		newFileName = "5_" + unprocessedYaml.SuggestionName + ".md"
	}

	if fileModified {
		utilities.ErrorChecker(os.WriteFile(_filePath, []byte(rawMarkdownContent), 0644))
	}

	newFilePath := filepath.Clean(filepath.Join(filepath.Dir(_filePath), newFileName))

	if _filePath != newFilePath {
		if os.Rename(_filePath, newFilePath) != nil {
			utilities.ErrorChecker(fmt.Errorf("cannot rename %s to %s (possible duplicate finding name/severity)", filepath.Base(_filePath), filepath.Base(newFilePath)))
		}
		_fileCache.RenameFile(_filePath, newFilePath, []byte(rawMarkdownContent))
		return newFilePath

	} else if fileModified {
		_fileCache.UpdateFile(_filePath, []byte(rawMarkdownContent))
	}

	return _filePath
}

/*
ModifyIdentifiers → Updates the FindingID, SuggestionID, RiskID YAML fields based on unique identifiers
  - Reads markdown file and normalises line endings to Unix format (Thanks to @bstlaurentnz for using windows)
  - Calls ValidateYamlFrontmatter()
  - Parses YAML into utilities.MarkdownYML
  - Determines appropriate identifier field and lock field based on directory:
    -- Findings: FindingID and FindingIDLocked
    -- Suggestions: SuggestionID and SuggestionIDLocked
    -- Risks: RiskID and RiskIDLocked
  - Handles identifier assignment based on document status:
    -- Draft status with unlocked identifier: Assigns new identifier using prefix and reserved ID
    -- Draft status with locked identifier: Assigns new identifier using prefix and reserved ID to unlocked field
    -- Release status: Assigns new identifier using prefix and reserved ID, Sets lock field to true for all identifiers
  - Writes modified content to file and updates file cache
  - Handles errors via utilities.ErrorChecker()
*/
func ModifyIdentifiers(_filePath, _identifierPrefix string, _reservedID int32, _fileCache *utilities.FileCache) {
	var unprocessedYaml utilities.MarkdownYML
	var identifier string
	var identifierLocked bool
	var identifierField string
	var identifierLockField string

	rawMarkdownContent, _, unprocessedYaml := parseFile(_filePath, _fileCache)
	directoryType := utilities.GetDirectoryType(_filePath)

	switch directoryType {
	case utilities.FindingsDirectory:
		identifier = strings.TrimSpace(unprocessedYaml.FindingID)
		identifierLocked = unprocessedYaml.FindingIDLocked
		identifierField = "FindingID"
		identifierLockField = "FindingIDLocked"
	case utilities.SuggestionsDirectory:
		identifier = strings.TrimSpace(unprocessedYaml.SuggestionID)
		identifierLocked = unprocessedYaml.SuggestionIDLocked
		identifierField = "SuggestionID"
		identifierLockField = "SuggestionIDLocked"
	case utilities.RisksDirectory:
		identifier = strings.TrimSpace(unprocessedYaml.RiskID)
		identifierLocked = unprocessedYaml.RiskIDLocked
		identifierField = "RiskID"
		identifierLockField = "RiskIDLocked"
	}

	if identifierLocked {
		if utilities.DocumentStatus != "Release" {
			return
		}

		if !strings.Contains(rawMarkdownContent, identifierLockField+": false") {
			return
		}

		rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+": false", identifierLockField+": true", 1)
		utilities.ErrorChecker(os.WriteFile(_filePath, []byte(rawMarkdownContent), 0644))
		_fileCache.UpdateFile(_filePath, []byte(rawMarkdownContent))
		return
	}

	generatedIdentifier := fmt.Sprintf("%s%d", _identifierPrefix, _reservedID)

	if identifier == "" {
		rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierField+":", identifierField+": "+generatedIdentifier, 1)
	} else {
		rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierField+": "+identifier, identifierField+": "+generatedIdentifier, 1)
	}

	if utilities.DocumentStatus == "Release" {
		if strings.Contains(rawMarkdownContent, identifierLockField+": false") {
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+": false", identifierLockField+": true", 1)
		} else if strings.Contains(rawMarkdownContent, identifierLockField+":") {
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+":", identifierLockField+": true", 1)
		}
	}

	utilities.ErrorChecker(os.WriteFile(_filePath, []byte(rawMarkdownContent), 0644))
	_fileCache.UpdateFile(_filePath, []byte(rawMarkdownContent))
}
