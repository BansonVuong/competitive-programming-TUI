package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	problemsDir := "./problems"
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
