package utilities

import (
	"regexp"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
)

// ======== Application Configuration Types ========

type Arguments struct {
	CustomPath string
	Debug      bool
	Watch      bool
}

type TemplateType struct {
	Technical bool
	Sra       bool
}

type ReportPaths struct {
	RootPath        string
	ConfigPath      string
	TemplatePath    string
	SummariesPath   string
	FindingsPath    string
	SuggestionsPath string
	RisksPath       string
	ControlsPath    string
	AppendicesPath  string
	ScreenshotsPath string
}

// ======== Regex Pattern Types ========

type YAMLPatterns struct {
	Frontmatter  *regexp.Regexp
	FindingID    *regexp.Regexp
	SuggestionID *regexp.Regexp
	RiskID       *regexp.Regexp
	ControlID    *regexp.Regexp
	Severity     *regexp.Regexp
	GrossRating  *regexp.Regexp
	TargetRating *regexp.Regexp
	CIDRef       *regexp.Regexp
}

type MarkdownPatterns struct {
	Token      *regexp.Regexp
	Retest     *regexp.Regexp
	Image      *regexp.Regexp
	ImageScale *regexp.Regexp
}

type RiskPatterns struct {
	Description         *regexp.Regexp
	Drivers             *regexp.Regexp
	Consequences        *regexp.Regexp
	RecommendedControls *regexp.Regexp
}

type ControlPatterns struct {
	Description *regexp.Regexp
	Reduces     *regexp.Regexp
}

// ======== Error Types ========

type ErrorSeverity int

type ErrorCategory int

type ReportError struct {
	Category ErrorCategory
	Severity ErrorSeverity
	Path     string
	Field    string
	Message  string
	Wrapped  error
}

// ======== YAML Configuration Types ========

