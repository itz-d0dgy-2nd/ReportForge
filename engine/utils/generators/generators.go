package generators

import (
	Utils "ReportForge/engine/utils"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"text/template"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/xuri/excelize/v2"
)

/*
GenerateHTML → Generate final HTML report from processed report data
  - Sorts findings and suggestions slice's by directory name (primary) and filename (secondary) for consistent ordering
  - Parses "template.html" file and creates "Report.html"
  - Handles errors via Utils.ErrorChecker()
*/
func GenerateHTML(_reportData Utils.ReportDataStruct, _reportPaths Utils.ReportPathsStruct) {

	sort.Slice(_reportData.Findings, func(i, j int) bool {
		if _reportData.Findings[i].Directory != _reportData.Findings[j].Directory {
			return _reportData.Findings[i].Directory < _reportData.Findings[j].Directory
		}
		return _reportData.Findings[i].FileName < _reportData.Findings[j].FileName
	})

	sort.Slice(_reportData.Suggestions, func(i, j int) bool {
		if _reportData.Suggestions[i].Directory != _reportData.Suggestions[j].Directory {
			return _reportData.Suggestions[i].Directory < _reportData.Suggestions[j].Directory
		}
		return _reportData.Suggestions[i].FileName < _reportData.Suggestions[j].FileName
	})

	templateHTML, errTemplateHTML := template.ParseFiles(_reportPaths.TemplatePath)
	Utils.ErrorChecker(errTemplateHTML)

	createHTML, errCreateHTML := os.Create("Report.html")
	Utils.ErrorChecker(errCreateHTML)
	defer createHTML.Close()

	Utils.ErrorChecker(templateHTML.Execute(createHTML, _reportData))
}

/*
GeneratePDF → Generate final PDF report from processed report data
  - Configures chromedp with headless browser options (disable GPU, extensions, audio, etc.)
  - Performs system-specific configuration on:
    -- Windows systems - Detects Microsoft Edge and uses as Chromium executable
    -- CICD / Action systems - Adds no-sandbox flag when running in CICD/Action
  - Parses "Report.html" file and creates "Report.pdf
  - Handles errors via Utils.ErrorChecker()
*/
func GeneratePDF(_reportPaths Utils.ReportPathsStruct) {
	PDFBuffer := []byte{}

	chromiumExecutionOptions := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-audio", true),
		chromedp.Flag("disable-webgl", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("no-first-run", true),
	}

	if runtime.GOOS == "windows" {
		edgePaths := []string{
			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
		}
		for _, path := range edgePaths {
			if _, err := os.Stat(path); err == nil {
				chromiumExecutionOptions = append(chromiumExecutionOptions, chromedp.ExecPath(path))
				break
			}
		}
	}

	if os.Getenv("ACTION") == "true" {
		chromiumExecutionOptions = append(chromiumExecutionOptions, chromedp.Flag("no-sandbox", true))
	}

	chromiumExecutionContext, chromiumExecutionContextCancel := chromedp.NewExecAllocator(context.Background(), chromiumExecutionOptions...)
	defer chromiumExecutionContextCancel()

	chromiumBrowserContext, chromiumBrowserContextCancel := chromedp.NewContext(chromiumExecutionContext)
	defer chromiumBrowserContextCancel()

	absoluteFilePath, errAbsoluteFilePath := filepath.Abs(filepath.Join("Report.html"))
	Utils.ErrorChecker(errAbsoluteFilePath)

	Utils.ErrorChecker(chromedp.Run(chromiumBrowserContext,
		chromedp.Navigate("file:///"+absoluteFilePath),
		chromedp.ActionFunc(func(context context.Context) error {
			print, _, errPrint := page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPreferCSSPageSize(true).
				Do(context)
			Utils.ErrorChecker(errPrint)

			PDFBuffer = print
			return nil
		}),
	))

	Utils.ErrorChecker(os.WriteFile(filepath.Join("Report.pdf"), PDFBuffer, 0o644))
}

/*
GenerateXLSX → Generate final XLSX spreadsheet report from processed report data
  - Creates new Excel workbook with bluemonday HTML sanitiser for content cleaning
  - Processes findings collection by:
    -- Creating worksheet per directory (truncated to 31 characters for Excel limits)
    -- Adding headers (Finding Name, Status, Impact, Likelihood, Details) to new sheets
    -- Sanitizing and populating finding data into appropriate rows and columns
  - Removes default "Sheet1" and saves as "Report.xlsx"
  - Handles errors via Utils.ErrorChecker()
*/
func GenerateXLSX(_findings []Utils.Markdown, suggestions []Utils.Markdown) {
	outputSpreadsheet := excelize.NewFile()
	sanitiser := bluemonday.StrictPolicy()

	for _, finding := range _findings {
		sheetName := finding.Directory

		if len(sheetName) > 31 {
			sheetName = sheetName[:31]
		}

		if sheetNameExists, _ := outputSpreadsheet.GetSheetIndex(sheetName); sheetNameExists == -1 {
			outputSpreadsheet.NewSheet(sheetName)

			headers := []string{"Finding Name", "Finding Status", "Finding Impact", "Finding Likelihood", "Finding Details"}
			for i, header := range headers {
				outputSpreadsheet.SetCellValue(sheetName, fmt.Sprintf("%c1", 'A'+i), header)
			}
		}

		rows, errRows := outputSpreadsheet.GetRows(sheetName)
		Utils.ErrorChecker(errRows)
		row := len(rows) + 1

		values := []string{
			sanitiser.Sanitize(finding.Headers.FindingTitle),
			sanitiser.Sanitize(finding.Headers.FindingStatus),
			sanitiser.Sanitize(finding.Headers.FindingImpact),
			sanitiser.Sanitize(finding.Headers.FindingLikelihood),
			sanitiser.Sanitize(finding.Body),
		}

		for i, value := range values {
			outputSpreadsheet.SetCellValue(sheetName, fmt.Sprintf("%c%d", 'A'+i, row), value)
		}
	}

	outputSpreadsheet.DeleteSheet("Sheet1")
	Utils.ErrorChecker(outputSpreadsheet.SaveAs("Report.xlsx"))

}
