package utilities

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"golang.org/x/image/draw"
	"gopkg.in/yaml.v3"
)

/*
ParseFile → Parses a markdown file from the cache, extracting YAML frontmatter and body content
*/
func ParseFile(_path string, _fileCache *FileCache) (string, []string, MarkdownYML) {
	var unprocessedYaml MarkdownYML

	rawFileContent := _fileCache.ReadFile(_path)

	rawMarkdownContent := strings.TrimRight(string(rawFileContent), "\t\r\n")
	regexMatches := YAMLPattern.Frontmatter.FindStringSubmatch(rawMarkdownContent)

	yaml.Unmarshal([]byte(regexMatches[1]), &unprocessedYaml)

	return rawMarkdownContent, regexMatches, unprocessedYaml
}

/*
IsRootLevelFile → Checks if file is directly in findings/suggestions/risks directory (not in subdirectory)
*/
func IsRootLevelFile(_path string) bool {
	return filepath.Base(filepath.Dir(_path)) == FindingsDirectory ||
		filepath.Base(filepath.Dir(_path)) == SuggestionsDirectory ||
		filepath.Base(filepath.Dir(_path)) == RisksDirectory ||
		filepath.Base(filepath.Dir(_path)) == ControlsDirectory ||
		filepath.Base(filepath.Dir(_path)) == AppendicesDirectory
}

/*
GetDirectoryType → Determines directory type from file path
*/
func GetDirectoryType(_path string) string {
	normalisedPath := filepath.ToSlash(_path)
	parts := strings.Split(normalisedPath, "/")

	for _, part := range parts {
		switch part {
		case SummariesDirectory:
			return SummariesDirectory
		case FindingsDirectory:
			return FindingsDirectory
		case SuggestionsDirectory:
			return SuggestionsDirectory
		case RisksDirectory:
			return RisksDirectory
		case ControlsDirectory:
			return ControlsDirectory
		case AppendicesDirectory:
			return AppendicesDirectory
		}
	}
	return ""
}

/*
SortSeverityMatrix → Sorts severity matrix finding IDs alphabetically within each cell
*/
func SortSeverityMatrix(_severityMatrix *SeverityMatrix) {
	for row := 0; row < len(_severityMatrix.Matrix); row++ {
		for column := 0; column < len(_severityMatrix.Matrix[row]); column++ {
			if _severityMatrix.Matrix[row][column] != "" {
				findings := strings.Split(_severityMatrix.Matrix[row][column], ", ")
				sort.Strings(findings)
				_severityMatrix.Matrix[row][column] = strings.Join(findings, ", ")
			}
		}
	}
}

/*
SortRiskMatrices → Sorts risk matrix risk IDs alphabetically within each cell
*/
func SortRiskMatrices(_riskMatrices *RiskMatrices) {
	for row := 0; row < len(_riskMatrices.GrossMatrix); row++ {
		for column := 0; column < len(_riskMatrices.GrossMatrix[row]); column++ {
			if _riskMatrices.GrossMatrix[row][column] != "" {
				risks := strings.Split(_riskMatrices.GrossMatrix[row][column], ", ")
				sort.Strings(risks)
				_riskMatrices.GrossMatrix[row][column] = strings.Join(risks, ", ")
			}
		}
	}

	for row := 0; row < len(_riskMatrices.TargetMatrix); row++ {
		for column := 0; column < len(_riskMatrices.TargetMatrix[row]); column++ {
			if _riskMatrices.TargetMatrix[row][column] != "" {
				risks := strings.Split(_riskMatrices.TargetMatrix[row][column], ", ")
				sort.Strings(risks)
				_riskMatrices.TargetMatrix[row][column] = strings.Join(risks, ", ")
			}
		}
	}
}

