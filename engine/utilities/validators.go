package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
validateConfigFiles → Validates all required config files exist and are readable.
  - Required: metadata.yml, content-order.yml
  - At least one of: severity-assessment.yml, risk-assessment.yml
*/
func validateConfigFiles(configPath string) {
	requiredConfigs := []string{
		ConfigFileMetadata,
		ConfigFileContentOrder,
	}

	for _, configName := range requiredConfigs {
		path := filepath.Join(configPath, configName+".yml")
		if _, errReadFile := os.ReadFile(path); errReadFile != nil {
			if os.IsNotExist(errReadFile) {
				Check(NewFileSystemError(
					path,
					fmt.Sprintf("required config file missing: %s.yml", configName),
					errReadFile,
				))
			}
		}
	}

	severityPath := filepath.Join(configPath, ConfigFileSeverityAssessment+".yml")
	riskPath := filepath.Join(configPath, ConfigFileRiskAssessment+".yml")

	_, errSeverity := os.Stat(severityPath)
	_, errRisk := os.Stat(riskPath)

	if errSeverity != nil && errRisk != nil {
		Check(NewConfigError(
			configPath,
			fmt.Sprintf("at least one assessment config required: %s.yml or %s.yml",
				ConfigFileSeverityAssessment, ConfigFileRiskAssessment),
		))
	}
}

/*
Validate → Validates metadata configuration values
*/
func (_metadata *MetadataYML) Validate(_path string) {
	if strings.TrimSpace(_metadata.Client) == "" {
		Check(NewValidationError(_path, "Client", "field is required"))
	}

	if strings.TrimSpace(_metadata.DocumentName) == "" {
		Check(NewValidationError(_path, "DocumentName", "field is required"))
	}

	if len(_metadata.DocumentInformation) == 0 {
		Check(NewValidationError(_path, "DocumentInformation", "at least one entry is required"))
	}

	for i, information := range _metadata.DocumentInformation {
		if strings.TrimSpace(information.DocumentVersioning["DocumentStatus"]) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("DocumentInformation[%d].DocumentStatus", i),
				"field is required",
			))
		}

		if status, exists := information.DocumentVersioning["DocumentStatus"]; exists {
			if status != ReportStatusDraft && status != ReportStatusRelease {
				Check(NewValidationError(
					_path,
					fmt.Sprintf("DocumentInformation[%d].DocumentStatus", i),
					fmt.Sprintf("must be '%s' or '%s', got '%s'", ReportStatusDraft, ReportStatusRelease, status),
				))
			}
		}

		if strings.TrimSpace(information.DocumentVersioning["DocumentVersion"]) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("DocumentInformation[%d].DocumentVersion", i),
				"field is required",
			))
		}

		if strings.TrimSpace(information.DocumentVersioning["DocumentDate"]) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("DocumentInformation[%d].DocumentDate", i),
				"field is required",
			))
		}

		if strings.TrimSpace(information.DocumentVersioning["DocumentStage"]) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("DocumentInformation[%d].DocumentStage", i),
				"field is required",
			))
		}

		if strings.TrimSpace(information.DocumentVersioning["DocumentContributor"]) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("DocumentInformation[%d].DocumentContributor", i),
				"field is required",
			))
		}
	}

	if len(_metadata.StakeholderInformation) == 0 {
		Check(NewValidationError(_path, "StakeholderInformation", "at least one stakeholder is required"))
	}

	for i, stakeholder := range _metadata.StakeholderInformation {
		if strings.TrimSpace(stakeholder.StakeholderName) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("StakeholderInformation[%d].StakeholderName", i),
				"field is required",
			))
		}

		if strings.TrimSpace(stakeholder.StakeholderRole) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("StakeholderInformation[%d].StakeholderRole", i),
				"field is required",
			))
		}

		if strings.TrimSpace(stakeholder.StakeholderCompany) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("StakeholderInformation[%d].StakeholderCompany", i),
				"field is required",
			))
		}
	}
}

