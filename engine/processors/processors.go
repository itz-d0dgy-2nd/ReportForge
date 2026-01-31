package processors

import (
	"ReportForge/engine/utilities"
	"ReportForge/engine/validators"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/russross/blackfriday/v2"
)

/*
parseFile → Decorator that wraps file reading with YAML validation
  - Calls utilities.ReadAndCleanMarkdownFile() to read, strip BOM, and normalise line endings
  - Validates YAML frontmatter and populates MarkdownYML struct
  - Returns raw markdown content, regex matches, and parsed YAML
*/
func parseFile(_filePath string, _fileCache *utilities.FileCache) (string, []string, utilities.MarkdownYML) {
	var unprocessedYaml utilities.MarkdownYML

	rawMarkdownContent, regexMatches, err := utilities.ReadAndCleanMarkdownFile(_filePath, _fileCache)
	utilities.ErrorChecker(err)

	validators.ValidateYamlFrontmatter(regexMatches, _filePath, &unprocessedYaml)

	return rawMarkdownContent, regexMatches, unprocessedYaml
}

/*
ProcessMarkdown → Process markdown files with YAML frontmatter for report content
  - Reads markdown file from specified file path and normalises line endings
  - Validates YAML frontmatter using validators.ValidateYamlFrontmatter()
  - Converts markdown content to HTML using blackfriday processor
  - Performs string replacements on:
    -- Custom tokens (!Client, !TargetAsset0, !TargetAsset1)
    -- Custom tags (<retest_fixed>, <retest_not_fixed>)
    -- Screenshot paths to include full report path
  - For findings specifically:
    -- Generates severity matrix update data
    -- Generates severity bar graph update data
  - Returns processed markdown file and optional severity updates
  - Handles errors via utilities.ErrorChecker()
*/
func ProcessMarkdown(_filePath string, _fileCache *utilities.FileCache) (utilities.MarkdownFile, *utilities.SeverityMatrixUpdate, *utilities.SeverityBarGraphUpdate) {
	var unprocessedYaml utilities.MarkdownYML
	var severityMatrixUpdate *utilities.SeverityMatrixUpdate
	var severityBarGraphUpdate *utilities.SeverityBarGraphUpdate

	_, regexMatches, unprocessedYaml := parseFile(_filePath, _fileCache)

	unprocessedMarkdown := string(blackfriday.Run([]byte(regexMatches[2])))

	unprocessedMarkdown = utilities.RegexTokenMatch.ReplaceAllStringFunc(unprocessedMarkdown, func(tokenMatch string) string {
		if tokenValue, exists := _fileCache.MetadataConfig.CustomVariables[strings.TrimPrefix(tokenMatch, "!")]; exists {
			return tokenValue
		}

		if strings.TrimPrefix(tokenMatch, "!") == "Client" {
			return _fileCache.MetadataConfig.Client
		}

		return tokenMatch
	})

	reportRoot := filepath.Dir(filepath.Dir(filepath.Dir(_filePath)))
	unprocessedMarkdown = utilities.RegexMarkdownRetestMatch.ReplaceAllString(unprocessedMarkdown, "<$1$2$3>")
	unprocessedMarkdown = utilities.RegexMarkdownImageMatchScale.ReplaceAllString(unprocessedMarkdown, `$1 src="`+reportRoot+`/$2"$3 style="$4"/>`)
	unprocessedMarkdown = utilities.RegexMarkdownImageMatch.ReplaceAllString(unprocessedMarkdown, `$1 src="`+reportRoot+`/$2"$3/>`)

	if strings.Contains(unprocessedMarkdown, "<qa>") {
		utilities.ErrorChecker(fmt.Errorf("%d QA comment(s) in ( %s )", strings.Count(unprocessedMarkdown, "<qa>"), _filePath))
	}

	markdownFile := utilities.MarkdownFile{
		Directory: filepath.Base(filepath.Dir(_filePath)),
		FileName:  filepath.Base(_filePath),
		Headers:   unprocessedYaml,
		Body:      unprocessedMarkdown,
	}

	if strings.Contains(_filePath, utilities.FindingsDirectory) {
		impactIndex := slices.Index(_fileCache.SeverityConfig.Impacts, unprocessedYaml.FindingImpact)
		likelihoodIndex := slices.Index(_fileCache.SeverityConfig.Likelihoods, unprocessedYaml.FindingLikelihood)

		rowIndex := impactIndex
		columnIndex := likelihoodIndex

		if _fileCache.SeverityConfig.SwapImpactLikelihoodAxis {
			rowIndex = likelihoodIndex
			columnIndex = impactIndex
		}

		if unprocessedYaml.FindingStatus != utilities.FindingsStatusResolved {
			severityMatrixUpdate = &utilities.SeverityMatrixUpdate{
				RowIndex:    rowIndex,
				ColumnIndex: columnIndex,
				FindingID:   unprocessedYaml.FindingID,
			}
		}

		if unprocessedYaml.FindingStatus != utilities.FindingsStatusResolved {
			severityBarGraphUpdate = &utilities.SeverityBarGraphUpdate{
				Severity: unprocessedYaml.FindingSeverity,
				Status:   unprocessedYaml.FindingStatus,
			}
		}
	}

	return markdownFile, severityMatrixUpdate, severityBarGraphUpdate
}
