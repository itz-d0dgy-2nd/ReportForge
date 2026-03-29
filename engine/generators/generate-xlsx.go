package generators

import (
	"ReportForge/engine/utilities"
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"github.com/xuri/excelize/v2"
)

/*
GenerateXLSX → Generate final XLSX spreadsheet report from processed report data
*/
func GenerateXLSX(_fileCache *utilities.FileCache) {
	metadataConfig := _fileCache.MetadataConfig()
	outputSpreadsheet := excelize.NewFile()
	sanitiser := bluemonday.StrictPolicy()
	rowCounters := make(map[string]int)

	for _, finding := range _fileCache.Findings {
		sheetName := finding.Directory

		if len(sheetName) > utilities.ExcelMaxSheetNameLength {
			sheetName = sheetName[:utilities.ExcelMaxSheetNameLength]
		}

		sheetIndex, errGetSheetIndex := outputSpreadsheet.GetSheetIndex(sheetName)
		if errGetSheetIndex != nil {
			utilities.Check(utilities.NewProcessingError(
				"",
				fmt.Sprintf("failed to check if sheet '%s' exists: %s", sheetName, errGetSheetIndex.Error()),
			))
		}

		if sheetIndex == -1 {
			_, errNewSheet := outputSpreadsheet.NewSheet(sheetName)
			if errNewSheet != nil {
				utilities.Check(utilities.NewProcessingError(
					"",
					fmt.Sprintf("failed to create sheet '%s': %s", sheetName, errNewSheet.Error()),
				))
			}

			headers := []string{"Finding Name", "Finding Status", "Finding Impact", "Finding Likelihood", "Finding Details"}
			for i, header := range headers {
				cell := fmt.Sprintf("%c1", 'A'+i)
				if errSetCellValue := outputSpreadsheet.SetCellValue(sheetName, cell, header); errSetCellValue != nil {
					utilities.Check(utilities.NewProcessingError(
						"",
						fmt.Sprintf("failed to set header '%s' in cell %s: %s", header, cell, errSetCellValue.Error()),
					))
				}
			}

			rowCounters[sheetName] = 1
		}

		rowCounters[sheetName]++
		row := rowCounters[sheetName]

		values := []string{
			sanitiser.Sanitize(finding.Headers.FindingTitle),
			sanitiser.Sanitize(finding.Headers.FindingStatus),
			sanitiser.Sanitize(finding.Headers.FindingImpact),
			sanitiser.Sanitize(finding.Headers.FindingLikelihood),
			sanitiser.Sanitize(finding.Body),
		}

		for i, value := range values {
			cell := fmt.Sprintf("%c%d", 'A'+i, row)
			if errSetCellValue := outputSpreadsheet.SetCellValue(sheetName, cell, value); errSetCellValue != nil {
				utilities.Check(utilities.NewProcessingError(
					"",
					fmt.Sprintf("failed to set value in cell %s for finding '%s': %s", cell, finding.Headers.FindingName, errSetCellValue.Error()),
				))
			}
		}
	}

	if errDeleteSheet := outputSpreadsheet.DeleteSheet("Sheet1"); errDeleteSheet != nil {
		utilities.Check(utilities.NewProcessingError(
			"",
			fmt.Sprintf("failed to delete default 'Sheet1': %s", errDeleteSheet.Error()),
		))
	}

	xlsxFileName := fmt.Sprintf("%s_%s.xlsx", metadataConfig.DocumentName, utilities.DocumentVersion)
	if errSaveAs := outputSpreadsheet.SaveAs(xlsxFileName); errSaveAs != nil {
		utilities.Check(utilities.NewFileSystemError(
			xlsxFileName,
			"failed to save Excel workbook",
			errSaveAs,
		))
	}
}
