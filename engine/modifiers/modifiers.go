package modifiers

import (
	"ReportForge/engine/utilities"
	"fmt"
	_ "image/png"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

/*
ModifyFindingFiles → Calculates and writes finding severity and constructs the filename
*/
func ModifyFindingFiles(_path string, _fileCache *utilities.FileCache) string {
	var calculatedSeverity string

	rawMarkdownContent, _, unprocessedYaml := utilities.ParseFile(_path, _fileCache)
	severityConfig := _fileCache.SeverityConfig()

	if severityConfig.ConductSeverityAssessment {
		impactIndex := slices.Index(severityConfig.Impacts, unprocessedYaml.FindingImpact)
		likelihoodIndex := slices.Index(severityConfig.Likelihoods, unprocessedYaml.FindingLikelihood)

		if impactIndex == -1 || likelihoodIndex == -1 {
			utilities.Check(utilities.NewValidationError(
				_path,
				"FindingImpact/FindingLikelihood",
				fmt.Sprintf("invalid impact '%s' or likelihood '%s' - must match configured values", unprocessedYaml.FindingImpact, unprocessedYaml.FindingLikelihood),
			))
			return _path
		}

		if severityConfig.SwapImpactLikelihoodAxis {
			calculatedSeverity = severityConfig.CalculatedMatrix[likelihoodIndex][impactIndex]
		} else {
			calculatedSeverity = severityConfig.CalculatedMatrix[impactIndex][likelihoodIndex]
		}

		if unprocessedYaml.FindingSeverity != calculatedSeverity {
			rawMarkdownContent = utilities.YAMLPattern.Severity.ReplaceAllString(rawMarkdownContent, "FindingSeverity: "+calculatedSeverity)
			if errWriteFile := os.WriteFile(_path, []byte(rawMarkdownContent), 0644); errWriteFile != nil {
				utilities.Check(utilities.NewFileSystemError(_path, "failed to write finding file", errWriteFile))
			}
			_fileCache.UpdateFile(_path, []byte(rawMarkdownContent))
		}

	} else {
		calculatedSeverity = unprocessedYaml.FindingSeverity
	}

	severityIndex := slices.Index(severityConfig.Severities, calculatedSeverity)
	if severityIndex == -1 {
		utilities.Check(utilities.NewValidationError(
			_path,
			"FindingSeverity",
			fmt.Sprintf("severity '%s' not found in configured severity levels", calculatedSeverity),
		))
		return _path
	}

	newFileName := strconv.Itoa(severityIndex) + "_" + unprocessedYaml.FindingName + ".md"
	return modifyFileName(_path, newFileName, []byte(rawMarkdownContent), _fileCache)
}

/*
ModifySuggestionFiles → Constructs the filename
*/
func ModifySuggestionFiles(_path string, _fileCache *utilities.FileCache) string {
	rawMarkdownContent, _, unprocessedYaml := utilities.ParseFile(_path, _fileCache)
	newFileName := unprocessedYaml.SuggestionName + ".md"
	return modifyFileName(_path, newFileName, []byte(rawMarkdownContent), _fileCache)
}

/*
ModifyRiskFiles → Calculates and writes gross and target risk ratings and constructs the filename
*/
func ModifyRiskFiles(_path string, _fileCache *utilities.FileCache) string {
	rawMarkdownContent, _, unprocessedYaml := utilities.ParseFile(_path, _fileCache)
	riskConfig := _fileCache.RiskConfig()

	grossImpactIndex := slices.Index(riskConfig.GrossImpacts, unprocessedYaml.RiskGrossImpact)
	grossLikelihoodIndex := slices.Index(riskConfig.GrossLikelihoods, unprocessedYaml.RiskGrossLikelihood)
	targetImpactIndex := slices.Index(riskConfig.TargetImpacts, unprocessedYaml.RiskTargetImpact)
	targetLikelihoodIndex := slices.Index(riskConfig.TargetLikelihoods, unprocessedYaml.RiskTargetLikelihood)

	if grossImpactIndex == -1 || grossLikelihoodIndex == -1 {
		utilities.Check(utilities.NewValidationError(
			_path,
			"RiskGrossImpact/RiskGrossLikelihood",
			fmt.Sprintf("invalid gross impact '%s' or likelihood '%s' - must match configured values", unprocessedYaml.RiskGrossImpact, unprocessedYaml.RiskGrossLikelihood),
		))
		return _path
	}

	if targetImpactIndex == -1 || targetLikelihoodIndex == -1 {
		utilities.Check(utilities.NewValidationError(
			_path,
			"RiskTargetImpact/RiskTargetLikelihood",
			fmt.Sprintf("invalid target impact '%s' or likelihood '%s' - must match configured values", unprocessedYaml.RiskTargetImpact, unprocessedYaml.RiskTargetLikelihood),
		))
		return _path
	}

	calculatedGrossRating := riskConfig.CalculatedGrossMatrix[grossImpactIndex][grossLikelihoodIndex]
	calculatedTargetRating := riskConfig.CalculatedTargetMatrix[targetImpactIndex][targetLikelihoodIndex]

	var fileModified bool
	if unprocessedYaml.RiskGrossRating != calculatedGrossRating {
		fileModified = true
		rawMarkdownContent = utilities.YAMLPattern.GrossRating.ReplaceAllString(rawMarkdownContent, "RiskGrossRating: "+calculatedGrossRating)
	}
	if unprocessedYaml.RiskTargetRating != calculatedTargetRating {
		fileModified = true
		rawMarkdownContent = utilities.YAMLPattern.TargetRating.ReplaceAllString(rawMarkdownContent, "RiskTargetRating: "+calculatedTargetRating)
	}

	if fileModified {
		if errWriteFile := os.WriteFile(_path, []byte(rawMarkdownContent), 0644); errWriteFile != nil {
			utilities.Check(utilities.NewFileSystemError(_path, "failed to write risk file", errWriteFile))
		}
		_fileCache.UpdateFile(_path, []byte(rawMarkdownContent))
	}

	grossRatingIndex := slices.Index(riskConfig.GrossRiskRatings, calculatedGrossRating)
	if grossRatingIndex == -1 {
		utilities.Check(utilities.NewValidationError(
			_path,
			"RiskGrossRating",
			fmt.Sprintf("gross rating '%s' not found in configured risk ratings", calculatedGrossRating),
		))
		return _path
	}

	targetRatingIndex := slices.Index(riskConfig.TargetRiskRatings, calculatedTargetRating)
	if targetRatingIndex == -1 {
		utilities.Check(utilities.NewValidationError(
			_path,
			"RiskTargetRating",
			fmt.Sprintf("target rating '%s' not found in configured risk ratings", calculatedTargetRating),
		))
		return _path
	}

	newFileName := strconv.Itoa(grossRatingIndex) + "_" + unprocessedYaml.RiskName + ".md"
	return modifyFileName(_path, newFileName, []byte(rawMarkdownContent), _fileCache)
}

/*
ModifyControlFiles → Constructs the filename
*/
func ModifyControlFiles(_path string, _fileCache *utilities.FileCache) string {
	rawMarkdownContent, _, unprocessedYaml := utilities.ParseFile(_path, _fileCache)
	newFileName := unprocessedYaml.ControlName + ".md"
	return modifyFileName(_path, newFileName, []byte(rawMarkdownContent), _fileCache)
}

func ModifyIdentifiers(_path, _identifierPrefix string, _reservedID int32, _fileCache *utilities.FileCache) {
	var identifierLocked bool
	var identifierLockField string
	var identifierRegex *regexp.Regexp

	rawMarkdownContent, _, unprocessedYaml := utilities.ParseFile(_path, _fileCache)

	switch utilities.GetDirectoryType(_path) {
	case utilities.FindingsDirectory:
		identifierLocked = unprocessedYaml.FindingIDLocked
		identifierLockField = "FindingIDLocked"
		identifierRegex = utilities.YAMLPattern.FindingID
	case utilities.SuggestionsDirectory:
		identifierLocked = unprocessedYaml.SuggestionIDLocked
		identifierLockField = "SuggestionIDLocked"
		identifierRegex = utilities.YAMLPattern.SuggestionID
	case utilities.RisksDirectory:
		identifierLocked = unprocessedYaml.RiskIDLocked
		identifierLockField = "RiskIDLocked"
		identifierRegex = utilities.YAMLPattern.RiskID
	case utilities.ControlsDirectory:
		identifierLocked = unprocessedYaml.ControlIDLocked
		identifierLockField = "ControlIDLocked"
		identifierRegex = utilities.YAMLPattern.ControlID
	}

	if identifierLocked {
		if utilities.DocumentStatus != utilities.ReportStatusRelease {
			return
		}

		if !strings.Contains(rawMarkdownContent, identifierLockField+": false") && !strings.Contains(rawMarkdownContent, identifierLockField+":false") {
			return
		}

		rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+": false", identifierLockField+": true", 1)
		if errWriteFile := os.WriteFile(_path, []byte(rawMarkdownContent), 0644); errWriteFile != nil {
			utilities.Check(utilities.NewFileSystemError(_path, "failed to write identifier lock changes", errWriteFile))
		}
		_fileCache.UpdateFile(_path, []byte(rawMarkdownContent))
		return
	}

	generatedIdentifier := fmt.Sprintf("%s%d", _identifierPrefix, _reservedID)
	rawMarkdownContent = identifierRegex.ReplaceAllString(rawMarkdownContent, "${1}"+generatedIdentifier)

	if utilities.DocumentStatus == utilities.ReportStatusRelease {
		if strings.Contains(rawMarkdownContent, identifierLockField+": false") {
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+": false", identifierLockField+": true", 1)
		} else if strings.Contains(rawMarkdownContent, identifierLockField+":") {
			rawMarkdownContent = strings.Replace(rawMarkdownContent, identifierLockField+":", identifierLockField+": true", 1)
		}
	}

	if errWriteFile := os.WriteFile(_path, []byte(rawMarkdownContent), 0644); errWriteFile != nil {
		utilities.Check(utilities.NewFileSystemError(_path, "failed to write identifier changes", errWriteFile))
	}
	_fileCache.UpdateFile(_path, []byte(rawMarkdownContent))
}

/*
modifyFileName → Renames file and updates cache if filename changes
*/
func modifyFileName(_path, _newFileName string, _fileContent []byte, _fileCache *utilities.FileCache) string {
	newFilePath := filepath.Clean(filepath.Join(filepath.Dir(_path), _newFileName))

	if _path != newFilePath {
		if errRename := os.Rename(_path, newFilePath); errRename != nil {
			utilities.Check(utilities.NewFileSystemError(
				_path,
				fmt.Sprintf("cannot rename to '%s' - file may already exist or path is invalid", filepath.Base(newFilePath)),
				errRename,
			))
		}
		_fileCache.RenameFile(_path, newFilePath, _fileContent)
		return newFilePath
	}

	return _path
}
