package handlers

import (
	"ReportForge/engine/modifiers"
	"ReportForge/engine/processors"
	"ReportForge/engine/utilities"
	"ReportForge/engine/validators"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

/*
HandleConfigProcessor → Recursively walks directory tree and processes YAML config files
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .yml file in report_config directory:
    -- Calls processors.ProcessConfigMetadata() or,
    -- Calls processors.ProcessConfigSeverityAssessment() or,
    -- Calls processors.ProcessConfigContentOrder()
  - Handles errors via utilities.ErrorChecker()
  - Returns processed yml of type utilities.MetadataYML, utilities.SeverityAssessmentYML, utilities.ContentOrderYML
*/
func HandleConfigs(_reportPaths utilities.ReportPaths, _fileCache *utilities.FileCache) {

	errDirectoryWalk := filepath.WalkDir(_reportPaths.ConfigPath, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil {
			return errAnonymousFunction
		}

		if directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".yml" {
			return nil
		}

		rawFileContent, errRawFileContent := _fileCache.ReadFile(filePath)
		utilities.ErrorChecker(errRawFileContent)

		if strings.Contains(directoryContents.Name(), utilities.ConfigFileMetadata) {
			var tempMetadataConfig utilities.MetadataYML
			errDecodeYAML := yaml.Unmarshal(rawFileContent, &tempMetadataConfig)
			utilities.ErrorChecker(errDecodeYAML)
			validators.ValidateConfigMetadata(&tempMetadataConfig, filePath)

			utilities.DocumentStatus = tempMetadataConfig.DocumentInformation[len(tempMetadataConfig.DocumentInformation)-1].DocumentVersioning["DocumentStatus"]

			_fileCache.MetadataConfig = tempMetadataConfig
		}

		if strings.Contains(directoryContents.Name(), utilities.ConfigFileSeverityAssessment) {
			var tempSeverityConfig utilities.SeverityAssessmentYML
			errDecodeYAML := yaml.Unmarshal(rawFileContent, &tempSeverityConfig)
			utilities.ErrorChecker(errDecodeYAML)
			validators.ValidateConfigSeverityAssessment(&tempSeverityConfig, filePath)

			slices.Reverse(tempSeverityConfig.Severities)

			if tempSeverityConfig.SwapImpactLikelihoodAxis {
				slices.Reverse(tempSeverityConfig.Likelihoods)
			} else {
				slices.Reverse(tempSeverityConfig.Impacts)
			}

			_fileCache.SeverityConfig = tempSeverityConfig
		}

		if strings.Contains(directoryContents.Name(), utilities.ConfigFileContentOrder) {
			var tempContentConfig utilities.ContentOrderYML
			errDecodeYAML := yaml.Unmarshal(rawFileContent, &tempContentConfig)
			utilities.ErrorChecker(errDecodeYAML)
			validators.ValidateConfigContentOrder(&tempContentConfig, filePath)

			_fileCache.ContentConfig = tempContentConfig
		}

		return nil
	})
	utilities.ErrorChecker(errDirectoryWalk)
}

/*
HandleModifications → Recursively walks directory tree and modifies markdown files
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file in findings/suggestions/risks directory:
    -- Skips root-level files
    -- Determines identifier prefix and reserves ID using atomic counter
    -- Calls modifiers.ModifySeverity()
    -- Calls modifiers.ModifyIdentifiers()
  - Uses waitgroup to coordinate all concurrent modifications
  - Handles errors via utilities.ErrorChecker()
*/
func HandleModifications(_reportPaths utilities.ReportPaths, _fileCache *utilities.FileCache) {
	var waitGroup sync.WaitGroup
	identifierPrefixMap, identifierCounterMap, lockedFiles := _fileCache.GetIdentifierMaps()

	directories := []struct {
		path string
	}{
		{_reportPaths.FindingsPath},
		{_reportPaths.SuggestionsPath},
		{_reportPaths.RisksPath},
	}

	for _, directory := range directories {
		waitGroup.Add(1)

		go func(dirPath string) {
			defer waitGroup.Done()

			errDirectoryWalk := filepath.WalkDir(dirPath, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
				if errAnonymousFunction != nil {
					return errAnonymousFunction
				}

				if directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" || utilities.IsRootLevelFile(filePath) {
					return nil
				}

				updatedPath := filePath

				if strings.Contains(filePath, utilities.FindingsDirectory) || strings.Contains(filePath, utilities.SuggestionsDirectory) {
					updatedPath = modifiers.ModifySeverity(filePath, _fileCache)
				}

				if len(identifierPrefixMap) == 0 {
					return nil
				}

				prefix := identifierPrefixMap[filepath.Dir(filePath)]
				if prefix == "" || lockedFiles[filePath] {
					return nil
				}

				modifiers.ModifyIdentifiers(updatedPath, prefix, atomic.AddInt32(identifierCounterMap[prefix], 1), _fileCache)

				return nil
			})

			utilities.ErrorChecker(errDirectoryWalk)

		}(directory.path)
	}

	waitGroup.Wait()
}

