package generators

import (
	"ReportForge/engine/utilities"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/xuri/excelize/v2"
)

/*
GenerateHTML → Generate final HTML report from processed report data
  - Parses "template.html" file and creates ".html"
  - Handles errors via utilities.ErrorChecker()
*/
func GenerateHTML(_fileCache *utilities.FileCache, _reportPaths utilities.ReportPaths) {
	templateFunctionMap := template.FuncMap{
		"inc":   func(i int) int { return i + 1 },
		"dec":   func(i int) int { return i - 1 },
		"add":   func(a, b int) int { return a + b },
		"sub":   func(a, b int) int { return a - b },
		"split": strings.Split,
	}

	templateHTML, errTemplateHTML := template.New("template.html").Funcs(templateFunctionMap).ParseFiles(_reportPaths.TemplatePath)
	utilities.ErrorChecker(errTemplateHTML)

	createHTML, errCreateHTML := os.Create(_fileCache.MetadataConfig.DocumentName + ".html")
	utilities.ErrorChecker(errCreateHTML)

	defer createHTML.Close()

	utilities.ErrorChecker(templateHTML.Execute(createHTML, _fileCache))
}

/*
GeneratePDF → Generate final PDF report from processed report data
  - Configures chromedp with headless browser options (disable GPU, extensions, audio, etc.)
  - Performs system-specific configuration on:
    -- Windows systems - Detects Microsoft Edge and uses as Chromium executable
    -- CICD / Action systems - Adds no-sandbox flag when running in CICD/Action
  - Parses ".html" file and creates ".pdf"
  - Handles errors via utilities.ErrorChecker()
  - TODO: Check if there is a way to get windows path from registry for known chromium binaries
*/
func GeneratePDF(_fileCache *utilities.FileCache, _reportPaths utilities.ReportPaths) {
	var buffer []byte

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
			if _, errStatCheck := os.Stat(path); errStatCheck == nil {
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

	absoluteFilePath, errAbsoluteFilePath := filepath.Abs(filepath.Join(_fileCache.MetadataConfig.DocumentName + ".html"))
	utilities.ErrorChecker(errAbsoluteFilePath)
	fileURL := "file:///" + filepath.ToSlash(absoluteFilePath)

	utilities.ErrorChecker(chromedp.Run(chromiumBrowserContext,
		chromedp.Navigate(fileURL),
		chromedp.ActionFunc(func(context context.Context) error {
			print, _, errPrintToPDF := page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPreferCSSPageSize(true).
				WithGenerateTaggedPDF(true).
				Do(context)
			utilities.ErrorChecker(errPrintToPDF)

			buffer = print
			return nil
		}),
	))

	utilities.ErrorChecker(os.WriteFile(filepath.Join(_fileCache.MetadataConfig.DocumentName+".pdf"), buffer, 0o644))
}

/*
GenerateXLSX → Generate final XLSX spreadsheet report from processed report data
  - Creates new Excel workbook with bluemonday HTML sanitiser for content cleaning
  - Processes findings collection by:
    -- Creating worksheet per directory (truncated to 31 characters for Excel limits)
    -- Adding headers (Finding Name, Status, Impact, Likelihood, Details) to new sheets
    -- Sanitising and populating finding data into appropriate rows and columns
  - Removes default "Sheet1" and saves as ".xlsx"
  - Handles errors via utilities.ErrorChecker()
*/
func GenerateXLSX(_fileCache *utilities.FileCache) {
	outputSpreadsheet := excelize.NewFile()
	sanitiser := bluemonday.StrictPolicy()

	for _, finding := range _fileCache.Findings {
		sheetName := finding.Directory

		if len(sheetName) > utilities.ExcelMaxSheetNameLength {
			sheetName = sheetName[:utilities.ExcelMaxSheetNameLength]
		}

		if sheetNameExists, _ := outputSpreadsheet.GetSheetIndex(sheetName); sheetNameExists == -1 {
			outputSpreadsheet.NewSheet(sheetName)

			headers := []string{"Finding Name", "Finding Status", "Finding Impact", "Finding Likelihood", "Finding Details"}
			for i, header := range headers {
				outputSpreadsheet.SetCellValue(sheetName, fmt.Sprintf("%c1", 'A'+i), header)
			}
		}

		rows, errRows := outputSpreadsheet.GetRows(sheetName)
		utilities.ErrorChecker(errRows)
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
	utilities.ErrorChecker(outputSpreadsheet.SaveAs(_fileCache.MetadataConfig.DocumentName + ".xlsx"))
}
