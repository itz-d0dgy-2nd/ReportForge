package processors

import (
	Utils "ReportForge/engine/utils"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

/*
ProcessConfigMetadata → Process YAML configuration file for report metadata
  - Reads YAML configuration file from specified file path
  - Unmarshals YAML content into pointer *Utils.MetadataYML struct
  - Draft or Release validation on:
    -- information.DocumentVersioning["DocumentStatus"]
  - Handles validation errors with descriptive messages including file path
  - Handles errors via Utils.ErrorChecker()
*/
func ProcessConfigMetadata(_filePath string, _metadata *Utils.MetadataYML) {
	rawFileContent, errRawFileContent := os.ReadFile(_filePath)
	Utils.ErrorChecker(errRawFileContent)

	errDecodeYML := yaml.Unmarshal(rawFileContent, &_metadata)
	Utils.ErrorChecker(errDecodeYML)

	for _, information := range _metadata.DocumentInformation {
		if status, exists := information.DocumentVersioning["DocumentStatus"]; exists {
			if status != "Draft" && status != "Release" {
				Utils.ErrorChecker(fmt.Errorf("invalid DocumentStatus found in '%s'", _filePath))
			}
		}
	}
}

/*
ProcessConfigSeverityAssessment → Process YAML configuration file for report severity assessment
  - Reads YAML configuration file from specified file path
  - Unmarshals YAML content into pointer *Utils.SeverityAssessmentYML struct
  - Empty or whitespace-only validation on:
    -- _severityAssessment.Impacts
    -- _severityAssessment.Likelihoods
    -- _severityAssessment.CalculatedMatrix
  - Handles validation errors with descriptive messages including file path
  - Handles errors via Utils.ErrorChecker()
*/
func ProcessConfigSeverityAssessment(_filePath string, _severityAssessment *Utils.SeverityAssessmentYML) {
	rawFileContent, errRawFileContent := os.ReadFile(_filePath)
	Utils.ErrorChecker(errRawFileContent)

	errDecodeYML := yaml.Unmarshal(rawFileContent, &_severityAssessment)
	Utils.ErrorChecker(errDecodeYML)

	for _, impact := range _severityAssessment.Impacts {
		if strings.TrimSpace(impact) == "" {
			Utils.ErrorChecker(fmt.Errorf("empty impact found in '%s'", _filePath))
		}
	}

	for _, likelihood := range _severityAssessment.Likelihoods {
		if strings.TrimSpace(likelihood) == "" {
			Utils.ErrorChecker(fmt.Errorf("empty likelihood found in '%s'", _filePath))
		}
	}

	for _, impact := range _severityAssessment.CalculatedMatrix {
		for _, severity := range impact {
			if strings.TrimSpace(severity) == "" {
				Utils.ErrorChecker(fmt.Errorf("empty severity found in '%s", _filePath))
			}
		}
	}
}

/*
ProcessSeverityMatrix → Process markdown files with YAML frontmatter for report severity assessment
  - Reads markdown file from specified file path
  - Extracts YAML frontmatter from markdown and unmarshals YAML content into Utils.MarkdownYML struct
  - Calculates matrix indices based on finding impact and likelihood values
  - Updates severity assessment matrix with finding ID at calculated position
  - Handles errors via Utils.ErrorChecker()
*/
func ProcessSeverityMatrix(_filePath string, _severityAssessment *Utils.SeverityAssessmentYML) {
	if _severityAssessment.ConductSeverityAssessment {
		unprocessedYaml := Utils.MarkdownYML{}

		rawFileContent, errReadMarkdown := os.ReadFile(_filePath)
		Utils.ErrorChecker(errReadMarkdown)

		rawMarkdownContent := strings.ReplaceAll(string(rawFileContent), "\r\n", "\n")
		regexMatches := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`).FindStringSubmatch(rawMarkdownContent)

		yaml.Unmarshal([]byte(regexMatches[1]), &unprocessedYaml)

		impactIndex := slices.Index(_severityAssessment.Impacts, unprocessedYaml.FindingImpact)
		likelihoodIndex := slices.Index(_severityAssessment.Likelihoods, unprocessedYaml.FindingLikelihood)

		if _severityAssessment.Matrix[impactIndex][likelihoodIndex] == "" {
			_severityAssessment.Matrix[impactIndex][likelihoodIndex] = unprocessedYaml.FindingID
		} else {
			_severityAssessment.Matrix[impactIndex][likelihoodIndex] += ", " + unprocessedYaml.FindingID
		}
	}
}

/*
ProcessMarkdown → Process markdown files with YAML frontmatter for report content
  - Reads markdown file from specified file path
  - Extracts YAML frontmatter from markdown and unmarshals YAML content into Utils.MarkdownYML struct
  - Converts markdown content to HTML using blackfriday processor
  - Performs strings checks and replace on:
    -- Custom Tags (<qa></qa>, <retest_fixed></retest_fixed>, <retest_not_fixed></retest_not_fixed>)
    -- Cutsom Tokens (!Client, !TargetAsset0, !TargetAsset1)
    -- Updates screenshot paths to include full report path
  - Thread-safely appends processed markdown to slice using mutex lock
  - Handles errors via Utils.ErrorChecker()
*/
func ProcessMarkdown(_reportPath string, _filePath string, _processedMarkdown *[]Utils.Markdown, _metadata Utils.MetadataYML, _mutexLock *sync.Mutex) {
	unprocessedYaml := Utils.MarkdownYML{}

	rawFileContent, errReadMD := os.ReadFile(_filePath)
	Utils.ErrorChecker(errReadMD)

	rawMarkdownContent := strings.ReplaceAll(string(rawFileContent), "\r\n", "\n")
	regexMatches := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`).FindStringSubmatch(rawMarkdownContent)

	errDecodeYML := yaml.Unmarshal([]byte(regexMatches[1]), &unprocessedYaml)
	Utils.ErrorChecker(errDecodeYML)

	unprocessedMarkdown := string(blackfriday.Run([]byte(regexMatches[2])))

	unprocessedMarkdown = regexp.MustCompile(`\B!([A-Za-z][A-Za-z0-9_]*)\b`).ReplaceAllStringFunc(unprocessedMarkdown, func(tokenMatch string) string {

		if tokenValue, exists := _metadata.CustomVariables[strings.TrimPrefix(tokenMatch, "!")]; exists {
			return tokenValue
		}

		if strings.TrimPrefix(tokenMatch, "!") == "Client" {
			return _metadata.Client
		}

		return tokenMatch
	})

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

	if strings.Contains(unprocessedMarkdown, "Screenshots/") {
		unprocessedMarkdown = strings.ReplaceAll(unprocessedMarkdown, "Screenshots/", _reportPath+"/Screenshots/")
	}

	_mutexLock.Lock()
	*_processedMarkdown = append(*_processedMarkdown, Utils.Markdown{
		Directory: filepath.Base(filepath.Dir(_filePath)),
		FileName:  filepath.Base(_filePath),
		Headers:   unprocessedYaml,
		Body:      unprocessedMarkdown,
	})

	_mutexLock.Unlock()
}