/*
Validate → Validates severity assessment configuration values
*/
func (_severityAssessment *SeverityAssessmentYML) Validate(_path string) {
	if len(_severityAssessment.Impacts) < 5 {
		Check(NewConfigError(
			_path,
			fmt.Sprintf("Impacts must have at least 5 entries, got %d", len(_severityAssessment.Impacts)),
		))
	}

	for i, impact := range _severityAssessment.Impacts {
		if strings.TrimSpace(impact) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("Impacts[%d]", i),
				"value cannot be empty",
			))
		}
	}

	if len(_severityAssessment.Likelihoods) < 5 {
		Check(NewConfigError(
			_path,
			fmt.Sprintf("Likelihoods must have at least 5 entries, got %d", len(_severityAssessment.Likelihoods)),
		))
	}

	for i, likelihood := range _severityAssessment.Likelihoods {
		if strings.TrimSpace(likelihood) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("Likelihoods[%d]", i),
				"value cannot be empty",
			))
		}
	}

	for row, impacts := range _severityAssessment.CalculatedMatrix {
		for col, severity := range impacts {
			if strings.TrimSpace(severity) == "" {
				Check(NewValidationError(
					_path,
					fmt.Sprintf("CalculatedMatrix[%d][%d]", row, col),
					"value cannot be empty",
				))
			}
		}
	}
}

/*
Validate → Validates risk assessment configuration values
*/
func (_riskAssessment *RiskAssessmentYML) Validate(_path string) {
	if len(_riskAssessment.GrossImpacts) < 5 {
		Check(NewConfigError(
			_path,
			fmt.Sprintf("GrossImpacts must have at least 5 entries, got %d", len(_riskAssessment.GrossImpacts)),
		))
	}

	for i, grossImpact := range _riskAssessment.GrossImpacts {
		if strings.TrimSpace(grossImpact) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("GrossImpacts[%d]", i),
				"value cannot be empty",
			))
		}
	}

	if len(_riskAssessment.GrossLikelihoods) < 5 {
		Check(NewConfigError(
			_path,
			fmt.Sprintf("GrossLikelihoods must have at least 5 entries, got %d", len(_riskAssessment.GrossLikelihoods)),
		))
	}

	for i, grossLikelihood := range _riskAssessment.GrossLikelihoods {
		if strings.TrimSpace(grossLikelihood) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("GrossLikelihoods[%d]", i),
				"value cannot be empty",
			))
		}
	}

	if len(_riskAssessment.TargetImpacts) < 5 {
		Check(NewConfigError(
			_path,
			fmt.Sprintf("TargetImpacts must have at least 5 entries, got %d", len(_riskAssessment.TargetImpacts)),
		))
	}

	for i, targetImpact := range _riskAssessment.TargetImpacts {
		if strings.TrimSpace(targetImpact) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("TargetImpacts[%d]", i),
				"value cannot be empty",
			))
		}
	}

	if len(_riskAssessment.TargetLikelihoods) < 5 {
		Check(NewConfigError(
			_path,
			fmt.Sprintf("TargetLikelihoods must have at least 5 entries, got %d", len(_riskAssessment.TargetLikelihoods)),
		))
	}

	for i, targetLikelihood := range _riskAssessment.TargetLikelihoods {
		if strings.TrimSpace(targetLikelihood) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("TargetLikelihoods[%d]", i),
				"value cannot be empty",
			))
		}
	}
}

