package Utils

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ProcessConfigMatrix(_directory string, _file os.DirEntry, _severityAssessment *SeverityAssessmentYML) {

	currentFileName := _file.Name()
	readYML, ErrReadYML := os.ReadFile(filepath.Join(_directory, currentFileName))
	ErrorChecker(ErrReadYML)

	ErrDecodeYML := yaml.Unmarshal(readYML, &_severityAssessment)
	ErrorChecker(ErrDecodeYML)

}
