package tui

import (
	"strings"

	// "github.com/charmbracelet/bubbles/list"
	"github.com/evertras/bubble-table/table"
)

// convertSearchResultsToListItems turns yay output lines into list items.
// func convertSearchResultsToListItems(results []string) []list.Item {
// 	var listItems []list.Item
// 	for _, line := range results {
// 		if strings.TrimSpace(line) == "" {
// 			continue
// 		}
// 		listItems = append(listItems, listItem(line))
// 	}
// 	return listItems
// }

// type listItem string

// func (i listItem) Title() string       { return string(i) }
// func (i listItem) Description() string { return "" }
// func (i listItem) FilterValue() string { return string(i) }

func convertSearchResultsToTableRows(results []string) []table.Row {
	var (
		rows           []table.Row
		aurRows        []table.Row
		currentRepo    string
		currentName    string
		currentVersion string
		currentDesc    string
		installed      string
		output         []table.Row
		before         string
		after          string
		cutBool        bool
	)

	for i := 0; i < len(results); i++ {
		line := results[i]
		if !strings.HasPrefix(line, " ") {

			if strings.Contains(line, "/aur") {
				before, after, cutBool = strings.Cut(line, "/")
				rows = append(rows, table.NewRow(table.RowData{
					"repo":        "Aur",
					"name":        before,
					"version":     after,
					"description": "line",
					"installed":   cutBool,
				}))
			} else {
				before, after, cutBool = strings.Cut(line, "/")
				if cutBool {
					currentRepo = before
				}
				before, after, cutBool = strings.Cut(after, " ")
				if cutBool {
					currentName = before
				}
				before, cutBool = strings.CutSuffix(after, " ")
				if cutBool {
					currentVersion = before
				}
				currentDesc = results[i+1] // description is next line
				if strings.Contains(line, "Installed") {
					installed = "Installed"
				} else {
					installed = ""
				}

				rows = append(rows, table.NewRow(table.RowData{
					"repo":        currentRepo,
					"name":        currentName,
					"version":     currentVersion,
					"description": currentDesc,
					"installed":   installed,
				}))

			}

		} else {
			continue
		}
	}
	output = append(rows, aurRows...)

	return output
}