/*
HandleProcessing → Recursively walks directory tree and processes markdown files
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file in summaries/findings/suggestions/risks/appendices directory:
    -- Skips root-level files
    -- Calls processors.ProcessMarkdown()
  - Uses channels and collector goroutines to safely aggregate results
  - Uses waitgroups to coordinate processing and collection phases
  - Handles errors via utilities.ErrorChecker()
  - Returns processed md of type utilities.SeverityMatrix, utilities.SeverityBarGraph, utilities.MarkdownFile
*/
func HandleProcessing(_reportPaths utilities.ReportPaths, _fileCache *utilities.FileCache) {
	var waitGroup sync.WaitGroup
	var waitGroupCollection sync.WaitGroup
	var processedSeverityMatrix utilities.SeverityMatrix
	var processedSeverityBarGraph utilities.SeverityBarGraph
	var processedSummaries, processedFindings, processedSuggestions, processedRisks, processedAppendices []utilities.MarkdownFile

	severityMatrixChannel := make(chan utilities.SeverityMatrixUpdate)
	severityBarGraphChannel := make(chan utilities.SeverityBarGraphUpdate)
	summariesChannel := make(chan utilities.MarkdownFile)
	findingsChannel := make(chan utilities.MarkdownFile)
	suggestionsChannel := make(chan utilities.MarkdownFile)
	risksChannel := make(chan utilities.MarkdownFile)
	appendicesChannel := make(chan utilities.MarkdownFile)

	directories := []struct {
		path    string
		channel chan utilities.MarkdownFile
		target  *[]utilities.MarkdownFile
	}{
		{_reportPaths.SummariesPath, summariesChannel, &processedSummaries},
		{_reportPaths.FindingsPath, findingsChannel, &processedFindings},
		{_reportPaths.SuggestionsPath, suggestionsChannel, &processedSuggestions},
		{_reportPaths.RisksPath, risksChannel, &processedRisks},
		{_reportPaths.AppendicesPath, appendicesChannel, &processedAppendices},
	}

	waitGroupCollection.Add(len(directories) + 2)

	go func() {
		defer waitGroupCollection.Done()
		for severity := range severityMatrixChannel {
			if processedSeverityMatrix.Matrix[severity.RowIndex][severity.ColumnIndex] == "" {
				processedSeverityMatrix.Matrix[severity.RowIndex][severity.ColumnIndex] = severity.FindingID
			} else {
				processedSeverityMatrix.Matrix[severity.RowIndex][severity.ColumnIndex] += ", " + severity.FindingID
			}
		}
	}()

	go func() {
		defer waitGroupCollection.Done()
		for update := range severityBarGraphChannel {
			processedSeverityBarGraph.Total++
			switch update.Status {
			case "Resolved":
				processedSeverityBarGraph.Resolved++
			case "Unresolved":
				processedSeverityBarGraph.Unresolved++

				if update.Severity != "" {
					if processedSeverityBarGraph.Severities == nil {
						processedSeverityBarGraph.Severities = make(map[string]int)
					}
					processedSeverityBarGraph.Severities[update.Severity]++
				}
			}
		}
	}()

	for _, directory := range directories {
		go func(channel chan utilities.MarkdownFile, target *[]utilities.MarkdownFile) {
			defer waitGroupCollection.Done()
			for markdown := range channel {
				*target = append(*target, markdown)
			}
		}(directory.channel, directory.target)
	}

	for _, directory := range directories {
		waitGroup.Add(1)

		go func(dirPath string, markdownChannel chan utilities.MarkdownFile) {
			defer waitGroup.Done()

			errDirectoryWalk := filepath.WalkDir(dirPath, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
				if errAnonymousFunction != nil {
					return errAnonymousFunction
				}

				if directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" || utilities.IsRootLevelFile(filePath) {
					return nil
				}

				markdownFile, severityMatrixUpdate, severityBarGraphUpdate := processors.ProcessMarkdown(filePath, _fileCache)

				markdownChannel <- markdownFile

				if strings.Contains(filePath, utilities.FindingsDirectory) {
					if _fileCache.SeverityConfig.ConductSeverityAssessment && severityMatrixUpdate != nil {
						severityMatrixChannel <- *severityMatrixUpdate
					}

					if _fileCache.SeverityConfig.DisplaySeverityBarGraph && severityBarGraphUpdate != nil {
						severityBarGraphChannel <- *severityBarGraphUpdate
					}
				}

				return nil
			})

			utilities.ErrorChecker(errDirectoryWalk)

		}(directory.path, directory.channel)
	}

	waitGroup.Wait()

	close(severityMatrixChannel)
	close(severityBarGraphChannel)
	for _, directory := range directories {
		close(directory.channel)
	}

	waitGroupCollection.Wait()

	_fileCache.SeverityMatrix = processedSeverityMatrix
	_fileCache.SeverityBarGraph = processedSeverityBarGraph
	_fileCache.Summaries = processedSummaries
	_fileCache.Findings = processedFindings
	_fileCache.Suggestions = processedSuggestions
	_fileCache.Risks = processedRisks
	_fileCache.Appendices = processedAppendices
}
