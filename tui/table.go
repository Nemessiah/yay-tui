package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true)
	cellStyle   = lipgloss.NewStyle().Padding(0, 1)
)

func NewPackageTable() table.Model {
	t := table.New([]table.Column{
		table.NewColumn("repo", "Repo", 15),
		table.NewColumn("name", "Package", 15),
		table.NewColumn("Version", "version", 10),
		table.NewColumn("description", "Description", 60),
		table.NewColumn("installed", "Installed", 10),
	}).WithRows([]table.Row{}).Focused(true)

	return t
}
