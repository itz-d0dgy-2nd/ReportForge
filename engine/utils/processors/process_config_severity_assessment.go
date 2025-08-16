package processors

import (
	"ReportForge/engine/utils"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

/*
ProcessConfigMatrix → Process yaml files
  - Reads yaml configuration file
  - Unmarshals yaml content into type *Utils.SeverityAssessmentYML
  - Validate _severityAssessment.Impacts, _severityAssessment.Likelihoods, _severityAssessment.Matrix, _severityAssessment.CalculatedMatrix
  - Updates _severityAssessment pointer with processed data
*/
func ProcessConfigMatrix(_directory string, _file os.DirEntry, _severityAssessment *Utils.SeverityAssessmentYML) {

	currentFileName := _file.Name()
	readYML, errReadYML := os.ReadFile(filepath.Clean(filepath.Join(_directory, currentFileName)))
	Utils.ErrorChecker(errReadYML)

	errDecodeYML := yaml.Unmarshal(readYML, &_severityAssessment)
	Utils.ErrorChecker(errDecodeYML)

	if !_severityAssessment.SeverityAssessmentEnabled {
		return
	}

	for _, impact := range _severityAssessment.Impacts {
		if impact == "" {
			Utils.ErrorChecker(fmt.Errorf("invalid impact in severity_config (%s/%s - %v) - please check that your impact is not nil", _directory, currentFileName, nil))
		}
	}

	for _, likelihood := range _severityAssessment.Likelihoods {
		if likelihood == "" {
			Utils.ErrorChecker(fmt.Errorf("invalid likelihood in severity_config (%s/%s - %v) - please check that your likelihood is not nil", _directory, currentFileName, nil))
		}
	}

	for rowIndex, row := range _severityAssessment.Matrix {
		for columnIndex, column := range row {
			if column != "" {
				Utils.ErrorChecker(fmt.Errorf("invalid value in severity_config Matrix[%d][%d] (%s/%s - %v) - please check that your Matrix is nil", rowIndex, columnIndex, _directory, currentFileName, nil))
			}
		}
	}

	for rowIndex, row := range _severityAssessment.CalculatedMatrix {
		for columnIndex, column := range row {
			if column == "" {
				Utils.ErrorChecker(fmt.Errorf("invalid severity in severity_config CalculatedMatrix[%d][%d] (%s/%s - %v) - please check that your severity is not nil", rowIndex, columnIndex, _directory, currentFileName, nil))
			}
		}
	}
}
