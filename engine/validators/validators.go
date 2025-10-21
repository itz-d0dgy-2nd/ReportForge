package validators

import (
	"ReportForge/engine/utilities"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

/*
ValidateYamlFrontmatter → Validates YAML frontmatter from regex matches
*/
func ValidateYamlFrontmatter(_regexMatches []string, _filePath string, _unprocessedYaml *utilities.MarkdownYML) {
	if len(_regexMatches) < 2 {
		utilities.ErrorChecker(fmt.Errorf("missing YAML frontmatter in ( %s )", _filePath))
	}

	if yaml.Unmarshal([]byte(_regexMatches[1]), &_unprocessedYaml) != nil {
		utilities.ErrorChecker(fmt.Errorf("invalid YAML frontmatter in ( %s )", _filePath))
	}
}

/*
ValidateImpactLikelihoodIndex → Validates that impact or likelihood index is valid
*/
func ValidateImpactLikelihoodIndex(_index int, _fieldName string, _filePath string) {
	if _index == -1 {
		utilities.ErrorChecker(fmt.Errorf("invalid %s found in ( %s )", _fieldName, _filePath))
	}
}

/*
ValidateSeverityExists → Validates that severity key is valid
*/
func ValidateSeverityKey(_severity string, _scales map[string]int, _filePath string) {
	if _severity == "" {
		utilities.ErrorChecker(fmt.Errorf("empty severity found in ( %s )", _filePath))
	}

	if _, exists := _scales[_severity]; !exists {
		utilities.ErrorChecker(fmt.Errorf("invalid severity found in ( %s )", _filePath))
	}
}

/*
ValidateConfigMetadata → Validates metadata configuration DocumentStatus values
*/
func ValidateConfigMetadata(_metadata *utilities.MetadataYML, _filePath string) {
	for _, information := range _metadata.DocumentInformation {
		if status, exists := information.DocumentVersioning["DocumentStatus"]; exists {
			if status != "Draft" && status != "Release" {
				utilities.ErrorChecker(fmt.Errorf("invalid DocumentStatus found in ( %s )", _filePath))
			}
		}
	}
}

/*
ValidateConfigSeverityAssessment → Validates severity assessment configuration values
*/
func ValidateConfigSeverityAssessment(_severityAssessment *utilities.SeverityAssessmentYML, _filePath string) {
	for _, impact := range _severityAssessment.Impacts {
		if strings.TrimSpace(impact) == "" {
			utilities.ErrorChecker(fmt.Errorf("empty impact found in ( %s )", _filePath))
		}
	}

	for _, likelihood := range _severityAssessment.Likelihoods {
		if strings.TrimSpace(likelihood) == "" {
			utilities.ErrorChecker(fmt.Errorf("empty likelihood found in ( %s )", _filePath))
		}
	}

	for _, impact := range _severityAssessment.CalculatedMatrix {
		for _, severity := range impact {
			if strings.TrimSpace(severity) == "" {
				utilities.ErrorChecker(fmt.Errorf("empty severity found in ( %s )", _filePath))
			}
		}
	}
}
