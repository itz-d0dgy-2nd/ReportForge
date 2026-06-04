package processors

import (
	"ReportForge/engine/utilities"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/russross/blackfriday/v2"
)

/*
ProcessMarkdown → Processes a markdown file with YAML frontmatter into a report-ready MarkdownFile
*/
func ProcessMarkdown(_path string, _fileCache *utilities.FileCache) (utilities.MarkdownFile, *utilities.SeverityMatrixUpdate, *utilities.SeverityBarGraphUpdate, *utilities.RiskMatricesUpdate) {
	var severityMatrixUpdate *utilities.SeverityMatrixUpdate
	var severityBarGraphUpdate *utilities.SeverityBarGraphUpdate
	var riskMatricesUpdate *utilities.RiskMatricesUpdate

	_, regexMatches, unprocessedYaml := utilities.ParseFile(_path, _fileCache)

	unprocessedMarkdown := string(blackfriday.Run([]byte(regexMatches[2])))

	metadataConfig := _fileCache.MetadataConfig()

	unprocessedMarkdown = utilities.MarkdownPattern.Token.ReplaceAllStringFunc(unprocessedMarkdown, func(tokenMatch string) string {
		if tokenValue, exists := metadataConfig.CustomVariables[strings.TrimPrefix(tokenMatch, "!")]; exists {
			return tokenValue
		}

		if strings.TrimPrefix(tokenMatch, "!") == "Client" {
			return metadataConfig.Client
		}

		return tokenMatch
	})

	reportRoot := filepath.Dir(filepath.Dir(filepath.Dir(_path)))
	unprocessedMarkdown = utilities.MarkdownPattern.Retest.ReplaceAllString(unprocessedMarkdown, "<$1$2$3>")
	unprocessedMarkdown = utilities.MarkdownPattern.ImageScale.ReplaceAllString(unprocessedMarkdown, `$1 src="`+reportRoot+`/$2"$3 style="$4"/>`)
	unprocessedMarkdown = utilities.MarkdownPattern.Image.ReplaceAllString(unprocessedMarkdown, `$1 src="`+reportRoot+`/$2"$3/>`)

	if strings.Contains(unprocessedMarkdown, "<qa>") {
		count := strings.Count(unprocessedMarkdown, "<qa>")
		utilities.Check(utilities.NewValidationWarning(
			_path,
			fmt.Sprintf("%d QA comment(s) found - remove before release", count),
		))
	}

	markdownFile := utilities.MarkdownFile{
		Directory: filepath.Base(filepath.Dir(_path)),
		FileName:  filepath.Base(_path),
		Headers:   unprocessedYaml,
		Body:      unprocessedMarkdown,
	}

	if strings.Contains(_path, utilities.FindingsDirectory) {
		severityConfig := _fileCache.SeverityConfig()

		impactIndex := slices.Index(severityConfig.Impacts, unprocessedYaml.FindingImpact)
		likelihoodIndex := slices.Index(severityConfig.Likelihoods, unprocessedYaml.FindingLikelihood)

		rowIndex := impactIndex
		columnIndex := likelihoodIndex

		if severityConfig.SwapImpactLikelihoodAxis {
			rowIndex = likelihoodIndex
			columnIndex = impactIndex
		}

		if unprocessedYaml.FindingStatus != utilities.FindingsStatusResolved {
			severityMatrixUpdate = &utilities.SeverityMatrixUpdate{
				RowIndex:    rowIndex,
				ColumnIndex: columnIndex,
				FindingID:   unprocessedYaml.FindingID,
			}

			severityBarGraphUpdate = &utilities.SeverityBarGraphUpdate{
				Severity: unprocessedYaml.FindingSeverity,
				Status:   unprocessedYaml.FindingStatus,
			}
		}
	}

	if strings.Contains(_path, utilities.RisksDirectory) {
		riskConfig := _fileCache.RiskConfig()

		grossImpactIndex := slices.Index(riskConfig.GrossImpacts, unprocessedYaml.RiskGrossImpact)
		grossLikelihoodIndex := slices.Index(riskConfig.GrossLikelihoods, unprocessedYaml.RiskGrossLikelihood)

		targetImpactIndex := slices.Index(riskConfig.TargetImpacts, unprocessedYaml.RiskTargetImpact)
		targetLikelihoodIndex := slices.Index(riskConfig.TargetLikelihoods, unprocessedYaml.RiskTargetLikelihood)

		riskMatricesUpdate = &utilities.RiskMatricesUpdate{
			GrossRowIndex:     grossImpactIndex,
			GrossColumnIndex:  grossLikelihoodIndex,
			TargetRowIndex:    targetImpactIndex,
			TargetColumnIndex: targetLikelihoodIndex,
			RiskID:            unprocessedYaml.RiskID,
		}
	}

	return markdownFile, severityMatrixUpdate, severityBarGraphUpdate, riskMatricesUpdate
}
