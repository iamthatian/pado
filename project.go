// merge this with projects? confusing filename?
package main

import (
	"os"
	"path/filepath"
)

type Project struct {
	Name        string
	Path        string
	Kind        string
	Description string
}

// And then, ruby path python git path etcetc (based on matchers)
type ProjectActions interface {
	Name()
	Path()
	Kind()
	Description()
	// Matchers() []string
	// Update()
	// Rename()
	// ChangeRename()
	Compile()
}

func getProjectFiles() []string {
	projectFiles := []string{
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
	return projectFiles
}

// TODO: return error
func searchAncestors(startPath string, matchList []string) string {
	cwd, err := NormalizePath(startPath)
	if err != nil {
		return ""
	}
	// TODO: Should loop for a fixed range just in case?
	for cwd != "/" {
		files, err := os.ReadDir(cwd)
		if err != nil {
			return ""
		}

		if hasProjectFiles(files, matchList) != "" {
			return filepath.Clean(cwd)
		}

		cwd = filepath.Dir(cwd)
	}

	return "/"
}

func hasProjectFiles(fileList []os.DirEntry, patterns []string) string {
	for _, pattern := range patterns {
		for _, file := range fileList {
			if file.Name() == pattern {
				return filepath.Clean(file.Name())
			}
		}
	}
	return ""
}

// TODO: I should just get the kind on matching type
// Given a matcher item, return the
// func parseLanguage(matcher string) string {
// }
