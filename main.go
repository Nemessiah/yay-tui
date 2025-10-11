package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"yay-tui/tui"
)

func main() {
	app := tea.NewProgram(tui.NewAppModel(), tea.WithAltScreen())

	if _, err := app.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
