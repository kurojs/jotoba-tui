package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kurojs/jotoba-tui/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
