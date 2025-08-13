package generators

import (
	"ReportForge/engine/utils"
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func SetupChromiumPrint(_PDFBuffer *[]byte, _reportPaths Utils.ReportPathsStruct) chromedp.Tasks {
	absoluteFilePath, errAbsoluteFilePath := filepath.Abs(filepath.Join("Report.html"))
	Utils.ErrorChecker(errAbsoluteFilePath)
	return chromedp.Tasks{

		chromedp.Navigate("file:///" + absoluteFilePath),
		chromedp.ActionFunc(func(context context.Context) error {

			// NOTE:
			//   Issue with A4 page dimensions (`WithPaperWidth(8.27).WithPaperHeight(11.69)`).
			//   A faint white line appears on the bottom and right side.
			//   Requires further investigation but fuck freedom units :P
			// UPDATE:
			//   WithPaperWidth() & WithPaperHeight() are dumb
			//   WithPreferCSSPageSize() works better

			print, _, errPrint := page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPreferCSSPageSize(true).
				Do(context)
			Utils.ErrorChecker(errPrint)

			*_PDFBuffer = print
			return nil

		}),
	}
}

func SetupChromiumBrowser() []chromedp.ExecAllocatorOption {
	chromiumExecutionOptions := []chromedp.ExecAllocatorOption{}
	chromiumExecutionOptions = append(chromiumExecutionOptions, chromedp.Flag("headless", true))

	if runtime.GOOS == "windows" {
		chromiumExecutionOptions = append(chromiumExecutionOptions, chromedp.ExecPath(`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`))
	}
	if os.Getenv("ACTION") == "true" {
		chromiumExecutionOptions = append(chromiumExecutionOptions, chromedp.Flag("no-sandbox", true))
		// This needs to be set due to either the restricted namespaces or apparmour in GitHub actions.
	}
	return chromiumExecutionOptions
}

func GeneratePDF(_reportPaths Utils.ReportPathsStruct) {

	PDFBuffer := []byte{}
	chromiumExecutionOptions := SetupChromiumBrowser()

	chromiumExecutionContext, errChromiumExecutionContext := chromedp.NewExecAllocator(context.Background(), chromiumExecutionOptions...)
	defer errChromiumExecutionContext()

	chromiumBrowserContext, errChromiumBrowserContext := chromedp.NewContext(chromiumExecutionContext)
	defer errChromiumBrowserContext()

	Utils.ErrorChecker(chromedp.Run(chromiumBrowserContext, SetupChromiumPrint(&PDFBuffer, _reportPaths)))
	Utils.ErrorChecker(os.WriteFile(filepath.Join("Report.pdf"), PDFBuffer, 0o777))
}
