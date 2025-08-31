package handlers

import (
	Utils "ReportForge/engine/utils"
	"ReportForge/engine/utils/modifiers"
	"ReportForge/engine/utils/processors"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

/*
ModifierFileHandler → Recursively walks directory tree and modifies markdown files YAML frontmatter
  - Initialises prefixMap and counterMap of modifiers.SetupPrefixMap()and map[string]*int
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file containing "2_findings" or "3_suggestions" in filename → modifiers.ModifyIdentifiers(), modifiers.ModifySeverity()
  - Handles errors via Utils.ErrorChecker()
  - TODO: Remove the need for the helper function ( SetupPrefixMap() )...
*/
func ModifierFileHandler(_directory string, _severityAssessment Utils.SeverityAssessmentYML) {
	identifierPrefixMap := Utils.SetupPrefixMap(_directory)
	identifierCounterMap := make(map[string]*int)

	if len(identifierPrefixMap) == 0 {
		return
	}

	for _, identifierPrefix := range identifierPrefixMap {
		identifierCounterMap[identifierPrefix] = new(int)
	}

	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" {
			return errAnonymousFunction
		}

		identifierPrefix := identifierPrefixMap[filepath.Dir(filePath)]

		if strings.Contains(filePath, "2_findings") || strings.Contains(filePath, "3_suggestions") {
			modifiers.ModifyIdentifiers(filePath, identifierPrefix, identifierCounterMap[identifierPrefix])
			modifiers.ModifySeverity(filePath, _severityAssessment)
		}

		return nil
	})

	Utils.ErrorChecker(errDirectoryWalk)
}

/*
ConfigFileHandler → Recursively walks directory tree and processes YAML config files
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .yml files containing "metadata" or "severity_assessment" in filename → calls processors.ProcessConfigMetadata(), processors.ProcessConfigSeverityAssessment()
  - Handles errors via Utils.ErrorChecker()
  - Returns processed yml of type Utils.MetadataYML{}, Utils.SeverityAssessmentYML{}
*/
func ConfigFileHandler(_directory string) (Utils.MetadataYML, Utils.SeverityAssessmentYML) {
	processedMetadata := Utils.MetadataYML{}
	processedSeverityAssessment := Utils.SeverityAssessmentYML{}

	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".yml" {
			return errAnonymousFunction
		}

		if strings.Contains(directoryContents.Name(), "metadata") {
			processors.ProcessConfigMetadata(filePath, &processedMetadata)
		}

		if strings.Contains(directoryContents.Name(), "severity_assessment") {
			processors.ProcessConfigSeverityAssessment(filePath, &processedSeverityAssessment)
		}

		return nil
	})

	Utils.ErrorChecker(errDirectoryWalk)

	return processedMetadata, processedSeverityAssessment
}

/*
SeverityFileHandler → Recursively walks directory tree and processes markdown files for severity assessment
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file → calls processors.ProcessSeverityMatrix()
  - Handles errors via Utils.ErrorChecker()
  - Returns processed YAML of type []Utils.SeverityAssessmentYML{}
*/
func SeverityFileHandler(_directory string, _severityAssessment Utils.SeverityAssessmentYML) Utils.SeverityAssessmentYML {
	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" {
			return errAnonymousFunction
		}

		processors.ProcessSeverityMatrix(filePath, &_severityAssessment)

		return nil
	})

	Utils.ErrorChecker(errDirectoryWalk)

	return _severityAssessment
}

/*
MarkdownFileHandler → Recursively walks directory tree and processes markdown files concurrently
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file → calls processors.ProcessMarkdown()
  - Handles errors via Utils.ErrorChecker()
  - Returns processed markdown of type []Utils.Markdown{}
*/
func MarkdownFileHandler(_reportPath string, _directory string, _metadata Utils.MetadataYML) []Utils.Markdown {
	processedMarkdown := []Utils.Markdown{}
	waitGroup := sync.WaitGroup{}
	mutexLock := sync.Mutex{}

	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" {
			return errAnonymousFunction
		}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			processors.ProcessMarkdown(_reportPath, filePath, &processedMarkdown, _metadata, &mutexLock)
		}()

		return nil
	})

	Utils.ErrorChecker(errDirectoryWalk)
	waitGroup.Wait()

	return processedMarkdown
}
