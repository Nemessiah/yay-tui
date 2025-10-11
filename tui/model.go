package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Constants that track which UI element currently has focus.
const (
	focusSearchInput = "search_input"
	focusPackageList = "package_list"
)

// AppModel holds all TUI state.
type AppModel struct {
	searchInputField   textinput.Model
	packageListDisplay list.Model
	isSearching        bool
	errorMessage       string
	selected           map[int]struct{}
	focusedComponent   string
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

// NewAppModel creates and initializes the app's state.
func NewAppModel() AppModel {
	searchInput := textinput.New()
	searchInput.Placeholder = "Enter package name to search..."
	searchInput.CharLimit = 64
	searchInput.Width = 40
	searchInput.Focus()

	listDisplay := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listDisplay.Title = "Search Results"

	return AppModel{
		searchInputField:   searchInput,
		packageListDisplay: listDisplay,
		isSearching:        false,
		errorMessage:       "",
		selected:           make(map[int]struct{}),
		focusedComponent:   focusSearchInput,
	}
}
