package modifiers

import (
	"ReportForge/engine/utilities"
	"ReportForge/engine/validators"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
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

			if _severityAssessment.FlipSeverityMatrix {
				calculatedSeverity = _severityAssessment.CalculatedMatrix[likelihoodIndex][impactIndex]
			} else {
				calculatedSeverity = _severityAssessment.CalculatedMatrix[impactIndex][likelihoodIndex]
			}

			if unprocessedYaml.FindingSeverity != calculatedSeverity {
				fileModified = true
				rawMarkdownContent = strings.Replace(rawMarkdownContent, "FindingSeverity: "+unprocessedYaml.FindingSeverity, "FindingSeverity: "+calculatedSeverity, 1)
			}
		} else {
			calculatedSeverity = unprocessedYaml.FindingSeverity
		}

		severityScaleIndex := slices.Index(_severityAssessment.Severities, calculatedSeverity)
		validators.ValidateSeverityIndex(severityScaleIndex, _filePath)
		newFileName = strconv.Itoa(severityScaleIndex) + "_" + unprocessedYaml.FindingName + ".md"
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
func ModifyIdentifiers(_filePath, _identifierPrefix string, _identifierCounter *int, _documentStatus string) {
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
			identifierSuffix := strings.TrimPrefix(identifier, _identifierPrefix)
			var lockedIdentifierNumber int

			if _, errParseIdentifier := fmt.Sscanf(identifierSuffix, "%d", &lockedIdentifierNumber); errParseIdentifier == nil {
				if lockedIdentifierNumber > *_identifierCounter {
					*_identifierCounter = lockedIdentifierNumber
				}
			}
		}

		if _documentStatus == "Release" && strings.Contains(rawMarkdownContent, identifierLockField+": false") {
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

		if _documentStatus == "Release" {
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+": false", identifierLockField+": true", 1)
		}
	}

	if fileModified {
		utilities.ErrorChecker(os.WriteFile(_filePath, []byte(rawMarkdownContent), 0644))
	}
}

/*
ModifyImage → Compresses images when document status is Release
  - Backs up original images to originals/ subdirectory
  - Creates JPEG compressed version (Quality 75) at original path
  - Skips processing if already compressed (backup exists) or file is in originals/ directory
*/
func ModifyImage(_filePath string, _documentStatus string) {
	if _documentStatus != "Release" {
		return
	}

	if filepath.Base(filepath.Dir(_filePath)) == "originals" {
		return
	}

	originalsDirectory := filepath.Join(filepath.Dir(_filePath), "originals")
	originalBackupPath := filepath.Join(originalsDirectory, filepath.Base(_filePath))

	if _, errStatCheck := os.Stat(originalBackupPath); errStatCheck == nil {
		return
	}

	errMakeDirectory := os.MkdirAll(originalsDirectory, 0755)
	utilities.ErrorChecker(errMakeDirectory)

	rawFileContent, errRawFileContent := os.Open(_filePath)
	utilities.ErrorChecker(errRawFileContent)
	defer rawFileContent.Close()

	decodedImage, _, errDecodedImage := image.Decode(rawFileContent)
	utilities.ErrorChecker(errDecodedImage)

	temporaryFilePath := _filePath + ".tmp"
	compressedFile, errCompressedFile := os.Create(temporaryFilePath)
	utilities.ErrorChecker(errCompressedFile)

	errEncodeImage := jpeg.Encode(compressedFile, decodedImage, &jpeg.Options{Quality: 75})
	compressedFile.Close()
	utilities.ErrorChecker(errEncodeImage)

	errRenameOriginalFile := os.Rename(_filePath, originalBackupPath)
	utilities.ErrorChecker(errRenameOriginalFile)

	errRenameCompressedFile := os.Rename(temporaryFilePath, _filePath)
	if errRenameCompressedFile != nil {
		os.Rename(originalBackupPath, _filePath)
		utilities.ErrorChecker(errRenameCompressedFile)
	}
}
