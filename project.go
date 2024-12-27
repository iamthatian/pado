package main

import (
	"os"
	"path/filepath"
)

// Project struct represents a project with metadata.
type Project struct {
	Name        string
	Path        string
	Kind        string
	Description string
	Priority    int
}

// IsEmpty checks if the project is not initialized.
func (p *Project) IsEmpty() bool {
	return len(p.Path) == 0
}

// FindProject traverses directories to locate project-specific files.
func (p *Project) FindProject(startPath string, maxDepth int) error {
	currentPath, err := NormalizePath(startPath)
	if err != nil {
		return err
	}

	// currentPath := normalizedPath
	depth := 0

	for currentPath != "/" && depth <= maxDepth {
		files, err := os.ReadDir(currentPath)
		if err != nil {
			return err
		}

		if matchedFile := matchProjectFiles(files, getProjectFiles()); matchedFile != "" {
			p.Path = filepath.Clean(currentPath)
			return nil
		}

		currentPath = filepath.Dir(currentPath)
		depth++
	}

	p.Path = "/"

	return nil
}

// matchProjectFiles checks if any file matches known project identifiers.
func matchProjectFiles(fileList []os.DirEntry, patterns []string) string {
	for _, pattern := range patterns {
		for _, file := range fileList {
			if file.Name() == pattern {
				return filepath.Clean(file.Name())
			}
		}
	}
	return ""
}

// getProjectFiles returns a list of file patterns associated with projects.
func getProjectFiles() []string {
	return []string{
		".git",
		".gitignore",
		"pyproject.toml",
		"compile_commands.json", // ccls clangd
		"compile_flags.txt",     // ccls clangd
		".ccls-cache",           // ccls
		".clangd",               // clangd
		"pubspec.yaml",          // dart
		"Dockerfile",            // docker - could be problematic
		"elm.json",              // elm
		".flowconfig",           // flow
		"fortls",                // fortran
		"project.godot",         // godot
		"stack.yaml",            // ghci, hie
		"hie-bios",              // ghci
		"BUILD.bazel",           // ghci
		"cabal.config",          // ghci
		"package.yaml",          // ghci, hie
		"go.mod",                // go
		"package.json",          // html, cssls, ocaml, typescript
		".envrc",                // nix flake
		"flake.nix",             // nix flake
		"composer.json",         // intelephense (php)
		"build.sbt",             // metals
		"build.sc",              // metals
		"build.gradle",          // metals
		"pom.xml",               // metals
		".merlin",               // ocaml
		"Cargo.toml",            // rust
		"Gemfile",               // ruby, solargraph
		"vue.config.js",         // vue
		"pyrightconfig.json",    // python
		"Makefile",              // misc
		".idea",                 // Editor
		".vscode",               // Editor
		".ensime_cache",         // scala
		".eunit",                // erlang
		".hg",                   // vc (mercurial)
		"_FOSSIL_",              // vc (fossil)
		".fslckout",             // vc (fossil)
		".bzr",                  // vc (bazaar)
		"arcs",                  // idk (was in projectile)
		".pijul",                // vc (pijul)
		".tox",                  // python
		".svn",                  // vc (svn)
		".stack-work",           // testing app
		".cache",                // idk
		".sl",                   // c#?
		".jj",                   // java
		"GTAGS",                 // c
		"TAGS",                  // c
		"configure.ac",          // c/c++
		"configure.in",          // c/c++
		"cscope.out",            // c
	}
}
