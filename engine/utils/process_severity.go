package Utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

func ProcessSeverityMatrix(directory string, file os.DirEntry, severityAssessment *SeverityAssessmentYML) {

	processedYML := MarkdownYML{}
	impact := 0
	likelihood := 0

	currentFileName := file.Name()
	readMD, ErrReadMD := os.ReadFile(filepath.Join(directory, currentFileName))
	ErrorChecker(ErrReadMD)

	regexYML := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
	regexMatches := regexYML.FindStringSubmatch(string(readMD))

	ErrDecodeYML := yaml.Unmarshal([]byte(regexMatches[1]), &processedYML)
	ErrorChecker(ErrDecodeYML)

	for key, value := range severityAssessment.Impacts {
		if value == processedYML.FindingImpact {
			impact = key
		}
	}

	for key, value := range severityAssessment.Likelihoods {
		if value == processedYML.FindingLikelihood {
			likelihood = key
		}
	}

	if _, validImpact := severityAssessment.Impacts[impact]; !validImpact {
		ErrorChecker(fmt.Errorf("invalid impact in finding (%s/%s - %s) - please check that your impact is supported", directory, processedYML.FindingName, processedYML.FindingImpact))
	}

	if _, validLikelihoods := severityAssessment.Likelihoods[likelihood]; !validLikelihoods {
		ErrorChecker(fmt.Errorf("invalid likelihood in finding (%s/%s - %s) - please check that your likelihood is supported", directory, processedYML.FindingName, processedYML.FindingLikelihood))
	}

	if severityAssessment.Matrix[impact][likelihood] == "" {
		severityAssessment.Matrix[impact][likelihood] = processedYML.FindingID
	} else {
		severityAssessment.Matrix[impact][likelihood] += ", " + processedYML.FindingID
	}

}
