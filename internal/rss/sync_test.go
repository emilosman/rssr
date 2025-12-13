package rss

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestSync(t *testing.T) {
	t.Run("updates list", func(t *testing.T) {
		l := NewListWithDefaults()
		rssItems := []*RssItem{
			{
				Ts:       0,
				Read:     false,
				Bookmark: false,
				Item: &gofeed.Item{
					GUID: "item-123",
				},
			},
			{
				Ts:       0,
				Read:     false,
				Bookmark: false,
				Item: &gofeed.Item{
					GUID: "item-456",
				},
			},
		}

		l.Feeds = []*RssFeed{{
			RssItems: rssItems,
		}}

		listState, err := l.SerializeList()
		if err != nil {
			t.Error(err)
		}

		response := []byte(`
{
  "ApiKey": "secret",
  "ItemIndex": {
    "item-123": {
      "Ts": 1700000000,
      "GUID": "item-123",
      "Read": true,
      "Bookmark": true
    },
    "item-456": {
      "Ts": 1700000000,
      "GUID": "item-123",
      "Read": true,
      "Bookmark": true
    }
  }
}
		`)

		server := Server(t, response)
		newListState, err := SyncState(server.URL, listState)
		if err != nil {
			t.Errorf("Sync error: %q", err)
		}

		if len(newListState.ItemIndex) != len(listState.ItemIndex) {
			t.Error("Wrong number of items parsed")
		}

		for i := range newListState.ItemIndex {
			li := newListState.ItemIndex[i]
			if li.Ts == 0 {
				t.Error("Timestamp not updated")
			}
			if !li.Bookmark {
				t.Error("Bookmark not updated")
			}
			if !li.Read {
				t.Error("Read not updated")
			}
		}
		l.SetListState(newListState)

		for _, rssItem := range l.ItemIndex {
			if rssItem.Ts == 0 {
				t.Error("Timestamp not updated")
			}
			if !rssItem.Bookmark {
				t.Error("Bookmark not updated")
			}
			if !rssItem.Read {
				t.Error("Read not updated")
			}
		}
	})
}
