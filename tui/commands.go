package tui

import (
	"yay-tui/yay"

	tea "github.com/charmbracelet/bubbletea"
)

// runYaySearchCommand launches yay -Ss <query> asynchronously
func runYaySearchCommand(searchQuery string) tea.Cmd {
	return func() tea.Msg {
		results, err := yay.Search(searchQuery)
		if err != nil {
			return yaySearchErrorMsg{errorText: err.Error()}
		}
		return yaySearchResultMsg{results: results}
	}
}

// runYayInspectCommand launches yay -Si <query> asynchronously
func runYayInspectCommand(searchQuery string) tea.Cmd {
	return func() tea.Msg {
		results, err := yay.Inspect(searchQuery)
		if err != nil {
			return yaySearchErrorMsg{errorText: err.Error()}
		}
		return yaySearchResultMsg{results: results}
	}
}
