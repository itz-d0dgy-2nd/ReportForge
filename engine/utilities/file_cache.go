package utilities

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

/*
cacheYMLFile → Processes YAML configuration files
*/
func (_fileCache *FileCache) cacheYMLFile(_path string, _rawFileContent []byte) {
	baseName := filepath.Base(_path)

	switch {
	case strings.Contains(baseName, ConfigFileMetadata):
		if errUnmarshal := yaml.Unmarshal(_rawFileContent, &_fileCache.metadataConfig); errUnmarshal != nil {
			Check(NewYAMLError(_path, "failed to unmarshal metadata config", errUnmarshal))
		}
		_fileCache.metadataConfig.Validate(_path)

		if len(_fileCache.metadataConfig.DocumentInformation) > 0 {
			lastDoc := _fileCache.metadataConfig.DocumentInformation[len(_fileCache.metadataConfig.DocumentInformation)-1]
			if status, exists := lastDoc.DocumentVersioning["DocumentStatus"]; exists {
				DocumentStatus = status
			}
			if version, exists := lastDoc.DocumentVersioning["DocumentVersion"]; exists {
				DocumentVersion = version
			}
		}

	case strings.Contains(baseName, ConfigFileSeverityAssessment):
		if errUnmarshal := yaml.Unmarshal(_rawFileContent, &_fileCache.severityConfig); errUnmarshal != nil {
			Check(NewYAMLError(_path, "failed to unmarshal severity config", errUnmarshal))
		}
		_fileCache.severityConfig.Validate(_path)

		slices.Reverse(_fileCache.severityConfig.Severities)
		if _fileCache.severityConfig.SwapImpactLikelihoodAxis {
			slices.Reverse(_fileCache.severityConfig.Likelihoods)
		} else {
			slices.Reverse(_fileCache.severityConfig.Impacts)
		}

	case strings.Contains(baseName, ConfigFileRiskAssessment):
		if errUnmarshal := yaml.Unmarshal(_rawFileContent, &_fileCache.riskConfig); errUnmarshal != nil {
			Check(NewYAMLError(_path, "failed to unmarshal risk config", errUnmarshal))
		}
		_fileCache.riskConfig.Validate(_path)

		slices.Reverse(_fileCache.riskConfig.GrossRiskRatings)
		slices.Reverse(_fileCache.riskConfig.GrossImpacts)
		slices.Reverse(_fileCache.riskConfig.TargetRiskRatings)
		slices.Reverse(_fileCache.riskConfig.TargetImpacts)

	case strings.Contains(baseName, ConfigFileContentOrder):
		if errUnmarshal := yaml.Unmarshal(_rawFileContent, &_fileCache.contentConfig); errUnmarshal != nil {
			Check(NewYAMLError(_path, "failed to unmarshal content order config", errUnmarshal))
		}
		_fileCache.contentConfig.Validate(_path)

	default:
		Check(NewConfigError(_path, "unknown config file - not recognised"))
	}
}

/*
cacheMDFile → Processes markdown files with YAML frontmatter
*/
func (_fileCache *FileCache) cacheMDFile(_path string, _rawFileContent []byte) {
	var fileName string
	var rawYaml map[string]any
	var unprocessedYaml MarkdownYML

	_rawFileContent = bytes.TrimPrefix(_rawFileContent, []byte{0xEF, 0xBB, 0xBF})
	_rawFileContent = bytes.TrimPrefix(_rawFileContent, []byte{0xFF, 0xFE})
	_rawFileContent = bytes.TrimPrefix(_rawFileContent, []byte{0xFE, 0xFF})
	normalisedContent := bytes.ReplaceAll(_rawFileContent, []byte("\r\n"), []byte("\n"))

	regexMatches := YAMLPattern.Frontmatter.FindSubmatch(normalisedContent)

	if len(regexMatches) < 3 || len(regexMatches[2]) == 0 {
		Check(NewValidationError(_path, "content", "missing markdown content after YAML frontmatter"))
	}

	if errUnmarshal := yaml.Unmarshal(regexMatches[1], &rawYaml); errUnmarshal != nil {
		Check(NewYAMLError(_path, "invalid frontmatter syntax", errUnmarshal))
	}

	if errUnmarshal := yaml.Unmarshal(regexMatches[1], &unprocessedYaml); errUnmarshal != nil {
		Check(NewYAMLError(_path, "frontmatter structure does not match expected schema", errUnmarshal))
	}

	unprocessedYaml.Validate(rawYaml, _path)

	switch GetDirectoryType(_path) {
	case SummariesDirectory:
		fileName = strings.TrimSpace(unprocessedYaml.ReportSummaryName)
	case FindingsDirectory:
		fileName = strings.TrimSpace(unprocessedYaml.FindingName)
	case SuggestionsDirectory:
		fileName = strings.TrimSpace(unprocessedYaml.SuggestionName)
	case RisksDirectory:
		fileName = strings.TrimSpace(unprocessedYaml.RiskName)
	case ControlsDirectory:
		fileName = strings.TrimSpace(unprocessedYaml.ControlName)
	case AppendicesDirectory:
		fileName = strings.TrimSpace(unprocessedYaml.AppendixName)
	}

	_fileCache.mutex.Lock()
	defer _fileCache.mutex.Unlock()

	if fileName != "" {
		if existing, found := _fileCache.fileNames[fileName]; found && existing != _path {
			Check(NewValidationWarning(
				_path,
				fmt.Sprintf("duplicate file name '%s' also exists in '%s' - hyperlinks may not work correctly", fileName, filepath.Base(filepath.Dir(existing))),
			))
		}
		_fileCache.fileNames[fileName] = _path
	}

	_fileCache.markdownFrontmatter[_path] = unprocessedYaml
	_fileCache.markdown[_path] = normalisedContent
}

