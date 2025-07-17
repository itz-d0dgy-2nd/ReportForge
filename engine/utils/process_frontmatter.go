package Utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ProcessFrontmatter(file string) FrontmatterYML {

	processedYML := FrontmatterYML{}

	readYML, ErrReadYML := os.ReadFile(file)
	ErrorChecker(ErrReadYML)

	ErrDecodeYML := yaml.Unmarshal(readYML, &processedYML)
	ErrorChecker(ErrDecodeYML)

	return processedYML
}
