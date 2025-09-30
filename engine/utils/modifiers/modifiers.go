package modifiers

import (
	Utils "ReportForge/engine/utils"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

/*
ModifySeverity → Updates the FindingSeverity YAML feild in markdown files based on calculated severity and renames files
  - Validates severity assessment is enabled, returns early if disabled
  - Reads markdown file and normalizes line endings to Unix format (Thanks to @bstlaurentnz for using windows)
  - Parses YAML into Utils.MarkdownYML
  - For "2_findings" files → calculates correct severity from impact/likelihood matrix, updates content and filename
  - For "3_suggestions" files → updates filename
*/
func ModifySeverity(_filePath string, _severityAssessment Utils.SeverityAssessmentYML) {
	newFileName := ""
	fileModified := false
	unprocessedYaml := Utils.MarkdownYML{}

	readMarkdownFile, errReadMarkdown := os.ReadFile(_filePath)
	Utils.ErrorChecker(errReadMarkdown)

	rawMarkdownContent := strings.ReplaceAll(string(readMarkdownFile), "\r\n", "\n")
	regexMatches := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`).FindStringSubmatch(rawMarkdownContent)

	if len(regexMatches) < 2 {
		Utils.ErrorChecker(fmt.Errorf("missing YAML frontmatter in ( %s )", _filePath))
	}

	if yaml.Unmarshal([]byte(regexMatches[1]), &unprocessedYaml) != nil {
		Utils.ErrorChecker(fmt.Errorf("invalid YAML frontmatter in ( %s )", _filePath))
	}

	if strings.Contains(_filePath, "2_findings") {
		if _severityAssessment.ConductSeverityAssessment {

			impactIndex := slices.Index(_severityAssessment.Impacts, unprocessedYaml.FindingImpact)
			likelihoodIndex := slices.Index(_severityAssessment.Likelihoods, unprocessedYaml.FindingLikelihood)

			if impactIndex == -1 || likelihoodIndex == -1 {
				Utils.ErrorChecker(fmt.Errorf("invalid impact or likelihood found in '%s'", _filePath))
			}

			if unprocessedYaml.FindingSeverity != _severityAssessment.CalculatedMatrix[impactIndex][likelihoodIndex] {
				fileModified = true
				rawMarkdownContent = strings.Replace(rawMarkdownContent, "FindingSeverity: "+unprocessedYaml.FindingSeverity, "FindingSeverity: "+_severityAssessment.CalculatedMatrix[impactIndex][likelihoodIndex], 1)
				newFileName = _severityAssessment.Scales[_severityAssessment.CalculatedMatrix[impactIndex][likelihoodIndex]] + "_" + unprocessedYaml.FindingName + ".md"
			}

		} else if !_severityAssessment.ConductSeverityAssessment {

			if _, exists := _severityAssessment.Scales[unprocessedYaml.FindingSeverity]; unprocessedYaml.FindingSeverity == "" || !exists {
				Utils.ErrorChecker(fmt.Errorf("invalid severity found in '%s'", _filePath))
			}

			fileModified = true
			newFileName = _severityAssessment.Scales[unprocessedYaml.FindingSeverity] + "_" + unprocessedYaml.FindingName + ".md"
		}
	}

	if strings.Contains(_filePath, "3_suggestions") {
		fileModified = true
		newFileName = "5_" + unprocessedYaml.SuggestionName + ".md"
	}

	if fileModified {
		Utils.ErrorChecker(os.WriteFile(_filePath, []byte(rawMarkdownContent), 0644))
		Utils.ErrorChecker(os.Rename(_filePath, filepath.Clean(filepath.Join(filepath.Dir(_filePath), newFileName))))
	}
}

/*
ModifyIdentifiers → Updates the FindingID and SuggestionID YAML fields based on unique identifiers
  - Reads markdown file and normalizes line endings to Unix format (Thanks to @bstlaurentnz for using windows)
  - Parses YAML into Utils.MarkdownYML
  - For "2_findings" files → generates FindingID if empty using prefix and counter
  - For "3_suggestions" files → generates SuggestionID if empty using prefix and counter
  - TODO: Track ID state to avoid conflicts
*/
func ModifyIdentifiers(_filePath, _identifierPrefix string, _identifierCounter *int) {
	fileModified := false
	unprocessedYaml := Utils.MarkdownYML{}

	readMarkdownFile, errReadMarkdown := os.ReadFile(_filePath)
	Utils.ErrorChecker(errReadMarkdown)

	rawMarkdownContent := strings.ReplaceAll(string(readMarkdownFile), "\r\n", "\n")
	regexMatches := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`).FindStringSubmatch(rawMarkdownContent)

	if len(regexMatches) < 2 {
		Utils.ErrorChecker(fmt.Errorf("missing YAML frontmatter '%s'", _filePath))
	}

	if yaml.Unmarshal([]byte(regexMatches[1]), &unprocessedYaml) != nil {
		Utils.ErrorChecker(fmt.Errorf("invalid YAML frontmatter '%s'", _filePath))
	}

	*_identifierCounter++

	if strings.Contains(_filePath, "2_findings") && strings.TrimSpace(unprocessedYaml.FindingID) == "" {
		fileModified = true
		generatedIdentifier := fmt.Sprintf("%s%d", _identifierPrefix, *_identifierCounter)
		rawMarkdownContent = strings.Replace(rawMarkdownContent, "FindingID:", "FindingID: "+generatedIdentifier, 1)
	}

	if strings.Contains(_filePath, "3_suggestions") && strings.TrimSpace(unprocessedYaml.SuggestionID) == "" {
		fileModified = true
		generatedIdentifier := fmt.Sprintf("%s%d", _identifierPrefix, *_identifierCounter)
		rawMarkdownContent = strings.Replace(rawMarkdownContent, "SuggestionID:", "SuggestionID: "+generatedIdentifier, 1)
	}

	if fileModified {
		Utils.ErrorChecker(os.WriteFile(_filePath, []byte(rawMarkdownContent), 0644))
	}
}
