package rss

import "errors"

var (
	ErrFeedHasNoUrl       = errors.New("feed has no URL")
	ErrNoFeedsInList      = errors.New("no feeds in list")
	ErrNoCategoryGiven    = errors.New("no category given")
	ErrNoBookmarkFeed     = errors.New("no bookmark feed found")
	ErrCooldown           = errors.New("5 second cooldown")
	ErrConfigDoesNotExist = "open urls.yaml: file does not exist"
	MsgFeedNotLoaded      = "Feed not loaded yet. Press shift+r"
	DefaultUrlsFile       = `# This file is written in YAML format.
# Each feed must be organized under a category.
# Feeds that are not assigned to a category will NOT appear in the app.
# Formatting Rules:
# - Use two spaces for indentation (no tabs allowed).
# - Follow proper YAML structure (see example below).
# Additional Information:
# - Refer to the README.md file for details on where this configuration file is stored
# on your operating system.
# - After saving and closing this file, the app will automatically load and display your feeds.
# Enjoy!
#
# Example (uncomment lines below to use):
#feeds:
#  - https://emilosman.com/feed
`
	DefaultConfigFile = `# This file is written in YAML format.
	# Below is the default config. Uncomment and change if needed.
	# render_markdown: true
`
)
