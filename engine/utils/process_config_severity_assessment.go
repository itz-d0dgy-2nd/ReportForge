package Utils

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ProcessConfigMatrix(_directory string, _file os.DirEntry, _severityAssessment *SeverityAssessmentYML) {

	currentFileName := _file.Name()
	readYML, errReadYML := os.ReadFile(filepath.Clean(filepath.Join(_directory, currentFileName)))
	ErrorChecker(errReadYML)

	errDecodeYML := yaml.Unmarshal(readYML, &_severityAssessment)
	ErrorChecker(errDecodeYML)

	// Add Error Handling for Impacts & Likelihoods

}
