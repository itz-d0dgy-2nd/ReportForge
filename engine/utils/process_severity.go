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
	impact := 0
	likelihood := 0

	currentFileName := file.Name()
	readMD, ErrReadMD := os.ReadFile(filepath.Join(directory, currentFileName))
	ErrorChecker(ErrReadMD)

	regex := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
	regexMatches := regex.FindStringSubmatch(string(readMD))

	ErrDecodeYML := yaml.Unmarshal([]byte(regexMatches[1]), &processedYML)
	ErrorChecker(ErrDecodeYML)

	for key, value := range storage.Impacts {
		if value == processedYML.FindingImpact {
			impact = key
		}
	}

	for key, value := range storage.Likelihoods {
		if value == processedYML.FindingLikelihood {
			likelihood = key
		}
	}

	if _, validImpact := storage.Impacts[impact]; !validImpact {
		ErrorChecker(fmt.Errorf("invalid impact in finding (%s/%s - %s) - please check that your impact is supported", directory, processedYML.FindingName, processedYML.FindingImpact))
	}

	if _, validLikelihoods := storage.Likelihoods[likelihood]; !validLikelihoods {
		ErrorChecker(fmt.Errorf("invalid likelihood in finding (%s/%s - %s) - please check that your likelihood is supported", directory, processedYML.FindingName, processedYML.FindingLikelihood))
	}

	if storage.Matrix[impact][likelihood] == "" {
		storage.Matrix[impact][likelihood] = processedYML.FindingID
	} else {
		storage.Matrix[impact][likelihood] += ", " + processedYML.FindingID
	}

}
