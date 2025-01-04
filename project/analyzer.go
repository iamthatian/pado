// TODO
package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/duckonomy/parkour/utils"
)

type ProjectAnalyzer interface {
	AnalyzeProject(path string) error
}

// Template for Project
// as hash? could be hash value?
type ProjectType struct {
	Name           string
	ProjectFile    string
	TestCommand    string
	BuildCommand   string
	InstallCommand string
	RunCommand     string
	SrcDir         string
	TestDir        string
	TestPrefix     string
	TestSuffix     string
}

// NOTE: Project must already have path and type
// TODO: Decide whether I should just add stuff to Project or embed separate ProjectType inside Project
// NOTE: only triggered when adding project "properly"
func (p *Project) AnalyzeProject() error {
	if p.Path == "" {
		return fmt.Errorf("project path is empty")
	}

	// Detect project type
	projType, err := detectProjectType(p.Path)
	if err != nil {
		return fmt.Errorf("failed to detect project type: %w", err)
	}

	// Update project with detected type information
	p.SimpleType = projType.Name
	p.BuildCommand = []string{projType.BuildCommand}
	p.TestCommand = []string{projType.TestCommand}
	p.RunCommand = []string{projType.RunCommand}

	// Walk through project directory to find children projects
	err = filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root path
		if path == p.Path {
			return nil
		}

		// Check if directory contains a project file
		if info.IsDir() {
			if childType, err := detectProjectType(path); err == nil {
				childProject := &Project{
					Path:       path,
					SimpleType: childType.Name,
					Name:       filepath.Base(path),
					Parent:     p,
					Dirty:      false,
				}
				p.Children = append(p.Children, childProject)
				// Skip walking into this directory since we've identified it as a child project
				return filepath.SkipDir
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk project directory: %w", err)
	}

	return nil
}

func detectProjectType(path string) (*ProjectType, error) {
	for _, projType := range PROJECT_TYPES() {
		switch projType.Name {
		// Node.js Package Managers
		case "yarn":
			pkgJSON := filepath.Join(path, "package.json")
			yarnLock := filepath.Join(path, "yarn.lock")
			yarnRc := filepath.Join(path, ".yarnrc.yml")
			if utils.FileExists(pkgJSON) && (utils.FileExists(yarnLock) || utils.FileExists(yarnRc)) {
				return &projType, nil
			}

		case "pnpm":
			pkgJSON := filepath.Join(path, "package.json")
			pnpmLock := filepath.Join(path, "pnpm-lock.yaml")
			if utils.FileExists(pkgJSON) && utils.FileExists(pnpmLock) {
				return &projType, nil
			}

		case "npm":
			pkgJSON := filepath.Join(path, "package.json")
			packageLock := filepath.Join(path, "package-lock.json")
			if utils.FileExists(pkgJSON) && utils.FileExists(packageLock) {
				return &projType, nil
			}

		// JavaScript/TypeScript Frameworks
		case "angular":
			angularJSON := filepath.Join(path, "angular.json")
			angularCLI := filepath.Join(path, ".angular-cli.json")
			pkgJSON := filepath.Join(path, "package.json")
			if utils.FileExists(angularJSON) || utils.FileExists(angularCLI) {
				return &projType, nil
			}
			if utils.FileExists(pkgJSON) && containsNodePackage(pkgJSON, "@angular/core") {
				return &projType, nil
			}

		case "vue":
			vueConfig := filepath.Join(path, "vue.config.js")
			pkgJSON := filepath.Join(path, "package.json")
			if utils.FileExists(vueConfig) || (utils.FileExists(pkgJSON) && containsNodePackage(pkgJSON, "vue")) {
				return &projType, nil
			}

		case "svelte":
			svelteConfig := filepath.Join(path, "svelte.config.js")
			pkgJSON := filepath.Join(path, "package.json")
			if utils.FileExists(svelteConfig) || (utils.FileExists(pkgJSON) && containsNodePackage(pkgJSON, "svelte")) {
				return &projType, nil
			}

		case "nextjs":
			nextConfig := filepath.Join(path, "next.config.js")
			pkgJSON := filepath.Join(path, "package.json")
			if utils.FileExists(nextConfig) || (utils.FileExists(pkgJSON) && containsNodePackage(pkgJSON, "next")) {
				return &projType, nil
			}

		case "astro":
			astroConfig := filepath.Join(path, "astro.config.mjs")
			pkgJSON := filepath.Join(path, "package.json")
			if utils.FileExists(astroConfig) || (utils.FileExists(pkgJSON) && containsNodePackage(pkgJSON, "astro")) {
				return &projType, nil
			}

		case "remix":
			remixConfig := filepath.Join(path, "remix.config.js")
			pkgJSON := filepath.Join(path, "package.json")
			if utils.FileExists(remixConfig) || (utils.FileExists(pkgJSON) && containsNodePackage(pkgJSON, "@remix-run/react")) {
				return &projType, nil
			}

		// Build Tools
		case "nx":
			nxJson := filepath.Join(path, "nx.json")
			if utils.FileExists(nxJson) {
				return &projType, nil
			}

		case "turborepo":
			turboJson := filepath.Join(path, "turbo.json")
			if utils.FileExists(turboJson) {
				return &projType, nil
			}

		// PHP Frameworks
		case "php-symfony":
			composerJSON := filepath.Join(path, "composer.json")
			appDir := filepath.Join(path, "app")
			srcDir := filepath.Join(path, "src")
			vendorDir := filepath.Join(path, "vendor")
			configDir := filepath.Join(path, "config")
			if utils.FileExists(composerJSON) {
				content, err := os.ReadFile(composerJSON)
				if err == nil && strings.Contains(string(content), "symfony/symfony") {
					if (utils.IsDirectory(appDir) && utils.IsDirectory(srcDir) && utils.IsDirectory(vendorDir)) ||
						(utils.IsDirectory(srcDir) && utils.IsDirectory(vendorDir) && utils.IsDirectory(configDir)) {
						return &projType, nil
					}
				}
			}

		// Python Frameworks
		case "django":
			managePy := filepath.Join(path, "manage.py")
			if utils.FileExists(managePy) {
				content, err := os.ReadFile(managePy)
				if err == nil && strings.Contains(string(content), "django") {
					return &projType, nil
				}
			}

		case "fastapi":
			paths := []string{
				filepath.Join(path, "main.py"),
				filepath.Join(path, "requirements.txt"),
				filepath.Join(path, "pyproject.toml"),
			}
			for _, p := range paths {
				if utils.FileExists(p) {
					content, err := os.ReadFile(p)
					if err == nil && strings.Contains(string(content), "fastapi") {
						return &projType, nil
					}
				}
			}

		// Ruby Frameworks
		case "rails-test":
			if isRailsProject(path) && utils.IsDirectory(filepath.Join(path, "test")) {
				return &projType, nil
			}

		case "rails-rspec":
			if isRailsProject(path) && utils.IsDirectory(filepath.Join(path, "spec")) {
				gemfile := filepath.Join(path, "Gemfile")
				if utils.FileExists(gemfile) {
					content, err := os.ReadFile(gemfile)
					if err == nil && strings.Contains(string(content), "rspec-rails") {
						return &projType, nil
					}
				}
			}

		// Java/Kotlin Frameworks
		case "spring-boot":
			pomXml := filepath.Join(path, "pom.xml")
			gradleFile := filepath.Join(path, "build.gradle")
			if utils.FileExists(pomXml) {
				content, err := os.ReadFile(pomXml)
				if err == nil && strings.Contains(string(content), "spring-boot") {
					return &projType, nil
				}
			}
			if utils.FileExists(gradleFile) {
				content, err := os.ReadFile(gradleFile)
				if err == nil && strings.Contains(string(content), "spring-boot") {
					return &projType, nil
				}
			}

		case "kotlin-gradle":
			buildGradleKts := filepath.Join(path, "build.gradle.kts")
			if utils.FileExists(buildGradleKts) {
				content, err := os.ReadFile(buildGradleKts)
				if err == nil && strings.Contains(string(content), "kotlin") {
					return &projType, nil
				}
			}

		// Mobile Development
		case "flutter":
			pubspec := filepath.Join(path, "pubspec.yaml")
			if utils.FileExists(pubspec) {
				content, err := os.ReadFile(pubspec)
				if err == nil && strings.Contains(string(content), "flutter:") {
					return &projType, nil
				}
			}

		// Rust Frameworks
		case "actix", "axum":
			cargoToml := filepath.Join(path, "Cargo.toml")
			if utils.FileExists(cargoToml) {
				content, err := os.ReadFile(cargoToml)
				if err == nil {
					framework := projType.Name
					if strings.Contains(string(content), framework) {
						return &projType, nil
					}
				}
			}

		// .NET
		case "dotnet", "dotnet-sln":
			patterns := []string{"*.csproj", "*.fsproj", "*.sln"}
			for _, pattern := range patterns {
				matches, err := filepath.Glob(filepath.Join(path, pattern))
				if err == nil && len(matches) > 0 {
					return &projType, nil
				}
			}

		// Build Systems
		case "bazel":
			workspace := filepath.Join(path, "WORKSPACE")
			if utils.FileExists(workspace) {
				return &projType, nil
			}

		case "cmake":
			cmakelists := filepath.Join(path, "CMakeLists.txt")
			if utils.FileExists(cmakelists) {
				return &projType, nil
			}

		case "meson":
			mesonBuild := filepath.Join(path, "meson.build")
			if utils.FileExists(mesonBuild) {
				return &projType, nil
			}

		// Swift
		case "swift-package":
			packageSwift := filepath.Join(path, "Package.swift")
			if utils.FileExists(packageSwift) {
				return &projType, nil
			}

		// Functional Languages
		case "haskell-stack":
			stackYaml := filepath.Join(path, "stack.yaml")
			if utils.FileExists(stackYaml) {
				return &projType, nil
			}

		case "ocaml-dune":
			duneProject := filepath.Join(path, "dune-project")
			if utils.FileExists(duneProject) {
				return &projType, nil
			}

		// Scala
		case "sbt":
			buildSbt := filepath.Join(path, "build.sbt")
			if utils.FileExists(buildSbt) {
				return &projType, nil
			}

		case "mill":
			buildSc := filepath.Join(path, "build.sc")
			if utils.FileExists(buildSc) {
				return &projType, nil
			}

		// Clojure
		case "lein-test", "lein-midje":
			projectClj := filepath.Join(path, "project.clj")
			if utils.FileExists(projectClj) {
				if projType.Name == "lein-midje" {
					midjeClj := filepath.Join(path, ".midje.clj")
					if utils.FileExists(midjeClj) {
						return &projType, nil
					}
				} else {
					return &projType, nil
				}
			}

		default:
			// For project types without special cases, check the project file directly
			projectFile := projType.ProjectFile

			// Handle glob patterns if present
			if strings.Contains(projectFile, "*") {
				matches, err := filepath.Glob(filepath.Join(path, projectFile))
				if err == nil && len(matches) > 0 {
					return &projType, nil
				}
			} else {
				// Direct file check
				if utils.FileExists(filepath.Join(path, projectFile)) {
					return &projType, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("unable to detect project type")
}

// PackageJSON represents the structure of a package.json file
type PackageJSON struct {
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
}

// isRailsProject checks if the given path contains a Rails project
func isRailsProject(path string) bool {
	// Check for required Rails directories
	requiredDirs := []string{
		"app",
		"config",
		"db",
	}

	for _, dir := range requiredDirs {
		if !utils.IsDirectory(filepath.Join(path, dir)) {
			return false
		}
	}

	// Check for typical Rails subdirectories in app/
	appSubdirs := []string{
		"controllers",
		"models",
		"views",
	}

	appPath := filepath.Join(path, "app")
	for _, subdir := range appSubdirs {
		if !utils.IsDirectory(filepath.Join(appPath, subdir)) {
			return false
		}
	}

	// Check for Gemfile with Rails
	gemfilePath := filepath.Join(path, "Gemfile")
	if !utils.FileExists(gemfilePath) {
		return false
	}

	content, err := os.ReadFile(gemfilePath)
	if err != nil {
		return false
	}

	gemfileContent := string(content)

	// Check for Rails gem in Gemfile
	// This handles various ways Rails might be specified:
	// gem 'rails'
	// gem 'rails', '6.1.0'
	// gem 'rails', '~> 6.1'
	// gem "rails"
	// etc.
	railsPatterns := []string{
		`gem 'rails'`,
		`gem "rails"`,
	}

	for _, pattern := range railsPatterns {
		if strings.Contains(gemfileContent, pattern) {
			return true
		}
	}

	// Additional checks for Rails-specific files
	railsFiles := []string{
		filepath.Join(path, "config", "routes.rb"),
		filepath.Join(path, "config", "application.rb"),
		filepath.Join(path, "config", "environment.rb"),
		filepath.Join(path, "Rakefile"),
	}

	// Count how many Rails-specific files exist
	railsFileCount := 0
	for _, file := range railsFiles {
		if utils.FileExists(file) {
			railsFileCount++
		}
	}

	// If we find at least 3 of the Rails-specific files, consider it a Rails project
	return railsFileCount >= 3
}

// isRailsAPI checks if it's a Rails API project
func isRailsAPI(path string) bool {
	if !isRailsProject(path) {
		return false
	}

	// Check application controller inheritance
	appControllerPath := filepath.Join(path, "app", "controllers", "application_controller.rb")
	if !utils.FileExists(appControllerPath) {
		return false
	}

	content, err := os.ReadFile(appControllerPath)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), "ActionController::API")
}

// getRailsVersion attempts to determine the Rails version
func getRailsVersion(path string) (string, error) {
	gemfileLockPath := filepath.Join(path, "Gemfile.lock")
	if !utils.FileExists(gemfileLockPath) {
		return "", fmt.Errorf("file Gemfile.lock not found")
	}

	content, err := os.ReadFile(gemfileLockPath)
	if err != nil {
		return "", err
	}

	// Look for rails version in Gemfile.lock
	// Example: rails (6.1.4.1)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "rails (") {
			version := strings.TrimPrefix(strings.TrimSpace(line), "rails (")
			version = strings.TrimSuffix(version, ")")
			return version, nil
		}
	}

	return "", fmt.Errorf("rails version not found in Gemfile.lock")
}

// containsPackage checks if a package.json file contains a specific package
// in any of its dependency sections (dependencies, devDependencies, or peerDependencies)
func containsNodePackage(packageJsonPath string, pkg string) bool {
	content, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return false
	}

	var packageJSON PackageJSON
	if err := json.Unmarshal(content, &packageJSON); err != nil {
		// Fallback to simple string matching if JSON parsing fails
		return strings.Contains(string(content), fmt.Sprintf(`"%s"`, pkg))
	}

	// Check in dependencies
	if _, ok := packageJSON.Dependencies[pkg]; ok {
		return true
	}

	// Check in devDependencies
	if _, ok := packageJSON.DevDependencies[pkg]; ok {
		return true
	}

	// Check in peerDependencies
	if _, ok := packageJSON.PeerDependencies[pkg]; ok {
		return true
	}

	// Check for scoped packages (e.g., @angular/core)
	// This handles cases where the package might be part of a scope
	if strings.Contains(pkg, "/") {
		scope := strings.Split(pkg, "/")[0]
		for dep := range packageJSON.Dependencies {
			if strings.HasPrefix(dep, scope) {
				return true
			}
		}
		for dep := range packageJSON.DevDependencies {
			if strings.HasPrefix(dep, scope) {
				return true
			}
		}
		for dep := range packageJSON.PeerDependencies {
			if strings.HasPrefix(dep, scope) {
				return true
			}
		}
	}

	return false
}

