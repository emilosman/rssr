package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) View() tea.View {
	var v tea.View
	switch {
	case m.i != nil:
		// Item view
		feedTitle := titleStyle.Render(m.i.FeedTitle)
		itemTitle := renderedItemTitle(m)
		status := renderedStatus(m)
		content := contentStyle.Render(m.v.View())
		help := helpStyle.Render(m.vh.View(viewKeyMap{}))
		view := lipgloss.JoinVertical(lipgloss.Left, feedTitle, itemTitle, status, content, help)
		v.SetContent(view)
		return v
	case m.f != nil:
		// Feed view
		title := renderedTitle(m)
		status := renderedStatus(m)
		list := m.li.View()
		view := lipgloss.JoinVertical(lipgloss.Left, title, status, list)
		v.SetContent(view)
		return v
	default:
		// List view
		tabs := renderedTabs(m)
		status := renderedStatus(m)
		list := m.lf.View()
		view := lipgloss.JoinVertical(lipgloss.Left, tabs, status, list)
		v.SetContent(view)
		return v
	}
}
