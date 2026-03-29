package handlers

import (
	"ReportForge/engine/processors"
	"ReportForge/engine/utilities"
	"io/fs"
	"path/filepath"
	"runtime"
	"sync"
)

/*
updateSeverityMatrix → Updates severity matrix with finding data
*/
func updateSeverityMatrix(_fileCache *utilities.FileCache, _update *utilities.SeverityMatrixUpdate) {
	row, col := _update.RowIndex, _update.ColumnIndex

	if _fileCache.SeverityMatrix.Matrix[row][col] == "" {
		_fileCache.SeverityMatrix.Matrix[row][col] = _update.FindingID
	} else {
		_fileCache.SeverityMatrix.Matrix[row][col] += ", " + _update.FindingID
	}
}

/*
updateSeverityBarGraph → Updates severity bar graph with finding data
*/
func updateSeverityBarGraph(_fileCache *utilities.FileCache, _update *utilities.SeverityBarGraphUpdate) {
	if _fileCache.SeverityBarGraph.Severities == nil {
		_fileCache.SeverityBarGraph.Severities = make(map[string]int)
	}

	_fileCache.SeverityBarGraph.Total++

	switch _update.Status {
	case "Resolved":
		_fileCache.SeverityBarGraph.Resolved++
	case "Unresolved":
		_fileCache.SeverityBarGraph.Unresolved++
		if _update.Severity != "" {
			_fileCache.SeverityBarGraph.Severities[_update.Severity]++
		}
	}
}

/*
updateRiskMatrices → Updates risk matrices with risk data
*/
func updateRiskMatrices(_fileCache *utilities.FileCache, _update *utilities.RiskMatricesUpdate) {
	grossRow, grossCol := _update.GrossRowIndex, _update.GrossColumnIndex
	if _fileCache.RiskMatrices.GrossMatrix[grossRow][grossCol] == "" {
		_fileCache.RiskMatrices.GrossMatrix[grossRow][grossCol] = _update.RiskID
	} else {
		_fileCache.RiskMatrices.GrossMatrix[grossRow][grossCol] += ", " + _update.RiskID
	}

	targetRow, targetCol := _update.TargetRowIndex, _update.TargetColumnIndex
	if _fileCache.RiskMatrices.TargetMatrix[targetRow][targetCol] == "" {
		_fileCache.RiskMatrices.TargetMatrix[targetRow][targetCol] = _update.RiskID
	} else {
		_fileCache.RiskMatrices.TargetMatrix[targetRow][targetCol] += ", " + _update.RiskID
	}
}

/*
processingCollector → Aggregates all processing results into the file cache
  - Routes markdown files to appropriate slices based on directory type
  - Aggregates severity matrix and bar graph data for findings
  - Aggregates risk matrix data for risks
*/
func processingCollector(_results <-chan utilities.ProcessingResult, _fileCache *utilities.FileCache, _waitGroup *sync.WaitGroup) {
	defer _waitGroup.Done()

	for result := range _results {
		switch result.Directory {
		case utilities.SummariesDirectory:
			_fileCache.Summaries = append(_fileCache.Summaries, result.Markdown)
		case utilities.FindingsDirectory:
			_fileCache.Findings = append(_fileCache.Findings, result.Markdown)
		case utilities.SuggestionsDirectory:
			_fileCache.Suggestions = append(_fileCache.Suggestions, result.Markdown)
		case utilities.RisksDirectory:
			_fileCache.Risks = append(_fileCache.Risks, result.Markdown)
		case utilities.ControlsDirectory:
			_fileCache.Controls = append(_fileCache.Controls, result.Markdown)
		case utilities.AppendicesDirectory:
			_fileCache.Appendices = append(_fileCache.Appendices, result.Markdown)
		}

		if result.Directory == utilities.FindingsDirectory {
			severityConfig := _fileCache.SeverityConfig()

			if severityConfig.ConductSeverityAssessment && result.SeverityMatrix != nil {
				updateSeverityMatrix(_fileCache, result.SeverityMatrix)
			}

			if severityConfig.DisplaySeverityBarGraph && result.SeverityBarGraph != nil {
				updateSeverityBarGraph(_fileCache, result.SeverityBarGraph)
			}
		}

		if result.Directory == utilities.RisksDirectory && result.RiskMatrices != nil {
			updateRiskMatrices(_fileCache, result.RiskMatrices)
		}
	}
}

/*
processingWorker → Processes jobs from the jobs channel and sends results to results channel
*/
func processingWorker(_jobs <-chan utilities.ProcessingJob, _results chan<- utilities.ProcessingResult, _fileCache *utilities.FileCache, _waitGroup *sync.WaitGroup) {
	defer _waitGroup.Done()

	for job := range _jobs {
		markdown, severityMatrix, severityBarGraph, riskMatrices := processors.ProcessMarkdown(job.Path, _fileCache)

		_results <- utilities.ProcessingResult{
			Markdown:         markdown,
			SeverityMatrix:   severityMatrix,
			SeverityBarGraph: severityBarGraph,
			RiskMatrices:     riskMatrices,
			Directory:        job.Directory,
		}
	}
}

/*
assignProcessorJobs → Walks a directory and sends processing jobs for each markdown file
*/
func assignProcessorJobs(_path, _directory string, _jobs chan<- utilities.ProcessingJob) {
	filepath.WalkDir(_path, func(path string, entry fs.DirEntry, errAnonymousFunction error) error {
		if errAnonymousFunction != nil {
			return errAnonymousFunction
		}

		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" || utilities.IsRootLevelFile(path) {
			return nil
		}

		_jobs <- utilities.ProcessingJob{
			Path:      path,
			Directory: _directory,
		}

		return nil
	})
}

/*
HandleProcessing → Walks all content directories concurrently using worker pool pattern
*/
func HandleProcessing(_reportPaths utilities.ReportPaths, _fileCache *utilities.FileCache) {
	var workersWaitgroup sync.WaitGroup
	var collectorWaitgroup sync.WaitGroup

	workers := runtime.NumCPU()
	jobs := make(chan utilities.ProcessingJob, workers*2)
	results := make(chan utilities.ProcessingResult, workers*2)

	collectorWaitgroup.Add(1)
	go processingCollector(results, _fileCache, &collectorWaitgroup)

	for i := 0; i < workers; i++ {
		workersWaitgroup.Add(1)
		go processingWorker(jobs, results, _fileCache, &workersWaitgroup)
	}

	directories := map[string]string{
		_reportPaths.SummariesPath:   utilities.SummariesDirectory,
		_reportPaths.FindingsPath:    utilities.FindingsDirectory,
		_reportPaths.SuggestionsPath: utilities.SuggestionsDirectory,
		_reportPaths.RisksPath:       utilities.RisksDirectory,
		_reportPaths.ControlsPath:    utilities.ControlsDirectory,
		_reportPaths.AppendicesPath:  utilities.AppendicesDirectory,
	}

	for path, directory := range directories {
		if path == "" {
			continue
		}
		assignProcessorJobs(path, directory, jobs)
	}

	close(jobs)
	workersWaitgroup.Wait()

	close(results)
	collectorWaitgroup.Wait()
}
