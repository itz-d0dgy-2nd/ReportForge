package generators

import (
	"ReportForge/engine/utils"
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"github.com/xuri/excelize/v2"
)

func GenerateXSLX(_findings []Utils.Markdown, suggestions []Utils.Markdown) {

	outputXLSX := excelize.NewFile()

	for _, findingMD := range _findings {

		sheetName := findingMD.Directory

		if len(sheetName) > 31 {
			sheetName = sheetName[:31]
		}

		if sheetNameExists, _ := outputXLSX.GetSheetIndex(sheetName); sheetNameExists == -1 {
			outputXLSX.NewSheet(sheetName)
			outputXLSX.SetCellValue(sheetName, "A1", "Finding Name")
			outputXLSX.SetCellValue(sheetName, "B1", "Finding Status")
			outputXLSX.SetCellValue(sheetName, "C1", "Finding Imapct")
			outputXLSX.SetCellValue(sheetName, "D1", "Finding Likelihood")
			outputXLSX.SetCellValue(sheetName, "E1", "Finding Details")
		}

		rows, _ := outputXLSX.GetRows(sheetName)
		row := len(rows) + 1

		outputXLSX.SetCellValue(findingMD.Directory, fmt.Sprintf("A%d", row), bluemonday.StrictPolicy().Sanitize(findingMD.Headers.FindingTitle))
		outputXLSX.SetCellValue(findingMD.Directory, fmt.Sprintf("B%d", row), bluemonday.StrictPolicy().Sanitize(findingMD.Headers.FindingStatus))
		outputXLSX.SetCellValue(findingMD.Directory, fmt.Sprintf("C%d", row), bluemonday.StrictPolicy().Sanitize(findingMD.Headers.FindingImpact))
		outputXLSX.SetCellValue(findingMD.Directory, fmt.Sprintf("D%d", row), bluemonday.StrictPolicy().Sanitize(findingMD.Headers.FindingLikelihood))
		outputXLSX.SetCellValue(findingMD.Directory, fmt.Sprintf("E%d", row), bluemonday.StrictPolicy().Sanitize(findingMD.Body))
	}

	outputXLSX.DeleteSheet("Sheet1")
	Utils.ErrorChecker(outputXLSX.SaveAs("Report.xlsx"))

}
