package Utils

import (
	"bytes"
	"encoding/json"
	"os"
)

func FrontmatterHandler(file string) FrontmatterJSON {

	projectFrontmatter := FrontmatterJSON{}

	readJSON, ErrReadJSON := os.ReadFile(file)
	ErrorChecker(ErrReadJSON)

	ErrDecodeJSON := json.NewDecoder(bytes.NewReader(readJSON)).Decode(&projectFrontmatter)
	ErrorChecker(ErrDecodeJSON)

	return projectFrontmatter

}
