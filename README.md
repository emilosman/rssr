# rssr

- **rssr** is a performant, terminal-based RSS reader written in Go, inspired by [Newsboat](https://github.com/newsboat/newsboat)
- It demonstrates idiomatic Go usage, concurrency, YAML-based configuration, and a TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)

Main view  
<img width="1082" height="893" alt="main" src="https://github.com/user-attachments/assets/39ebff9f-6803-4fac-a2a2-d475c5da988c" />


Feed view  
<img width="1251" height="891" alt="feed-view" src="https://github.com/user-attachments/assets/a0933c7e-730b-471e-9201-7b2d5ab85f8c" />


Item view  
<img width="830" height="893" alt="viewport" src="https://github.com/user-attachments/assets/fea95c67-540d-4bb6-99b5-17a61b996caa" />


Bookmarks view  
<img width="435" height="239" alt="bookmarks" src="https://github.com/user-attachments/assets/5cc7d9ca-f59d-4806-b25a-0e2e809a9dd4" />

## Features
- Vim-style navigation for feeds and articles
- Keyboard shortcuts inspired by Newsboat (`?` for help)
- Asynchronous feed fetching using Go routines and channels
- Read/unread tracking for feed items
- Configurable via a simple YAML file

## Installation
```bash
git clone https://github.com/emilosman/rssr.git
cd rssr
go install ./cmd/rssr

rssr
```

## Usage
- Use arrow keys or Vim-style shortcuts to navigate
- Press `?` for the full help menu
- Edit the feed list with your preferred editor (vi by default)
- Mark feeds or items as read/unread

## Configuration (MacOS)
- Config file: `~/Library/Application\ Support/rssr/urls.yaml`
- Cache file: `~/Library/Caches/rssr/data.json`

Example urls.yaml:
```
Tech:
  - https://example.com/tech.rss
  - https://example.com/golang.rss
News:
  - https://example.com/worldnews.rss
```

## Development
- This is a hobby project, exploring Go and terminal UI development
- See the [TODO list](./docs/todo.md) for planned features and improvements

## License
- Licensed under [GPLv3](./LICENSE)
