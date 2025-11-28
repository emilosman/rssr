package rss

import (
	"fmt"
	"net/url"

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
	date := i.Item.PublishedParsed
	link := i.Link()
	desc := i.Description()
	content := i.Item.Content
	enclosuers := i.Enclosures()

	if content != "" {
		desc = ""
	}

	return fmt.Sprintf(
		"%s\n%s\n\n%s\n\n%s\n\n%s\n\n",
		date,
		link,
		desc,
		content,
		enclosuers,
	)
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
	desc := i.Item.Description
	if desc == "" {
		desc = i.Item.Content
	}
	return desc
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

func sanitizeItem(i *gofeed.Item) {
	i.Title = clean(i.Title)
	i.Description = clean(i.Description)
	i.Content = clean(i.Content)
}
