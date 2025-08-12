package processors

import (
	"ReportForge/engine/utils"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ProcessConfigMatrix(_directory string, _file os.DirEntry, _severityAssessment *Utils.SeverityAssessmentYML) {

	currentFileName := _file.Name()
	readYML, errReadYML := os.ReadFile(filepath.Clean(filepath.Join(_directory, currentFileName)))
	Utils.ErrorChecker(errReadYML)

	errDecodeYML := yaml.Unmarshal(readYML, &_severityAssessment)
	Utils.ErrorChecker(errDecodeYML)

	// Add Error Handling for Impacts & Likelihoods

}
