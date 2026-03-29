package generators

import (
	"ReportForge/engine/utilities"
	"regexp"
	"strings"
)

func extractRiskSection(body, sectionName string) string {
	var regex *regexp.Regexp

	switch sectionName {
	case "Risk Description":
		regex = utilities.RiskPattern.Description
	case "Risk Drivers":
		regex = utilities.RiskPattern.Drivers
	case "Risk Consequences":
		regex = utilities.RiskPattern.Consequences
	case "Recommended Controls":
		regex = utilities.RiskPattern.RecommendedControls
	default:
		return ""
	}

	matches := regex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

func extractControlSection(body, sectionName string) string {
	var regex *regexp.Regexp

	switch sectionName {
	case "Control Description":
		regex = utilities.ControlPattern.Description
	case "Control Reduces":
		regex = utilities.ControlPattern.Reduces
	default:
		return ""
	}

	matches := regex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}
