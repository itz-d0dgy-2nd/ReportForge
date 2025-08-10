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

func ProcessMarkdown(_reportTemplatePath string, _frontmatter FrontmatterYML, _severityAssessment SeverityAssessmentYML, _directory string, _file os.DirEntry, _storage *[]Markdown) []Markdown {

	processedYML := MarkdownYML{}
	impact := -1
	likelihood := -1

	currentFileName := _file.Name()
	readMD, ErrReadMD := os.ReadFile(filepath.Join(_directory, currentFileName))
	ErrorChecker(ErrReadMD)

	regexYML := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
	regexMatches := regexYML.FindStringSubmatch(string(readMD))

	ErrDecodeYML := yaml.Unmarshal([]byte(regexMatches[1]), &processedYML)
	ErrorChecker(ErrDecodeYML)

	if strings.Contains(_directory, "findings") {

		for key, value := range _severityAssessment.Impacts {
			if value == processedYML.FindingImpact {
				impact = key
			}
		}

		for key, value := range _severityAssessment.Likelihoods {
			if value == processedYML.FindingLikelihood {
				likelihood = key
			}
		}

		if _, validImpact := _severityAssessment.Impacts[impact]; !validImpact {
			ErrorChecker(fmt.Errorf("invalid impact in finding (%s/%s - %s) - please check that your impact is supported", _directory, processedYML.FindingName, processedYML.FindingImpact))
		}

		if _, validLikelihoods := _severityAssessment.Likelihoods[likelihood]; !validLikelihoods {
			ErrorChecker(fmt.Errorf("invalid likelihood in finding (%s/%s - %s) - please check that your likelihood is supported", _directory, processedYML.FindingName, processedYML.FindingLikelihood))
		}

		processedYML.FindingSeverity = _severityAssessment.CalculatedMatrix[impact][likelihood]
		currentFileName = processedYML.FindingSeverity + "_" + processedYML.FindingName + ".md"
		ErrRename := os.Rename(filepath.Join(_directory, _file.Name()), filepath.Join(_directory, currentFileName))
		ErrorChecker(ErrRename)
	}

	if strings.Contains(_directory, "suggestions") {
		currentFileName = "Suggestion_" + processedYML.SuggestionName + ".md"
		ErrRename := os.Rename(filepath.Join(_directory, _file.Name()), filepath.Join(_directory, currentFileName))
		ErrorChecker(ErrRename)
	}

	processedMD := string(blackfriday.Run([]byte(regexMatches[2])))

	if strings.Contains(processedMD, "<qa>") {
		fmt.Printf("::warning:: %s: %d QA Comment Present In File \n", currentFileName, strings.Count(processedMD, "<qa>"))
	}

	if strings.Contains(processedMD, "!Client") {
		processedMD = strings.ReplaceAll(processedMD, "!Client", _frontmatter.Client)
	}

	if strings.Contains(processedMD, "!TargetAsset0") {
		processedMD = strings.ReplaceAll(processedMD, "!TargetAsset0", _frontmatter.TargetInformation["TargetAsset0"])
	}

	if strings.Contains(processedMD, "!TargetAsset1") {
		processedMD = strings.ReplaceAll(processedMD, "!TargetAsset1", _frontmatter.TargetInformation["TargetAsset1"])
	}

	if strings.Contains(processedMD, "<p><retest_fixed></p>") {
		processedMD = strings.ReplaceAll(processedMD, "<p><retest_fixed></p>", "<retest_fixed>")
		processedMD = strings.ReplaceAll(processedMD, "<p></retest_fixed></p>", "</retest_fixed>")
	}

	if strings.Contains(processedMD, "<p><retest_not_fixed></p>") {
		processedMD = strings.ReplaceAll(processedMD, "<p><retest_not_fixed></p>", "<retest_not_fixed>")
		processedMD = strings.ReplaceAll(processedMD, "<p></retest_not_fixed></p>", "</retest_not_fixed>")
	}

	if strings.Contains(processedMD, "Screenshots/") {
		processedMD = strings.ReplaceAll(processedMD, "Screenshots/", _reportTemplatePath+"/Screenshots/")
	}

	*_storage = append(*_storage, Markdown{
		Directory: filepath.Base(_directory),
		FileName:  strings.TrimSuffix(currentFileName, filepath.Ext(currentFileName)),
		Headers:   processedYML,
		Body:      processedMD,
	})

	return *_storage

}
