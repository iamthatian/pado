// TODO
package project

import "fmt"

type ProjectAnalyzer interface {
	AnalyzeProject(path string) error
}

// Template for Project
// as hash? could be hash value?
type ProjectType struct {
	Name         string
	ProjectFile  string
	TestCommand  string
	BuildCommand string
	RunCommand   string
}

// NOTE: Project must already have path and type
// TODO: Decide whether I should just add stuff to Project or embed separate ProjectType inside Project
// NOTE: only triggered when adding project "properly"
func (p *Project) AnaylzeProject() error {
	return fmt.Errorf("TODO")
}

// TODO: AI to add functionality
func PROJECT_TYPES() []ProjectType {
	return []ProjectType{
		{
			Name:         "bazel",
			ProjectFile:  "WORKSPACE",
			BuildCommand: "bazel build",
			TestCommand:  "bazel test",
			RunCommand:   "bazel run",
		},
		{
			Name:         "haskell-cabal",
			ProjectFile:  "WORKSPACE",
			BuildCommand: "cabal build",
			TestCommand:  "bazel test",
			RunCommand:   "bazel run",
		},
		{
			Name:         "bazel",
			ProjectFile:  "WORKSPACE",
			BuildCommand: "bazel build",
			TestCommand:  "bazel test",
			RunCommand:   "bazel run",
		},
		{
			Name:         "bazel",
			ProjectFile:  "WORKSPACE",
			BuildCommand: "bazel build",
			TestCommand:  "bazel test",
			RunCommand:   "bazel run",
		},
		{
			Name:         "bazel",
			ProjectFile:  "WORKSPACE",
			BuildCommand: "bazel build",
			TestCommand:  "bazel test",
			RunCommand:   "bazel run",
		},
		{
			Name:         "bazel",
			ProjectFile:  "WORKSPACE",
			BuildCommand: "bazel build",
			TestCommand:  "bazel test",
			RunCommand:   "bazel run",
		},
		{
			Name:         "bazel",
			ProjectFile:  "WORKSPACE",
			BuildCommand: "bazel build",
			TestCommand:  "bazel test",
			RunCommand:   "bazel run",
		},
	}
}

