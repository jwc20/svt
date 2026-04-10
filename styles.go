package svt

import "charm.land/lipgloss/v2"

var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170")).MarginBottom(1)
	SelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	DimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	PanelStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).Padding(0, 1)
	FocusLabel    = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	PlainLabel    = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	GoodStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	BadStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	WarnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	EventStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	PromptStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
)
