package utilities

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

func NewFileCache(_directory string) *FileCache {
	var waitGroup sync.WaitGroup

	fileCache := &FileCache{
		cache: make(map[string][]byte),
	}

	errDirectoryWalk := filepath.WalkDir(_directory, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil {
			return errAnonymousFunction
		}

		if directoryContents.IsDir() || (filepath.Ext(directoryContents.Name()) != ".md" && filepath.Ext(directoryContents.Name()) != ".yml") {
			return nil
		}

		waitGroup.Add(1)
		go func(path string, extension string) {
			defer waitGroup.Done()

			rawFileContent, errRawFileContent := os.ReadFile(path)
			ErrorChecker(errRawFileContent)

			if extension == ".yml" && strings.Contains(path, ConfigDirectory) {
				if strings.Contains(filepath.Base(path), ConfigFileMetadata) {
					yaml.Unmarshal(rawFileContent, &fileCache.MetadataConfig)
				} else if strings.Contains(filepath.Base(path), ConfigFileSeverityAssessment) {
					yaml.Unmarshal(rawFileContent, &fileCache.SeverityConfig)
				} else if strings.Contains(filepath.Base(path), ConfigFileDirectoryOrder) {
					yaml.Unmarshal(rawFileContent, &fileCache.DirectoryConfig)
				}
			} else {
				fileCache.mutex.Lock()
				fileCache.cache[path] = rawFileContent
				fileCache.mutex.Unlock()
			}

		}(filePath, filepath.Ext(directoryContents.Name()))
		return nil
	})

	ErrorChecker(errDirectoryWalk)
	waitGroup.Wait()
	return fileCache
}

func (_fileCache *FileCache) ReadFile(_filePath string) ([]byte, error) {
	_fileCache.mutex.RLock()
	cachedContent, exists := _fileCache.cache[_filePath]
	_fileCache.mutex.RUnlock()

	if exists {
		return cachedContent, nil
	}

	rawFileContent, errRawFileContent := os.ReadFile(_filePath)
	if errRawFileContent != nil {
		return nil, errRawFileContent
	}

	_fileCache.mutex.Lock()
	_fileCache.cache[_filePath] = rawFileContent
	_fileCache.mutex.Unlock()

	return rawFileContent, nil
}

func (_fileCache *FileCache) UpdateFile(_filePath string, _content []byte) {
	_fileCache.mutex.Lock()
	defer _fileCache.mutex.Unlock()
	_fileCache.cache[_filePath] = _content
}

func (_fileCache *FileCache) RenameFile(_oldPath, _newPath string, _content []byte) {
	_fileCache.mutex.Lock()
	defer _fileCache.mutex.Unlock()
	delete(_fileCache.cache, _oldPath)
	_fileCache.cache[_newPath] = _content
}

func (_fileCache *FileCache) IterateCachedFiles(_function func(filePath string, content []byte)) {
	_fileCache.mutex.RLock()
	defer _fileCache.mutex.RUnlock()

	for path, content := range _fileCache.cache {
		_function(path, content)
	}
}

func (_fileCache *FileCache) GetIdentifierMaps() (map[string]string, map[string]*int32, map[string]bool) {
	identifierPrefixMap := make(map[string]string)
	identifierCounterMap := make(map[string]*int32)
	lockedFiles := make(map[string]bool)

	subdirectories := make(map[string]bool)
	_fileCache.IterateCachedFiles(func(filePath string, _ []byte) {
		if filepath.Ext(filePath) == ".md" && !IsRootLevelFile(filePath) {
			subdirectories[filepath.Dir(filePath)] = true
		}
	})

	directories := make([]string, 0, len(subdirectories))
	for directory := range subdirectories {
		directories = append(directories, directory)
	}

	sort.Slice(directories, func(i, j int) bool {
		parentDirectoryI := filepath.Base(filepath.Dir(directories[i]))
		parentDirectoryJ := filepath.Base(filepath.Dir(directories[j]))

		if parentDirectoryI != parentDirectoryJ {
			return parentDirectoryI < parentDirectoryJ
		}

		subdirectoryI := filepath.Base(directories[i])
		subdirectoryJ := filepath.Base(directories[j])

		var directoryOrderList []string
		switch parentDirectoryI {
		case FindingsDirectory:
			directoryOrderList = _fileCache.DirectoryConfig.Findings
		case SuggestionsDirectory:
			directoryOrderList = _fileCache.DirectoryConfig.Suggestions
		case RisksDirectory:
			directoryOrderList = _fileCache.DirectoryConfig.Risks
		}

		if len(directoryOrderList) > 0 {
			positionI := slices.Index(directoryOrderList, subdirectoryI)
			positionJ := slices.Index(directoryOrderList, subdirectoryJ)

			if positionI != -1 && positionJ != -1 {
				return positionI < positionJ
			}
			if positionI != -1 {
				return true
			}
			if positionJ != -1 {
				return false
			}
		}

		return subdirectoryI < subdirectoryJ
	})

	for index, directory := range directories {
		if index >= MaxIdentifierPrefixes {
			ErrorChecker(fmt.Errorf("too many subdirectories: maximum %d allowed (A-Z)", MaxIdentifierPrefixes))
		}
		identifierPrefix := string(rune('A' + index))
		identifierPrefixMap[directory] = identifierPrefix
		identifierCounterMap[identifierPrefix] = new(int32)
	}

	_fileCache.IterateCachedFiles(func(filePath string, rawFileContent []byte) {
		identifierPrefix := identifierPrefixMap[filepath.Dir(filePath)]
		if identifierPrefix == "" {
			return
		}

		rawMarkdownContent := strings.ReplaceAll(string(rawFileContent), "\r\n", "\n")
		regexMatches := RegexYamlMatch.FindStringSubmatch(rawMarkdownContent)
		if len(regexMatches) < 2 {
			return
		}

		var unprocessedYaml MarkdownYML
		if yaml.Unmarshal([]byte(regexMatches[1]), &unprocessedYaml) != nil {
			return
		}

		var identifier string
		var identifierLocked bool

		parentDirectory := filepath.Base(filepath.Dir(filepath.Dir(filePath)))
		switch parentDirectory {
		case FindingsDirectory:
			identifier = strings.TrimSpace(unprocessedYaml.FindingID)
			identifierLocked = unprocessedYaml.FindingIDLocked
		case SuggestionsDirectory:
			identifier = strings.TrimSpace(unprocessedYaml.SuggestionID)
			identifierLocked = unprocessedYaml.SuggestionIDLocked
		case RisksDirectory:
			identifier = strings.TrimSpace(unprocessedYaml.RiskID)
			identifierLocked = unprocessedYaml.RiskIDLocked
		default:
			return
		}

		if identifierLocked {
			lockedFiles[filePath] = true

			if strings.HasPrefix(identifier, identifierPrefix) {
				var identifierNumber int
				if _, err := fmt.Sscanf(strings.TrimPrefix(identifier, identifierPrefix), "%d", &identifierNumber); err == nil {
					if int32(identifierNumber) > *identifierCounterMap[identifierPrefix] {
						*identifierCounterMap[identifierPrefix] = int32(identifierNumber)
					}
				}
			}
		}
	})

	return identifierPrefixMap, identifierCounterMap, lockedFiles
}
