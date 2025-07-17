package Utils

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ProcessConfigFrontmatter(directory string, file os.DirEntry, frontMatter *FrontmatterYML) {

	currentFileName := file.Name()
	readYML, ErrReadYML := os.ReadFile(filepath.Join(directory, currentFileName))
	ErrorChecker(ErrReadYML)

	ErrDecodeYML := yaml.Unmarshal(readYML, &frontMatter)
	ErrorChecker(ErrDecodeYML)

	for _, information := range frontMatter.DocumentInformation {
		for metadataKey, metadata := range information.DocumentMetadata {
			if metadataKey == "DocumentStatus" {
				if metadata != "Draft" && metadata != "Release" {
					ErrorChecker(fmt.Errorf("invalid documentStatus in frontmatter (%s/%s - %s) - please check that your documentStatus is 'Draft' or 'Release'", directory, currentFileName, metadata))
				}
			}
		}
	}
}
