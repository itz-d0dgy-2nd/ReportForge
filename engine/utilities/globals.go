package utilities

import (
	"regexp"
	"sync"
)

const (
	RootDirectory                 string = "report"
	ConfigDirectory               string = "0_report_config"
	TemplateDirectory             string = "0_report_template"
	SummariesDirectory            string = "1_summaries"
	FindingsDirectory             string = "2_findings"
	SuggestionsDirectory          string = "3_suggestions"
	RisksDirectory                string = "4_risks"
	AppendicesDirectory           string = "5_appendices"
	ScreenshotsDirectory          string = "Screenshots"
	ScreenshotsOriginalsDirectory string = "originals"
)

const (
	ConfigFileMetadata           string = "metadata"
	ConfigFileSeverityAssessment string = "severity_assessment"
	ConfigFileContentOrder       string = "content_order"
)

const (
	ReportStatusDraft   string = "Draft"
	ReportStatusRelease string = "Release"
)

const (
	ImageCompressionQuality int = 75
	MaxIdentifierPrefixes   int = 26
	ExcelMaxSheetNameLength int = 31
	PDFOptimalImageWidth    int = 1200
)

const (
	FindingsStatusUnresolved string = "Unresolved"
	FindingsStatusResolved   string = "Resolved"
)

var (
	RegexYamlMatch               = regexp.MustCompile(`(?s)^---[ \t]*\r?\n(.*?)\r?\n---[ \t]*(?:\r?\n(.*))?$`)
	RegexFindingSeverity         = regexp.MustCompile(`FindingSeverity:[^\n]*`)
	RegexTokenMatch              = regexp.MustCompile(`\B!([A-Za-z][A-Za-z0-9_]*)\b`)
	RegexMarkdownRetestMatch     = regexp.MustCompile(`<p><(/?)(retest_)(fixed|not_fixed)></p>`)
	RegexMarkdownImageMatch      = regexp.MustCompile(`(<img\s+)src="(?:\.*\/)*(Screenshots/[^"]+)"([^>]*)\s*/>`)
	RegexMarkdownImageMatchScale = regexp.MustCompile(`(<img\s+)src="(?:\.*\/)*(Screenshots/[^"]+)"([^>]*)\s*/>\{([^}]*)\}`)
)

var DocumentStatus string

type Arguments struct {
	RebuildCache    bool
	DevelopmentMode bool
	CustomPath      string
}

type ReportPaths struct {
	RootPath        string
	ConfigPath      string
	TemplatePath    string
	SummariesPath   string
	FindingsPath    string
	SuggestionsPath string
	RisksPath       string
	AppendicesPath  string
	ScreenshotsPath string
}

type FileCache struct {
	cache            map[string][]byte
	mutex            sync.RWMutex
	Path             string
	MetadataConfig   MetadataYML
	SeverityConfig   SeverityAssessmentYML
	ContentConfig    ContentOrderYML
	SeverityMatrix   SeverityMatrix
	SeverityBarGraph SeverityBarGraph
	Summaries        []MarkdownFile
	Findings         []MarkdownFile
	Suggestions      []MarkdownFile
	Risks            []MarkdownFile
	Appendices       []MarkdownFile
}

type MetadataYML struct {
	Client              string            `yaml:"Client"`
	TargetInformation   map[string]string `yaml:"TargetInformation"`
	DocumentName        string            `yaml:"DocumentName"`
	DocumentInformation []struct {
		DocumentCurrent    bool              `yaml:"DocumentCurrent"`
		DocumentVersioning map[string]string `yaml:"DocumentVersioning"`
	} `yaml:"DocumentInformation"`
	StakeholderInformation []map[string]any  `yaml:"StakeholderInformation"`
	CustomVariables        map[string]string `yaml:"CustomVariables"`
}

type SeverityAssessmentYML struct {
	ConductSeverityAssessment bool         `yaml:"ConductSeverityAssessment"`
	DisplaySeverityMatrix     bool         `yaml:"DisplaySeverityMatrix"`
	SwapImpactLikelihoodAxis  bool         `yaml:"SwapImpactLikelihoodAxis"`
	DisplaySeverityBarGraph   bool         `yaml:"DisplaySeverityBarGraph"`
	Impacts                   []string     `yaml:"Impacts"`
	Likelihoods               []string     `yaml:"Likelihoods"`
	Severities                []string     `yaml:"Severities"`
	CalculatedMatrix          [5][5]string `yaml:"CalculatedMatrix"`
}