// (projectile-register-project-type 'dotnet #'projectile-dotnet-project-p
//                                   :project-file '("?*.csproj" "?*.fsproj")
//                                   :compile "dotnet build"
//                                   :run "dotnet run"
//                                   :test "dotnet test")
// (projectile-register-project-type 'dotnet-sln #'projectile-dotnet-sln-project-p
//                                   :project-file "?*.sln"
//                                   :compile "dotnet build"
//                                   :run "dotnet run"
//                                   :test "dotnet test")
// (projectile-register-project-type 'nim-nimble #'projectile-nimble-project-p
//                                   :project-file "?*.nimble"
//                                   :compile "nimble --noColor build --colors:off"
//                                   :install "nimble --noColor install --colors:off"
//                                   :test "nimble --noColor test -d:nimUnittestColor:off --colors:off"
//                                   :run "nimble --noColor run --colors:off"
//                                   :src-dir "src"
//                                   :test-dir "tests")
// ;; File-based detection project types
//
// ;; Universal
// (projectile-register-project-type 'xmake '("xmake.lua")
//                                   :project-file "xmake.lua"
//                                   :compile "xmake build"
//                                   :test "xmake test"
//                                   :run "xmake run"
//                                   :install "xmake install")
// (projectile-register-project-type 'scons '("SConstruct")
//                                   :project-file "SConstruct"
//                                   :compile "scons"
//                                   :test "scons test"
//                                   :test-suffix "test")
// (projectile-register-project-type 'meson '("meson.build")
//                                   :project-file "meson.build"
//                                   :compilation-dir "build"
//                                   :configure "meson %s"
//                                   :compile "ninja"
//                                   :test "ninja test")
// (projectile-register-project-type 'nix '("default.nix")
//                                   :project-file "default.nix"
//                                   :compile "nix-build"
//                                   :test "nix-build")
// (projectile-register-project-type 'nix-flake '("flake.nix")
//                                   :project-file "flake.nix"
//                                   :compile "nix build"
//                                   :test "nix flake check"
//                                   :run "nix run")
// (projectile-register-project-type 'bazel '("WORKSPACE")
//                                   :project-file "WORKSPACE"
//                                   :compile "bazel build"
//                                   :test "bazel test"
//                                   :run "bazel run")
// (projectile-register-project-type 'debian '("debian/control")
//                                   :project-file "debian/control"
//                                   :compile "debuild -uc -us")
//
// ;; Make & CMake
// (projectile-register-project-type 'make '("Makefile")
//                                   :project-file "Makefile"
//                                   :compile "make"
//                                   :test "make test"
//                                   :install "make install")
// (projectile-register-project-type 'gnumake '("GNUMakefile")
//                                   :project-file "GNUMakefile"
//                                   :compile "make"
//                                   :test "make test"
//                                   :install "make install")
// (projectile-register-project-type 'cmake '("CMakeLists.txt")
//                                   :project-file "CMakeLists.txt"
//                                   :configure #'projectile--cmake-configure-command
//                                   :compile #'projectile--cmake-compile-command
//                                   :test #'projectile--cmake-test-command
//                                   :install #'projectile--cmake-install-command
//                                   :package #'projectile--cmake-package-command)
// ;; go-task/task
// (projectile-register-project-type 'go-task '("Taskfile.yml")
//                                   :project-file "Taskfile.yml"
//                                   :compile "task build"
//                                   :test "task test"
//                                   :install "task install")
// ;; Go should take higher precedence than Make because Go projects often have a Makefile.
// (projectile-register-project-type 'go projectile-go-project-test-function
//                                   :compile "go build"
//                                   :test "go test ./..."
//                                   :test-suffix "_test")
// ;; PHP
// (projectile-register-project-type 'php-symfony '("composer.json" "app" "src" "vendor")
//                                   :project-file "composer.json"
//                                   :compile "app/console server:run"
//                                   :test "phpunit -c app "
//                                   :test-suffix "Test")
// ;; Erlang & Elixir
// (projectile-register-project-type 'rebar '("rebar.config")
//                                   :project-file "rebar.config"
//                                   :compile "rebar3 compile"
//                                   :test "rebar3 do eunit,ct"
//                                   :test-suffix "_SUITE")
// (projectile-register-project-type 'elixir '("mix.exs")
//                                   :project-file "mix.exs"
//                                   :compile "mix compile"
//                                   :src-dir "lib/"
//                                   :test "mix test"
//                                   :test-suffix "_test")
// ;; JavaScript
// (projectile-register-project-type 'grunt '("Gruntfile.js")
//                                   :project-file "Gruntfile.js"
//                                   :compile "grunt"
//                                   :test "grunt test")
// (projectile-register-project-type 'gulp '("gulpfile.js")
//                                   :project-file "gulpfile.js"
//                                   :compile "gulp"
//                                   :test "gulp test")
// (projectile-register-project-type 'npm '("package.json" "package-lock.json")
//                                   :project-file "package.json"
//                                   :compile "npm install && npm run build"
//                                   :test "npm test"
//                                   :test-suffix ".test")
// (projectile-register-project-type 'yarn '("package.json" "yarn.lock")
//                                   :project-file "package.json"
//                                   :compile "yarn && yarn build"
//                                   :test "yarn test"
//                                   :test-suffix ".test")
// (projectile-register-project-type 'pnpm '("package.json" "pnpm-lock.yaml")
//                                   :project-file "package.json"
//                                   :compile "pnpm install && pnpm build"
//                                   :test "pnpm test"
//                                   :test-suffix ".test")
// ;; Angular
// (projectile-register-project-type 'angular '("angular.json" ".angular-cli.json")
//                                   :project-file "angular.json"
//                                   :compile "ng build"
//                                   :run "ng serve"
//                                   :test "ng test"
//                                   :test-suffix ".spec")
// ;; Python
// (projectile-register-project-type 'django '("manage.py")
//                                   :project-file "manage.py"
//                                   :compile "python manage.py runserver"
//                                   :test "python manage.py test"
//                                   :test-prefix "test_"
//                                   :test-suffix"_test")
// (projectile-register-project-type 'python-pip '("requirements.txt")
//                                   :project-file "requirements.txt"
//                                   :compile "python setup.py build"
//                                   :test "python -m unittest discover"
//                                   :test-prefix "test_"
//                                   :test-suffix"_test")
// (projectile-register-project-type 'python-pkg '("setup.py")
//                                   :project-file "setup.py"
//                                   :compile "python setup.py build"
//                                   :test "python -m unittest discover"
//                                   :test-prefix "test_"
//                                   :test-suffix"_test")
// (projectile-register-project-type 'python-tox '("tox.ini")
//                                   :project-file "tox.ini"
//                                   :compile "tox -r --notest"
//                                   :test "tox"
//                                   :test-prefix "test_"
//                                   :test-suffix"_test")
// (projectile-register-project-type 'python-pipenv '("Pipfile")
//                                   :project-file "Pipfile"
//                                   :compile "pipenv run build"
//                                   :test "pipenv run test"
//                                   :test-prefix "test_"
//                                   :test-suffix "_test")
// (projectile-register-project-type 'python-poetry '("poetry.lock")
//                                   :project-file "poetry.lock"
//                                   :compile "poetry build"
//                                   :test "poetry run python -m unittest discover"
//                                   :test-prefix "test_"
//                                   :test-suffix "_test")
// (projectile-register-project-type 'python-toml '("pyproject.toml")
//                                   :project-file "pyproject.toml"
//                                   :compile "python -m build"
//                                   :test "python -m unittest discover"
//                                   :test-prefix "test_"
//                                   :test-suffix "_test")
// ;; Java & friends
// (projectile-register-project-type 'maven '("pom.xml")
//                                   :project-file "pom.xml"
//                                   :compile "mvn -B clean install"
//                                   :test "mvn -B test"
//                                   :test-suffix "Test"
//                                   :src-dir "src/main/"
//                                   :test-dir "src/test/")
// (projectile-register-project-type 'gradle '("build.gradle")
//                                   :project-file "build.gradle"
//                                   :compile "gradle build"
//                                   :test "gradle test"
//                                   :test-suffix "Spec")
// (projectile-register-project-type 'gradlew '("gradlew")
//                                   :project-file "gradlew"
//                                   :compile "./gradlew build"
//                                   :test "./gradlew test"
//                                   :test-suffix "Spec")
// (projectile-register-project-type 'grails '("application.yml" "grails-app")
//                                   :project-file "application.yml"
//                                   :compile "grails package"
//                                   :test "grails test-app"
//                                   :test-suffix "Spec")
// ;; Scala
// (projectile-register-project-type 'sbt '("build.sbt")
//                                   :project-file "build.sbt"
//                                   :src-dir "main"
//                                   :test-dir "test"
//                                   :compile "sbt compile"
//                                   :test "sbt test"
//                                   :test-suffix "Spec")
//
// (projectile-register-project-type 'mill '("build.sc")
//                                   :project-file "build.sc"
//                                   :src-dir "src/"
//                                   :test-dir "test/src/"
//                                   :compile "mill __.compile"
//                                   :test "mill __.test"
//                                   :test-suffix "Test")
//
// ;; Clojure
// (projectile-register-project-type 'lein-test '("project.clj")
//                                   :project-file "project.clj"
//                                   :compile "lein compile"
//                                   :test "lein test"
//                                   :test-suffix "_test")
// (projectile-register-project-type 'lein-midje '("project.clj" ".midje.clj")
//                                   :project-file "project.clj"
//                                   :compile "lein compile"
//                                   :test "lein midje"
//                                   :test-prefix "t_")
// (projectile-register-project-type 'boot-clj '("build.boot")
//                                   :project-file "build.boot"
//                                   :compile "boot aot"
//                                   :test "boot test"
//                                   :test-suffix "_test")
// (projectile-register-project-type 'clojure-cli '("deps.edn")
//                                   :project-file "deps.edn"
//                                   :test-suffix "_test")
// (projectile-register-project-type 'bloop '(".bloop")
//                                   :project-file ".bloop"
//                                   :compile "bloop compile root"
//                                   :test "bloop test --propagate --reporter scalac root"
//                                   :src-dir "src/main/"
//                                   :test-dir "src/test/"
//                                   :test-suffix "Spec")
// ;; Ruby
// (projectile-register-project-type 'ruby-rspec '("Gemfile" "lib" "spec")
//                                   :project-file "Gemfile"
//                                   :compile "bundle exec rake"
//                                   :src-dir "lib/"
//                                   :test "bundle exec rspec"
//                                   :test-dir "spec/"
//                                   :test-suffix "_spec")
// (projectile-register-project-type 'ruby-test '("Gemfile" "lib" "test")
//                                   :project-file "Gemfile"
//                                   :compile"bundle exec rake"
//                                   :src-dir "lib/"
//                                   :test "bundle exec rake test"
//                                   :test-suffix "_test")
// ;; Rails needs to be registered after npm, otherwise `package.json` makes it `npm`.
// ;; https://github.com/bbatsov/projectile/pull/1191
// (projectile-register-project-type 'rails-test '("Gemfile" "app" "lib" "db" "config" "test")
//                                   :project-file "Gemfile"
//                                   :compile "bundle exec rails server"
//                                   :src-dir "app/"
//                                   :test "bundle exec rake test"
//                                   :test-suffix "_test")
// (projectile-register-project-type 'rails-rspec '("Gemfile" "app" "lib" "db" "config" "spec")
//                                   :project-file "Gemfile"
//                                   :compile "bundle exec rails server"
//                                   :src-dir "app/"
//                                   :test "bundle exec rspec"
//                                   :test-dir "spec/"
//                                   :test-suffix "_spec")
// ;; Crystal
// (projectile-register-project-type 'crystal-spec '("shard.yml")
//                                   :project-file "shard.yml"
//                                   :src-dir "src/"
//                                   :test "crystal spec"
//                                   :test-dir "spec/"
//                                   :test-suffix "_spec")
//
// ;; Emacs
// (projectile-register-project-type 'emacs-cask '("Cask")
//                                   :project-file "Cask"
//                                   :compile "cask install"
//                                   :test-prefix "test-"
//                                   :test-suffix "-test")
//
// (projectile-register-project-type 'emacs-eask '("Eask")
//                                   :project-file "Eask"
//                                   :compile "eask install"
//                                   :test-prefix "test-"
//                                   :test-suffix "-test")
//
// (projectile-register-project-type 'emacs-eldev #'projectile-eldev-project-p
//                                   :project-file "Eldev"
//                                   :compile "eldev compile"
//                                   :test "eldev test"
//                                   :run "eldev emacs"
//                                   :package "eldev package")
//
// ;; R
// (projectile-register-project-type 'r '("DESCRIPTION")
//                                   :project-file "DESCRIPTION"
//                                   :compile "R CMD INSTALL --with-keep.source ."
//                                   :test (concat "R CMD check -o " temporary-file-directory " ."))
//
// ;; Haskell
// (projectile-register-project-type 'haskell-stack '("stack.yaml")
//                                   :project-file "stack.yaml"
//                                   :compile "stack build"
//                                   :test "stack build --test"
//                                   :test-suffix "Spec")
//
// ;; Rust
// (projectile-register-project-type 'rust-cargo '("Cargo.toml")
//                                   :project-file "Cargo.toml"
//                                   :compile "cargo build"
//                                   :test "cargo test"
//                                   :run "cargo run")
//
// ;; Racket
// (projectile-register-project-type 'racket '("info.rkt")
//                                   :project-file "info.rkt"
//                                   :test "raco test ."
//                                   :install "raco pkg install"
//                                   :package "raco pkg create --source $(pwd)")
//
// ;; Dart
// (projectile-register-project-type 'dart '("pubspec.yaml")
//                                   :project-file "pubspec.yaml"
//                                   :compile "pub get"
//                                   :test "pub run test"
//                                   :run "dart"
//                                   :test-suffix "_test.dart")
//
// ;; Elm
// (projectile-register-project-type 'elm '("elm.json")
//                                   :project-file "elm.json"
//                                   :compile "elm make")
//
// ;; Julia
// (projectile-register-project-type 'julia '("Project.toml")
//                                   :project-file "Project.toml"
//                                   :compile "julia --project=@. -e 'import Pkg; Pkg.precompile(); Pkg.build()'"
//                                   :test "julia --project=@. -e 'import Pkg; Pkg.test()' --check-bounds=yes"
//                                   :src-dir "src"
//                                   :test-dir "test")
//
// ;; OCaml
// (projectile-register-project-type 'ocaml-dune '("dune-project")
//                                   :project-file "dune-project"
//                                   :compile "dune build"
//                                   :test "dune runtest")
//
// ;; Zig
// (projectile-register-project-type 'zig '("build.zig.zon")
//                                   :project-file "build.zig.zon"
//                                   :compile "zig build"
//                                   :test "zig build test"
//                                   :run "zig build run")
