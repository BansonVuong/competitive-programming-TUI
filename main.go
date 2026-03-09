package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	problems, err := DiscoverProblems("./problems")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to discover problems: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(NewModel(problems), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}
