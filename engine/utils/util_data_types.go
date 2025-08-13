package Utils

type ArgumentsStruct struct {
	DevelopmentMode bool
	CustomMode      string
}

type ReportPathsStruct struct {
	RootPath        string
	ConfigPath      string
	TemplatePath    string
	SummariesPath   string
	FindingsPath    string
	SuggestionsPath string
	AppendicesPath  string
}

type Markdown struct {
	Directory string
	FileName  string
	Headers   MarkdownYML
	Body      string
}

type ReportDataStruct struct {
	Frontmatter FrontmatterYML
	Severity    SeverityAssessmentYML
	Summaries   []Markdown
	Findings    []Markdown
	Suggestions []Markdown
	Appendices  []Markdown
	Path        string
}

type FrontmatterYML struct {
	Client              string            `yaml:"Client"`
	TargetInformation   map[string]string `yaml:"TargetInformation"`
	DocumentInformation []struct {
		DocumentCurrent  bool              `yaml:"DocumentCurrent"`
		DocumentMetadata map[string]string `yaml:"DocumentMetadata"`
	} `yaml:"DocumentInformation"`
	StakeholderInformation []map[string]any `yaml:"StakeholderInformation"`
}

type SeverityAssessmentYML struct {
	Impacts          map[int]string `yaml:"Impacts"`
	Likelihoods      map[int]string `yaml:"Likelihoods"`
	Matrix           [5][5]string   `yaml:"Matrix"`
	CalculatedMatrix [5][5]string   `yaml:"CalculatedMatrix"`
}

type MarkdownYML struct {
	ReportSummariesAuthor    string `yaml:"ReportSummariesAuthor"`
	ReportSummariesReviewers string `yaml:"ReportSummariesReviewers"`
	FindingID                string `yaml:"FindingID"`
	FindingName              string `yaml:"FindingName"`
	FindingTitle             string `yaml:"FindingTitle"`
	FindingStatus            string `yaml:"FindingStatus"`
	FindingImpact            string `yaml:"FindingImpact"`
	FindingLikelihood        string `yaml:"FindingLikelihood"`
	FindingSeverity          string `yaml:"FindingSeverity"`
	FindingAuthor            string `yaml:"FindingAuthor"`
	FindingReviewers         string `yaml:"FindingReviewers"`
	SuggestionID             string `yaml:"SuggestionID"`
	SuggestionName           string `yaml:"SuggestionName"`
	SuggestionTitle          string `yaml:"SuggestionTitle"`
	SuggestionStatus         string `yaml:"SuggestionStatus"`
	SuggestionAuthor         string `yaml:"SuggestionAuthor"`
	SuggestionReviewers      string `yaml:"SuggestionReviewers"`
	AppendixName             string `yaml:"AppendixName"`
	AppendixTitle            string `yaml:"AppendixTitle"`
	AppendixStatus           string `yaml:"AppendixStatus"`
	AppendixAuthor           string `yaml:"AppendixAuthor"`
	AppendixReviewers        string `yaml:"AppendixReviewers"`
}
