package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssr/internal/rss"
	"github.com/muesli/reflow/wordwrap"
)

type keyHandler func(*model) tea.Cmd

var (
	feedKeyHandlers = map[string]keyHandler{
		"A":      handleMarkFeedRead,
		"b":      handlePrevUnreadFeed,
		"B":      handleViewBookmarks,
		"E":      handleEdit,
		"h":      handlePrevTab,
		"l":      handleNextTab,
		"n":      handleNextUnreadFeed,
		"o":      handleOpenLatest,
		"O":      handleOpenFeed,
		"p":      handlePrevUnreadFeed,
		"q":      handleQuit,
		"r":      handleUpdateFeed,
		"R":      handleUpdateAllFeeds,
		"S":      handleSync,
		"enter":  handleEnterFeed,
		"esc":    handleQuit,
		"tab":    handleNextTab,
		"left":   handlePrevTab,
		"right":  handleNextTab,
		"ctrl+a": handleMarkTabAsRead,
		"ctrl+c": handleInterrupt,
		"ctrl+r": handleTabUpdate,
	}

	itemKeyHandlers = map[string]keyHandler{
		"a":      handleToggleRead,
		"A":      handleMarkItemsRead,
		"b":      handleBack,
		"B":      handleViewBookmarks,
		"c":      handleToggleBookmark,
		"n":      handleNextUnreadItem,
		"o":      handleOpenItem,
		"O":      handleOpenFeed,
		"p":      handlePrevUnreadItem,
		"q":      handleBack,
		"r":      handleUpdateFeed,
		"R":      handleUpdateAllFeeds,
		"S":      handleSync,
		"enter":  handleViewItem,
		"esc":    handleBack,
		"ctrl+c": handleInterrupt,
	}

	viewKeyHandlers = map[string]keyHandler{
		"a":      handleToggleRead,
		"b":      handleBack,
		"B":      handleViewBookmarks,
		"c":      handleToggleBookmark,
		"g":      handleGoToStart,
		"h":      handleViewPrev,
		"l":      handleViewNext,
		"n":      handleNextUnreadItem,
		"o":      handleOpenItem,
		"p":      handlePrevUnreadItem,
		"q":      handleBack,
		"?":      handleViewHelp,
		"ctrl+c": handleInterrupt,
		"enter":  handleOpenItem,
		"esc":    handleBack,
		"left":   handleViewPrev,
		"right":  handleViewNext,
	}
)

func handleSync(m *model) tea.Cmd {
	err := m.l.SyncList()
	if err != nil {
		m.UpdateStatus(err.Error())
	} else {
		m.UpdateStatus("Sync ok")
		m.SaveState()
	}

	if m.f != nil {
		rebuildItemsList(m)
	}
	return rebuildFeedList(m)
}

func handleEdit(m *model) tea.Cmd {
	urlsFilePath, err := rss.UrlsFilePath()
	if err != nil {
		fmt.Println("Error opening config dir", err)
		return nil
	}
	urlsFile := filepath.Join(urlsFilePath, "urls.yaml")

	editor := os.Getenv("EDITOR")
	if editor == "" {
		switch runtime.GOOS {
		case "windows":
			editor = "notepad"
		default:
			editor = "nvim"
		}
	}

	m.SaveState()
	m.prog.ReleaseTerminal()
	cmd := exec.Command(editor, urlsFile)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		m.UpdateStatus(err.Error())
		return nil
	}

	m.prog.RestoreTerminal()
	filesystem := os.DirFS(urlsFilePath)
	l, err := rss.LoadList(filesystem)
	if err != nil {
		m.UpdateStatus(err.Error())
		return nil
	}

	m.l = l
	m.tabs = l.Categories()

	if len(m.tabs) == 0 {
		m.activeTab = 0
	} else if m.activeTab > len(m.tabs)-1 {
		m.activeTab = len(m.tabs) - 1
	}

	m.UpdateStatus("URLs file edited")

	return rebuildFeedList(m)
}

func handleNextTab(m *model) tea.Cmd {
	m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
	m.lf.ResetSelected()
	return rebuildFeedList(m)
}

func handlePrevTab(m *model) tea.Cmd {
	m.activeTab = max(m.activeTab-1, 0)
	m.lf.ResetSelected()
	return rebuildFeedList(m)
}

func handleToggleRead(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		i.item.ToggleRead()
		rebuildItemsList(m)
		if i.item.Read {
			m.UpdateStatus(MsgMarkItemRead)
		} else {
			m.UpdateStatus(MsgMarkItemUnread)
		}
	}
	return nil
}

