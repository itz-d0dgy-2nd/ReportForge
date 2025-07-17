package Utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

func ProcessMarkdown(frontmatter FrontmatterYML, severityAssessment SeverityAssessmentYML, directory string, file os.DirEntry, storage *[]Markdown) []Markdown {

	processedYML := MarkdownYML{}
	impact := 0
	likelihood := 0

	currentFileName := file.Name()
	readMD, ErrReadMD := os.ReadFile(filepath.Join(directory, currentFileName))
	ErrorChecker(ErrReadMD)

	regexYML := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
	regexMatches := regexYML.FindStringSubmatch(string(readMD))

	ErrDecodeYML := yaml.Unmarshal([]byte(regexMatches[1]), &processedYML)
	ErrorChecker(ErrDecodeYML)

	if strings.Contains(directory, "findings") {

		for key, value := range severityAssessment.Impacts {
			if value == processedYML.FindingImpact {
				impact = key
			}
		}

		for key, value := range severityAssessment.Likelihoods {
			if value == processedYML.FindingLikelihood {
				likelihood = key
			}
		}

		if _, validImpact := severityAssessment.Impacts[impact]; !validImpact {
			ErrorChecker(fmt.Errorf("invalid impact in finding (%s/%s - %s) - please check that your impact is supported", directory, processedYML.FindingName, processedYML.FindingImpact))
		}

		if _, validLikelihoods := severityAssessment.Likelihoods[likelihood]; !validLikelihoods {
			ErrorChecker(fmt.Errorf("invalid likelihood in finding (%s/%s - %s) - please check that your likelihood is supported", directory, processedYML.FindingName, processedYML.FindingLikelihood))
		}

		processedYML.FindingSeverity = severityAssessment.CalculatedMatrix[impact][likelihood]
		currentFileName = processedYML.FindingSeverity + "_" + processedYML.FindingName + ".md"
		ErrRename := os.Rename(filepath.Join(directory, file.Name()), filepath.Join(directory, currentFileName))
		ErrorChecker(ErrRename)
	}

	if strings.Contains(directory, "suggestions") {
		currentFileName = "Suggestion_" + processedYML.SuggestionName + ".md"
		ErrRename := os.Rename(filepath.Join(directory, file.Name()), filepath.Join(directory, currentFileName))
		ErrorChecker(ErrRename)
	}

	processedMD := string(blackfriday.Run([]byte(regexMatches[2])))

	if strings.Contains(processedMD, "<qa>") {
		fmt.Printf("::warning:: %s: %d QA Comment Present In File \n", currentFileName, strings.Count(processedMD, "<qa>"))
	}

	if strings.Contains(processedMD, "!Client") {
		processedMD = strings.ReplaceAll(processedMD, "!Client", frontmatter.Client)
	}

	if strings.Contains(processedMD, "!TargetAsset0") {
		processedMD = strings.ReplaceAll(processedMD, "!TargetAsset0", frontmatter.TargetInformation["TargetAsset0"])
	}

	if strings.Contains(processedMD, "!TargetAsset1") {
		processedMD = strings.ReplaceAll(processedMD, "!TargetAsset1", frontmatter.TargetInformation["TargetAsset1"])
	}

	if strings.Contains(processedMD, "<p><retest_fixed></p>") {
		processedMD = strings.ReplaceAll(processedMD, "<p><retest_fixed></p>", "<retest_fixed>")
		processedMD = strings.ReplaceAll(processedMD, "<p></retest_fixed></p>", "</retest_fixed>")
	}

	if strings.Contains(processedMD, "<p><retest_not_fixed></p>") {
		processedMD = strings.ReplaceAll(processedMD, "<p><retest_not_fixed></p>", "<retest_not_fixed>")
		processedMD = strings.ReplaceAll(processedMD, "<p></retest_not_fixed></p>", "</retest_not_fixed>")
	}

	*storage = append(*storage, Markdown{
		Directory: filepath.Base(directory),
		FileName:  strings.TrimSuffix(currentFileName, filepath.Ext(currentFileName)),
		Headers:   processedYML,
		Body:      processedMD,
	})

	return *storage

}
