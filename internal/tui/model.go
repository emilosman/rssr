package tui

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/emilosman/rssr/internal/rss"
	"github.com/muesli/reflow/wordwrap"
)

type feedItem struct {
	title, desc string
	rssFeed     *rss.RssFeed
}

func (f feedItem) Title() string       { return f.title }
func (f feedItem) Description() string { return f.desc }
func (f feedItem) FilterValue() string { return f.title }

type rssListItem struct {
	title, desc string
	item        *rss.RssItem
}

func (r rssListItem) Title() string       { return r.title }
func (r rssListItem) Description() string { return r.desc }
func (r rssListItem) FilterValue() string { return r.item.FilterContent() }

type model struct {
	prog       *tea.Program
	title      string
	status     string
	clearTimer *time.Timer
	l          *rss.List
	f          *rss.RssFeed
	i          *rss.RssItem
	lf         list.Model
	li         list.Model
	v          viewport.Model
	vh         help.Model
	tabs       []string
	activeTab  int
}

func initialModel() *model {
	urlsFilePath, err := rss.UrlsFilePath()
	if err != nil {
		fmt.Println("Error opening config dir", err)
	}
	filesystem := os.DirFS(urlsFilePath)
	l, err := rss.LoadList(filesystem)
	t := l.Categories()

	df := list.NewDefaultDelegate()
	df.ShortHelpFunc = listShortHelp
	df.FullHelpFunc = listFullHelp

	di := list.NewDefaultDelegate()
	di.ShortHelpFunc = itemsShortHelp
	di.FullHelpFunc = itemsFullHelp

	m := &model{
		l:         l,
		lf:        list.New(nil, df, 0, 0),
		li:        list.New(nil, di, 0, 0),
		tabs:      t,
		activeTab: 0,
		v:         viewport.New(),
		vh:        help.New(),
	}

	rebuildFeedList(m)

	m.lf.DisableQuitKeybindings()
	m.li.DisableQuitKeybindings()
	m.lf.SetShowTitle(false)
	m.li.SetShowTitle(false)
	m.lf.SetShowStatusBar(false)
	m.li.SetShowStatusBar(true)

	if err != nil {
		m.UpdateStatus(err.Error())
	}

	if len(m.lf.Items()) == 0 {
		m.UpdateStatus(MsgNoFeedsInList)
	}

	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case feedUpdatedMsg:
		if msg.Err != nil {
			m.UpdateStatus(fmt.Sprintf("Error updating: %v", msg.Err))
		} else {
			m.UpdateStatus(fmt.Sprintf("Updated %s", msg.Feed.Title()))
		}
		rebuildFeedList(m)
		return m, nil
	case feedsDoneMsg:
		m.UpdateStatus(MsgAllFeedsUpdated)
		return m, nil
	case statusClearMsg:
		m.status = ""
		return m, nil
	case tea.KeyPressMsg:
		var handlers map[string]keyHandler
		lfState := m.lf.FilterState().String()
		liState := m.li.FilterState().String()

		if lfState == "filtering" || liState == "filtering" {
			break
		}

		switch {
		case m.i != nil:
			handlers = viewKeyHandlers
			if i, err := strconv.Atoi(msg.String()); err == nil {
				return m, handleEnclosureNumber(m, i)
			}
		case m.f != nil:
			handlers = itemKeyHandlers
			if i, err := strconv.Atoi(msg.String()); err == nil {
				return m, handleItemNumber(m, i)
			}
		default:
			handlers = feedKeyHandlers
			if i, err := strconv.Atoi(msg.String()); err == nil {
				return m, handleTabNumber(m, i)
			}
		}

		if handler, ok := handlers[msg.String()]; ok {
			return m, handler(m)
		}

	case tea.WindowSizeMsg:
		topBarHeigh := 2
		m.lf.SetSize(msg.Width, msg.Height-topBarHeigh)
		m.li.SetSize(msg.Width, msg.Height-topBarHeigh)

		m.v.SetWidth(msg.Width)
		m.v.SetHeight(msg.Height - topBarHeigh)

		if m.i != nil {
			m.v.SetContent(wordwrap.String(m.i.Content(), 80))
		}
	}

	var cmd tea.Cmd

	switch {
	case m.i != nil:
		m.v, cmd = m.v.Update(msg)
	case m.f != nil:
		m.li, cmd = m.li.Update(msg)
	default:
		m.lf, cmd = m.lf.Update(msg)
	}

	return m, cmd
}
