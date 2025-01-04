package project

// TODO: Split this and ProjectType per file?
//
//	Get values from config
type ProjectRootIndicator struct {
	Type     string
	Priority int
	Patterns []string
}

func VC_INDICATORS() []ProjectRootIndicator {
	return []ProjectRootIndicator{
		{
			Type:     "git",
			Priority: 300,
			Patterns: []string{
				".git",
				".gitignore",
				".gitmodules",
				".gitattributes",
				".gitkeep",
			},
		},
		// Mercurial
		{
			Type:     "mercurial",
			Priority: 300,
			Patterns: []string{
				".hg",
				".hgignore",
				".hgtags",
				".hgeol",
				".hgsub",
				".hgsubstate",
			},
		},
		// Subversion
		{
			Type:     "svn",
			Priority: 300,
			Patterns: []string{
				".svn",
				".svnignore",
				"_svn",
			},
		},
		// Bazaar
		{
			Type:     "bzr",
			Priority: 300,
			Patterns: []string{
				".bzr",
				".bzrignore",
				".bzrtags",
			},
		},
		// CVS
		{
			Type:     "cvs",
			Priority: 300,
			Patterns: []string{
				"CVS",
				".cvsignore",
				".cvsrc",
			},
		},
		// Darcs
		{
			Type:     "darcs",
			Priority: 300,
			Patterns: []string{
				"_darcs",
				".darcs-ignore",
			},
		},
		// Fossil
		{
			Type:     "fossil",
			Priority: 300,
			Patterns: []string{
				".fossil",
				".fossil-settings",
			},
		},
		// Monotone
		{
			Type:     "monotone",
			Priority: 300,
			Patterns: []string{
				"_MTN",
				".mtn-ignore",
			},
		},
		// Perforce
		{
			Type:     "perforce",
			Priority: 300,
			Patterns: []string{
				"p4config",
				".p4ignore",
				"p4env.txt",
			},
		},
		// Pijul
		{
			Type:     "pijul",
			Priority: 300,
			Patterns: []string{
				".pijul",
			},
		},
		// RCS
		{
			Type:     "rcs",
			Priority: 300,
			Patterns: []string{
				"RCS",
				",v",
			},
		},
	}
}

func WORKSPACE_INDICATORS() []ProjectRootIndicator {
	return []ProjectRootIndicator{
		{
			Type:     "pnpm-workspace",
			Priority: 250,
			Patterns: []string{"pnpm-workspace.yaml"},
		},
		{
			Type:     "yarn-workspace",
			Priority: 250,
			Patterns: []string{"package.json"},
		},
		{
			Type:     "nx-workspace",
			Priority: 250,
			Patterns: []string{"nx.json", "workspace.json"},
		},
		{
			Type:     "rush",
			Priority: 250,
			Patterns: []string{"rush.json"},
		},
		{
			Type:     "turborepo",
			Priority: 250,
			Patterns: []string{"turbo.json"},
		},
		{
			Type:     "lerna",
			Priority: 250,
			Patterns: []string{"lerna.json"},
		},

		// Rust
		{
			Type:     "cargo-workspace",
			Priority: 250,
			Patterns: []string{"Cargo.toml"},
		},

		// Go
		{
			Type:     "go-workspace",
			Priority: 250,
			Patterns: []string{"go.work"},
		},

		// Python
		{
			Type:     "python-namespace",
			Priority: 250,
			Patterns: []string{"setup.cfg", "pyproject.toml"},
		},

		// Java/Gradle
		{
			Type:     "gradle-multiproject",
			Priority: 250,
			Patterns: []string{"settings.gradle", "settings.gradle.kts"},
		},

		// Maven
		{
			Type:     "maven-multimodule",
			Priority: 250,
			Patterns: []string{"pom.xml"},
		},

		// .NET
		{
			Type:     "dotnet-solution",
			Priority: 250,
			Patterns: []string{"*.sln"},
		},

		// Generic Build Systems
		{
			Type:     "bazel-workspace",
			Priority: 250,
			Patterns: []string{"WORKSPACE", "WORKSPACE.bazel"},
		},
		{
			Type:     "buck-workspace",
			Priority: 250,
			Patterns: []string{".buckconfig"},
		},
		{
			Type:     "pants-workspace",
			Priority: 250,
			Patterns: []string{"pants.toml", "BUILD"},
		},
	}
}

