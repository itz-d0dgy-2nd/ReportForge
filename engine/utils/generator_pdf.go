package Utils

import (
	"context"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ChromePDFPrint(bytes *[]byte) chromedp.Tasks {

	absPath, ErrAbsPath := filepath.Abs("Report.html")
	ErrorChecker(ErrAbsPath)

	return chromedp.Tasks{

		chromedp.Navigate("file:///" + absPath),
		chromedp.ActionFunc(func(context context.Context) error {

			// NOTE:
			//   Issue with A4 page dimensions (`WithPaperWidth(8.27).WithPaperHeight(11.69)`).
			//   A faint white line appears on the bottom and right side.
			//   Requires further investigation but fuck freedom units :P
			// UPDATE:
			//   WithPaperWidth() & WithPaperHeight() are dumb
			//   WithPreferCSSPageSize() works better

			print, _, ErrPrint := page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPreferCSSPageSize(true).
				Do(context)
			ErrorChecker(ErrPrint)

			*bytes = print
			return nil

		}),
	}
}

func GeneratePDF() {

	bufferPDF := []byte{}

	executionContext, ErrExecutionContext := chromedp.NewExecAllocator(context.Background(), chromedp.Flag("headless", true), chromedp.Flag("no-sandbox", true))
	defer ErrExecutionContext()

	browserContext, ErrBrowserContext := chromedp.NewContext(executionContext)
	defer ErrBrowserContext()

	ErrorChecker(chromedp.Run(browserContext, ChromePDFPrint(&bufferPDF)))
	ErrorChecker(os.WriteFile("Report.pdf", bufferPDF, 0644))
	// ErrorChecker(os.Remove("Report.html"))

}