func handleToggleBookmark(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if !ok {
		return nil
	}

	i.item.ToggleBookmark()
	err := m.l.SetBookmark(i.item.Bookmark, i.item)
	if err != nil {
		m.UpdateStatus(err.Error())
		return nil
	}

	if i.item.Bookmark {
		m.UpdateStatus(MsgBookmarkAdded)
	} else {
		m.UpdateStatus(MsgBookmarkRemoved)
	}

	rebuildItemsList(m)
	return nil
}

func handleViewBookmarks(m *model) tea.Cmd {
	if bookmarks := m.l.Bookmarks(); bookmarks != nil {
		m.f = bookmarks
		m.title = "Bookmarks"
		m.i = nil

		rebuildItemsList(m)

		if len(bookmarks.RssItems) > 0 {
			m.li.ResetSelected()
		}
	}

	return nil
}

func handleNextUnreadItem(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		prev := i.item
		index, next := m.f.NextUnreadItem(prev)
		if next != nil {
			m.li.Select(index)
		}
	}
	if m.i != nil {
		handleViewItem(m)
	}
	return nil
}

func handleMarkFeedRead(m *model) tea.Cmd {
	if i, ok := m.lf.SelectedItem().(feedItem); ok {
		f := i.rssFeed
		f.MarkAllItemsRead()
		rebuildFeedList(m)
		m.UpdateStatus(MsgMarkFeedRead)
	}
	return nil
}

func handleMarkItemsRead(m *model) tea.Cmd {
	if m.f != nil {
		m.f.MarkAllItemsRead()
		rebuildItemsList(m)
		m.UpdateStatus(MsgMarkFeedRead)
	}
	return nil
}

func handleMarkTabAsRead(m *model) tea.Cmd {
	feeds, err := m.l.GetCategory(activeTab(m.tabs, m.activeTab))
	if err != nil {
		m.UpdateStatus(err.Error())
	}

	rss.MarkFeedsAsRead(feeds...)
	rebuildFeedList(m)
	m.UpdateStatus(MsgMakrTabAsRead)

	return nil
}

func handleBack(m *model) tea.Cmd {
	if m.i != nil {
		m.i = nil
	} else {
		m.lf.ResetFilter()
		m.li.ResetFilter()
		rebuildFeedList(m)
		m.f = nil
	}
	return nil
}

func handleOpenFeed(m *model) tea.Cmd {
	i, ok := m.lf.SelectedItem().(feedItem)
	if ok {
		f := i.rssFeed
		url, err := f.Link()
		if err != nil {
			m.UpdateStatus(err.Error())
		}

		err = openInBrowser(url)
		if err != nil {
			m.UpdateStatus(err.Error())
		}
	}
	return nil
}

func handleOpenLatest(m *model) tea.Cmd {
	i, ok := m.lf.SelectedItem().(feedItem)
	if ok {
		f := i.rssFeed

		latest := f.LatestItem()
		if latest == nil {
			return nil
		}

		url := latest.Link()

		err := openInBrowser(url)
		if err != nil {
			m.UpdateStatus(err.Error())
		}

		latest.MarkRead()
		rebuildFeedList(m)
	}
	return nil
}

func handleOpenItem(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		rssItem := i.item
		if rssItem.Item != nil {
			err := openInBrowser(rssItem.Link())
			if err != nil {
				errorMessage := fmt.Sprintf("Error opening item, %q", err)
				m.UpdateStatus(errorMessage)
			}
			rssItem.MarkRead()
			rebuildItemsList(m)
		}
	}
	return nil
}

func handlePrevUnreadItem(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		next := i.item
		index, prev := m.f.PrevUnreadItem(next)
		if prev != nil {
			m.li.Select(index)
		}
	}
	if m.i != nil {
		handleViewItem(m)
	}
	return nil
}

func handleUpdateFeed(m *model) tea.Cmd {
	if len(m.l.Feeds) == 0 {
		m.UpdateStatus(ErrUpdatingFeed)
		return nil
	}

	feed := m.f
	if m.f == nil {
		if i, ok := m.lf.SelectedItem().(feedItem); ok {
			feed = i.rssFeed
		}
	}

	message := fmt.Sprintf("%s %s", MsgUpdatingFeed, feed.Url)
	m.UpdateStatus(message)
	return updateFeedCmd(m, feed)
}

func handleUpdateAllFeeds(m *model) tea.Cmd {
	m.UpdateStatus(MsgUpdatingAllFeeds)
	return updateAllFeedsCmd(m)
}

