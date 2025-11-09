package handlers

import (
	"ReportForge/engine/helpers"
	"ReportForge/engine/modifiers"
	"ReportForge/engine/processors"
	"ReportForge/engine/utilities"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

/*
HandleConfigProcessor → Recursively walks directory tree and processes YAML config files
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .yml file containing "metadata" or "severity_assessment" in filename → calls processors.ProcessConfigMetadata() or processors.ProcessConfigSeverityAssessment()
  - Handles errors via utilities.ErrorChecker()
  - Returns processed yml of type utilities.MetadataYML, utilities.SeverityAssessmentYML
*/
func HandleConfigProcessor(_directory string) (utilities.MetadataYML, utilities.SeverityAssessmentYML) {
	var processedMetadata utilities.MetadataYML
	var processedSeverityAssessment utilities.SeverityAssessmentYML

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

	utilities.ErrorChecker(errDirectoryWalk)

	return processedMetadata, processedSeverityAssessment
}

/*
HandleSeverityModifier → Recursively walks directory tree and modifies markdown files YAML frontmatter
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file within "2_findings" or "3_suggestions" directories → modifiers.ModifySeverity()
    -- Uses waitgroup to wait for all goroutines to finish before exiting
    -- Uses semaphore to limit concurrent goroutines to runtime.NumCPU()*2
  - Handles errors via utilities.ErrorChecker()
*/
func HandleSeverityModifier(_directory string, _severityAssessment utilities.SeverityAssessmentYML) {
	var waitGroup sync.WaitGroup
	var semaphore = make(chan struct{}, runtime.NumCPU()*2)

	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" {
			return errAnonymousFunction
		}

		if strings.Contains(filePath, "2_findings") || strings.Contains(filePath, "3_suggestions") {
			waitGroup.Add(1)
			semaphore <- struct{}{}

			go func(path string) {
				defer waitGroup.Done()
				defer func() { <-semaphore }()
				modifiers.ModifySeverity(path, _severityAssessment)
			}(filePath)
		}
		return nil
	})

	utilities.ErrorChecker(errDirectoryWalk)
	waitGroup.Wait()
}

/*
HandleIdentifierModifier → Recursively walks directory tree and modifies markdown files YAML frontmatter
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file within "2_findings" or "3_suggestions" directories → modifiers.ModifyIdentifiers()
  - Handles errors via utilities.ErrorChecker()
  - TODO: Remove the need for the helper function ( PrefixMapHelper() )...
  - TODO: Remove the need for the helper function ( TrackedLockedHelper() )...
*/
func HandleIdentifierModifier(_directory string, _metadata utilities.MetadataYML) {
	var documentStatus string
	var identifierPrefixMap = helpers.PrefixMapHelper(_directory)
	var identifierCounterMap = make(map[string]*int)

	if len(identifierPrefixMap) == 0 {
		return
	}

	for _, identifierPrefix := range identifierPrefixMap {
		identifierCounterMap[identifierPrefix] = new(int)
	}

	helpers.TrackedLockedHelper(_directory, identifierPrefixMap, identifierCounterMap)

	for _, documentInformation := range _metadata.DocumentInformation {
		if documentInformation.DocumentCurrent {
			documentStatus = documentInformation.DocumentVersioning["DocumentStatus"]
			break
		}
	}
	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" {
			return errAnonymousFunction
		}

		if strings.Contains(filePath, "2_findings") || strings.Contains(filePath, "3_suggestions") {
			identifierPrefix := identifierPrefixMap[filepath.Dir(filePath)]
			modifiers.ModifyIdentifiers(filePath, identifierPrefix, identifierCounterMap[identifierPrefix], documentStatus)
		}

		return nil
	})

	utilities.ErrorChecker(errDirectoryWalk)
}

/*
HandleImageModifier → Recursively walks directory tree and compresses images when document status is Release
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .jpg, .jpeg, .png file within "Screenshots" directories → modifiers.ModifyImage()
    -- Uses waitgroup to wait for all goroutines to finish before exiting
    -- Uses semaphore to limit concurrent goroutines to runtime.NumCPU()*2
  - Handles errors via utilities.ErrorChecker()
*/
func HandleImageModifier(_directory string, _metadata utilities.MetadataYML) {
	var documentStatus string
	var waitGroup sync.WaitGroup
	var semaphore = make(chan struct{}, runtime.NumCPU()*2)

	for _, documentInformation := range _metadata.DocumentInformation {
		if documentInformation.DocumentCurrent {
			documentStatus = documentInformation.DocumentVersioning["DocumentStatus"]
			break
		}
	}

	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || (strings.ToLower(filepath.Ext(directoryContents.Name())) != ".jpg" && strings.ToLower(filepath.Ext(directoryContents.Name())) != ".jpeg" && strings.ToLower(filepath.Ext(directoryContents.Name())) != ".png") {
			return errAnonymousFunction
		}

		if strings.Contains(filePath, "Screenshots") {
			waitGroup.Add(1)
			semaphore <- struct{}{}

			go func(path string) {
				defer waitGroup.Done()
				defer func() { <-semaphore }()
				modifiers.ModifyImage(path, documentStatus)
			}(filePath)
		}
		return nil
	})

	utilities.ErrorChecker(errDirectoryWalk)
	waitGroup.Wait()
}

/*
HandleSeverityProcessor → Recursively walks directory tree and processes markdown files for severity assessment
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file → calls processors.ProcessSeverityMatrix()
  - Handles errors via utilities.ErrorChecker()
  - Returns processed YAML of type utilities.SeverityAssessmentYML
*/
func HandleSeverityProcessor(_directory string, _severityAssessment utilities.SeverityAssessmentYML) utilities.SeverityAssessmentYML {
	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" {
			return errAnonymousFunction
		}

		processors.ProcessSeverityMatrix(filePath, &_severityAssessment)

		return nil
	})

	utilities.ErrorChecker(errDirectoryWalk)

	return _severityAssessment
}

/*
HandleMarkdownProcessor → Recursively walks directory tree and processes markdown files concurrently
  - Walks through all subdirectories using filepath.WalkDir()
  - For each .md file → calls processors.ProcessMarkdown()
    -- Uses waitgroup to wait for all goroutines to finish before exiting
    -- Uses mutex to safely update shared data []utilities.Markdown
    -- Uses semaphore to limit concurrent goroutines to runtime.NumCPU()*2
  - Handles errors via utilities.ErrorChecker()
  - Returns processed markdown of type []utilities.Markdown
*/
func HandleMarkdownProcessor(_reportPath string, _directory string, _metadata utilities.MetadataYML) []utilities.Markdown {
	var processedMarkdown []utilities.Markdown
	var waitGroup sync.WaitGroup
	var mutexLock sync.Mutex
	var semaphore = make(chan struct{}, runtime.NumCPU()*2)

	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() || filepath.Ext(directoryContents.Name()) != ".md" {
			return errAnonymousFunction
		}

		waitGroup.Add(1)
		semaphore <- struct{}{}

		go func(path string) {
			defer waitGroup.Done()
			defer func() { <-semaphore }()
			processors.ProcessMarkdown(_reportPath, path, &processedMarkdown, _metadata, &mutexLock)
		}(filePath)

		return nil
	})

	utilities.ErrorChecker(errDirectoryWalk)
	waitGroup.Wait()

	return processedMarkdown
}
