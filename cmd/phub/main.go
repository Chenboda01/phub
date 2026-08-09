package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"phub/internal/ui"
)

func main() {
	program := tea.NewProgram(ui.New())
	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "phub could not start:", err)
		os.Exit(1)
	}
}