type ContentOrderYML struct {
	Summaries                    []string          `yaml:"Summaries"`
	Findings                     []string          `yaml:"Findings"`
	Suggestions                  []string          `yaml:"Suggestions"`
	Risks                        []string          `yaml:"Risks"`
	FindingIdentifierPrefixes    map[string]string `yaml:"FindingIdentifierPrefixes"`    // folder_name: "A"
	SuggestionIdentifierPrefixes map[string]string `yaml:"SuggestionIdentifierPrefixes"` // folder_name: "B"
	RiskIdentifierPrefixes       map[string]string `yaml:"RiskIdentifierPrefixes"`       // folder_name: "C"
}

type SeverityMatrixUpdate struct {
	RowIndex    int
	ColumnIndex int
	FindingID   string
}

type SeverityMatrix struct {
	Matrix [5][5]string
}

type SeverityBarGraphUpdate struct {
	Severity string
	Status   string
}

type SeverityBarGraph struct {
	Severities map[string]int
	Total      int
	Resolved   int
	Unresolved int
}

type MarkdownFile struct {
	Directory string
	FileName  string
	Headers   MarkdownYML
	Body      string
}

type MarkdownYML struct {
	ReportSummaryName        string `yaml:"ReportSummaryName"`
	ReportSummaryTitle       string `yaml:"ReportSummaryTitle"`
	ReportSummariesAuthor    string `yaml:"ReportSummariesAuthor"`
	ReportSummariesReviewers string `yaml:"ReportSummariesReviewers"`
	FindingID                string `yaml:"FindingID"`
	FindingIDLocked          bool   `yaml:"FindingIDLocked"`
	FindingName              string `yaml:"FindingName"`
	FindingTitle             string `yaml:"FindingTitle"`
	FindingStatus            string `yaml:"FindingStatus"`
	FindingImpact            string `yaml:"FindingImpact"`
	FindingLikelihood        string `yaml:"FindingLikelihood"`
	FindingSeverity          string `yaml:"FindingSeverity"`
	FindingAuthor            string `yaml:"FindingAuthor"`
	FindingReviewers         string `yaml:"FindingReviewers"`
	SuggestionID             string `yaml:"SuggestionID"`
	SuggestionIDLocked       bool   `yaml:"SuggestionIDLocked"`
	SuggestionName           string `yaml:"SuggestionName"`
	SuggestionTitle          string `yaml:"SuggestionTitle"`
	SuggestionStatus         string `yaml:"SuggestionStatus"`
	SuggestionAuthor         string `yaml:"SuggestionAuthor"`
	SuggestionReviewers      string `yaml:"SuggestionReviewers"`
	RiskID                   string `yaml:"RiskID"`
	RiskIDLocked             bool   `yaml:"RiskIDLocked"`
	RiskName                 string `yaml:"RiskName"`
	RiskTitle                string `yaml:"RiskTitle"`
	RiskStatus               string `yaml:"RiskStatus"`
	RiskGrossImpact          string `yaml:"RiskGrossImpact"`
	RiskGrossLikelihood      string `yaml:"RiskGrossLikelihood"`
	RiskGrossRating          string `yaml:"RiskGrossRating"`
	RiskTargetImpact         string `yaml:"RiskTargetImpact"`
	RiskTargetLikelihood     string `yaml:"RiskTargetLikelihood"`
	RiskTargetRating         string `yaml:"RiskTargetRating"`
	RiskAuthor               string `yaml:"RiskAuthor"`
	RiskReviewers            string `yaml:"RiskReviewers"`
	AppendixName             string `yaml:"AppendixName"`
	AppendixTitle            string `yaml:"AppendixTitle"`
	AppendixStatus           string `yaml:"AppendixStatus"`
	AppendixAuthor           string `yaml:"AppendixAuthor"`
	AppendixReviewers        string `yaml:"AppendixReviewers"`
}
