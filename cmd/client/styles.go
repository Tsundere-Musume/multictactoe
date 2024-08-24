package main

import "github.com/charmbracelet/lipgloss"

const (
	active = lipgloss.Color("#f6c177")
	// inactive     = lipgloss.Color("#c4a7e7")
	invalid      = lipgloss.Color("#eb6f92")
	valid        = lipgloss.Color("#a6e3a1")
	borderColor  = lipgloss.Color("#31748f")
	TEXT         = lipgloss.Color("#313244")
	NOT_EDITABLE = lipgloss.Color("#ebbcba")
)

var (
	base        = lipgloss.NewStyle().Padding(0, 1).Foreground(active).BorderForeground(borderColor)
	borderStyle = lipgloss.NewStyle().Foreground(borderColor).BorderForeground(borderColor)
	boardBorder = borderStyle.UnsetPadding().BorderStyle(lipgloss.ThickBorder())
)
