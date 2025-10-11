package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ===========================================================
// MODEL STRUCTURE
// ===========================================================

// Focus constants to indicate which component currently receives key events.
const (
	focusSearchInput = "search_input"
	focusPackageList = "package_list"
)

// AppModel is the main data structure that holds all the state for our TUI program.
// Bubble Tea uses a "model" to track what should be displayed and what happens on input.
type AppModel struct {
	searchInputField   textinput.Model  // text input field where user types package name to search
	packageListDisplay list.Model       // list UI element to display yay search results
	isSearching        bool             // flag to track if we’re currently waiting for yay command output
	errorMessage       string           // stores an error message if yay fails
	selected           map[int]struct{} // which item is selected
	focusedComponent   string           // which UI component currently receives keys

}

// ===========================================================
// INITIALIZATION
// ===========================================================

// NewAppModel initializes our TUI's starting state.
func NewAppModel() AppModel {
	// Create and configure a text input field for the search bar.
	searchInputField := textinput.New()
	searchInputField.Placeholder = "Enter package name to search..."
	searchInputField.Focus() // Give focus to the input so we can type right away.
	searchInputField.CharLimit = 64
	searchInputField.Width = 40

	// Create a list for displaying package search results.
	// Each list item will be a package name and short description from yay output.
	packageListDisplay := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	packageListDisplay.Title = "Search Results"

	return AppModel{
		searchInputField:   searchInputField,
		packageListDisplay: packageListDisplay,
		isSearching:        false,
		errorMessage:       "",
		selected:           make(map[int]struct{}),
		focusedComponent:   focusSearchInput, // start with search focused
	}
}

// ===========================================================
// BUBBLE TEA CORE FUNCTIONS
// ===========================================================

// Init runs once at the beginning. It’s used for setup.
// We don’t need any background startup tasks, so return nil.
func (model AppModel) Init() tea.Cmd {
	return nil
}

// Update is called whenever the user presses a key or an event occurs.
// It modifies the model’s state and can trigger new commands (like running yay).
func (model AppModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd

	switch msg := message.(type) {

	case tea.KeyMsg:
		k := msg.String()

		// Global quit keys
		if k == "ctrl+c" || k == "esc" || k == "q" {
			return model, tea.Quit
		}

		// Focus switching keys (handled regardless of focused component)
		switch k {
		case "s":
			// focus search input
			model.focusedComponent = focusSearchInput
			// ensure search input has cursor so user can type
			model.searchInputField.Focus()
			return model, nil

		case "tab":
			// toggle focus to package list
			model.focusedComponent = focusPackageList
			return model, nil
		}

		// Now route key events to the component that has focus.
		if model.focusedComponent == focusSearchInput {
			// If Enter is pressed while search input is focused, start search.
			if k == "enter" {
				if !model.isSearching && len(strings.TrimSpace(model.searchInputField.Value())) > 0 {
					model.isSearching = true
					model.packageListDisplay.Title = "Searching..."
					// blur the input by moving focus, but we manage blur by changing focusedComponent
					model.focusedComponent = focusPackageList
					model.searchInputField.Blur() // textinput has Blur()
					return model, runYaySearchCommand(model.searchInputField.Value())
				}
			}

			// Otherwise forward normal editing keys to the input component
			model.searchInputField, command = model.searchInputField.Update(msg)
			return model, command
		}

		// If the package list has focus, handle its navigation keys here.
		if model.focusedComponent == focusPackageList {
			// Let the list component handle keys (it will respond to up/down, etc.)
			model.packageListDisplay, command = model.packageListDisplay.Update(msg)

			// Optionally add custom keys for list-level actions:
			switch k {
			case "enter":
				// user pressed Enter on a list item - do something (e.g., show details)
				// You can inspect model.packageListDisplay.Index() or SelectedItem() here.
				return model, nil
			case "j", "down", "k", "up":
				// list.Update already handles movement, so no-op here if using list.Update
			case "s":
				// already handled above; kept for clarity
			}
			return model, command
		}

		// Fallback: just return
		return model, nil

	case yaySearchResultMsg:
		model.isSearching = false
		model.errorMessage = ""
		model.packageListDisplay.Title = fmt.Sprintf("Results for \"%s\"", model.searchInputField.Value())
		model.packageListDisplay.SetItems(convertSearchResultsToListItems(msg.results))
		model.packageListDisplay.SetSize(80, 20)

		// set focus to the package list so up/down keys will work
		model.focusedComponent = focusPackageList
		// blur the search input (textinput supports Blur)
		model.searchInputField.Blur()

		// select the first item for convenience
		if len(model.packageListDisplay.Items()) > 0 {
			model.packageListDisplay.Select(0)
		}
		return model, nil

	// Handle the message when yay search fails.
	case yaySearchErrorMsg:
		model.isSearching = false
		model.errorMessage = msg.errorText
		model.packageListDisplay.Title = "Error"
		model.packageListDisplay.SetItems(nil)

		// return focus to the search input so user can try again
		model.focusedComponent = focusSearchInput
		model.searchInputField.Focus()

		return model, nil

	}

	// Always update the list view (so it responds to scrolling, etc.)
	model.packageListDisplay, command = model.packageListDisplay.Update(message)
	return model, command
}

