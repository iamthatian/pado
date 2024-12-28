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
func (p *Project) FindProject(startPath string) error {
	maxDepth := 100
	depth := 0
	currentPath, err := NormalizePath(startPath)
	if err != nil {
		return err
	}

	completePatterns := getCategorizedProjectFiles()

	for currentPath != "/" && depth <= maxDepth {
		files, err := os.ReadDir(currentPath)
		if err != nil {
			return err
		}

		// This is an unsorted map so might have to keep track of depth
		// TODO: priority of projects (i.e. custom takes priority over project and lang takes priority project over git)
		for category, patterns := range completePatterns {
			if matched := matchProjectFiles(files, patterns); matched {
				p.Path = filepath.Clean(currentPath)
				p.Kind = category
				return nil
			}
		}
		currentPath = filepath.Dir(currentPath)
		depth++
	}

	p.Path = "/"
	p.Kind = ""
	return nil
}

func matchProjectFiles(fileList []os.DirEntry, patterns []string) bool {
	for _, pattern := range patterns {
		for _, file := range fileList {
			if file.Name() == pattern {
				return true
			}
		}
	}
	return false
}

func getCategorizedProjectFiles() map[string][]string {
	return map[string][]string{
		// Version Control
		"git": {
			".git",
			".gitignore",
		},
		"mercurial": {
			".hg",
		},
		"svn": {
			".svn",
		},
		"bazaar": {
			".bzr",
		},
		"fossil": {
			"_FOSSIL_",
			".fslckout",
		},
		"pijul": {
			".pijul",
		},

		// C and C++
		"c": {
			"compile_commands.json",
			"compile_flags.txt",
			"Makefile",
			"configure.ac",
			"configure.in",
			"cscope.out",
			"GTAGS",
			"TAGS",
		},
		"cpp": {
			"compile_commands.json",
			"compile_flags.txt",
			"Makefile",
			".clangd",
			".ccls-cache",
		},

		// Python
		"python": {
			"pyproject.toml",
			"requirements.txt",
			"setup.py",
			"tox.ini",
			".tox",
			"pyrightconfig.json",
		},

		// JavaScript/Node.js
		"nodejs": {
			"package.json",
			"yarn.lock",
			"pnpm-lock.yaml",
			"webpack.config.js",
			"rollup.config.js",
			"vite.config.js",
		},

		// Go
		"go": {
			"go.mod",
			"go.sum",
		},

		// Rust
		"rust": {
			"Cargo.toml",
			"Cargo.lock",
		},

		// Java
		"java": {
			"pom.xml",
			"build.gradle",
			"build.gradle.kts",
			".classpath",
			".project",
		},

		// Haskell
		"haskell": {
			"stack.yaml",
			"cabal.config",
			"package.yaml",
			"hie-bios",
		},

		// Dart/Flutter
		"dart": {
			"pubspec.yaml",
		},

		// Ruby
		"ruby": {
			"Gemfile",
			"Gemfile.lock",
		},

		// PHP
		"php": {
			"composer.json",
			"composer.lock",
		},

		// Docker
		"docker": {
			"Dockerfile",
			"docker-compose.yml",
		},

		// Elm
		"elm": {
			"elm.json",
		},

		// Fortran
		"fortran": {
			"fortls",
		},

		// Nix
		"nix": {
			"flake.nix",
			".envrc",
		},

		// Scala
		"scala": {
			"build.sbt",
			".ensime_cache",
		},

		// Vue
		"vue": {
			"vue.config.js",
		},

		// Godot
		"godot": {
			"project.godot",
		},

		// Editor Configurations
		"editor": {
			".idea",
			".vscode",
		},

		// Miscellaneous
		"make": {
			"Makefile",
		},
		"ocaml": {
			".merlin",
		},
		"erlang": {
			".eunit",
		},
		"metals": {
			"metals.sbt",
			"build.sc",
		},
		"environment": {
			".env",
			".envrc",
		},
		"cache": {
			".cache",
		},
	}
}
