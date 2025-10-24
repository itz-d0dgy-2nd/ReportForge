package modifiers

import (
	"ReportForge/engine/utilities"
	"ReportForge/engine/validators"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

/*
ModifySeverity → Updates the FindingSeverity YAML field in markdown files based on calculated severity and renames files
  - Reads markdown file and normalises line endings to Unix format (Thanks to @bstlaurentnz for using windows)
    -- For each file conduct validation for errors → calls ValidateYamlFrontmatter()
  - Parses YAML into utilities.MarkdownYML
  - For "2_findings" files → calculates correct severity from impact/likelihood matrix, updates content and filename
    -- For each file conduct validation for errors → calls ValidateImpactLikelihoodIndex()
    -- For each file conduct validation for errors → calls ValidateSeverityKey()
  - For "3_suggestions" files → updates filename
  - TODO: Clean this up, it can be improved.
*/
func ModifySeverity(_filePath string, _severityAssessment utilities.SeverityAssessmentYML) {
	var calculatedSeverity string
	var newFileName string
	var fileModified bool
	var unprocessedYaml utilities.MarkdownYML

	rawFileContent, errRawFileContent := os.ReadFile(_filePath)
	utilities.ErrorChecker(errRawFileContent)

	rawMarkdownContent := strings.ReplaceAll(string(rawFileContent), "\r\n", "\n")
	regexMatches := utilities.RegexYamlMatch.FindStringSubmatch(rawMarkdownContent)

	validators.ValidateYamlFrontmatter(regexMatches, _filePath, &unprocessedYaml)

	if strings.Contains(_filePath, "2_findings") {
		if _severityAssessment.ConductSeverityAssessment {
			impactIndex := slices.Index(_severityAssessment.Impacts, unprocessedYaml.FindingImpact)
			validators.ValidateImpactIndex(impactIndex, "impact", _filePath)

			likelihoodIndex := slices.Index(_severityAssessment.Likelihoods, unprocessedYaml.FindingLikelihood)
			validators.ValidateLikelihoodIndex(likelihoodIndex, "likelihood", _filePath)

			if _severityAssessment.FlipSeverityAssessment {
				calculatedSeverity = _severityAssessment.CalculatedMatrix[likelihoodIndex][impactIndex]
			} else {
				calculatedSeverity = _severityAssessment.CalculatedMatrix[impactIndex][likelihoodIndex]
			}

			severityScaleIndex := slices.Index(_severityAssessment.Severities, calculatedSeverity)
			validators.ValidateSeverityIndex(severityScaleIndex, _filePath)

			if unprocessedYaml.FindingSeverity != calculatedSeverity {
				fileModified = true
				rawMarkdownContent = strings.Replace(rawMarkdownContent, "FindingSeverity: "+unprocessedYaml.FindingSeverity, "FindingSeverity: "+calculatedSeverity, 1)
			}

			newFileName = strconv.Itoa(severityScaleIndex) + "_" + unprocessedYaml.FindingName + ".md"

		} else {
			severityScaleIndex := slices.Index(_severityAssessment.Severities, calculatedSeverity)
			validators.ValidateSeverityIndex(severityScaleIndex, _filePath)

			newFileName = strconv.Itoa(severityScaleIndex) + "_" + unprocessedYaml.FindingName + ".md"
		}
	}

	if strings.Contains(_filePath, "3_suggestions") {
		newFileName = "4_" + unprocessedYaml.SuggestionName + ".md"
	}

	finalFilePath := filepath.Clean(filepath.Join(filepath.Dir(_filePath), newFileName))

	if fileModified {
		utilities.ErrorChecker(os.WriteFile(_filePath, []byte(rawMarkdownContent), 0644))
	}

	if _filePath != finalFilePath {
		utilities.ErrorChecker(os.Rename(_filePath, finalFilePath))
	}
}

/*
ModifyIdentifiers → Updates the FindingID and SuggestionID YAML fields based on unique identifiers
  - Reads markdown file and normalises line endings to Unix format (Thanks to @bstlaurentnz for using windows)
    -- For each file conduct validation for errors → calls ValidateYamlFrontmatter()
  - Parses YAML into utilities.MarkdownYML
  - For "2_findings" or "3_suggestions" files → assigns Finding/Suggestions fields
    -- In Draft status → overwrites FindingID and SuggestionID if locked fields are false
    -- In Release status → sets all lock fields to true
  - TODO: Clean this up, it can be improved.
*/
func ModifyIdentifiers(_filePath, _identifierPrefix string, _identifierCounter *int, _isRelease bool) {
	var fileModified bool
	var unprocessedYaml utilities.MarkdownYML
	var identifier string
	var identifierLocked bool
	var identifierField string
	var identifierLockField string

	rawFileContent, errRawFileContent := os.ReadFile(_filePath)
	utilities.ErrorChecker(errRawFileContent)

	rawMarkdownContent := strings.ReplaceAll(string(rawFileContent), "\r\n", "\n")
	regexMatches := utilities.RegexYamlMatch.FindStringSubmatch(rawMarkdownContent)

	validators.ValidateYamlFrontmatter(regexMatches, _filePath, &unprocessedYaml)

	if strings.Contains(_filePath, "2_findings") {
		identifier = strings.TrimSpace(unprocessedYaml.FindingID)
		identifierLocked = unprocessedYaml.FindingIDLocked
		identifierField = "FindingID"
		identifierLockField = "FindingIDLocked"

	} else {
		identifier = strings.TrimSpace(unprocessedYaml.SuggestionID)
		identifierLocked = unprocessedYaml.SuggestionIDLocked
		identifierField = "SuggestionID"
		identifierLockField = "SuggestionIDLocked"

	}

	if identifierLocked {
		if identifier != "" && strings.HasPrefix(identifier, _identifierPrefix) {
			var lockedIDNumber int
			if _, err := fmt.Sscanf(strings.TrimPrefix(identifier, _identifierPrefix), "%d", &lockedIDNumber); err == nil {
				if lockedIDNumber > *_identifierCounter {
					*_identifierCounter = lockedIDNumber
				}
			}
		}

		if _isRelease && strings.Contains(rawMarkdownContent, identifierLockField+": false") {
			fileModified = true
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+": false", identifierLockField+": true", 1)
		}

	} else {
		*_identifierCounter++
		generatedIdentifier := fmt.Sprintf("%s%d", _identifierPrefix, *_identifierCounter)
		fileModified = true

		if identifier == "" {
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierField+":", identifierField+": "+generatedIdentifier, 1)
		} else {
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierField+": "+identifier, identifierField+": "+generatedIdentifier, 1)
		}

		if _isRelease {
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+": false", identifierLockField+": true", 1)
		}
	}

	if fileModified {
		utilities.ErrorChecker(os.WriteFile(_filePath, []byte(rawMarkdownContent), 0644))
	}
}
