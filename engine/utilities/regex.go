package utilities

import (
	"regexp"
)

// Global doesnt feel right... Must refactor down the line
func init() {
	YAMLPattern = YAMLPatterns{
		Frontmatter:  regexp.MustCompile(`(?s)^---[ \t]*\r?\n(.*?)\r?\n---[ \t]*(?:\r?\n(.*))?$`),
		FindingID:    regexp.MustCompile(`(?m)^(\s*FindingID\s*:\s*).*$`),
		SuggestionID: regexp.MustCompile(`(?m)^(\s*SuggestionID\s*:\s*).*$`),
		RiskID:       regexp.MustCompile(`(?m)^(\s*RiskID\s*:\s*).*$`),
		ControlID:    regexp.MustCompile(`(?m)^(\s*ControlID\s*:\s*).*$`),
		Severity:     regexp.MustCompile(`FindingSeverity:[^\n]*`),
		GrossRating:  regexp.MustCompile(`RiskGrossRating:[^\n]*`),
		TargetRating: regexp.MustCompile(`RiskTargetRating:[^\n]*`),
		CIDRef:       regexp.MustCompile(`^CID:\d+$`),
	}

	MarkdownPattern = MarkdownPatterns{
		Token:      regexp.MustCompile(`\B!([A-Za-z][A-Za-z0-9_]*)\b`),
		Retest:     regexp.MustCompile(`<p><(/?)(retest_)(fixed|not_fixed)></p>`),
		Image:      regexp.MustCompile(`(<img\s+)src="(?:\.*\/)*(Screenshots/[^"]+)"([^>]*)\s*/>`),
		ImageScale: regexp.MustCompile(`(<img\s+)src="(?:\.*\/)*(Screenshots/[^"]+)"([^>]*)\s*/>\{([^}]*)\}`),
	}

	RiskPattern = RiskPatterns{
		Description:         regexp.MustCompile(`(?s)<h3[^>]*>Risk Description</h3>\s*(.+?)(?:<h3|$)`),
		Drivers:             regexp.MustCompile(`(?s)<h3[^>]*>Risk Drivers</h3>\s*(.+?)(?:<h3|$)`),
		Consequences:        regexp.MustCompile(`(?s)<h3[^>]*>Risk Consequences</h3>\s*(.+?)(?:<h3|$)`),
		RecommendedControls: regexp.MustCompile(`(?s)<h3[^>]*>Recommended Controls</h3>\s*(.+?)(?:<h3|$)`),
	}

	ControlPattern = ControlPatterns{
		Description: regexp.MustCompile(`(?s)<h3[^>]*>Control Description</h3>\s*(.+?)(?:<h3|$)`),
		Reduces:     regexp.MustCompile(`(?s)<h3[^>]*>Reduces</h3>\s*(.+?)(?:<h3|$)`),
	}
}
