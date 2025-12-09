package rss

import (
	"fmt"
	"net/url"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/mmcdole/gofeed"
)

type RssItem struct {
	Item     *gofeed.Item
	Bookmark bool
	Read     bool
}

func (i *RssItem) Link() string {
	var raw string
	if i.Item == nil {
		return raw
	}
	if len(i.Item.Enclosures) > 0 {
		raw = i.Item.Enclosures[0].URL
	}
	if i.Item.Link != "" {
		raw = i.Item.Link
	}

	url, err := url.ParseRequestURI(raw)
	if err != nil {
		return ""
	}

	return url.String()
}

func (i *RssItem) Title() string {
	title := i.Item.Title
	if i.Bookmark {
		title = fmt.Sprintf("* %s", title)
	}
	if !i.Read {
		title = fmt.Sprintf("+ %s", title)
	}
	return title
}

func (i *RssItem) FilterContent() string {
	return fmt.Sprintf("%s %s", i.Title(), i.Description())
}

func (i *RssItem) Content() string {
	time := i.Timestamp()
	link := i.Link()
	content := i.Item.Content
	enclosuers := i.Enclosures()

	if content == "" {
		content = i.Description()
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err == nil {
		render, err := r.Render(content)
		if err == nil {
			content = render
		}
	}

	return fmt.Sprintf(
		"%s\n%s\n\n%s\n\n%s\n\n",
		time,
		link,
		content,
		enclosuers,
	)
}

func (i *RssItem) Timestamp() *time.Time {
	if i.Item == nil {
		return nil
	}
	p := i.Item.Published
	u := i.Item.Updated

	var pp, up *time.Time

	if p != "" {
		pp = i.Item.PublishedParsed
	}

	if u != "" {
		up = i.Item.UpdatedParsed
	}

	if pp != nil && up != nil {
		if up.After(*pp) {
			return up
		}
	}

	return pp
}

func (i *RssItem) Enclosures() string {
	var text string
	if len(i.Item.Enclosures) == 0 {
		return text
	}

	text = "Enclosed links:\n"
	for i, enc := range i.Item.Enclosures {
		text += fmt.Sprintf("- [%d] %s\n", i, enc.URL)
	}

	return text
}

func (i *RssItem) Description() string {
	if i.Item.Description != "" {
		return i.Item.Description
	}
	if i.Item.Content != "" {
		return i.Item.Content
	}
	if i.Item.PublishedParsed != nil {
		return i.Item.PublishedParsed.String()
	}
	return ""
}

func (i *RssItem) ToggleRead() {
	i.Read = !i.Read
}

func (i *RssItem) ToggleBookmark() {
	i.Bookmark = !i.Bookmark
}

func (i *RssItem) MarkRead() {
	i.Read = true
}

func sanitizeItem(item *gofeed.Item) {
	item.Title = clean(item.Title)
	item.Description = clean(item.Description)
	item.Content = toMarkdown(item.Content)

	for i, enc := range item.Enclosures {
		u, err := url.ParseRequestURI(enc.URL)
		if err != nil {
			item.Enclosures = append(item.Enclosures[:i], item.Enclosures[i+1:]...)
		} else {
			enc.URL = u.String()
		}
	}
}
