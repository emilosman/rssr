package rss

import (
	"html"
	"os"
	"path/filepath"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/microcosm-cc/bluemonday"
)

func DataFilePath() (string, error) {
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
	err = createDefaultFile(urlsFile, DefaultUrlsFile)

	return appDir, err
}

func ConfigFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dir, "rssr")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	configFile := filepath.Join(appDir, "config.yaml")
	err = createDefaultFile(configFile, DefaultConfigFile)

	return appDir, err
}

func createDefaultFile(filePath, fileContent string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.WriteString(fileContent); err != nil {
			return err
		}
	}
	return nil
}

func clean(input string) string {
	p := bluemonday.StrictPolicy()
	s := p.Sanitize(input)
	s = html.UnescapeString(s)
	s = normalizeSpaces(s)
	s = strings.ToValidUTF8(s, "")
	return s
}

func toMarkdown(input string) string {
	markdown, err := htmltomarkdown.ConvertString(input)
	if err != nil {
		return input
	}
	return markdown
}

func normalizeSpaces(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	return strings.Join(strings.Fields(s), " ")
}
