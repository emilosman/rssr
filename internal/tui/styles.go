package tui

import "github.com/charmbracelet/lipgloss"

var (
	unreadStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")) // green

	bookmarkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")) // yellow

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")) // red

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("13")). // magenta
		Margin(0, 0, 0, 1)

	itemTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13")).
		Margin(0, 0, 0, 1)

	bookmarkItemTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11")).
		Margin(0, 0, 0, 1)

	unreadItemTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("2")).
		Margin(1, 0, 0, 1)

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")) // gray

	contentStyle = lipgloss.NewStyle().
		Margin(0, 0, 0, 1)

	helpStyle = lipgloss.NewStyle().
		Margin(0, 0, 0, 2)

	activeTabStyle = lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.Color("13")).
		Margin(0, 1, 0, 1)

	unreadTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("2")).
		Margin(0, 1, 0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")). // light gray/white
		Margin(0, 1, 0, 1)
)
