package tui

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssr/internal/rss"
	"github.com/muesli/reflow/truncate"
)

type feedUpdatedMsg struct {
	Feed *rss.RssFeed
	Err  error
}

type feedsDoneMsg struct{}
type statusClearMsg struct{}

func updateAllFeedsCmd(m *model) tea.Cmd {
	return func() tea.Msg {
		results, err := m.l.UpdateAllFeeds()
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		go func() {
			for res := range results {
				m.prog.Send(feedUpdatedMsg{Feed: res.Feed, Err: res.Err})
			}
			m.prog.Send(feedsDoneMsg{})
		}()

		return MsgUpdatingAllFeeds
	}
}

func updateTabFeedsCmd(m *model) tea.Cmd {
	return func() tea.Msg {
		feeds, err := m.l.GetCategory(activeTab(m.tabs, m.activeTab))
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		results, err := rss.UpdateFeeds(feeds...)
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		go func() {
			for res := range results {
				m.prog.Send(feedUpdatedMsg{Feed: res.Feed, Err: res.Err})
			}
			m.prog.Send(feedsDoneMsg{})
		}()

		return MsgUpdatingAllFeeds
	}
}

func updateFeedCmd(m *model, feed *rss.RssFeed) tea.Cmd {
	return func() tea.Msg {
		results, err := rss.UpdateFeeds(feed)
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		go func() {
			for res := range results {
				m.prog.Send(feedUpdatedMsg{Feed: res.Feed, Err: res.Err})
			}
		}()

		return fmt.Sprintf("%s %s", MsgUpdatingFeed, feed.Url)
	}
}

// Builds the feed list and sets the items
func rebuildFeedList(m *model) tea.Cmd {
	items := buildFeedList(m)
	m.lf.SetItems(items)
	return nil
}

func rebuildItemsList(m *model) tea.Cmd {
	if m.li.FilterState().String() != "filter applied" {
		items := buildItemsList(m)
		m.li.SetItems(items)
	}
	return nil
}

// Builds the feed list
func buildFeedList(m *model) []list.Item {
	var listItems []list.Item

	activeTab := activeTab(m.tabs, m.activeTab)
	feeds, err := m.l.GetCategory(activeTab)
	if err != nil {
		feeds = m.l.Feeds
	}

	if len(feeds) != 0 {
		for _, feed := range feeds {
			if feed.Category == "" {
				continue
			}

			title := feed.Title()
			description := feed.Latest()

			if feed.HasUnread() {
				title = unreadStyle.Render(title)
			}

			if feed.Error != "" {
				description = errorStyle.Render(description)
			}

			width := uint(m.lf.Width() - 3)
			title = truncate.StringWithTail(title, width, "...")
			description = truncate.StringWithTail(description, width, "...")

			listItems = append(listItems, feedItem{
				title:   title,
				desc:    description,
				rssFeed: feed,
			})
		}
	}

	return listItems
}

func activeTab(t []string, a int) string {
	var activeTab string
	if len(t) != 0 {
		activeTab = t[a]
	}
	return activeTab
}

func buildItemsList(m *model) []list.Item {
	feed := m.f
	listItems := make([]list.Item, 0, len(feed.RssItems))
	for idx := range feed.RssItems {
		ri := feed.RssItems[idx]
		title := ri.Title()
		description := ri.Description()

		if ri.Bookmark {
			title = bookmarkStyle.Render(title)
		}

		if !ri.Read {
			title = unreadStyle.Render(title)
		}

		width := uint(m.li.Width() - 3)
		title = truncate.StringWithTail(title, width, "...")
		description = truncate.StringWithTail(description, width, "...")

		listItems = append(listItems, rssListItem{
			title: title,
			desc:  description,
			item:  ri,
		})
	}
	return listItems
}

func renderedTabs(m *model) string {
	var renderedTabs string
	for i, tab := range m.tabs {
		if i == m.activeTab {
			renderedTabs += activeTabStyle.Render(tab)
		} else {
			feeds, _ := m.l.GetCategory(tab)
			hasUnread := false
			for _, f := range feeds {
				if f.HasUnread() {
					hasUnread = true
					break
				}
			}
			if hasUnread {
				renderedTabs += unreadTabStyle.Render(tab)
			} else {
				renderedTabs += inactiveTabStyle.Render(tab)
			}
		}
	}

	return renderedTabs
}

func renderedTitle(m *model) string {
	return titleStyle.Render(m.title)
}

func renderedStatus(m *model) string {
	return statusStyle.Render(m.status)
}

func renderedItemTitle(m *model) string {
	if m.i.Bookmark {
		return bookmarkItemTitleStyle.Render(m.i.Title())
	}
	if !m.i.Read {
		return unreadItemTitleStyle.Render(m.i.Title())
	}
	return itemTitleStyle.Render(m.i.Title())
}

func openInBrowser(raw string) error {
	var cmd *exec.Cmd

	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return err
	}

	url := parsed.String()

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func (m *model) UpdateTitle(title string) {
	m.title = title
}

func (m *model) UpdateStatus(msg string) {
	m.status = msg

	if m.clearTimer != nil {
		m.clearTimer.Stop()
	}

	m.clearTimer = time.AfterFunc(8*time.Second, func() {
		m.prog.Send(statusClearMsg{})
	})
}

func (m *model) SaveState() error {
	dataFilePath, err := rss.DataFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(dataFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return m.l.Save(f, time.Now())
}

func BuildApp() {
	m := initialModel()
	p := tea.NewProgram(m)
	m.prog = p

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
