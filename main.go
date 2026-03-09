package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	problemsDirFlag := flag.String("problems-dir", "./problems", "directory containing .cpp/.md/.in/.out problem files")
	initFlag := flag.Bool("init", false, "run first-time setup")
	flag.Parse()

	problemsDir := filepath.Clean(*problemsDirFlag)
	if err := maybeRunInit(problemsDir, *initFlag); err != nil {
		fmt.Fprintf(os.Stderr, "initialization failed: %v\n", err)
		os.Exit(1)
	}
	if err := validateProblemsDir(problemsDir); err != nil {
		fmt.Fprintf(os.Stderr, "invalid problems directory %q: %v\n", problemsDir, err)
		fmt.Fprintf(os.Stderr, "Run `cptui --init` or specify `--problems-dir`.\n")
		os.Exit(1)
	}

	problems, err := DiscoverProblems(problemsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to discover problems in %q: %v\n", problemsDir, err)
		os.Exit(1)
	}

	p := tea.NewProgram(NewModel(problems), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}

func validateProblemsDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".cpp") {
			return nil
		}
	}

	return fmt.Errorf("no .cpp files found")
}
