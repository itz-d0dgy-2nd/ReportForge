package utilities

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

/*
Error → Implements error interface for ReportError
*/
func (_reportError *ReportError) Error() string {
	if _reportError.Path != "" && _reportError.Field != "" {
		return fmt.Sprintf("%s in field '%s' at ( %s ): %s", _reportError.categoryName(), _reportError.Field, _reportError.Path, _reportError.Message)
	}
	if _reportError.Path != "" {
		return fmt.Sprintf("%s at ( %s ): %s", _reportError.categoryName(), _reportError.Path, _reportError.Message)
	}
	return fmt.Sprintf("%s: %s", _reportError.categoryName(), _reportError.Message)
}

/*
Unwrap → Supports error unwrapping for errors.Is and errors.As
*/
func (_reportError *ReportError) Unwrap() error {
	return _reportError.Wrapped
}

/*
IsWarning → Returns true if error should be treated as a warning
*/
func (_reportError *ReportError) IsWarning() bool {
	return _reportError.Severity == Warning
}

/*
categoryName → Returns category name
*/
func (_reportError *ReportError) categoryName() string {
	switch _reportError.Category {
	case Validation:
		return "Validation failed"
	case FileSystem:
		return "File system error"
	case YAML:
		return "YAML parsing failed"
	case Configuration:
		return "Configuration error"
	case Processing:
		return "Processing failed"
	case External:
		return "External dependency error"
	default:
		return "Error"
	}
}

// ======== Custom Error's ========

/*
NewValidationError → Creates a validation error (fatal by default)
*/
func NewValidationError(_path, _field, _message string) *ReportError {
	return &ReportError{
		Category: Validation,
		Severity: Error,
		Path:     _path,
		Field:    _field,
		Message:  _message,
	}
}

/*
NewValidationWarning → Creates a validation warning (non-fatal)
*/
func NewValidationWarning(_path, _message string) *ReportError {
	return &ReportError{
		Category: Validation,
		Severity: Warning,
		Path:     _path,
		Message:  _message,
	}
}

/*
NewFileSystemError → Creates a file system error
*/
func NewFileSystemError(_path, _message string, _wrapped error) *ReportError {
	return &ReportError{
		Category: FileSystem,
		Severity: Error,
		Path:     _path,
		Message:  _message,
		Wrapped:  _wrapped,
	}
}

/*
NewYAMLError → Creates a YAML parsing error
*/
func NewYAMLError(_path, _message string, _wrapped error) *ReportError {
	return &ReportError{
		Category: YAML,
		Severity: Error,
		Path:     _path,
		Message:  _message,
		Wrapped:  _wrapped,
	}
}

/*
NewConfigError → Creates a configuration error
*/
func NewConfigError(_path, _message string) *ReportError {
	return &ReportError{
		Category: Configuration,
		Severity: Error,
		Path:     _path,
		Message:  _message,
	}
}

/*
NewProcessingError → Creates a processing error
*/
func NewProcessingError(_path, _message string) *ReportError {
	return &ReportError{
		Category: Processing,
		Severity: Error,
		Path:     _path,
		Message:  _message,
	}
}

/*
NewExternalError → Creates an external dependency error
*/
func NewExternalError(_message string, _wrapped error) *ReportError {
	return &ReportError{
		Category: External,
		Severity: Error,
		Message:  _message,
		Wrapped:  _wrapped,
	}
}

/*
Check → Handles error reporting and exits on fatal errors.
  - Warnings: are printed to stderr and execution continues.
  - Fatal: errors are printed to stderr and exit with code 1.
*/
func Check(_uhm error) {
	if _uhm == nil {
		return
	}

	// Handle custom ReportError types
	var errReport *ReportError
	if errors.As(_uhm, &errReport) {
		if errReport.IsWarning() {
			fmt.Fprintf(os.Stderr, "::warning:: %s\n", errReport.Error())
			return
		}
		fmt.Fprintf(os.Stderr, "::error:: %s\n", errReport.Error())
		os.Exit(1)
	}

	// Handle standard errors
	switch {
	case errors.Is(_uhm, fs.ErrNotExist):
		fmt.Fprintf(os.Stderr, "::error:: File not found: %s\n", _uhm.Error())
		os.Exit(1)

	case errors.Is(_uhm, fs.ErrPermission):
		fmt.Fprintf(os.Stderr, "::error:: Permission denied: %s\n", _uhm.Error())
		os.Exit(1)

	default:
		fmt.Fprintf(os.Stderr, "::error:: %s\n", _uhm.Error())
		os.Exit(1)
	}
}
