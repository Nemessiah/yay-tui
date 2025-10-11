package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

// convertSearchResultsToListItems turns yay output lines into list items.
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

type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return string(i) }
