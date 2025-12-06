package rss

import (
	"html"
	"os"
	"path/filepath"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/microcosm-cc/bluemonday"
)

func CacheFilePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "rssr")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "data.json"), nil
}

func UrlsFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dir, "rssr")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	urlsFile := filepath.Join(appDir, "urls.yaml")
	err = defaultUrlsFile(urlsFile)

	return appDir, err
}

func defaultUrlsFile(urlsFile string) error {
	if _, err := os.Stat(urlsFile); os.IsNotExist(err) {
		f, err := os.Create(urlsFile)
		if err != nil {
			return err
		}
		defer f.Close()

		defaultYAML := ExampleConfigFile

		if _, err := f.WriteString(defaultYAML); err != nil {
			return err
		}
	}
	return nil
}

func clean(input string) string {
	p := bluemonday.StrictPolicy()
	s := p.Sanitize(input)
	s = html.UnescapeString(s)
	s = fixMojibake(s)
	s = normalizeSpaces(s)
	return s
}

func toMarkdown(input string) string {
	markdown, err := htmltomarkdown.ConvertString(input)
	if err != nil {
		return input
	}
	return markdown
}

func fixMojibake(s string) string {
	b := make([]byte, len(s))
	for i, r := range s {
		if r > 255 {
			return s
		}
		b[i] = byte(r)
	}
	return string(b)
}

func normalizeSpaces(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	return strings.Join(strings.Fields(s), " ")
}
