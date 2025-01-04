// This module implements everything related to finding the project root
// And provides default matchers
package project

type ProjectMetadata struct {
	License     string
	VCS         string
	Languages   []string
	BuildSystem string
}

type Project struct {
	Path         string
	Type         ProjectType
	SimpleType   string
	Name         string
	Priority     int
	Parent       *Project
	Children     []*Project
	Dirty        bool
	Metadata     ProjectMetadata
	BuildCommand []string
	RunCommand   []string
	TestCommand  []string
}

type RegexString string

func (p *Project) InitProject(path string) error {
	err := p.FindProjectRoot(path)
	if err != nil {
		return err
	}

	err = p.AnaylzeProject()
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) IsEmpty() bool {
	return len(p.Path) == 0
}

// Merge this later with values from config
func DEFAULT_IGNORE() []RegexString {
	return []RegexString{
		// Version Control
		`^\.git$`,
		`^\.hg$`,
		`^\.svn$`,

		// Build outputs
		`^target$`,
		`^build$`,
		`^dist$`,
		`^out$`,

		// Dependencies
		`^node_modules$`,
		`^vendor$`,
		`^\.venv$`,

		// IDE
		`^\.idea$`,
		`^\.vscode$`,
		`^\.vs$`,

		// OS
		`^\.DS_Store$`,
		`^Thumbs\.db$`,

		// Logs
		`^\.log$`,
		`^logs$`,
		`^npm-debug\.log.*`,

		// Temp files
		`^\.tmp$`,
		`^\.temp$`,
		`^\.cache$`,

		// Build artifacts
		`^\.class$`,
		`^\.pyc$`,
		`^\.pyo$`,
		`^\.o$`,
		`^\.obj$`,
	}
}
