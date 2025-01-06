package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ProjectFinder struct {
	editor string
}

func NewProjectFinder() *ProjectFinder {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // because it is everywhere?
	}
	return &ProjectFinder{
		editor: editor,
	}
}

// Return selected directory instead of changing to it
func (pf *ProjectFinder) FindProject() (string, error) {
	pkCmd := exec.Command("pk", "list")
	fzfCmd := exec.Command("fzf", "--bind=ctrl-c:abort")

	pipe, err := pkCmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	fzfCmd.Stdin = pipe
	fzfCmd.Stderr = os.Stderr

	if err := pkCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start pk command: %w", err)
	}

	projectOutput, err := fzfCmd.Output()
	if err != nil {
		pkCmd.Wait() // Clean up pk command
		return "", err
	}

	if err := pkCmd.Wait(); err != nil {
		return "", fmt.Errorf("pk command failed: %w", err)
	}

	selection := strings.TrimSpace(string(projectOutput))
	if selection == "" {
		return "", nil
	}

	menuOptions := []string{
		"find(edit)",
		"find(show)",
		"grep(edit)",
		"go to project",
	}

	menuCmd := exec.Command("fzf", "--bind=ctrl-c:abort")
	menuInput := strings.Join(menuOptions, "\n")
	menuCmd.Stdin = strings.NewReader(menuInput)
	menuCmd.Stderr = os.Stderr

	menuOutput, err := menuCmd.Output()
	if err != nil {
		return "", err
	}

	choice := strings.TrimSpace(string(menuOutput))
	switch choice {
	case "find(edit)":
		result, err := pf.FindFile(selection)
		fmt.Println(result)
		return "", err
	case "find(show)":
		fdCmd := exec.Command("fd", ".", selection)
		fdCmd.Stdout = os.Stdout
		fdCmd.Stderr = os.Stderr
		return "", fdCmd.Run()
	case "grep(edit)":
		err = pf.GrepEdit(selection)
		return "", err
	case "go to project":
		return selection, nil
	}

	return "", nil
}

func (pf *ProjectFinder) GrepEdit(projectPath string) error {
	os.Remove("/tmp/rg-fzf-r")
	os.Remove("/tmp/rg-fzf-f")

	rgPrefix := fmt.Sprintf("rg --column --line-number --no-heading --color=always --smart-case {q} %s", projectPath)

	fzfCmd := exec.Command("fzf",
		"--ansi",
		"--disabled",
		"--bind", fmt.Sprintf("start:reload:%s", rgPrefix),
		"--bind", fmt.Sprintf("change:reload:sleep 0.1; %s || true", rgPrefix),
		"--bind", `ctrl-t:transform:[[ ! $FZF_PROMPT =~ ripgrep ]] &&
			echo "rebind(change)+change-prompt(1. ripgrep> )+disable-search+transform-query:echo {q} > /tmp/rg-fzf-f; cat /tmp/rg-fzf-r" ||
			echo "unbind(change)+change-prompt(2. fzf> )+enable-search+transform-query:echo {q} > /tmp/rg-fzf-r; cat /tmp/rg-fzf-f"`,
		"--color", "hl:-1:underline,hl+:-1:underline:reverse",
		"--prompt", "1. ripgrep> ",
		"--delimiter", ":",
		"--header", "CTRL-T: Switch between ripgrep/fzf",
		"--preview", "bat --color=always {1} --highlight-line {2}",
		"--preview-window", "up,60%,border-bottom,+{2}+3/3,~3",
		"--bind", fmt.Sprintf("enter:become(%s {1} +{2})", pf.editor),
	)

	fzfCmd.Stdin = os.Stdin
	fzfCmd.Stdout = os.Stdout
	fzfCmd.Stderr = os.Stderr

	return fzfCmd.Run()
}

func (pf *ProjectFinder) FindFile(projectArg string) (string, error) {
	// Get project path using pk
	pkCmd := exec.Command("pk", projectArg)
	pkOutput, err := pkCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get project path: %w", err)
	}
	projectPath := strings.TrimSpace(string(pkOutput))

	// Run fd command
	fdCmd := exec.Command("fd", ".", projectPath, "-tf")
	fzfCmd := exec.Command("fzf")

	// Pipe fd output to fzf
	pipe, err := fdCmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	fzfCmd.Stdin = pipe
	fzfCmd.Stderr = os.Stderr

	// Start fd
	if err := fdCmd.Start(); err != nil {
		return "", err
	}

	// Get fzf output
	output, err := fzfCmd.Output()
	if err != nil {
		fdCmd.Wait() // Clean up fd command
		return "", err
	}

	if err := fdCmd.Wait(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return "", nil
	}

	fileInfo, err := os.Stat(result)
	if err != nil {
		return "", err
	}

	if fileInfo.IsDir() {
		return result, nil
	}

	// For files, launch editor and return empty string
	editorCmd := exec.Command(pf.editor, result)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	err = editorCmd.Run()
	return "", err
}