func handleTabUpdate(m *model) tea.Cmd {
	m.UpdateStatus(MsgUpdatingAllFeeds)
	return updateTabFeedsCmd(m)
}

func handleQuit(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}

func handleEnterFeed(m *model) tea.Cmd {
	if i, ok := m.lf.SelectedItem().(feedItem); ok {
		if i.rssFeed.Feed != nil {
			m.f = i.rssFeed
			m.UpdateTitle(i.title)
			m.li.ResetSelected()
			rebuildItemsList(m)
		}
	}
	return nil
}

func handleViewItem(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		m.i = i.item
		if m.i.Item != nil {
			m.v.YOffset = 0
			m.v.SetContent(wordwrap.String(m.i.Content(), 80))
			m.i.MarkRead()
			rebuildItemsList(m)
		}
	}
	return nil
}

func handleInterrupt(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}

func handleTabNumber(m *model, i int) tea.Cmd {
	if i > len(m.tabs)-1 {
		return nil
	}
	m.lf.ResetSelected()
	m.activeTab = i
	return rebuildFeedList(m)
}

func handleItemNumber(m *model, i int) tea.Cmd {
	if i == 0 {
		m.li.ResetSelected()
	}
	return nil
}

func handleEnclosureNumber(m *model, i int) tea.Cmd {
	if len(m.i.Item.Enclosures) < (i + 1) {
		return nil
	}

	selected := m.i.Item.Enclosures[i].URL
	err := openInBrowser(selected)
	if err != nil {
		errorMessage := fmt.Sprintf("Error opening item, %q", err)
		m.UpdateStatus(errorMessage)
	}

	return nil
}

func handleViewNext(m *model) tea.Cmd {
	index, next := m.f.NextAfter(m.i)
	if next != nil {
		m.i = next
		m.li.Select(index)
		m.v.YOffset = 0
		m.v.SetContent(wordwrap.String(next.Content(), 80))
		next.MarkRead()
		rebuildItemsList(m)
	}
	return nil
}

func handleViewPrev(m *model) tea.Cmd {
	index, prev := m.f.PrevBefore(m.i)
	if prev != nil {
		m.i = prev
		m.li.Select(index)
		m.v.YOffset = 0
		m.v.SetContent(wordwrap.String(prev.Content(), 80))
		prev.MarkRead()
		rebuildItemsList(m)
	}
	return nil
}

func handleNextUnreadFeed(m *model) tea.Cmd {
	i, ok := m.lf.SelectedItem().(feedItem)
	if ok {
		prev := i.rssFeed
		feeds, err := m.l.GetCategory(activeTab(m.tabs, m.activeTab))
		if err != nil {
			return nil
		}
		index, next := rss.NextUnreadFeed(feeds, prev)
		if next != nil {
			m.lf.Select(index)
			return nil
		}
		nextUnreadTab(m)
	}
	return nil
}

func nextUnreadTab(m *model) tea.Cmd {
	currentTab := m.activeTab
	for i := currentTab; i < len(m.tabs); i++ {
		category := m.tabs[i]
		feeds, err := m.l.GetCategory(category)
		if err != nil {
			return nil
		}
		for j, feed := range feeds {
			if feed.HasUnread() {
				m.activeTab = i
				rebuildFeedList(m)
				m.lf.Select(j)
				break
			}
		}
		if currentTab != m.activeTab {
			break
		}
	}
	return nil
}

func handlePrevUnreadFeed(m *model) tea.Cmd {
	i, ok := m.lf.SelectedItem().(feedItem)
	if ok {
		next := i.rssFeed
		feeds, err := m.l.GetCategory(activeTab(m.tabs, m.activeTab))
		if err != nil {
			return nil
		}
		index, prev := rss.PrevUnreadFeed(feeds, next)
		if prev != nil {
			m.lf.Select(index)
			return nil
		}
		prevUnreadTab(m)
	}
	return nil
}

func prevUnreadTab(m *model) tea.Cmd {
	currentTab := m.activeTab
	for i := currentTab; i >= 0; i-- {
		category := m.tabs[i]
		feeds, err := m.l.GetCategory(category)
		if err != nil {
			return nil
		}
		for j, feed := range slices.Backward(feeds) {
			if feed.HasUnread() {
				m.activeTab = i
				rebuildFeedList(m)
				m.lf.Select(j)
				break
			}
		}
		if currentTab != m.activeTab {
			break
		}
	}
	return nil
}

func handleViewHelp(m *model) tea.Cmd {
	m.vh.ShowAll = !m.vh.ShowAll
	return nil
}

func handleGoToStart(m *model) tea.Cmd {
	m.v.YOffset = 0
	return nil
}