func PROJECT_TYPES() []ProjectType {
	return []ProjectType{
		// Kotlin
		{
			Name:         "kotlin-gradle",
			ProjectFile:  "build.gradle.kts",
			BuildCommand: "gradle build",
			TestCommand:  "gradle test",
			RunCommand:   "gradle run",
			TestSuffix:   "Test",
		},
		// Flutter/Dart
		{
			Name:         "flutter",
			ProjectFile:  "pubspec.yaml",
			BuildCommand: "flutter build",
			TestCommand:  "flutter test",
			RunCommand:   "flutter run",
			TestSuffix:   "_test.dart",
		},
		// Swift
		{
			Name:         "swift-package",
			ProjectFile:  "Package.swift",
			BuildCommand: "swift build",
			TestCommand:  "swift test",
			RunCommand:   "swift run",
			TestSuffix:   "Tests",
		},
		// Deno
		{
			Name:         "deno",
			ProjectFile:  "deno.json",
			BuildCommand: "deno compile",
			TestCommand:  "deno test",
			RunCommand:   "deno run",
			TestSuffix:   "_test.ts",
		},
		// Bun
		{
			Name:         "bun",
			ProjectFile:  "bun.lockb",
			BuildCommand: "bun build",
			TestCommand:  "bun test",
			RunCommand:   "bun run",
			TestSuffix:   ".test.ts",
		},
		// Vue.js
		{
			Name:         "vue",
			ProjectFile:  "vue.config.js",
			BuildCommand: "vue-cli-service build",
			TestCommand:  "vue-cli-service test:unit",
			RunCommand:   "vue-cli-service serve",
			TestSuffix:   ".spec.js",
		},
		// Svelte
		{
			Name:         "svelte",
			ProjectFile:  "svelte.config.js",
			BuildCommand: "vite build",
			TestCommand:  "vitest",
			RunCommand:   "vite dev",
			TestSuffix:   ".test.js",
		},
		// Next.js
		{
			Name:         "nextjs",
			ProjectFile:  "next.config.js",
			BuildCommand: "next build",
			TestCommand:  "jest",
			RunCommand:   "next dev",
			TestSuffix:   ".test.js",
		},
		// Astro
		{
			Name:         "astro",
			ProjectFile:  "astro.config.mjs",
			BuildCommand: "astro build",
			TestCommand:  "astro check",
			RunCommand:   "astro dev",
		},
		// Remix
		{
			Name:         "remix",
			ProjectFile:  "remix.config.js",
			BuildCommand: "remix build",
			TestCommand:  "remix test",
			RunCommand:   "remix dev",
			TestSuffix:   ".test.ts",
		},
		// SolidJS
		{
			Name:         "solid",
			ProjectFile:  "vite.config.ts",
			BuildCommand: "vite build",
			TestCommand:  "vitest",
			RunCommand:   "vite serve",
			TestSuffix:   ".test.tsx",
		},
		// Qwik
		{
			Name:         "qwik",
			ProjectFile:  "qwik.config.ts",
			BuildCommand: "qwik build",
			TestCommand:  "vitest",
			RunCommand:   "qwik dev",
			TestSuffix:   ".test.tsx",
		},
		// Nx Monorepo
		{
			Name:         "nx",
			ProjectFile:  "nx.json",
			BuildCommand: "nx build",
			TestCommand:  "nx test",
			RunCommand:   "nx serve",
			TestSuffix:   ".spec.ts",
		},
		// Turborepo
		{
			Name:         "turborepo",
			ProjectFile:  "turbo.json",
			BuildCommand: "turbo build",
			TestCommand:  "turbo test",
			RunCommand:   "turbo dev",
		},
		// Vite
		{
			Name:         "vite",
			ProjectFile:  "vite.config.ts",
			BuildCommand: "vite build",
			TestCommand:  "vitest",
			RunCommand:   "vite",
			TestSuffix:   ".test.ts",
		},
		// Phoenix (Elixir)
		{
			Name:         "phoenix",
			ProjectFile:  "mix.exs",
			BuildCommand: "mix phx.server",
			TestCommand:  "mix test",
			RunCommand:   "iex -S mix phx.server",
			TestSuffix:   "_test.exs",
		},
		// FastAPI
		{
			Name:         "fastapi",
			ProjectFile:  "main.py",
			BuildCommand: "pip install -e .",
			TestCommand:  "pytest",
			RunCommand:   "uvicorn main:app --reload",
			TestPrefix:   "test_",
		},
		// Tauri
		{
			Name:         "tauri",
			ProjectFile:  "tauri.conf.json",
			BuildCommand: "tauri build",
			TestCommand:  "cargo test",
			RunCommand:   "tauri dev",
		},
		// Actix (Rust)
		{
			Name:         "actix",
			ProjectFile:  "Cargo.toml",
			BuildCommand: "cargo build",
			TestCommand:  "cargo test",
			RunCommand:   "cargo run",
			TestSuffix:   "_test",
		},
		// Axum (Rust)
		{
			Name:         "axum",
			ProjectFile:  "Cargo.toml",
			BuildCommand: "cargo build",
			TestCommand:  "cargo test",
			RunCommand:   "cargo run",
			TestSuffix:   "_test",
		},
		// Spring Boot
		{
			Name:         "spring-boot",
			ProjectFile:  "pom.xml",
			BuildCommand: "mvn spring-boot:run",
			TestCommand:  "mvn test",
			RunCommand:   "java -jar target/*.jar",
			TestSuffix:   "Test",
		},
		// Quarkus
		{
			Name:         "quarkus",
			ProjectFile:  "pom.xml",
			BuildCommand: "mvn quarkus:dev",
			TestCommand:  "mvn test",
			RunCommand:   "java -jar target/*-runner.jar",
			TestSuffix:   "Test",
		},
		// Micronaut
		{
			Name:         "micronaut",
			ProjectFile:  "build.gradle",
			BuildCommand: "gradle build",
			TestCommand:  "gradle test",
			RunCommand:   "gradle run",
			TestSuffix:   "Spec",
		},
		// .NET
		{
			Name:         "dotnet",
			ProjectFile:  "*.csproj", // Will need glob handling
			BuildCommand: "dotnet build",
			RunCommand:   "dotnet run",
			TestCommand:  "dotnet test",
		},
		{
			Name:         "dotnet-sln",
			ProjectFile:  "*.sln", // Will need glob handling
			BuildCommand: "dotnet build",
			RunCommand:   "dotnet run",
			TestCommand:  "dotnet test",
		},
		// Nim
		{
			Name:           "nim-nimble",
			ProjectFile:    "*.nimble", // Will need glob handling
			BuildCommand:   "nimble --noColor build --colors:off",
			InstallCommand: "nimble --noColor install --colors:off",
			TestCommand:    "nimble --noColor test -d:nimUnittestColor:off --colors:off",
			RunCommand:     "nimble --noColor run --colors:off",
			SrcDir:         "src",
			TestDir:        "tests",
		},
		// Universal build systems
		{
			Name:           "xmake",
			ProjectFile:    "xmake.lua",
			BuildCommand:   "xmake build",
			TestCommand:    "xmake test",
			RunCommand:     "xmake run",
			InstallCommand: "xmake install",
		},
		{
			Name:         "scons",
			ProjectFile:  "SConstruct",
			BuildCommand: "scons",
			TestCommand:  "scons test",
			TestSuffix:   "test",
		},
		{
			Name:         "meson",
			ProjectFile:  "meson.build",
			BuildCommand: "ninja",
			TestCommand:  "ninja test",
		},
		{
			Name:         "bazel",
			ProjectFile:  "WORKSPACE",
			BuildCommand: "bazel build",
			TestCommand:  "bazel test",
			RunCommand:   "bazel run",
		},
		// Make & CMake
		{
			Name:           "make",
			ProjectFile:    "Makefile",
			BuildCommand:   "make",
			TestCommand:    "make test",
			InstallCommand: "make install",
		},
		{
			Name:         "cmake",
			ProjectFile:  "CMakeLists.txt",
			BuildCommand: "cmake --build .",
			TestCommand:  "ctest",
		},
		// Go
		{
			Name:         "go",
			ProjectFile:  "go.mod",
			BuildCommand: "go build",
			TestCommand:  "go test ./...",
			TestSuffix:   "_test",
		},
		// Rust
		{
			Name:         "rust-cargo",
			ProjectFile:  "Cargo.toml",
			BuildCommand: "cargo build",
			TestCommand:  "cargo test",
			RunCommand:   "cargo run",
		},
		// Node.js
		{
			Name:         "npm",
			ProjectFile:  "package.json",
			BuildCommand: "npm install && npm run build",
			TestCommand:  "npm test",
			TestSuffix:   ".test",
		},
		{
			Name:         "yarn",
			ProjectFile:  "package.json", // With yarn.lock check in analyzer
			BuildCommand: "yarn && yarn build",
			TestCommand:  "yarn test",
			TestSuffix:   ".test",
		},
		// Python
		{
			Name:         "python-pip",
			ProjectFile:  "requirements.txt",
			BuildCommand: "python setup.py build",
			TestCommand:  "python -m unittest discover",
			TestPrefix:   "test_",
			TestSuffix:   "_test",
		},
		{
			Name:         "python-poetry",
			ProjectFile:  "poetry.lock",
			BuildCommand: "poetry build",
			TestCommand:  "poetry run python -m unittest discover",
			TestPrefix:   "test_",
			TestSuffix:   "_test",
		},
		// Java
		{
			Name:         "maven",
			ProjectFile:  "pom.xml",
			BuildCommand: "mvn -B clean install",
			TestCommand:  "mvn -B test",
			TestSuffix:   "Test",
			SrcDir:       "src/main/",
			TestDir:      "src/test/",
		},
		{
			Name:         "gradle",
			ProjectFile:  "build.gradle",
			BuildCommand: "gradle build",
			TestCommand:  "gradle test",
			TestSuffix:   "Spec",
		},
		// Haskell
		{
			Name:         "haskell-stack",
			ProjectFile:  "stack.yaml",
			BuildCommand: "stack build",
			TestCommand:  "stack build --test",
			TestSuffix:   "Spec",
		},
		// Nix
		{
			Name:         "nix",
			ProjectFile:  "default.nix",
			BuildCommand: "nix-build",
			TestCommand:  "nix-build",
		},
		{
			Name:         "nix-flake",
			ProjectFile:  "flake.nix",
			BuildCommand: "nix build",
			TestCommand:  "nix flake check",
			RunCommand:   "nix run",
		},
		// Debian
		{
			Name:         "debian",
			ProjectFile:  "debian/control",
			BuildCommand: "debuild -uc -us",
		},
		// PHP
		{
			Name:         "php-symfony",
			ProjectFile:  "composer.json",
			BuildCommand: "app/console server:run",
			TestCommand:  "phpunit -c app",
			TestSuffix:   "Test",
		},
		// Erlang & Elixir
		{
			Name:         "rebar",
			ProjectFile:  "rebar.config",
			BuildCommand: "rebar3 compile",
			TestCommand:  "rebar3 do eunit,ct",
			TestSuffix:   "_SUITE",
		},
		{
			Name:         "elixir",
			ProjectFile:  "mix.exs",
			BuildCommand: "mix compile",
			TestCommand:  "mix test",
			TestSuffix:   "_test",
			SrcDir:       "lib/",
		},
		// JavaScript
		{
			Name:         "grunt",
			ProjectFile:  "Gruntfile.js",
			BuildCommand: "grunt",
			TestCommand:  "grunt test",
		},
		{
			Name:         "gulp",
			ProjectFile:  "gulpfile.js",
			BuildCommand: "gulp",
			TestCommand:  "gulp test",
		},
		{
			Name:         "pnpm",
			ProjectFile:  "package.json", // With pnpm-lock.yaml check
			BuildCommand: "pnpm install && pnpm build",
			TestCommand:  "pnpm test",
			TestSuffix:   ".test",
		},
		// Angular
		{
			Name:         "angular",
			ProjectFile:  "angular.json", // Also check .angular-cli.json
			BuildCommand: "ng build",
			RunCommand:   "ng serve",
			TestCommand:  "ng test",
			TestSuffix:   ".spec",
		},
		// Python additional frameworks
		{
			Name:         "django",
			ProjectFile:  "manage.py",
			BuildCommand: "python manage.py runserver",
			TestCommand:  "python manage.py test",
			TestPrefix:   "test_",
			TestSuffix:   "_test",
		},
		{
			Name:         "python-tox",
			ProjectFile:  "tox.ini",
			BuildCommand: "tox -r --notest",
			TestCommand:  "tox",
			TestPrefix:   "test_",
			TestSuffix:   "_test",
		},
		{
			Name:         "python-pipenv",
			ProjectFile:  "Pipfile",
			BuildCommand: "pipenv run build",
			TestCommand:  "pipenv run test",
			TestPrefix:   "test_",
			TestSuffix:   "_test",
		},
		// Java additional frameworks
		{
			Name:         "grails",
			ProjectFile:  "application.yml", // Also check grails-app dir
			BuildCommand: "grails package",
			TestCommand:  "grails test-app",
			TestSuffix:   "Spec",
		},
		// Scala
		{
			Name:         "sbt",
			ProjectFile:  "build.sbt",
			BuildCommand: "sbt compile",
			TestCommand:  "sbt test",
			TestSuffix:   "Spec",
			SrcDir:       "main",
			TestDir:      "test",
		},
		{
			Name:         "mill",
			ProjectFile:  "build.sc",
			BuildCommand: "mill __.compile",
			TestCommand:  "mill __.test",
			TestSuffix:   "Test",
			SrcDir:       "src/",
			TestDir:      "test/src/",
		},
		{
			Name:         "bloop",
			ProjectFile:  ".bloop",
			BuildCommand: "bloop compile root",
			TestCommand:  "bloop test --propagate --reporter scalac root",
			TestSuffix:   "Spec",
			SrcDir:       "src/main/",
			TestDir:      "src/test/",
		},
		// Clojure
		{
			Name:         "lein-test",
			ProjectFile:  "project.clj",
			BuildCommand: "lein compile",
			TestCommand:  "lein test",
			TestSuffix:   "_test",
		},
		{
			Name:         "lein-midje",
			ProjectFile:  "project.clj", // Also check .midje.clj
			BuildCommand: "lein compile",
			TestCommand:  "lein midje",
			TestPrefix:   "t_",
		},
		{
			Name:         "boot-clj",
			ProjectFile:  "build.boot",
			BuildCommand: "boot aot",
			TestCommand:  "boot test",
			TestSuffix:   "_test",
		},
		{
			Name:        "clojure-cli",
			ProjectFile: "deps.edn",
			TestSuffix:  "_test",
		},
		// Ruby
		{
			Name:         "ruby-rspec",
			ProjectFile:  "Gemfile", // Also check lib and spec dirs
			BuildCommand: "bundle exec rake",
			TestCommand:  "bundle exec rspec",
			TestSuffix:   "_spec",
			SrcDir:       "lib/",
			TestDir:      "spec/",
		},
		{
			Name:         "ruby-test",
			ProjectFile:  "Gemfile", // Also check lib and test dirs
			BuildCommand: "bundle exec rake",
			TestCommand:  "bundle exec rake test",
			SrcDir:       "lib/",
			TestSuffix:   "_test",
		},
		{
			Name:         "rails-test",
			ProjectFile:  "Gemfile", // Also check app, lib, db, config, test dirs
			BuildCommand: "bundle exec rails server",
			TestCommand:  "bundle exec rake test",
			SrcDir:       "app/",
			TestSuffix:   "_test",
		},
		{
			Name:         "rails-rspec",
			ProjectFile:  "Gemfile", // Also check app, lib, db, config, spec dirs
			BuildCommand: "bundle exec rails server",
			TestCommand:  "bundle exec rspec",
			SrcDir:       "app/",
			TestDir:      "spec/",
			TestSuffix:   "_spec",
		},
		// Crystal
		{
			Name:        "crystal-spec",
			ProjectFile: "shard.yml",
			TestCommand: "crystal spec",
			TestSuffix:  "_spec",
			SrcDir:      "src/",
			TestDir:     "spec/",
		},
		// R
		{
			Name:         "r",
			ProjectFile:  "DESCRIPTION",
			BuildCommand: "R CMD INSTALL --with-keep.source .",
			TestCommand:  "R CMD check -o /tmp .", // Note: Temp dir handling needed
		},
		// OCaml
		{
			Name:         "ocaml-dune",
			ProjectFile:  "dune-project",
			BuildCommand: "dune build",
			TestCommand:  "dune runtest",
		},
		// Dart
		{
			Name:         "dart",
			ProjectFile:  "pubspec.yaml",
			BuildCommand: "pub get",
			TestCommand:  "pub run test",
			RunCommand:   "dart",
			TestSuffix:   "_test.dart",
		},
		// Elm
		{
			Name:         "elm",
			ProjectFile:  "elm.json",
			BuildCommand: "elm make",
		},
		// Julia
		{
			Name:         "julia",
			ProjectFile:  "Project.toml",
			BuildCommand: "julia --project=@. -e 'import Pkg; Pkg.precompile(); Pkg.build()'",
			TestCommand:  "julia --project=@. -e 'import Pkg; Pkg.test()' --check-bounds=yes",
			SrcDir:       "src",
			TestDir:      "test",
		},
		// Zig
		{
			Name:         "zig",
			ProjectFile:  "build.zig.zon",
			BuildCommand: "zig build",
			TestCommand:  "zig build test",
			RunCommand:   "zig build run",
		},
	}
}

