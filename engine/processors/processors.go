package processors

import (
	"ReportForge/engine/utilities"
	"ReportForge/engine/validators"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

/*
ProcessConfigMetadata → Process YAML configuration file for report metadata
  - Reads YAML configuration file from specified file path
  - Unmarshals YAML content into pointer *utilities.MetadataYML struct
    -- For each file conduct validation for errors → calls validators.ValidateConfigMetadata()
*/
func ProcessConfigMetadata(_filePath string, _metadata *utilities.MetadataYML) {
	rawFileContent, errRawFileContent := os.ReadFile(_filePath)
	utilities.ErrorChecker(errRawFileContent)

	errDecodeYML := yaml.Unmarshal(rawFileContent, &_metadata)
	utilities.ErrorChecker(errDecodeYML)

	validators.ValidateConfigMetadata(_metadata, _filePath)
}

/*
ProcessConfigSeverityAssessment → Process YAML configuration file for report severity assessment
  - Reads YAML configuration file from specified file path
  - Unmarshals YAML content into pointer *utilities.SeverityAssessmentYML struct
    -- For each file handles validation errors → calls validators.ValidateConfigSeverityAssessment()
*/
func ProcessConfigSeverityAssessment(_filePath string, _severityAssessment *utilities.SeverityAssessmentYML) {
	rawFileContent, errRawFileContent := os.ReadFile(_filePath)
	utilities.ErrorChecker(errRawFileContent)

	errDecodeYML := yaml.Unmarshal(rawFileContent, &_severityAssessment)
	utilities.ErrorChecker(errDecodeYML)

	validators.ValidateConfigSeverityAssessment(_severityAssessment, _filePath)
}

/*
ProcessSeverityMatrix → Process markdown files with YAML frontmatter for report severity assessment
  - Reads markdown file from specified file path
    -- For each file handles validation errors → calls validators.ValidateYamlFrontmatter()
  - Calculates matrix indices based on finding impact and likelihood values
  - Updates severity assessment matrix with finding ID at calculated position
  - Handles errors via utilities.ErrorChecker()
*/
func ProcessSeverityMatrix(_filePath string, _severityAssessment *utilities.SeverityAssessmentYML) {
	if _severityAssessment.ConductSeverityAssessment {
		var unprocessedYaml utilities.MarkdownYML

		rawFileContent, errRawFileContent := os.ReadFile(_filePath)
		utilities.ErrorChecker(errRawFileContent)

		rawMarkdownContent := strings.ReplaceAll(string(rawFileContent), "\r\n", "\n")
		regexMatches := utilities.RegexYamlMatch.FindStringSubmatch(rawMarkdownContent)

		validators.ValidateYamlFrontmatter(regexMatches, _filePath, &unprocessedYaml)

		impactIndex := slices.Index(_severityAssessment.Impacts, unprocessedYaml.FindingImpact)
		likelihoodIndex := slices.Index(_severityAssessment.Likelihoods, unprocessedYaml.FindingLikelihood)

		rowIndex, columnIndex := impactIndex, likelihoodIndex
		if _severityAssessment.FlipSeverityAssessment {
			rowIndex, columnIndex = likelihoodIndex, impactIndex
		}

		if _severityAssessment.Matrix[rowIndex][columnIndex] == "" {
			_severityAssessment.Matrix[rowIndex][columnIndex] = unprocessedYaml.FindingID
		} else {
			_severityAssessment.Matrix[rowIndex][columnIndex] += ", " + unprocessedYaml.FindingID
		}
	}
}

/*
ProcessMarkdown → Process markdown files with YAML frontmatter for report content
  - Reads markdown file from specified file path
    -- For each file handles validation errors → calls validators.ValidateYamlFrontmatter()
  - Converts markdown content to HTML using blackfriday processor
  - Performs strings checks and replace on:
    -- Custom Tags (<qa></qa>, <retest_fixed></retest_fixed>, <retest_not_fixed></retest_not_fixed>)
    -- Custom Tokens (!Client, !TargetAsset0, !TargetAsset1)
    -- Updates screenshot paths to include full report path
  - Thread-safely appends processed markdown to slice using mutex lock
  - Handles errors via utilities.ErrorChecker()
*/
func ProcessMarkdown(_reportPath string, _filePath string, _processedMarkdown *[]utilities.Markdown, _metadata utilities.MetadataYML, _mutexLock *sync.Mutex) {
	var unprocessedYaml utilities.MarkdownYML

	rawFileContent, errRawFileContent := os.ReadFile(_filePath)
	utilities.ErrorChecker(errRawFileContent)

	rawMarkdownContent := strings.ReplaceAll(string(rawFileContent), "\r\n", "\n")
	regexMatches := utilities.RegexYamlMatch.FindStringSubmatch(rawMarkdownContent)

	validators.ValidateYamlFrontmatter(regexMatches, _filePath, &unprocessedYaml)

	unprocessedMarkdown := string(blackfriday.Run([]byte(regexMatches[2])))

	unprocessedMarkdown = utilities.RegexTokenMatch.ReplaceAllStringFunc(unprocessedMarkdown, func(tokenMatch string) string {
		if tokenValue, exists := _metadata.CustomVariables[strings.TrimPrefix(tokenMatch, "!")]; exists {
			return tokenValue
		}

		if strings.TrimPrefix(tokenMatch, "!") == "Client" {
			return _metadata.Client
		}

		return tokenMatch
	})

	unprocessedMarkdown = utilities.RegexMarkdownImageMatchScale.ReplaceAllString(unprocessedMarkdown, `$1 src="`+_reportPath+`/$2"$3 style="$4"/></p>`)
	unprocessedMarkdown = utilities.RegexMarkdownImageMatch.ReplaceAllString(unprocessedMarkdown, `$1 src="`+_reportPath+`/$2"$3/></p>`)

	if strings.Contains(unprocessedMarkdown, "<qa>") {
		fmt.Printf("::warning:: %s: %d QA Comment Present In File \n", _filePath, strings.Count(unprocessedMarkdown, "<qa>"))
	}

	if strings.Contains(unprocessedMarkdown, "<p><retest_fixed></p>") {
		unprocessedMarkdown = strings.ReplaceAll(unprocessedMarkdown, "<p><retest_fixed></p>", "<retest_fixed>")
		unprocessedMarkdown = strings.ReplaceAll(unprocessedMarkdown, "<p></retest_fixed></p>", "</retest_fixed>")
	}

	if strings.Contains(unprocessedMarkdown, "<p><retest_not_fixed></p>") {
		unprocessedMarkdown = strings.ReplaceAll(unprocessedMarkdown, "<p><retest_not_fixed></p>", "<retest_not_fixed>")
		unprocessedMarkdown = strings.ReplaceAll(unprocessedMarkdown, "<p></retest_not_fixed></p>", "</retest_not_fixed>")
	}

	_mutexLock.Lock()
	*_processedMarkdown = append(*_processedMarkdown, utilities.Markdown{
		Directory: filepath.Base(filepath.Dir(_filePath)),
		FileName:  filepath.Base(_filePath),
		Headers:   unprocessedYaml,
		Body:      unprocessedMarkdown,
	})
	_mutexLock.Unlock()
}
