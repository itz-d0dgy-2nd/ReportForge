package utilities

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

/*
ErrorChecker → Checks error and handles it based on severity level
*/
func ErrorChecker(_uhm error) {
	var isError bool
	var message string

	if _uhm == nil {
		return
	}

	errorMessage := _uhm.Error()

	switch {
	case errors.Is(_uhm, fs.ErrNotExist):
		isError = false
		message = fmt.Sprintf("File not found - %s", errorMessage)

	case strings.Contains(errorMessage, "QA") && strings.Contains(errorMessage, "comment"):
		isError = false
		message = fmt.Sprintf("QA comment found - %s", errorMessage)

	case errors.Is(_uhm, fs.ErrPermission):
		isError = true
		message = fmt.Sprintf("Permission denied - %s", errorMessage)

	case strings.Contains(errorMessage, "executable file not found"):
		isError = true
		message = "Chromium-based browser not found in $PATH (required for PDF generation)"

	case strings.Contains(errorMessage, "invalid") && strings.Contains(errorMessage, "path"):
		isError = true
		message = fmt.Sprintf("Invalid file path - %s", errorMessage)

	case strings.Contains(errorMessage, "yaml:") || strings.Contains(errorMessage, "unmarshal"):
		isError = true
		message = fmt.Sprintf("YAML parsing failed - %s", errorMessage)

	case strings.Contains(errorMessage, "missing required key"):
		isError = true
		message = fmt.Sprintf("Frontmatter validation failed - %s", errorMessage)

	default:
		isError = true
		message = errorMessage
	}

	if isError {
		fmt.Fprintf(os.Stderr, "::error:: %s\n", message)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stderr, "::warning:: %s\n", message)
	}
}
