package Utils

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ChromePDFPrint(_bytes *[]byte) chromedp.Tasks {

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

			print, _, errPrint := page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPreferCSSPageSize(true).
				Do(context)
			ErrorChecker(errPrint)

			*_bytes = print
			return nil

		}),
	}
}

func GeneratePDF() {

	bufferPDF := []byte{}
	execAllocatorOpts := []chromedp.ExecAllocatorOption{}

	execAllocatorOpts = append(execAllocatorOpts, chromedp.Flag("headless", true))

	if runtime.GOOS == "windows" {
		execAllocatorOpts = append(execAllocatorOpts,
			chromedp.ExecPath(`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`),
		)
	}

	if os.Getenv("ACTION") == "true" {
		execAllocatorOpts = append(execAllocatorOpts,
			chromedp.Flag("no-sandbox", true),
		)
	}

	executionContext, ErrExecutionContext := chromedp.NewExecAllocator(context.Background(), execAllocatorOpts...)
	defer ErrExecutionContext()

	browserContext, ErrBrowserContext := chromedp.NewContext(executionContext)
	defer ErrBrowserContext()

	ErrorChecker(chromedp.Run(browserContext, ChromePDFPrint(&bufferPDF)))
	ErrorChecker(os.WriteFile("Report.pdf", bufferPDF, 0644))
	// ErrorChecker(os.Remove("Report.html"))

}
