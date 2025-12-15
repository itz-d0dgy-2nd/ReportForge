package utilities

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

/*
IsRootLevelFile → Checks if file is directly in findings/suggestions/risks directory (not in subdirectory)
*/
func IsRootLevelFile(_filePath string) bool {
	return filepath.Base(filepath.Dir(_filePath)) == FindingsDirectory ||
		filepath.Base(filepath.Dir(_filePath)) == SuggestionsDirectory ||
		filepath.Base(filepath.Dir(_filePath)) == RisksDirectory ||
		filepath.Base(filepath.Dir(_filePath)) == AppendicesDirectory
}

/*
SortSeverityMatrix → Sorts ReportForge severity matrix finding IDs alphabetically within each cell
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
SortReportData → Sorts ReportForge data alphabetically by directory and then by filename
*/
func SortReportData(_markdown []MarkdownFile, _parentDirectory string, _directoryOrder DirectoryOrderYML) {
	var directoryOrderList []string
	var directoryOrderMap map[string]int

	switch _parentDirectory {
	case SummariesDirectory:
		directoryOrderList = _directoryOrder.Summaries
	case FindingsDirectory:
		directoryOrderList = _directoryOrder.Findings
	case SuggestionsDirectory:
		directoryOrderList = _directoryOrder.Suggestions
	case RisksDirectory:
		directoryOrderList = _directoryOrder.Risks
	}

	if len(directoryOrderList) > 0 {
		directoryOrderMap = make(map[string]int, len(directoryOrderList))
		for index, directory := range directoryOrderList {
			directoryOrderMap[directory] = index
		}
	}

	sort.Slice(_markdown, func(i, j int) bool {
		directoryI := _markdown[i].Directory
		directoryJ := _markdown[j].Directory

		if directoryI == directoryJ {
			return _markdown[i].FileName < _markdown[j].FileName
		}

		if directoryOrderMap == nil {
			return directoryI < directoryJ
		}

		priorityI, orderedI := directoryOrderMap[directoryI]
		priorityJ, orderedJ := directoryOrderMap[directoryJ]

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
GetDirectoryType → Determines directory type from file path
*/
func GetDirectoryType(_filePath string) string {
	switch {
	case strings.Contains(_filePath, FindingsDirectory):
		return FindingsDirectory
	case strings.Contains(_filePath, SuggestionsDirectory):
		return SuggestionsDirectory
	case strings.Contains(_filePath, RisksDirectory):
		return RisksDirectory
	default:
		return ""
	}
}

/*
OptimiseImagesForPDF → Resizes images for PDF optimisation in Release mode only
*/
func OptimiseImagesForPDF(_screenshotsPath string) {
	var waitGroup sync.WaitGroup

	if DocumentStatus != ReportStatusRelease {
		return
	}

	errDirectoryWalk := filepath.WalkDir(_screenshotsPath, func(filePath string, directoryContents fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil || directoryContents.IsDir() {
			return errAnonymousFunction
		}

		if filepath.Base(filepath.Dir(filePath)) == ScreenshotsOriginalsDirectory {
			return nil
		}

		currentFileExtension := strings.ToLower(filepath.Ext(directoryContents.Name()))

		if currentFileExtension == ".jpg" || currentFileExtension == ".jpeg" || currentFileExtension == ".png" {
			waitGroup.Add(1)

			go func(path string) {
				defer waitGroup.Done()

				originalsDirectory := filepath.Join(filepath.Dir(path), ScreenshotsOriginalsDirectory)
				originalBackupPath := filepath.Join(originalsDirectory, filepath.Base(path))

				_, errOriginalBackupStatistics := os.Stat(originalBackupPath)
				if errOriginalBackupStatistics == nil {
					return
				}

				rawFileContent, errReadFile := os.ReadFile(path)
				ErrorChecker(errReadFile)

				decodedImage, imageFormat, errDecodedImage := image.Decode(bytes.NewReader(rawFileContent))
				ErrorChecker(errDecodedImage)

				originalBounds := decodedImage.Bounds()
				originalWidth := originalBounds.Dx()

				if originalWidth <= PDFOptimalImageWidth {
					return
				}

				errMakeDirectory := os.MkdirAll(originalsDirectory, 0755)
				ErrorChecker(errMakeDirectory)

				errWriteBackup := os.WriteFile(originalBackupPath, rawFileContent, 0644)
				ErrorChecker(errWriteBackup)

				targetWidth := PDFOptimalImageWidth
				targetHeight := (originalBounds.Dy() * PDFOptimalImageWidth) / originalWidth
				resizedImage := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

				for resizedImageY := 0; resizedImageY < targetHeight; resizedImageY++ {
					for resizedImageX := 0; resizedImageX < targetWidth; resizedImageX++ {
						resizedImage.Set(resizedImageX, resizedImageY, decodedImage.At((resizedImageX*originalWidth)/targetWidth, (resizedImageY*originalBounds.Dy())/targetHeight))
					}
				}

				temporaryFilePath := path + ".tmp"
				outputFile, errOutputFile := os.Create(temporaryFilePath)
				ErrorChecker(errOutputFile)

				if imageFormat == "png" {
					ErrorChecker((&png.Encoder{CompressionLevel: png.BestCompression}).Encode(outputFile, resizedImage))
				} else {
					ErrorChecker(jpeg.Encode(outputFile, resizedImage, &jpeg.Options{Quality: ImageCompressionQuality}))
				}

				errCloseFile := outputFile.Close()
				ErrorChecker(errCloseFile)

				errRenameFile := os.Rename(temporaryFilePath, path)
				if errRenameFile != nil {
					os.Remove(temporaryFilePath)
					os.Rename(originalBackupPath, path)
					ErrorChecker(errRenameFile)
				}

			}(filePath)
		}

		return nil
	})

	ErrorChecker(errDirectoryWalk)
	waitGroup.Wait()
}
