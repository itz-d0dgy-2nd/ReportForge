package Utils

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ProcessConfigFrontmatter(_directory string, _file os.DirEntry, _frontMatter *FrontmatterYML) {

	currentFileName := _file.Name()
	readYML, ErrReadYML := os.ReadFile(filepath.Join(_directory, currentFileName))
	ErrorChecker(ErrReadYML)

	ErrDecodeYML := yaml.Unmarshal(readYML, &_frontMatter)
	ErrorChecker(ErrDecodeYML)

	for _, information := range _frontMatter.DocumentInformation {
		for metadataKey, metadata := range information.DocumentMetadata {
			if metadataKey == "DocumentStatus" {
				if metadata != "Draft" && metadata != "Release" {
					ErrorChecker(fmt.Errorf("invalid documentStatus in frontmatter (%s/%s - %s) - please check that your documentStatus is 'Draft' or 'Release'", _directory, currentFileName, metadata))
				}
			}
		}
	}
}
