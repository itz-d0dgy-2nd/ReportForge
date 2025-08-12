package Utils

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ProcessConfigFrontmatter(_directory string, _file os.DirEntry, _frontmatter *FrontmatterYML) {

	currentFileName := _file.Name()
	readYML, errReadYML := os.ReadFile(filepath.Clean(filepath.Join(_directory, currentFileName)))
	ErrorChecker(errReadYML)

	errDecodeYML := yaml.Unmarshal(readYML, &_frontmatter)
	ErrorChecker(errDecodeYML)

	for _, information := range _frontmatter.DocumentInformation {
		if status, exists := information.DocumentMetadata["DocumentStatus"]; exists {
			if status != "Draft" && status != "Release" {
				ErrorChecker(fmt.Errorf("invalid documentStatus in frontmatter (%s/%s - %s) - please check that your documentStatus is 'Draft' or 'Release'", _directory, currentFileName, status))
			}
		}
	}
}
