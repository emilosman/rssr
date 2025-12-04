package rss

import (
	"html"
	"os"
	"path/filepath"
	"strings"

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

func ConfigFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dir, "rssr")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	configFile := filepath.Join(appDir, "urls.yaml")
	err = defaultConfigFile(configFile)

	return appDir, err
}

func defaultConfigFile(configFile string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		f, err := os.Create(configFile)
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