func LANGUAGE_INDICATORS() []ProjectRootIndicator {
	return []ProjectRootIndicator{
		// Node.js/JavaScript
		{
			Type:     "nodejs",
			Priority: 200,
			Patterns: []string{"package.json", "package-lock.json", "yarn.lock", "pnpm-lock.yaml"},
		},
		// React Specific
		{
			Type:     "react",
			Priority: 195,
			Patterns: []string{"next.config.js", "react-scripts", "craco.config.js"},
		},
		// Vue Specific
		{
			Type:     "vue",
			Priority: 195,
			Patterns: []string{"vue.config.js", "nuxt.config.js"},
		},

		// Rust
		{
			Type:     "rust",
			Priority: 200,
			Patterns: []string{"Cargo.toml", "Cargo.lock", "rust-toolchain.toml"},
		},

		// Go
		{
			Type:     "go",
			Priority: 200,
			Patterns: []string{"go.mod", "go.sum"},
		},

		// Python
		{
			Type:     "python",
			Priority: 200,
			Patterns: []string{
				"requirements.txt",
				"setup.py",
				"pyproject.toml",
				"Pipfile",
				"poetry.lock",
				"conda.yaml",
				"environment.yml",
			},
		},

		// Haskell
		{
			Type:     "haskell",
			Priority: 200,
			Patterns: []string{
				"stack.yaml",
				"cabal.config",
				"package.yaml",
				"hie-bios",
			},
		},

		// Java/Kotlin/JVM
		{
			Type:     "jvm",
			Priority: 200,
			Patterns: []string{
				"pom.xml",
				"build.gradle",
				"build.gradle.kts",
				"gradlew",
				"settings.gradle",
				".mvn",
			},
		},

		// Dart/Flutter
		{
			Type:     "dart",
			Priority: 200,
			Patterns: []string{"pubspec.yaml"},
		},

		// Elm
		{
			Type:     "elm",
			Priority: 200,
			Patterns: []string{"elm.json"},
		},

		// Fortran
		{
			Type:     "fortran",
			Priority: 200,
			Patterns: []string{"fortls"},
		},

		// Nix (maybe move this out?)
		{
			Type:     "Nix",
			Priority: 100,
			Patterns: []string{"flake.nix", ".envrc"},
		},

		// Scala
		{
			Type:     "Scala",
			Priority: 200,
			Patterns: []string{"built.sbt", ".ensime_cache"},
		},

		// Godot
		{
			Type:     "godot",
			Priority: 200,
			Patterns: []string{"project.godot"},
		},

		// Ocaml
		{
			Type:     "ocaml",
			Priority: 200,
			Patterns: []string{".merlin"},
		},

		// Erlang
		{
			Type:     "erlang",
			Priority: 200,
			Patterns: []string{".eunit"},
		},

		// Metals
		{
			Type:     "metals",
			Priority: 200,
			Patterns: []string{"metals.sbt", "build.sc"},
		},

		// Swift
		{
			Type:     "swift",
			Priority: 200,
			Patterns: []string{"Package.swift"},
		},

		// Kotlin
		{
			Type:     "kotlin",
			Priority: 200,
			Patterns: []string{"build.gradle.kts", "settings.gradle.kts"},
		},

		// Julia
		{
			Type:     "julia",
			Priority: 200,
			Patterns: []string{"Project.toml", "Manifest.toml"},
		},

		// R
		{
			Type:     "r",
			Priority: 200,
			Patterns: []string{"DESCRIPTION", ".Rproj"},
		},

		// Elixir
		{
			Type:     "elixir",
			Priority: 200,
			Patterns: []string{"mix.exs"},
		},

		// Clojure
		{
			Type:     "clojure",
			Priority: 200,
			Patterns: []string{"project.clj", "deps.edn"},
		},

		// Ruby
		{
			Type:     "ruby",
			Priority: 200,
			Patterns: []string{
				"Gemfile",
				"Gemfile.lock",
				"config.ru",
				".ruby-version",
				"Rakefile",
			},
		},

		// PHP
		{
			Type:     "php",
			Priority: 200,
			Patterns: []string{
				"composer.json",
				"composer.lock",
				"artisan",
				"index.php",
			},
		},

		// C/C++
		{
			Type:     "cpp",
			Priority: 200,
			Patterns: []string{
				"CMakeLists.txt",
				"Makefile",
				"configure.ac",
				"meson.build",
			},
		},

		// .NET
		{
			Type:     "dotnet",
			Priority: 200,
			Patterns: []string{
				"*.csproj",
				"*.fsproj",
				"*.sln",
				"global.json",
				"nuget.config",
			},
		},
	}
}
