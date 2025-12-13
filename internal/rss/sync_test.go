package rss

import (
	"log/slog"
	"testing"
)

func TestSync(t *testing.T) {
	t.Run("syncs list", func(t *testing.T) {
		mock := &ListState{
			ApiKey: "mock-api-key",
			ItemIndex: map[string]*ItemState{
				"item-123": {
					Ts:       123456789,
					GUID:     "item-123",
					Read:     false,
					Bookmark: true,
				},
				"item-456": {
					Ts:       123456790,
					GUID:     "item-456",
					Read:     true,
					Bookmark: false,
				},
			},
		}

		data := []byte{}
		server := Server(t, data)

		listState, err := SyncState(server.URL, mock)
		if err != nil {
			t.Errorf("Sync error: %q", err)
		}

		slog.Info("state", "state", listState)
	})
}
