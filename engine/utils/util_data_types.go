package Utils

type FrontmatterJSON struct {
	Client                 string            `json:"Client"`
	TargetInformation      map[string]string `json:"TargetInformation"`
	DocumentInformation    []map[string]any  `json:"DocumentInformation"`
	StakeholderInformation []map[string]any  `json:"StakeholderInformation"`
}

type SeverityMatrix struct {
	Impacts     map[string]int
	Likelihoods map[string]int
	Matrix      [5][5]string
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

type Markdown struct {
	Directory string
	FileName  string
	Headers   MarkdownYML
	Body      string
}

type Report struct {
	Frontmatter     FrontmatterJSON
	Severity        [5][5]string
	ReportSummaries []Markdown
	Findings        []Markdown
	Suggestions     []Markdown
	Appendices      []Markdown
}
