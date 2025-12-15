# rssr

- `rssr` is a performant, terminal-based RSS reader written in Go
- Features YAML based configuration, tabs, `vim`-style navigation, high-speed concurrent updates, and syncing across devices

## Main view
- Feeds lists are organized in tabs, and unread feeds are highlighted
- `shift+e` edits the URLs file
- `shift+r` refreshes all feeds

<img width="1082" height="893" alt="main" src="https://github.com/user-attachments/assets/39ebff9f-6803-4fac-a2a2-d475c5da988c" />

## Feed view
- Feed items can be bookmarked by pressing `c`

<img width="1251" height="891" alt="feed-view" src="https://github.com/user-attachments/assets/a0933c7e-730b-471e-9201-7b2d5ab85f8c" />

## Item view
- HTML content is shown as markdown and highlighted

<img width="830" height="893" alt="viewport" src="https://github.com/user-attachments/assets/fea95c67-540d-4bb6-99b5-17a61b996caa" />

## Bookmarks view
- Shows saved items for later reading
- `shift+b` opens the bookmark view

<img width="435" height="239" alt="bookmarks" src="https://github.com/user-attachments/assets/5cc7d9ca-f59d-4806-b25a-0e2e809a9dd4" />

## Usage help
- Use arrow keys or Vim-style shortcuts to navigate
- Press `?` for full keyboard shortcut help

## Requirements
- [_Go_ should be installed on the system to build and install the app](https://go.dev/dl/)
- Binaries will be available in the near future for more convenient installs

## Installation
```bash
git clone https://github.com/emilosman/rssr.git
cd rssr
go install ./cmd/rssr

rssr
```

## Syncing across devices
- Syncing can be done with the [rssr-sync](https://github.com/emilosman/rssr-sync) server

## Configuration files (MacOS)
- URLs file: `~/Library/Application\ Support/rssr/urls.yaml`
- Cache file: `~/Library/Caches/rssr/data.json`

## Configuration files (Linux)
- URLs file: `~/.config/urls.yanl`
- Config file: `~/.cache/data.json`

## Development
- See the [TODO list](./docs/todo.md) for planned features and improvements

## License
- Licensed under [GPLv3](./LICENSE)
