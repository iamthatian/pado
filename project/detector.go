package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/duckonomy/parkour/utils"
)

// NOTE: Don't have separate ProjectRoot type because
// we might just want to get Project root value in some cases
func (p *Project) FindProjectRoot(path string) error {
	maxDepth := 100
	depth := 0
	currentPath, err := utils.CanonicalizePath(path)
	if err != nil {
		return err
	}

	vcIndicators := VC_INDICATORS()
	workspaceIndicators := WORKSPACE_INDICATORS()
	languageIndicators := LANGUAGE_INDICATORS()

	for depth <= maxDepth {
		// list files
		files, err := os.ReadDir(currentPath)
		if err != nil {
			return err
		}

		if indicator, matched := containsRootIndicator(files, vcIndicators); matched {
			// What does this mean though lol
			p.SimpleType = indicator.Type
			p.Path = currentPath
			return nil
		}

		if indicator, matched := containsRootIndicator(files, workspaceIndicators); matched {
			p.SimpleType = indicator.Type
			p.Path = currentPath
			return nil
		}

		if indicator, matched := containsRootIndicator(files, languageIndicators); matched {
			p.SimpleType = indicator.Type
			p.Path = currentPath
			return nil
		}

		currentPath = filepath.Dir(currentPath)
		if utils.IsHomeDir(currentPath) || utils.IsRootDir(currentPath) {
			break
		}
		depth++
	}

	return fmt.Errorf("no project root found")
}

func containsRootIndicator(files []os.DirEntry, indicators []ProjectRootIndicator) (*ProjectRootIndicator, bool) {
	for _, indicator := range indicators {
		for _, pattern := range indicator.Patterns {
			for _, file := range files {
				if file.Name() == pattern {
					return &indicator, true
				}
			}
		}
	}
	return &ProjectRootIndicator{}, false
}
