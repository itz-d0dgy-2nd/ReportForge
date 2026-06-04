package handlers

import (
	"ReportForge/engine/modifiers"
	"ReportForge/engine/utilities"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

/*
modificationWorker → Processes modification jobs from the jobs channel
*/
func modificationWorker(_jobs <-chan utilities.ModificationJob, _fileCache *utilities.FileCache, _waitGroup *sync.WaitGroup) {
	defer _waitGroup.Done()

	for job := range _jobs {
		updatedPath := job.Path

		switch job.DirectoryType {
		case utilities.FindingsDirectory:
			updatedPath = modifiers.ModifyFindingFiles(job.Path, _fileCache)
		case utilities.SuggestionsDirectory:
			updatedPath = modifiers.ModifySuggestionFiles(job.Path, _fileCache)
		case utilities.RisksDirectory:
			updatedPath = modifiers.ModifyRiskFiles(job.Path, _fileCache)
		case utilities.ControlsDirectory:
			updatedPath = modifiers.ModifyControlFiles(job.Path, _fileCache)
		}

		if !job.IsLocked && job.IdentifierPrefix != "" {
			modifiers.ModifyIdentifiers(
				updatedPath,
				job.IdentifierPrefix,
				job.AssignedID,
				_fileCache,
			)
		}
	}
}

/*
assignModificationJobs → Walks a directory and sends modification jobs with pre-assigned sequential identifiers
*/
func assignModificationJobs(_path string, _directory string, identifierPrefixMap map[string]string, identifierCounterMap map[string]*int32, lockedFiles map[string]bool, jobs chan<- utilities.ModificationJob) {
	type fileInfo struct {
		path      string
		directory string
		prefix    string
		isLocked  bool
	}

	var allFiles []fileInfo

	errDirectoryWalk := filepath.WalkDir(_path, func(path string, entry fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil {
			return errAnonymousFunction
		}

		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" || utilities.IsRootLevelFile(path) {
			return nil
		}

		directory := filepath.Dir(path)
		prefix := identifierPrefixMap[directory]
		isLocked := lockedFiles[path]

		allFiles = append(allFiles, fileInfo{
			path:      path,
			directory: directory,
			prefix:    prefix,
			isLocked:  isLocked,
		})

		return nil
	})

	for _, file := range allFiles {
		var assignedID int32

		if !file.isLocked && file.prefix != "" {
			counter := identifierCounterMap[file.prefix]
			*counter++
			assignedID = *counter
		}

		jobs <- utilities.ModificationJob{
			Path:             file.path,
			DirectoryType:    _directory,
			AssignedID:       assignedID,
			IdentifierPrefix: file.prefix,
			IsLocked:         file.isLocked,
		}
	}
	utilities.Check(errDirectoryWalk)
}

/*
buildPrefixMaps → Builds identifier prefix maps for markdown files
*/
func buildPrefixMaps(_fileCache *utilities.FileCache) (map[string]string, map[string]*int32, map[string]bool) {
	identifierPrefixMap := make(map[string]string)
	identifierCounterMap := make(map[string]*int32)
	lockedFiles := make(map[string]bool)

	contentConfig := _fileCache.ContentConfig()

	_fileCache.IterateCachedFiles(func(path string, rawFileContent []byte) {
		if filepath.Ext(path) != ".md" || utilities.IsRootLevelFile(path) {
			return
		}

		directory := filepath.Dir(path)
		subdirectoryName := filepath.Base(directory)
		parentDirectory := filepath.Base(filepath.Dir(directory))

		type directoryConfig struct {
			prefixMap  map[string]string
			identifier string
			locked     bool
		}

		unprocessedYaml := _fileCache.ReadFileFrontmatter(path)

		var directorySettings directoryConfig
		switch parentDirectory {
		case utilities.FindingsDirectory:
			directorySettings = directoryConfig{contentConfig.FindingIdentifierPrefixes, strings.TrimSpace(unprocessedYaml.FindingID), unprocessedYaml.FindingIDLocked}
		case utilities.SuggestionsDirectory:
			directorySettings = directoryConfig{contentConfig.SuggestionIdentifierPrefixes, strings.TrimSpace(unprocessedYaml.SuggestionID), unprocessedYaml.SuggestionIDLocked}
		case utilities.RisksDirectory:
			directorySettings = directoryConfig{contentConfig.RiskIdentifierPrefixes, strings.TrimSpace(unprocessedYaml.RiskID), unprocessedYaml.RiskIDLocked}
		case utilities.ControlsDirectory:
			directorySettings = directoryConfig{contentConfig.ControlIdentifierPrefixes, strings.TrimSpace(unprocessedYaml.ControlID), unprocessedYaml.ControlIDLocked}
		default:
			return
		}

		if _, exists := identifierPrefixMap[directory]; !exists {
			prefix, exists := directorySettings.prefixMap[subdirectoryName]
			if !exists {
				utilities.Check(utilities.NewConfigError(
					_fileCache.Path,
					fmt.Sprintf("no identifier prefix configured for subdirectory '%s' in %s", subdirectoryName, parentDirectory),
				))
				return
			}
			identifierPrefixMap[directory] = prefix
			identifierCounterMap[prefix] = new(int32)
		}

		if directorySettings.locked {
			lockedFiles[path] = true
			identifierPrefix := identifierPrefixMap[directory]

			if strings.HasPrefix(directorySettings.identifier, identifierPrefix) {
				var identifierNumber int
				if _, err := fmt.Sscanf(strings.TrimPrefix(directorySettings.identifier, identifierPrefix), "%d", &identifierNumber); err == nil {
					if int32(identifierNumber) > *identifierCounterMap[identifierPrefix] {
						*identifierCounterMap[identifierPrefix] = int32(identifierNumber)
					}
				}
			}
		}
	})

	return identifierPrefixMap, identifierCounterMap, lockedFiles
}

/*
HandleModifications → Walks findings, suggestions, and risks directories concurrently using worker pool
*/
func HandleModifications(_reportPaths utilities.ReportPaths, _fileCache *utilities.FileCache) {
	var workersWaitgroup sync.WaitGroup

	workers := runtime.NumCPU()
	jobs := make(chan utilities.ModificationJob, workers*2)
	identifierPrefixMap, identifierCounterMap, lockedFiles := buildPrefixMaps(_fileCache)

	for i := 0; i < workers; i++ {
		workersWaitgroup.Add(1)
		go modificationWorker(jobs, _fileCache, &workersWaitgroup)
	}

	directories := map[string]string{
		_reportPaths.FindingsPath:    utilities.FindingsDirectory,
		_reportPaths.SuggestionsPath: utilities.SuggestionsDirectory,
		_reportPaths.RisksPath:       utilities.RisksDirectory,
		_reportPaths.ControlsPath:    utilities.ControlsDirectory,
	}

	for path, directory := range directories {
		if path == "" {
			continue
		}
		assignModificationJobs(path, directory, identifierPrefixMap, identifierCounterMap, lockedFiles, jobs)
	}

	close(jobs)
	workersWaitgroup.Wait()
}
