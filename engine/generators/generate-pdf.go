package generators

import (
	"ReportForge/engine/utilities"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func findChromiumBrowser() string {
	chromiumBrowsers := []string{"google-chrome", "chromium", "chromium-browser", "microsoft-edge"}
	if runtime.GOOS == "windows" {
		chromiumBrowsers = []string{"msedge.exe", "chrome.exe", "chromium.exe"}
	}

	for _, chromiumBrowser := range chromiumBrowsers {
		if path, errLookPath := exec.LookPath(chromiumBrowser); errLookPath == nil {
			return path
		}
	}

	utilities.Check(utilities.NewExternalError(
		"no Chromium-based browser found in $PATH — install Chrome, Chromium, or Edge",
		nil,
	))
	return ""
}

func htmlFileURL(_htmlFileName string) string {
	absoluteFilePath, errAbsoluteFilePath := filepath.Abs(_htmlFileName)
	if errAbsoluteFilePath != nil {
		utilities.Check(utilities.NewFileSystemError(_htmlFileName, "failed to resolve absolute path", errAbsoluteFilePath))
	}

	fileURL := &url.URL{Scheme: "file", Path: filepath.ToSlash(absoluteFilePath)}
	return fileURL.String()
}

func buildChromiumOptions(_browserPath string) []chromedp.ExecAllocatorOption {
	chromiumExecutionOptions := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath(_browserPath),
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

	if os.Getenv("ACTION") == "true" {
		chromiumExecutionOptions = append(chromiumExecutionOptions, chromedp.Flag("no-sandbox", true))
	}

	return chromiumExecutionOptions
}

func renderPDF(_fileURL string, _browserPath string) []byte {
	chromiumExecutionOptions := buildChromiumOptions(_browserPath)
	chromiumExecutionContext, chromiumExecutionContextCancel := chromedp.NewExecAllocator(context.Background(), chromiumExecutionOptions...)
	defer chromiumExecutionContextCancel()

	chromiumBrowserContext, chromiumBrowserContextCancel := chromedp.NewContext(chromiumExecutionContext)
	defer chromiumBrowserContextCancel()

	var buffer []byte
	errChromedp := chromedp.Run(chromiumBrowserContext,
		chromedp.Navigate(_fileURL),
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
			if errPrintToPDF != nil {
				return errPrintToPDF
			}
			buffer = print
			return nil
		}),
	)
	if errChromedp != nil {
		utilities.Check(utilities.NewExternalError(
			fmt.Sprintf("Chrome DevTools Protocol error rendering '%s'", _fileURL),
			errChromedp,
		))
	}
	if len(buffer) == 0 {
		utilities.Check(utilities.NewExternalError("PDF rendering produced empty output", nil))
	}

	return buffer
}

/*
GeneratePDF → Generate final PDF report from processed report data
*/
func GeneratePDF(_fileCache *utilities.FileCache, _reportPaths utilities.ReportPaths) {
	metadataConfig := _fileCache.MetadataConfig()
	chromiumPath := findChromiumBrowser()
	fileURL := htmlFileURL(fmt.Sprintf("%s_%s.html", metadataConfig.DocumentName, utilities.DocumentVersion))
	pdfBytes := renderPDF(fileURL, chromiumPath)

	pdfFileName := fmt.Sprintf("%s_%s.pdf", metadataConfig.DocumentName, utilities.DocumentVersion)
	if errWriteFile := os.WriteFile(pdfFileName, pdfBytes, 0o644); errWriteFile != nil {
		utilities.Check(utilities.NewFileSystemError(pdfFileName, "failed to write PDF output file", errWriteFile))
	}
}