/*
Validate → Validates content order configuration values
*/
func (_contentOrder *ContentOrderYML) Validate(_path string) {
	prefixUsage := make(map[string]string)

	for subdirName, prefix := range _contentOrder.FindingIdentifierPrefixes {
		if strings.TrimSpace(prefix) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("FindingIdentifierPrefixes[%s]", subdirName),
				"prefix cannot be empty",
			))
		}
		if existingLocation, exists := prefixUsage[prefix]; exists {
			Check(NewValidationError(
				_path,
				"FindingIdentifierPrefixes",
				fmt.Sprintf("duplicate prefix '%s' - already used by '%s', cannot reuse for '%s'", prefix, existingLocation, subdirName),
			))
		}
		prefixUsage[prefix] = subdirName + " (findings)"
	}

	for subdirName, prefix := range _contentOrder.SuggestionIdentifierPrefixes {
		if strings.TrimSpace(prefix) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("SuggestionIdentifierPrefixes[%s]", subdirName),
				"prefix cannot be empty",
			))
		}
		if existingLocation, exists := prefixUsage[prefix]; exists {
			Check(NewValidationError(
				_path,
				"SuggestionIdentifierPrefixes",
				fmt.Sprintf("duplicate prefix '%s' - already used by '%s', cannot reuse for '%s'", prefix, existingLocation, subdirName),
			))
		}
		prefixUsage[prefix] = subdirName + " (suggestions)"
	}

	for subdirName, prefix := range _contentOrder.RiskIdentifierPrefixes {
		if strings.TrimSpace(prefix) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("RiskIdentifierPrefixes[%s]", subdirName),
				"prefix cannot be empty",
			))
		}
		if existingLocation, exists := prefixUsage[prefix]; exists {
			Check(NewValidationError(
				_path,
				"RiskIdentifierPrefixes",
				fmt.Sprintf("duplicate prefix '%s' - already used by '%s', cannot reuse for '%s'", prefix, existingLocation, subdirName),
			))
		}
		prefixUsage[prefix] = subdirName + " (risks)"
	}

	for subdirName, prefix := range _contentOrder.ControlIdentifierPrefixes {
		if strings.TrimSpace(prefix) == "" {
			Check(NewValidationError(
				_path,
				fmt.Sprintf("ControlIdentifierPrefixes[%s]", subdirName),
				"prefix cannot be empty",
			))
		}
		if existingLocation, exists := prefixUsage[prefix]; exists {
			Check(NewValidationError(
				_path,
				"ControlIdentifierPrefixes",
				fmt.Sprintf("duplicate prefix '%s' - already used by '%s', cannot reuse for '%s'", prefix, existingLocation, subdirName),
			))
		}
		prefixUsage[prefix] = subdirName + " (controls)"
	}
}

/*
Validate → Validates markdown frontmatter fields based on the file's directory context
*/
func (_markdown *MarkdownYML) Validate(_rawYaml map[string]any, _path string) {
	requiredFieldsMap := map[string][]string{
		SummariesDirectory:   {"ReportSummaryName", "ReportSummaryTitle", "ReportSummariesAuthor", "ReportSummariesReviewers"},
		FindingsDirectory:    {"FindingID", "FindingIDLocked", "FindingName", "FindingTitle", "FindingStatus", "FindingImpact", "FindingLikelihood", "FindingSeverity", "FindingAuthor", "FindingReviewers"},
		SuggestionsDirectory: {"SuggestionID", "SuggestionIDLocked", "SuggestionName", "SuggestionTitle", "SuggestionStatus", "SuggestionAuthor", "SuggestionReviewers"},
		RisksDirectory:       {"RiskID", "RiskIDLocked", "RiskName", "RiskTitle", "RiskGrossImpact", "RiskGrossLikelihood", "RiskGrossRating", "RiskTargetImpact", "RiskTargetLikelihood", "RiskTargetRating", "RiskAuthor", "RiskReviewers"},
		ControlsDirectory:    {"ControlID", "ControlIDLocked", "ControlName", "ControlTitle", "ControlNZISMReferences", "ControlAuthor", "ControlReviewers"},
		AppendicesDirectory:  {"AppendixName", "AppendixTitle", "AppendixStatus", "AppendixAuthor", "AppendixReviewers"},
	}

	var missingKeys []string
	for directory, fields := range requiredFieldsMap {
		if strings.Contains(_path, directory) {
			for _, field := range fields {
				if _, exists := _rawYaml[field]; !exists {
					missingKeys = append(missingKeys, field)
				}
			}
			break
		}
	}

	if len(missingKeys) > 0 {
		Check(NewValidationError(
			_path,
			strings.Join(missingKeys, ", "),
			"required frontmatter field(s) missing",
		))
	}

	if strings.Contains(_path, ControlsDirectory) {
		for _, cid := range _markdown.ControlNZISMReferences {
			if !YAMLPattern.CIDRef.MatchString(cid) {
				Check(NewValidationWarning(
					_path,
					fmt.Sprintf("NZISM reference '%s' does not match format 'CID:<number>'", cid),
				))
			}
		}
	}
}
