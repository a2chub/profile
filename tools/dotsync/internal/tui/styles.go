package tui

import "charm.land/lipgloss/v2"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle()

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")) // green

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")) // yellow

	dangerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")) // red

	managedStyle = lipgloss.NewStyle().
			Faint(true)

	statusBarStyle = lipgloss.NewStyle().
			Reverse(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Faint(true)
)
