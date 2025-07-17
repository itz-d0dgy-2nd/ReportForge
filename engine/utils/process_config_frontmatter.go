package Utils

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ProcessConfigFrontmatter(directory string, file os.DirEntry, storage *FrontmatterYML) {

	currentFileName := file.Name()
	readYML, ErrReadYML := os.ReadFile(filepath.Join(directory, currentFileName))
	ErrorChecker(ErrReadYML)

	ErrDecodeYML := yaml.Unmarshal(readYML, &storage)
	ErrorChecker(ErrDecodeYML)

	// for _, information := range processedYML.DocumentInformation {
	// 	for metadataKey, metadata := range information.DocumentMetadata {
	// 		if metadataKey == "DocumentStatus" {
	// 			if metadata != "Draft" && metadata != "Release" {
	// 				ErrorChecker(fmt.Errorf("invalid document statu in frontmatter (%s/%s - %s) - please check that your document status is supported", directory, , processedYML.FindingLikelihood))
	// 			}
	// 		}
	// 	}
	// }
}
