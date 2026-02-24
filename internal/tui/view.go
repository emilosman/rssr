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
		title := renderedTitle(m)
		itemTitle := renderedItemTitle(m)
		status := renderedStatus(m)
		feedTitle := titleStyle.Render(m.i.FeedTitle)
		content := contentStyle.Render(m.v.View())
		help := helpStyle.Render(m.vh.View(viewKeyMap{}))
		view := lipgloss.JoinVertical(lipgloss.Left, title, itemTitle, status, feedTitle, content, help)
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
