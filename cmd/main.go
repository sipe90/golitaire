package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sipe90/golitaire/internal/golitaire"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Printf("Failed to open log file %v:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	p := tea.NewProgram(golitaire.CreateModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to start application: %v", err)
		os.Exit(1)
	}
}
