package Utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

func ProcessSeverityMatrix(directory string, file os.DirEntry, storage *SeverityMatrix) {

	processedYML := MarkdownYML{}

	readMD, ErrReadMD := os.ReadFile(filepath.Join(directory, file.Name()))
	ErrorChecker(ErrReadMD)

	regex := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
	regexMatches := regex.FindStringSubmatch(string(readMD))

	ErrDecodeYML := yaml.Unmarshal([]byte(regexMatches[1]), &processedYML)
	ErrorChecker(ErrDecodeYML)

	if _, validImpact := storage.Impacts[processedYML.FindingImpact]; !validImpact {
		ErrorChecker(fmt.Errorf("invalid impact in finding (%s/%s - %s) - please check that your impact is supported", directory, processedYML.FindingName, processedYML.FindingImpact))
	}

	if _, validLikelihoods := storage.Likelihoods[processedYML.FindingLikelihood]; !validLikelihoods {
		ErrorChecker(fmt.Errorf("invalid likelihood in finding (%s/%s) - %s", directory, processedYML.FindingName, processedYML.FindingLikelihood))
	}

	if storage.Matrix[storage.Impacts[processedYML.FindingImpact]][storage.Likelihoods[processedYML.FindingLikelihood]] == "" {
		storage.Matrix[storage.Impacts[processedYML.FindingImpact]][storage.Likelihoods[processedYML.FindingLikelihood]] = processedYML.FindingID
	} else {
		storage.Matrix[storage.Impacts[processedYML.FindingImpact]][storage.Likelihoods[processedYML.FindingLikelihood]] += ", " + processedYML.FindingID
	}

}

func CalculateSeverity(impact string, likelihood string) string {

	// [^DIA]: https://www.digital.govt.nz/standards-and-guidance/privacy-security-and-risk/risk-management/risk-assessments/analyse/initial-risk-ratings#table-1
	intersections := SeverityMatrix{
		map[string]int{
			"Severe":      0,
			"Significant": 1,
			"Moderate":    2,
			"Minor":       3,
			"Minimal":     4,
		},

		map[string]int{
			"Almost Never":          0,
			"Possible but Unlikely": 1,
			"Possible":              2,
			"Highly Probable":       3,
			"Almost Certain":        4,
		},

		[5][5]string{
			{"1", "1", "0", "0", "0"},
			{"2", "1", "1", "1", "0"},
			{"2", "2", "2", "1", "1"},
			{"3", "2", "2", "2", "1"},
			{"3", "3", "2", "2", "2"},
		},
	}

	return intersections.Matrix[intersections.Impacts[impact]][intersections.Likelihoods[likelihood]]

}
