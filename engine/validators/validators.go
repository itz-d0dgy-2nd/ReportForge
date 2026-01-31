package validators

import (
	"ReportForge/engine/utilities"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

/*
ValidateConfigMetadata → Validates metadata configuration DocumentStatus values
*/
func ValidateConfigMetadata(_metadata *utilities.MetadataYML, _filePath string) {
	if len(_metadata.DocumentInformation) == 0 {
		utilities.ErrorChecker(fmt.Errorf("empty DocumentInformation in ( %s )", _filePath))
	}

	for _, information := range _metadata.DocumentInformation {
		if status, exists := information.DocumentVersioning["DocumentStatus"]; exists {
			if status != utilities.ReportStatusDraft && status != utilities.ReportStatusRelease {
				utilities.ErrorChecker(fmt.Errorf("invalid DocumentStatus in ( %s )", _filePath))
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
			utilities.ErrorChecker(fmt.Errorf("empty impact in ( %s )", _filePath))
		}
	}

	for _, likelihood := range _severityAssessment.Likelihoods {
		if strings.TrimSpace(likelihood) == "" {
			utilities.ErrorChecker(fmt.Errorf("empty likelihood in ( %s )", _filePath))
		}
	}

	for _, impact := range _severityAssessment.CalculatedMatrix {
		for _, severity := range impact {
			if strings.TrimSpace(severity) == "" {
				utilities.ErrorChecker(fmt.Errorf("empty severity in ( %s )", _filePath))
			}
		}
	}
}

/*
ValidateConfigContentOrder → Validates content order configuration values
*/
func ValidateConfigContentOrder(_contentOrder *utilities.ContentOrderYML, _filePath string) {
	for _, prefix := range _contentOrder.FindingIdentifierPrefixes {
		if strings.TrimSpace(prefix) == "" {
			utilities.ErrorChecker(fmt.Errorf("empty finding prefix in ( %s )", _filePath))
		}
	}

	for _, prefix := range _contentOrder.SuggestionIdentifierPrefixes {
		if strings.TrimSpace(prefix) == "" {
			utilities.ErrorChecker(fmt.Errorf("empty suggestion prefix in ( %s )", _filePath))
		}
	}

	for _, prefix := range _contentOrder.RiskIdentifierPrefixes {
		if strings.TrimSpace(prefix) == "" {
			utilities.ErrorChecker(fmt.Errorf("empty risk prefix in ( %s )", _filePath))
		}
	}
}

/*
ValidateYamlFrontmatter → Validates YAML frontmatter from regex matches
*/
func ValidateYamlFrontmatter(_regexMatches []string, _filePath string, _unprocessedYaml *utilities.MarkdownYML) {
	var rawYaml map[string]string
	var missingKeys []string

	if len(_regexMatches) < 2 {
		utilities.ErrorChecker(fmt.Errorf("missing YAML frontmatter in ( %s )", _filePath))
		return
	}

	if errDecodeYAML := yaml.Unmarshal([]byte(_regexMatches[1]), &rawYaml); errDecodeYAML != nil {
		utilities.ErrorChecker(fmt.Errorf("invalid YAML frontmatter in ( %s )", _filePath))
		return
	}

	requiredFieldsMap := map[string][]string{
		utilities.SummariesDirectory:   {"ReportSummaryName", "ReportSummaryTitle", "ReportSummariesAuthor", "ReportSummariesReviewers"},
		utilities.FindingsDirectory:    {"FindingID", "FindingIDLocked", "FindingName", "FindingTitle", "FindingStatus", "FindingImpact", "FindingLikelihood", "FindingSeverity", "FindingAuthor", "FindingReviewers"},
		utilities.SuggestionsDirectory: {"SuggestionID", "SuggestionIDLocked", "SuggestionName", "SuggestionTitle", "SuggestionStatus", "SuggestionAuthor", "SuggestionReviewers"},
		utilities.RisksDirectory:       {"RiskID", "RiskIDLocked", "RiskName", "RiskTitle", "RiskStatus", "RiskGrossImpact", "RiskGrossLikelihood", "RiskGrossRating", "RiskTargetImpact", "RiskTargetLikelihood", "RiskTargetRating", "RiskAuthor", "RiskReviewers"},
		utilities.AppendicesDirectory:  {"AppendixName", "AppendixTitle", "AppendixStatus", "AppendixAuthor", "AppendixReviewers"},
	}

	for directory, fields := range requiredFieldsMap {
		if strings.Contains(_filePath, directory) {
			for _, field := range fields {
				if _, exists := rawYaml[field]; !exists {
					missingKeys = append(missingKeys, field)
				}
			}
			break
		}
	}

	if len(missingKeys) > 0 {
		utilities.ErrorChecker(fmt.Errorf("missing required key(s) [%s] in ( %s )", strings.Join(missingKeys, ", "), _filePath))
		return
	}

	if errDecodeYAML := yaml.Unmarshal([]byte(_regexMatches[1]), _unprocessedYaml); errDecodeYAML != nil {
		utilities.ErrorChecker(fmt.Errorf("invalid YAML frontmatter in ( %s )", _filePath))
	}
}

/*
ValidateImpactIndex → Validates that impact index is valid
*/
func ValidateImpactIndex(_index int, _filePath string) {
	if _index == -1 {
		utilities.ErrorChecker(fmt.Errorf("invalid impact in ( %s )", _filePath))
	}
}

/*
ValidateLikelihoodIndex → Validates that likelihood index is valid
*/
func ValidateLikelihoodIndex(_index int, _filePath string) {
	if _index == -1 {
		utilities.ErrorChecker(fmt.Errorf("invalid likelihood in ( %s )", _filePath))
	}
}

/*
ValidateSeverityIndex → Validates that severity index is valid
*/
func ValidateSeverityIndex(_index int, _filePath string) {
	if _index == -1 {
		utilities.ErrorChecker(fmt.Errorf("invalid severity in ( %s )", _filePath))
	}
}
