package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI layout.
func (model AppModel) View() string {
	layout := lipgloss.NewStyle().Margin(1, 2)

	if model.isSearching {
		return layout.Render(fmt.Sprintf(
			"%s\n\nSearching for \"%s\"...",
			model.searchInputField.View(),
			model.searchInputField.Value(),
		))
	}

	if model.errorMessage != "" {
		return layout.Render(fmt.Sprintf(
			"%s\n\nError: %s",
			model.searchInputField.View(),
			model.errorMessage,
		))
	}

	if model.searchComplete {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			model.searchInputField.View(),
			model.packageTable.View(),
		)
	}

	return layout.Render(fmt.Sprintf(
		"%s\n\n%s",
		model.searchInputField.View(),
		model.packageListDisplay.View(),
	))
}
