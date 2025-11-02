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
		table.NewColumn("repo", "Repo", 10),
		table.NewColumn("name", "Package", 30),
		table.NewColumn("Version", "version", 20),
		table.NewColumn("description", "Description", 100),
		table.NewColumn("installed", "Installed", 10),
	}).WithRows([]table.Row{}).Focused(true)

	return t
}
