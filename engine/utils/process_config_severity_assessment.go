package Utils

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ProcessConfigMatrix(directory string, file os.DirEntry, severityAssessment *SeverityAssessmentYML) {

	currentFileName := file.Name()
	readYML, ErrReadYML := os.ReadFile(filepath.Join(directory, currentFileName))
	ErrorChecker(ErrReadYML)

	ErrDecodeYML := yaml.Unmarshal(readYML, &severityAssessment)
	ErrorChecker(ErrDecodeYML)

}
