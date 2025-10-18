package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles keypresses and messages.
func (model AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m := msg.(type) {

	case tea.WindowSizeMsg:
		model.Width = m.Width
		model.Height = m.Height
		model.packageListDisplay.SetSize(model.Width-4, model.Height-4)
	case tea.KeyMsg:
		k := m.String()

		// Global quit keys
		if k == "ctrl+c" || k == "esc" {
			return model, tea.Quit
		}

		// Focus switching
		switch k {
		case "s":
			model.focusedComponent = focusSearchInput
			model.searchInputField.Focus()
			return model, nil
		case "tab":
			model.focusedComponent = focusPackageTable
			model.searchInputField.Blur()
			return model, nil
		}

		// Search input focused
		if model.focusedComponent == focusSearchInput {
			if k == "enter" && !model.isSearching && strings.TrimSpace(model.searchInputField.Value()) != "" {
				model.isSearching = true
				model.packageListDisplay.Title = "Searching..."
				model.focusedComponent = focusPackageList
				model.searchInputField.Blur()
				return model, runYaySearchCommand(model.searchInputField.Value())
			}
			model.searchInputField, cmd = model.searchInputField.Update(m)
			return model, cmd
		}

		// Package list focused
		if model.focusedComponent == focusPackageList {
			model.packageListDisplay, cmd = model.packageListDisplay.Update(m)
			return model, cmd
		}
		switch model.focusedComponent {
		case focusPackageTable:
			var cmd tea.Cmd
			model.packageTable, cmd = model.packageTable.Update(m)
			return model, cmd
		}

	case yaySearchResultMsg:
		model.isSearching = false
		model.errorMessage = ""
		// model.packageListDisplay.Title = fmt.Sprintf("Results for \"%s\"", model.searchInputField.Value())
		// model.packageListDisplay.SetItems(convertSearchResultsToListItems(m.results))
		// model.focusedComponent = focusPackageList
		rows := convertSearchResultsToTableRows(m.results)
		model.packageTable = model.packageTable.WithRows(rows)
		model.focusedComponent = focusPackageTable
		model.searchInputField.Blur()
		model.searchComplete = true
		return model, nil

	case yaySearchErrorMsg:
		model.isSearching = false
		model.errorMessage = m.errorText
		model.packageListDisplay.Title = "Error"
		model.packageListDisplay.SetItems(nil)
		model.focusedComponent = focusSearchInput
		model.searchInputField.Focus()
		return model, nil
	}

	return model, cmd
}
