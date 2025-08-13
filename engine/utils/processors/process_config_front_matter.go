package processors

import (
	"ReportForge/engine/utils"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

/*
ProcessConfigFrontmatter → Process yaml files
  - Reads yaml configuration file
  - Unmarshals yaml content into type *Utils.FrontmatterYML
  - Validate _frontmatter.DocumentInformation
  - Updates frontmatter pointer with processed data
*/
func ProcessConfigFrontmatter(_directory string, _file os.DirEntry, _frontmatter *Utils.FrontmatterYML) {

	currentFileName := _file.Name()
	readYML, errReadYML := os.ReadFile(filepath.Clean(filepath.Join(_directory, currentFileName)))
	Utils.ErrorChecker(errReadYML)

	errDecodeYML := yaml.Unmarshal(readYML, &_frontmatter)
	Utils.ErrorChecker(errDecodeYML)

	for _, information := range _frontmatter.DocumentInformation {
		if status, exists := information.DocumentMetadata["DocumentStatus"]; exists {
			if status != "Draft" && status != "Release" {
				Utils.ErrorChecker(fmt.Errorf("invalid documentStatus in frontmatter (%s/%s - %s) - please check that your documentStatus is 'Draft' or 'Release'", _directory, currentFileName, status))
			}
		}
	}
}
