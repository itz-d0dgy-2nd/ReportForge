<p align="center">
    <img src="engine/template/html/images/logo.png">
</p>

> [!Note]
> As someone who doesn’t work with Go on a regular basis, this codebase may not follow standard design patterns or best practices. 
> It was primarily developed to address a specific problem, and while it gets the job done, there's certainly room for improvement.
> 
> If you're familiar with Go and have suggestions for better practices or optimisations, please feel free to open a pull request! I'm always happy to improve the project and welcome contributions.

For detailed information on this project, please read the [wiki](https://github.com/itz-d0dgy-2nd/ReportForge/wiki)

## Roadmap:

### Improvements
- [ ] Improve finding cross referencing
- [ ] Improve XLSX generation 

### Implement
- [ ] Footnotes
- [ ] Better error handling in most if not all areas
    - [ ] File Handlers - What if the files dont exist? - Proposed format ("missing XYZ.XYZ file (%s/.) - please check that you have not deleted this document")
        - [ ] `file_handler_report_config.go`
        - [ ] `file_handler_markdown.go`
        - [ ] `file_handler_severity.go`
    - [ ] Generators - What if generation failed? - Proposed format ("")
        - [ ] `generator_html.go` 
        - [ ] `generator_pdf`
        - [ ] `generator_xslx.go`
    - [ ] Processors - What if data is not correct / expected - Proposed format ("invalid XYZ in XYZ (%s/%s - %s) - please check that your XYZ is ")
        - [ ] `process_config_front_matter.go`
        - [ ] `process_config_severity_assessment.go`
        - [ ] `process_config_markdown.go`
        - [ ] `process_severity`

### Refactor
- [ ] Refactor HTML
- [ ] Refactor CSS
- [ ] Refactor GoLang
- [ ] Refactor GoLang Templating

### Stretch Goals:
- [ ] Improve the `report.html` output so that it could be used
- [ ] Create ReportForge binary