// TODO
// // ;; File-based detection project types
// // ;; Make & CMake
// // (projectile-register-project-type 'gnumake '("GNUMakefile")
// //                                   :project-file "GNUMakefile"
// //                                   :compile "make"
// //                                   :test "make test"
// //                                   :install "make install")
// // ;; go-task/task
// // (projectile-register-project-type 'go-task '("Taskfile.yml")
// //                                   :project-file "Taskfile.yml"
// //                                   :compile "task build"
// //                                   :test "task test"
// //                                   :install "task install")
// // ;; Go should take higher precedence than Make because Go projects often have a Makefile.
// // (projectile-register-project-type 'go projectile-go-project-test-function
// //                                   :compile "go build"
// //                                   :test "go test ./..."
// //                                   :test-suffix "_test")

// // Need laravel!!

// // (projectile-register-project-type 'python-pkg '("setup.py")
// //                                   :project-file "setup.py"
// //                                   :compile "python setup.py build"
// //                                   :test "python -m unittest discover"
// //                                   :test-prefix "test_"
// //                                   :test-suffix"_test")

// // (projectile-register-project-type 'python-toml '("pyproject.toml")
// //                                   :project-file "pyproject.toml"
// //                                   :compile "python -m build"
// //                                   :test "python -m unittest discover"
// //                                   :test-prefix "test_"
// //                                   :test-suffix "_test")

// // (projectile-register-project-type 'gradlew '("gradlew")
// //                                   :project-file "gradlew"
// //                                   :compile "./gradlew build"
// //                                   :test "./gradlew test"
// //                                   :test-suffix "Spec")
// // ;; Racket

// // (projectile-register-project-type 'racket '("info.rkt")
// //                                   :project-file "info.rkt"
// //                                   :test "raco test ."
// //                                   :install "raco pkg install"
// //                                   :package "raco pkg create --source $(pwd)")
// //