/*
SortReportData → Sorts data alphabetically by directory and then by filename
*/
func SortReportData(_markdown []MarkdownFile, _parentDirectory string, _contentOrder *ContentOrderYML) {
	var contentOrderList []string

	switch _parentDirectory {
	case SummariesDirectory:
		contentOrderList = _contentOrder.Summaries
	case FindingsDirectory:
		contentOrderList = _contentOrder.Findings
	case SuggestionsDirectory:
		contentOrderList = _contentOrder.Suggestions
	case RisksDirectory:
		contentOrderList = _contentOrder.Risks
	case ControlsDirectory:
		contentOrderList = _contentOrder.Controls
	}

	if len(contentOrderList) == 0 {
		sort.Slice(_markdown, func(i, j int) bool {
			if _markdown[i].Directory == _markdown[j].Directory {
				return _markdown[i].FileName < _markdown[j].FileName
			}
			return _markdown[i].Directory < _markdown[j].Directory
		})
		return
	}

	contentOrderMap := make(map[string]int, len(contentOrderList))
	for index, directory := range contentOrderList {
		contentOrderMap[directory] = index
	}

	sort.Slice(_markdown, func(i, j int) bool {
		directoryI := _markdown[i].Directory
		directoryJ := _markdown[j].Directory

		if directoryI == directoryJ {
			return _markdown[i].FileName < _markdown[j].FileName
		}

		priorityI, orderedI := contentOrderMap[directoryI]
		priorityJ, orderedJ := contentOrderMap[directoryJ]

		if orderedI && orderedJ {
			return priorityI < priorityJ
		}

		if orderedI {
			return true
		}

		if orderedJ {
			return false
		}

		return directoryI < directoryJ
	})
}

/*
OptimiseImagesForPDF → Resizes images for PDF optimisation in Release mode only.
*/
func OptimiseImagesForPDF(_path string) {
	if DocumentStatus != ReportStatusRelease {
		return
	}

	var waitGroup sync.WaitGroup

	errDirectoryWalk := filepath.WalkDir(_path, func(path string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() {
			return errAnonymousFunction
		}

		if filepath.Base(filepath.Dir(path)) == ScreenshotsOriginalsDirectory {
			return nil
		}

		extension := strings.ToLower(filepath.Ext(directoryContents.Name()))
		if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
			return nil
		}

		waitGroup.Add(1)
		go func(path string) {
			defer waitGroup.Done()

			originalsDirectory := filepath.Join(filepath.Dir(path), ScreenshotsOriginalsDirectory)
			backupPath := filepath.Join(originalsDirectory, filepath.Base(path))

			if _, errStat := os.Stat(backupPath); errStat == nil {
				return
			}

			rawFileContent, errRawFileContent := os.ReadFile(path)
			if errRawFileContent != nil {
				Check(NewFileSystemError(path, "failed to read image file", errRawFileContent))
			}

			decodedImage, imageFormat, errDecoding := image.Decode(bytes.NewReader(rawFileContent))
			if errDecoding != nil {
				Check(NewProcessingError(path, fmt.Sprintf("failed to decode image: %s", errDecoding.Error())))
			}

			originalBounds := decodedImage.Bounds()
			originalWidth := originalBounds.Dx()

			if originalWidth <= PDFOptimalImageWidth {
				return
			}

			if errMkdir := os.MkdirAll(originalsDirectory, 0755); errMkdir != nil {
				Check(NewFileSystemError(originalsDirectory, "failed to create originals directory", errMkdir))
			}

			if errWriteBackup := os.WriteFile(backupPath, rawFileContent, 0644); errWriteBackup != nil {
				Check(NewFileSystemError(backupPath, "failed to write backup image", errWriteBackup))
			}

			targetWidth := PDFOptimalImageWidth
			targetHeight := (originalBounds.Dy() * PDFOptimalImageWidth) / originalWidth

			resizedImage := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
			draw.CatmullRom.Scale(resizedImage, resizedImage.Bounds(), decodedImage, originalBounds, draw.Over, nil)

			tmpPath := path + ".tmp"
			outputFile, errCreate := os.Create(tmpPath)
			if errCreate != nil {
				Check(NewFileSystemError(tmpPath, "failed to create temporary image file", errCreate))
			}

			var errEncoding error
			if imageFormat == "png" {
				errEncoding = (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(outputFile, resizedImage)
			} else {
				errEncoding = jpeg.Encode(outputFile, resizedImage, &jpeg.Options{Quality: ImageCompressionQuality})
			}

			if errClose := outputFile.Close(); errEncoding == nil && errClose != nil {
				errEncoding = errClose
			}

			if errEncoding != nil {
				os.Remove(tmpPath)
				Check(NewProcessingError(tmpPath, fmt.Sprintf("failed to encode image: %s", errEncoding.Error())))
				return
			}

			if errRename := os.Rename(tmpPath, path); errRename != nil {
				os.Remove(tmpPath)
				os.Rename(backupPath, path)
				Check(NewFileSystemError(path, "failed to replace original image with optimised version", errRename))
			}
		}(path)

		return nil
	})

	waitGroup.Wait()
	Check(errDirectoryWalk)
}