// View defines how everything should look on the screen.
func (model AppModel) View() string {
	mainLayoutStyle := lipgloss.NewStyle().Margin(1, 2)

	if model.isSearching {
		return mainLayoutStyle.Render(fmt.Sprintf(
			"%s\n\nSearching for \"%s\"...",
			model.searchInputField.View(),
			model.searchInputField.Value(),
		))
	}

	if model.errorMessage != "" {
		return mainLayoutStyle.Render(fmt.Sprintf(
			"%s\n\nError: %s",
			model.searchInputField.View(),
			model.errorMessage,
		))
	}

	return mainLayoutStyle.Render(fmt.Sprintf(
		"%s\n\n%s",
		model.searchInputField.View(),
		model.packageListDisplay.View(),
	))
}

// ===========================================================
// YAY COMMAND EXECUTION
// ===========================================================

// yaySearchResultMsg and yaySearchErrorMsg are message types sent back to Update()
// when yay finishes running. Bubble Tea uses this for async work.
type yaySearchResultMsg struct {
	results []string
}

type yaySearchErrorMsg struct {
	errorText string
}

// runYaySearchCommand starts yay -Ss <query> as a background process
// and returns a command that Bubble Tea will handle asynchronously.
func runYaySearchCommand(searchQuery string) tea.Cmd {
	return func() tea.Msg {
		// Build the command for yay search.
		command := exec.Command("yay", "-Ss", searchQuery)
		fmt.Println("Running yay command:", searchQuery)
		// Get the stdout pipe (output from yay).
		outputPipe, err := command.StdoutPipe()
		if err != nil {
			return yaySearchErrorMsg{errorText: fmt.Sprintf("Failed to capture output: %v", err)}
		}

		// Start the command execution.
		if err := command.Start(); err != nil {
			return yaySearchErrorMsg{errorText: fmt.Sprintf("Failed to run yay: %v", err)}
		}

		// Read yay’s output line by line.
		scanner := bufio.NewScanner(outputPipe)
		var results []string
		for scanner.Scan() {
			results = append(results, scanner.Text())
		}

		// Wait for yay to finish.
		if err := command.Wait(); err != nil {
			return yaySearchErrorMsg{errorText: fmt.Sprintf("yay failed: %v", err)}
		}

		return yaySearchResultMsg{results: results}
	}
}

// convertSearchResultsToListItems turns yay’s output lines into list items
// so they can be displayed in the Bubble Tea list UI.
func convertSearchResultsToListItems(results []string) []list.Item {
	var listItems []list.Item
	for _, line := range results {
		if strings.TrimSpace(line) == "" {
			continue
		}
		listItems = append(listItems, listItem(line))
	}
	return listItems
}

// listItem is a simple string wrapper so Bubble Tea's list can display it.
type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return string(i) }

// ===========================================================
// MAIN FUNCTION
// ===========================================================

func main() {
	app := tea.NewProgram(NewAppModel(), tea.WithAltScreen())

	// Run the Bubble Tea program — this starts the interactive TUI loop.
	if _, err := app.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