/*
MetadataConfig → Returns pointer to metadata configuration
*/
func (_fileCache *FileCache) MetadataConfig() *MetadataYML {
	_fileCache.mutex.RLock()
	defer _fileCache.mutex.RUnlock()
	return &_fileCache.metadataConfig
}

/*
SeverityConfig → Returns pointer to severity assessment configuration
*/
func (_fileCache *FileCache) SeverityConfig() *SeverityAssessmentYML {
	_fileCache.mutex.RLock()
	defer _fileCache.mutex.RUnlock()
	return &_fileCache.severityConfig
}

/*
RiskConfig → Returns pointer to risk assessment configuration
*/
func (_fileCache *FileCache) RiskConfig() *RiskAssessmentYML {
	_fileCache.mutex.RLock()
	defer _fileCache.mutex.RUnlock()
	return &_fileCache.riskConfig
}

/*
ContentConfig → Returns pointer to content order configuration
*/
func (_fileCache *FileCache) ContentConfig() *ContentOrderYML {
	_fileCache.mutex.RLock()
	defer _fileCache.mutex.RUnlock()
	return &_fileCache.contentConfig
}

/*
ReadFile → Retrieves cached file content by path
*/
func (_fileCache *FileCache) ReadFile(_path string) []byte {
	_fileCache.mutex.RLock()
	defer _fileCache.mutex.RUnlock()
	return _fileCache.markdown[_path]
}

/*
ReadFileFrontmatter → Returns cached parsed MarkdownYML for a given path
*/
func (_fileCache *FileCache) ReadFileFrontmatter(_path string) MarkdownYML {
	_fileCache.mutex.RLock()
	defer _fileCache.mutex.RUnlock()
	return _fileCache.markdownFrontmatter[_path]
}

/*
UpdateFile → Updates cached file content
*/
func (_fileCache *FileCache) UpdateFile(_path string, _content []byte) {
	_fileCache.mutex.Lock()
	defer _fileCache.mutex.Unlock()
	_fileCache.markdown[_path] = _content
}

/*
RenameFile → Renames file in cache
*/
func (_fileCache *FileCache) RenameFile(_oldPath, _newPath string, _content []byte) {
	_fileCache.mutex.Lock()
	defer _fileCache.mutex.Unlock()
	delete(_fileCache.markdown, _oldPath)
	_fileCache.markdown[_newPath] = _content
}

/*
IterateCachedFiles → Executes function on each cached file
*/
func (_fileCache *FileCache) IterateCachedFiles(_function func(path string, content []byte)) {
	_fileCache.mutex.RLock()
	defer _fileCache.mutex.RUnlock()

	for path, content := range _fileCache.markdown {
		_function(path, content)
	}
}

/*
ReloadFile → Reloads and re-parses a file from disk (used by watcher)
*/
func (_fileCache *FileCache) ReloadFile(path string) {
	rawFileContent, errRawFileContent := os.ReadFile(path)
	if errRawFileContent != nil {
		Check(NewFileSystemError(path, "failed to read file", errRawFileContent))
	}

	switch filepath.Ext(path) {
	case ".yml", ".yaml":
		_fileCache.cacheYMLFile(path, rawFileContent)
	case ".md":
		_fileCache.cacheMDFile(path, rawFileContent)
	default:
		Check(NewValidationError(path, "extension", "unsupported file type"))
	}

	Logger.Debug("file reloaded in cache", "file", filepath.Base(path))
}

/*
ClearProcessedData → Clears all processed markdown slices and resets severity and risk matrices
*/
func (_fileCache *FileCache) ClearProcessedData() {
	_fileCache.mutex.Lock()
	defer _fileCache.mutex.Unlock()

	_fileCache.Summaries = nil
	_fileCache.Findings = nil
	_fileCache.Suggestions = nil
	_fileCache.Risks = nil
	_fileCache.Controls = nil
	_fileCache.Appendices = nil

	_fileCache.SeverityMatrix = SeverityMatrix{}
	_fileCache.SeverityBarGraph = SeverityBarGraph{}
	_fileCache.RiskMatrices = RiskMatrices{}

	Logger.Debug("processed data cleared for regeneration")
}

/*
NewFileCache → Creates and initialises a file cache for report processing
*/
func NewFileCache(_path string) *FileCache {
	var waitGroup sync.WaitGroup

	fileCache := &FileCache{
		markdown:            make(map[string][]byte),
		markdownFrontmatter: make(map[string]MarkdownYML),
		fileNames:           make(map[string]string),
		Path:                _path,
	}

	validateConfigFiles(filepath.Join(_path, ConfigDirectory))

	errDirectoryWalk := filepath.WalkDir(_path, func(path string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
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
			Check(errRawFileContent)

			switch extension {
			case ".yml":
				fileCache.cacheYMLFile(path, rawFileContent)
			case ".md":
				fileCache.cacheMDFile(path, rawFileContent)
			}

		}(path, filepath.Ext(directoryContents.Name()))
		return nil
	})

	Check(errDirectoryWalk)
	waitGroup.Wait()
	return fileCache
}
