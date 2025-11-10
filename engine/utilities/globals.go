package utilities

import "regexp"

var RegexYamlMatch = regexp.MustCompile(`(?s)^---\n(.*?)\n---(?:\n(.*))?$`)
var RegexTokenMatch = regexp.MustCompile(`\B!([A-Za-z][A-Za-z0-9_]*)\b`)
var RegexMarkdownImageMatch = regexp.MustCompile(`(<img\s+)src="(?:\.*\/)*(Screenshots/[^"]+)"([^>]*)\s*/>`)
var RegexMarkdownImageMatchScale = regexp.MustCompile(`(<img\s+)src="(?:\.*\/)*(Screenshots/[^"]+)"([^>]*)\s*/>\{([^}]*)\}`)

type Arguments struct {
	DevelopmentMode bool
	CustomMode      string
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
}

type MarkdownFile struct {
	Directory string
	FileName  string
	Headers   MarkdownYML
	Body      string
}

type ReportData struct {
	Metadata    MetadataYML
	Severity    SeverityAssessmentYML
	Summaries   []MarkdownFile
	Findings    []MarkdownFile
	Suggestions []MarkdownFile
	Risks       []MarkdownFile
	Appendices  []MarkdownFile
	Path        string
}

type MetadataYML struct {
	Client              string            `yaml:"Client"`
	TargetInformation   map[string]string `yaml:"TargetInformation"`
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
	FlipSeverityMatrix        bool         `yaml:"FlipSeverityMatrix"`
	DisplaySeverityBarGraph   bool         `yaml:"DisplaySeverityBarGraph"`
	Impacts                   []string     `yaml:"Impacts"`
	Likelihoods               []string     `yaml:"Likelihoods"`
	Severities                []string     `yaml:"Severities"`
	CalculatedMatrix          [5][5]string `yaml:"CalculatedMatrix"`
	Matrix                    [5][5]string `yaml:"Matrix"`
}

type MarkdownYML struct {
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
	RiskDescription          string `yaml:"RiskDescription"`
	RiskDrivers              string `yaml:"RiskDrivers"`
	RiskConsequences         string `yaml:"RiskConsequences"`
	RiskGrossLikelihood      string `yaml:"RiskGrossLikelihood"`
	RiskGrossImpact          string `yaml:"RiskGrossImpact"`
	RiskGrossRating          string `yaml:"RiskGrossRating"`
	RiskRecommendedControls  string `yaml:"RiskRecommendedControls"`
	RiskOwner                string `yaml:"RiskOwner"`
	RiskTargetLikelihood     string `yaml:"RiskTargetLikelihood"`
	RiskTargetImpact         string `yaml:"RiskTargetImpact"`
	RiskTargetRating         string `yaml:"RiskTargetRating"`
	AppendixName             string `yaml:"AppendixName"`
	AppendixTitle            string `yaml:"AppendixTitle"`
	AppendixStatus           string `yaml:"AppendixStatus"`
	AppendixAuthor           string `yaml:"AppendixAuthor"`
	AppendixReviewers        string `yaml:"AppendixReviewers"`
}
