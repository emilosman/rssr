package rss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ListState struct {
	ApiKey    string
	ItemIndex map[string]*ItemState
}

type ItemState struct {
	Ts       int64
	GUID     string
	Read     bool
	Bookmark bool
}

func (l *List) SerializeList() (*ListState, error) {
	ls := &ListState{
		ApiKey:    "localhost",
		ItemIndex: make(map[string]*ItemState),
	}
	if len(l.Feeds) == 0 {
		return nil, ErrNoFeedsInList
	}

	for _, feed := range l.Feeds {
		for _, item := range feed.RssItems {
			if item.Item != nil {
				ls.ItemIndex[item.GUID()] = &ItemState{
					Ts:       item.Ts,
					GUID:     item.GUID(),
					Read:     item.Read,
					Bookmark: item.Bookmark,
				}
			}
		}
	}

	return ls, nil
}

func (l *List) SyncList() error {
	ls, err := l.SerializeList()
	if err != nil {
		return err
	}

	ls, err = SyncState("http://192.168.1.52:8080", ls)
	//ls, err = SyncState("http://localhost:8080", ls)
	if err != nil {
		return err
	}

	return l.SetListState(ls)
}

func SyncState(url string, ls *ListState) (*ListState, error) {
	body, err := json.Marshal(ls)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	resp, err := http.Post(url+"/sync", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("post error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var merged ListState
	if err := json.NewDecoder(resp.Body).Decode(&merged); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return &merged, nil
}

func (l *List) SetListState(ls *ListState) error {
	for guid, is := range ls.ItemIndex {
		item := l.ItemIndex[guid]
		if item != nil {
			item.Ts = is.Ts
			item.Read = is.Read
			item.Bookmark = is.Bookmark
		}
	}
	return nil
}
