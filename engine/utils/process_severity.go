package Utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

func ProcessSeverityMatrix(_directory string, _file os.DirEntry, _severityAssessment *SeverityAssessmentYML) {

	processedYML := MarkdownYML{}
	impact := -1
	likelihood := -1

	currentFileName := _file.Name()
	readMD, errReadMD := os.ReadFile(filepath.Join(_directory, currentFileName))
	ErrorChecker(errReadMD)

	regexYML := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
	regexMatches := regexYML.FindStringSubmatch(string(readMD))

	errDecodeYML := yaml.Unmarshal([]byte(regexMatches[1]), &processedYML)
	ErrorChecker(errDecodeYML)

	for key, value := range _severityAssessment.Impacts {
		if value == processedYML.FindingImpact {
			impact = key
		}
	}

	for key, value := range _severityAssessment.Likelihoods {
		if value == processedYML.FindingLikelihood {
			likelihood = key
		}
	}

	if _, validImpact := _severityAssessment.Impacts[impact]; !validImpact {
		ErrorChecker(fmt.Errorf("invalid impact in finding (%s/%s - %s) - please check that your impact is supported", _directory, processedYML.FindingName, processedYML.FindingImpact))
	}

	if _, validLikelihoods := _severityAssessment.Likelihoods[likelihood]; !validLikelihoods {
		ErrorChecker(fmt.Errorf("invalid likelihood in finding (%s/%s - %s) - please check that your likelihood is supported", _directory, processedYML.FindingName, processedYML.FindingLikelihood))
	}

	if _severityAssessment.Matrix[impact][likelihood] == "" {
		_severityAssessment.Matrix[impact][likelihood] = processedYML.FindingID
	} else {
		_severityAssessment.Matrix[impact][likelihood] += ", " + processedYML.FindingID
	}

}