type MetadataYML struct {
	Client              string            `yaml:"Client"`
	TargetInformation   map[string]string `yaml:"TargetInformation"`
	DocumentName        string            `yaml:"DocumentName"`
	DocumentInformation []struct {
		DocumentCurrent    bool              `yaml:"DocumentCurrent"`
		DocumentVersioning map[string]string `yaml:"DocumentVersioning"`
	} `yaml:"DocumentInformation"`
	StakeholderInformation []struct {
		StakeholderName    string `yaml:"StakeholderName"`
		StakeholderRole    string `yaml:"StakeholderRole"`
		StakeholderCompany string `yaml:"StakeholderCompany"`
	} `yaml:"StakeholderInformation"`
	CustomVariables map[string]string `yaml:"CustomVariables"`
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

type RiskAssessmentYML struct {
	GrossImpacts           []string     `yaml:"GrossImpacts"`
	GrossLikelihoods       []string     `yaml:"GrossLikelihoods"`
	GrossRiskRatings       []string     `yaml:"GrossRiskRatings"`
	CalculatedGrossMatrix  [5][5]string `yaml:"CalculatedGrossMatrix"`
	TargetImpacts          []string     `yaml:"TargetImpacts"`
	TargetLikelihoods      []string     `yaml:"TargetLikelihoods"`
	TargetRiskRatings      []string     `yaml:"TargetRiskRatings"`
	CalculatedTargetMatrix [5][5]string `yaml:"CalculatedTargetMatrix"`
}

type ContentOrderYML struct {
	Summaries                    []string          `yaml:"Summaries"`
	Findings                     []string          `yaml:"Findings"`
	Suggestions                  []string          `yaml:"Suggestions"`
	Risks                        []string          `yaml:"Risks"`
	Controls                     []string          `yaml:"Controls"`
	FindingIdentifierPrefixes    map[string]string `yaml:"FindingIdentifierPrefixes"`
	SuggestionIdentifierPrefixes map[string]string `yaml:"SuggestionIdentifierPrefixes"`
	RiskIdentifierPrefixes       map[string]string `yaml:"RiskIdentifierPrefixes"`
	ControlIdentifierPrefixes    map[string]string `yaml:"ControlIdentifierPrefixes"`
}

// ======== Markdown Types ========

type MarkdownYML struct {
	ReportSummaryName        string `yaml:"ReportSummaryName"`
	ReportSummaryTitle       string `yaml:"ReportSummaryTitle"`
	ReportSummariesAuthor    string `yaml:"ReportSummariesAuthor"`
	ReportSummariesReviewers string `yaml:"ReportSummariesReviewers"`

	FindingID         string `yaml:"FindingID"`
	FindingIDLocked   bool   `yaml:"FindingIDLocked"`
	FindingName       string `yaml:"FindingName"`
	FindingTitle      string `yaml:"FindingTitle"`
	FindingStatus     string `yaml:"FindingStatus"`
	FindingImpact     string `yaml:"FindingImpact"`
	FindingLikelihood string `yaml:"FindingLikelihood"`
	FindingSeverity   string `yaml:"FindingSeverity"`
	FindingAuthor     string `yaml:"FindingAuthor"`
	FindingReviewers  string `yaml:"FindingReviewers"`

	SuggestionID        string `yaml:"SuggestionID"`
	SuggestionIDLocked  bool   `yaml:"SuggestionIDLocked"`
	SuggestionName      string `yaml:"SuggestionName"`
	SuggestionTitle     string `yaml:"SuggestionTitle"`
	SuggestionStatus    string `yaml:"SuggestionStatus"`
	SuggestionAuthor    string `yaml:"SuggestionAuthor"`
	SuggestionReviewers string `yaml:"SuggestionReviewers"`

	RiskID               string `yaml:"RiskID"`
	RiskIDLocked         bool   `yaml:"RiskIDLocked"`
	RiskName             string `yaml:"RiskName"`
	RiskTitle            string `yaml:"RiskTitle"`
	RiskGrossImpact      string `yaml:"RiskGrossImpact"`
	RiskGrossLikelihood  string `yaml:"RiskGrossLikelihood"`
	RiskGrossRating      string `yaml:"RiskGrossRating"`
	RiskTargetImpact     string `yaml:"RiskTargetImpact"`
	RiskTargetLikelihood string `yaml:"RiskTargetLikelihood"`
	RiskTargetRating     string `yaml:"RiskTargetRating"`
	RiskAuthor           string `yaml:"RiskAuthor"`
	RiskReviewers        string `yaml:"RiskReviewers"`

	ControlID              string   `yaml:"ControlID"`
	ControlIDLocked        bool     `yaml:"ControlIDLocked"`
	ControlName            string   `yaml:"ControlName"`
	ControlTitle           string   `yaml:"ControlTitle"`
	ControlNZISMReferences []string `yaml:"ControlNZISMReferences"`
	ControlAuthor          string   `yaml:"ControlAuthor"`
	ControlReviewers       string   `yaml:"ControlReviewers"`

	AppendixName      string `yaml:"AppendixName"`
	AppendixTitle     string `yaml:"AppendixTitle"`
	AppendixStatus    string `yaml:"AppendixStatus"`
	AppendixAuthor    string `yaml:"AppendixAuthor"`
	AppendixReviewers string `yaml:"AppendixReviewers"`
}

type MarkdownFile struct {
	Directory string
	FileName  string
	Headers   MarkdownYML
	Body      string
}

// ======== Cache Types ========

type FileCache struct {
	markdown            map[string][]byte
	markdownFrontmatter map[string]MarkdownYML
	fileNames           map[string]string
	mutex               sync.RWMutex
	Path                string
	metadataConfig      MetadataYML
	severityConfig      SeverityAssessmentYML
	riskConfig          RiskAssessmentYML
	contentConfig       ContentOrderYML
	SeverityMatrix      SeverityMatrix
	SeverityBarGraph    SeverityBarGraph
	RiskMatrices        RiskMatrices
	Summaries           []MarkdownFile
	Findings            []MarkdownFile
	Suggestions         []MarkdownFile
	Risks               []MarkdownFile
	Controls            []MarkdownFile
	Appendices          []MarkdownFile
}

// === Worker Types

type ModificationJob struct {
	Path             string
	DirectoryType    string
	AssignedID       int32
	IdentifierPrefix string
	IsLocked         bool
}

type ProcessingJob struct {
	Path      string
	Directory string
}

type ProcessingResult struct {
	Markdown         MarkdownFile
	SeverityMatrix   *SeverityMatrixUpdate
	SeverityBarGraph *SeverityBarGraphUpdate
	RiskMatrices     *RiskMatricesUpdate
	Directory        string
}

// ======== Matrix Types ========

type SeverityMatrix struct {
	Matrix [5][5]string
}

type SeverityMatrixUpdate struct {
	RowIndex    int
	ColumnIndex int
	FindingID   string
}

type SeverityBarGraph struct {
	Severities map[string]int
	Total      int
	Resolved   int
	Unresolved int
}

type SeverityBarGraphUpdate struct {
	Severity string
	Status   string
}

type RiskMatrices struct {
	GrossMatrix  [5][5]string
	TargetMatrix [5][5]string
}

type RiskMatricesUpdate struct {
	GrossRowIndex     int
	GrossColumnIndex  int
	TargetRowIndex    int
	TargetColumnIndex int
	RiskID            string
}

// ======== Template Types ========

type DirectoryGroup struct {
	Directory string
	Items     []MarkdownFile
}

type ControlRiskMapping struct {
	ControlID    string
	ControlTitle string
	ControlName  string
	RiskIDs      []string
	RiskName     []string
}

type TemplateData struct {
	*FileCache
	AppendixGroups      []DirectoryGroup
	RiskGroups          []DirectoryGroup
	ControlGroups       []DirectoryGroup
	ControlRiskMappings []ControlRiskMapping
	FindingGroups       []DirectoryGroup
	SuggestionGroups    []DirectoryGroup
}

type Watcher struct {
	paused          atomic.Bool
	skipNextTrigger atomic.Bool
	watcher         *fsnotify.Watcher
}
