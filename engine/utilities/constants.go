package utilities

import "time"

// ======== Error Constants ========
const (
	Validation ErrorCategory = iota
	FileSystem
	YAML
	Configuration
	Processing
	External
)

const (
	Warning ErrorSeverity = iota
	Error
)

// ======== Directory Structure Constants ========
const (
	RootDirectory                 string = "report"
	ConfigDirectory               string = "0_report_config"
	TemplateDirectory             string = "0_report_template"
	SummariesDirectory            string = "1_summaries"
	AppendicesDirectory           string = "Appendices"
	ScreenshotsDirectory          string = "Screenshots"
	ScreenshotsOriginalsDirectory string = "originals"

	FindingsDirectory    string = "2_findings"
	SuggestionsDirectory string = "3_suggestions"

	RisksDirectory    string = "2_risks"
	ControlsDirectory string = "3_controls"
)

// ======== Configuration File Constants ========
const (
	ConfigFileMetadata           string = "metadata"
	ConfigFileSeverityAssessment string = "severity_assessment"
	ConfigFileRiskAssessment     string = "risk_assessment"
	ConfigFileContentOrder       string = "content_order"
)

// ======== Report Status Constants ========
const (
	ReportStatusDraft   string = "Draft"
	ReportStatusRelease string = "Release"
)

// ======== Finding Status Constants ========
const (
	FindingsStatusUnresolved string = "Unresolved"
	FindingsStatusResolved   string = "Resolved"
)

// ======== Processing Constants ========
const (
	ImageCompressionQuality int           = 75
	MaxIdentifierPrefixes   int           = 26
	ExcelMaxSheetNameLength int           = 31
	PDFOptimalImageWidth    int           = 1200
	DebounceInterval        time.Duration = 500 * time.Millisecond
)

// ======== Global States ========
var (
	DocumentStatus  string
	DocumentVersion string
)

var (
	YAMLPattern     YAMLPatterns
	MarkdownPattern MarkdownPatterns
	RiskPattern     RiskPatterns
	ControlPattern  ControlPatterns
)
